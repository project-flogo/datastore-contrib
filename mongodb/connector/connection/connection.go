package mongodbconnection

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	b64 "encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/connection"
	"github.com/project-flogo/core/support/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var logCache = log.ChildLogger(log.RootLogger(), "mongodb-connection")
var factory = &mongodbFactory{}

// Settings struct
type Settings struct {
	Name          string `md:"Name,required"`
	Description   string `md:"Description"`
	ConnectionURI string `md:"ConnectionURI"`
	Database      string `md:"Database"`
	DocsMetadata  string `md:"DocsMetadata"`
	CredType      string `md:"CredType"`
	UserName      string `md:"UserName"`
	Password      string `md:"Password"`
	Ssl           bool   `md:"Ssl"`
	TrustCert     string `md:"TrustCert"`
	ClientKey     string `md:"ClientKey"`
	ClientCert    string `md:"ClientCert"`
	KeyPass       string `md:"KeyPass"`
	X509          bool   `md:"X509"`
}

// type MongodbClientConfig struct {
// 	Database    string
// 	MongoClient *mongo.Client
// }

func init() {
	if os.Getenv(log.EnvKeyLogLevel) == "DEBUG" {
		// Enable debug logs for sarama lib
		// sarama.Logger = dlog.New(os.Stderr, "[flogo-mongodb]", dlog.LstdFlags)
		// todo
	}
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
		fmt.Println("returning cache connection===")
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
			//	panic("failed to parse root certificate")
			fmt.Println("failed to parse root certificate")
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
			//	fmt.Println("ClientKey", ClientKey)
			//	fmt.Println("ClientCert", ClientCert)
			fmt.Println("KeyPass", KeyPass)

			keyPEMBlock := parseCert(ClientKey)
			certPEMBlock := parseCert(ClientCert)
			var cert tls.Certificate
			if KeyPass != "" || len(KeyPass) != 0 {
				var pkey []byte
				v, _ := pem.Decode([]byte(keyPEMBlock))
				if v == nil {
					fmt.Println("not able to decode key")
				}
				if v.Type == "RSA PRIVATE KEY" {
					fmt.Println("key type is rsa private")

					if x509.IsEncryptedPEMBlock(v) {
						fmt.Println("IsEncryptedPEMBlock")
						pkey, _ = x509.DecryptPEMBlock(v, []byte(KeyPass))
						pkey = pem.EncodeToMemory(&pem.Block{
							Type:  v.Type,
							Bytes: pkey,
						})
					} else {
						fmt.Println("else IsEncryptedPEMBlock")
						pkey = pem.EncodeToMemory(v)
					}
				}
				cert, err = tls.X509KeyPair([]byte(certPEMBlock), pkey)
			} else {
				cert, err = tls.X509KeyPair([]byte(certPEMBlock), []byte(keyPEMBlock))
			}
			if err != nil {
				//	log.Fatal(err)
				fmt.Println("error", err)
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
	//	fmt.Println("url====", url)
	client, err := mongo.NewClient(opts.ApplyURI(url))
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		fmt.Println("===connect error==", err)
		return nil, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Println("===ping error==", err)
	} else {
		fmt.Println("Ping success")
		sharedConn.mclient = client
		fmt.Println("returning new connection===")
	}
	return sharedConn, nil
}

type MongodbSharedConfigManager struct {
	config  *Settings
	name    string
	mclient *mongo.Client
}

func (k *MongodbSharedConfigManager) Type() string {
	return "mongodb"
}

// GetConnection ss
func (k *MongodbSharedConfigManager) GetConnection() interface{} {
	return k
}

// GetClient type
func (k *MongodbSharedConfigManager) GetClient() *mongo.Client {
	return k.mclient
}

func (k *MongodbSharedConfigManager) GetClientConfiguration() *Settings {
	return k.config
}

// ReleaseConnection ss
func (k *MongodbSharedConfigManager) ReleaseConnection(connection interface{}) {

}

func (k *MongodbSharedConfigManager) Start() error {
	return nil
}

func (k *MongodbSharedConfigManager) Stop() error {
	logCache.Debug("Cleaning up client cache")
	//	k.config.MongoClient.Disconnect(context.Background())

	return nil
}

func GetSharedConfiguration(conn interface{}) (connection.Manager, error) {
	var cManager connection.Manager
	var err error
	//	_, ok := conn.(map[string]interface{})
	// if ok {
	// 	//	cManager, err = handleLegacyConnection(conn)
	// } else {
	cManager, err = coerce.ToConnection(conn)
	//	}

	if err != nil {
		return nil, err
	}
	return cManager, nil
}

// func handleLegacyConnection(conn interface{}) (connection.Manager, error) {

// 	connectionObject, _ := coerce.ToObject(conn)
// 	if connectionObject == nil {
// 		return nil, errors.New("Connection object is nil")
// 	}

// 	id := connectionObject["id"].(string)

// 	cManager := connection.GetManager(id)
// 	if cManager == nil {

// 		connObject, err := generic.NewConnection(connectionObject)
// 		if err != nil {
// 			return nil, err
// 		}
// 		cManager, err = factory.NewManager(connObject.Settings())
// 		if err != nil {
// 			return nil, err
// 		}

// 		err = connection.RegisterManager(id, cManager)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	return cManager, nil

// }

func getmongodbClientConfig(settings map[string]interface{}) (*Settings, error) {
	connectionConfig := &Settings{}

	s := &Settings{}

	err := metadata.MapToStruct(settings, s, false)

	if err != nil {
		return nil, err
	}

	connectionConfig = s
	return connectionConfig, nil
}

// parse cert

func parseCert(cert string) string {
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(cert), &m)
	if err != nil {
		fmt.Println("=======Error Parsing Json=====", err)

	}
	//	fmt.Println(m["content"])
	content := m["content"].(string)
	lastBin := strings.LastIndex(content, "base64,")
	sEnc := content[lastBin+7 : len(content)]
	//	fmt.Printf("substring is: %v\n", sEnc)
	sDec, _ := b64.StdEncoding.DecodeString(sEnc)
	//	fmt.Println(string(sDec))
	return (string(sDec))
}
