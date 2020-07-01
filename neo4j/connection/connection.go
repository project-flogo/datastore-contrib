package neo4jconnection

import (
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/connection"
	"github.com/project-flogo/core/support/log"
)

var logneo4jconn = log.ChildLogger(log.RootLogger(), "neo4j-connection")
var factory = &neo4jFactory{}

// Settings struct
type Settings struct {
	Name          string `md:"name,required"`
	Description   string `md:"description"`
	ConnectionURI string `md:"connectionURI,required"`
	CredType      string `md:"credType,required"`
	UserName      string `md:"username"`
	Password      string `md:"password"`
}

func init() {
	err := connection.RegisterManagerFactory(factory)
	if err != nil {
		panic(err)
	}
}

type neo4jFactory struct {
}

func (*neo4jFactory) Type() string {
	return "neo4j"
}

func (*neo4jFactory) NewManager(settings map[string]interface{}) (connection.Manager, error) {
	sharedConn := &Neo4jSharedConfigManager{}
	var err error
	sharedConn.config, err = getNeo4jClientConfig(settings)
	if err != nil {
		return nil, err
	}
	if sharedConn.driver != nil {
		return sharedConn, nil
	}

	url := sharedConn.config.ConnectionURI
	credType := sharedConn.config.CredType
	username := sharedConn.config.UserName
	password := sharedConn.config.Password

	configForNeo4j40 := func(conf *neo4j.Config) { conf.Encrypted = false }

	auth := neo4j.NoAuth()
	if credType != "None" {
		auth = neo4j.BasicAuth(username, password, "")
	}

	driver, err := neo4j.NewDriver(url, auth, configForNeo4j40)
	if err != nil {
		logneo4jconn.Errorf("===driver error==", err)
		return nil, err
	}

	sharedConn.driver = driver
	logneo4jconn.Debugf("Returning neo4j connection")
	return sharedConn, nil
}

// Neo4jSharedConfigManager Structure
type Neo4jSharedConfigManager struct {
	config *Settings
	name   string
	driver neo4j.Driver
}

// Type of SharedConfigManager
func (k *Neo4jSharedConfigManager) Type() string {
	return "neo4j"
}

// GetConnection ss
func (k *Neo4jSharedConfigManager) GetConnection() interface{} {
	return k.driver
}

// ReleaseConnection ss
func (k *Neo4jSharedConfigManager) ReleaseConnection(connection interface{}) {

}

// Start connection manager
func (k *Neo4jSharedConfigManager) Start() error {
	return nil
}

// Stop connection manager
func (k *Neo4jSharedConfigManager) Stop() error {
	logneo4jconn.Debug("Cleaning up client connection cache")
	k.driver.Close()
	return nil
}

// GetSharedConfiguration function to return Neo4j connection manager
func GetSharedConfiguration(conn interface{}) (connection.Manager, error) {
	var cManager connection.Manager
	var err error
	cManager, err = coerce.ToConnection(conn)
	if err != nil {
		return nil, err
	}
	return cManager, nil
}

func getNeo4jClientConfig(settings map[string]interface{}) (*Settings, error) {
	connectionConfig := &Settings{}
	err := metadata.MapToStruct(settings, connectionConfig, false)
	if err != nil {
		return nil, err
	}
	return connectionConfig, nil
}
