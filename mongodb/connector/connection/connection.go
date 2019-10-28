package mongodb

import (
	"context"
	"fmt"
	"os"
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

type Settings struct {
	Name          string `md:"Name,required"`
	Description   string `md:"Description"`
	ConnectionURI string `md:"ConnectionURI"`
	Database      string `md:"Database"`
	DocsMetadata  string `md:"DocsMetadata"`
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

func (k *MongodbSharedConfigManager) GetConnection() interface{} {
	return k
}
func (k *MongodbSharedConfigManager) GetClient() *mongo.Client {
	return k.mclient
}

func (k *MongodbSharedConfigManager) GetClientConfiguration() *Settings {
	return k.config
}

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
