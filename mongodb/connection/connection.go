package mongodbconnection

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	b64 "encoding/base64"
	"encoding/json"
	"encoding/pem"
	"strings"
	"time"

	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/connection"
	"github.com/project-flogo/core/support/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var logmongoconn = log.ChildLogger(log.RootLogger(), "mongodb-connection")
var factory = &mongodbFactory{}

// Settings struct
type Settings struct {
	Name          string `md:"name,required"`
	Description   string `md:"description"`
	ConnectionURI string `md:"connectionURI"`
	DocsMetadata  string `md:"DocsMetadata"`
	CredType      string `md:"credType"`
	UserName      string `md:"username"`
	Password      string `md:"password"`
	Ssl           bool   `md:"ssl"`
	TrustCert     string `md:"trustCert"`
	ClientKey     string `md:"clientKey"`
	ClientCert    string `md:"clientCert"`
	KeyPass       string `md:"keyPassword"`
	X509          bool   `md:"x509"`
}

func init() {
	err := connection.RegisterManagerFactory(factory)
	if err != nil {
		panic(err)
	}
}

type mongodbFactory struct {
}

func (*mongodbFactory) Type() string {
	return "mongodb"
}

func (*mongodbFactory) NewManager(settings map[string]interface{}) (connection.Manager, error) {
	sharedConn := &MongodbSharedConfigManager{}
	var err error
	sharedConn.config, err = getmongodbClientConfig(settings)
	if err != nil {
		return nil, err
	}
	if sharedConn.mclient != nil {
		return sharedConn, nil
	}
	opts := options.Client()

	url := sharedConn.config.ConnectionURI
	credType := sharedConn.config.CredType
	ssl := sharedConn.config.Ssl

	if credType != "None" {
		userName := sharedConn.config.UserName
		password := sharedConn.config.Password
		opts.SetAuth(options.Credential{
			AuthMechanism: credType,
			Username:      userName,
			Password:      password,
		})
	}
	//ssl
	if ssl {
		var tlsConfig *tls.Config
		trustCert := sharedConn.config.TrustCert
		rootCert := parseCert(trustCert)
		roots := x509.NewCertPool()
		ok := roots.AppendCertsFromPEM([]byte(rootCert))
		if !ok {
			logmongoconn.Errorf("Failed to parse root certificate for SSL")
			//	return nil, err
		}
		ClientKey := sharedConn.config.ClientKey
		ClientCert := sharedConn.config.ClientCert
		if ClientKey == "" || len(ClientKey) == 0 || ClientCert == "" || len(ClientCert) == 0 {
			tlsConfig = &tls.Config{
				RootCAs: roots,
			}
		} else {
			KeyPass := sharedConn.config.KeyPass //TODO need to check how password protected keys will be handled by platform
			keyPEMBlock := parseCert(ClientKey)
			certPEMBlock := parseCert(ClientCert)
			var cert tls.Certificate
			if KeyPass != "" || len(KeyPass) != 0 {
				var pkey []byte
				v, _ := pem.Decode([]byte(keyPEMBlock))
				if v == nil {
					logmongoconn.Warnf("Not able to decode client key")
				}
				if v.Type == "RSA PRIVATE KEY" {
					if x509.IsEncryptedPEMBlock(v) {
						pkey, _ = x509.DecryptPEMBlock(v, []byte(KeyPass))
						pkey = pem.EncodeToMemory(&pem.Block{
							Type:  v.Type,
							Bytes: pkey,
						})
					} else {
						pkey = pem.EncodeToMemory(v)
					}
				}
				cert, err = tls.X509KeyPair([]byte(certPEMBlock), pkey)
			} else {
				cert, err = tls.X509KeyPair([]byte(certPEMBlock), []byte(keyPEMBlock))
			}
			if err != nil {
				logmongoconn.Errorf("Error while creating client certs for establishing 2 way SSL", err)
				return nil, err
			}
			tlsConfig = &tls.Config{
				RootCAs: roots,
				//	ClientAuth: tls.RequireAndVerifyClientCert,
				//	ClientCAs:  clients,
				Certificates: []tls.Certificate{cert},
			}
			// x-509 implementation
			x509 := sharedConn.config.X509
			if x509 {
				opts.SetAuth(options.Credential{
					AuthMechanism: "MONGODB-X509",
					AuthSource:    "$external",
				})
			}
		}
		opts.SetTLSConfig(tlsConfig)

	}
	client, err := mongo.NewClient(opts.ApplyURI(url))
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		logmongoconn.Errorf("===connect error==", err)
		return nil, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		logmongoconn.Errorf("===ping error===", err)
		return nil, err
	} else {
		logmongoconn.Debugf("===Ping success===")
		sharedConn.mclient = client
		logmongoconn.Debugf("Returning new mongodb connection")
	}
	return sharedConn, nil
}

// MongodbSharedConfigManager Structure
type MongodbSharedConfigManager struct {
	config  *Settings
	name    string
	mclient *mongo.Client
}

// Type of SharedConfigManager
func (k *MongodbSharedConfigManager) Type() string {
	return "mongodb"
}

// GetConnection ss
func (k *MongodbSharedConfigManager) GetConnection() interface{} {
	return k.mclient
}

// ReleaseConnection ss
func (k *MongodbSharedConfigManager) ReleaseConnection(connection interface{}) {

}

// Start connection manager
func (k *MongodbSharedConfigManager) Start() error {
	return nil
}

// Stop connection manager
func (k *MongodbSharedConfigManager) Stop() error {
	logmongoconn.Debug("Cleaning up client connection cache")
	return nil
}

// GetSharedConfiguration function to return MongoDB connection manager
func GetSharedConfiguration(conn interface{}) (connection.Manager, error) {
	var cManager connection.Manager
	var err error
	cManager, err = coerce.ToConnection(conn)
	if err != nil {
		return nil, err
	}
	return cManager, nil
}

func getmongodbClientConfig(settings map[string]interface{}) (*Settings, error) {
	connectionConfig := &Settings{}
	err := metadata.MapToStruct(settings, connectionConfig, false)
	if err != nil {
		return nil, err
	}
	return connectionConfig, nil
}

// parse cert

func parseCert(cert string) string {
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(cert), &m)
	if err != nil {
		logmongoconn.Errorf("=======Error Parsing Certificate for SSL handshake=====", err)
	}
	content := m["content"].(string)
	lastBin := strings.LastIndex(content, "base64,")
	sEnc := content[lastBin+7 : len(content)]
	sDec, _ := b64.StdEncoding.DecodeString(sEnc)
	return (string(sDec))
}
