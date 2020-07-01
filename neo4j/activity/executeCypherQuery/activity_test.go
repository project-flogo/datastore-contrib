package executeCypherQuery

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/mapper"
	"github.com/project-flogo/core/data/resolve"
	"github.com/project-flogo/core/support"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/support/test"
	_ "github.com/project-flogo/datastore-contrib/neo4j/connection"
	"github.com/stretchr/testify/assert"
)

var activityMetadata *activity.Metadata

var settingsRead = `{
	"settings": {
		"connection": {
			"id": "e1e890d0-de91-11e9-aef0-13201957902e",
			"name": "neo4jcon",
			"ref": "github.com/project-flogo/datastore-contrib/neo4j/connection",
			"settings": {
				"name": "neo4jcon",
				"description": "",
				"connectionURI": "bolt://localhost:7687",
				"credType": "None",
				"username": "",
				"password": ""
			}
		},
		"databaseName": "neo4j",
		"accessMode": "Read"
	}
}`

var settingsWrite = `{
	"settings": {
		"connection": {
			"id": "e1e890d0-de91-11e9-aef0-13201957902e",
			"name": "neo4jcon",
			"ref": "github.com/project-flogo/datastore-contrib/neo4j/connection",
			"settings": {
				"name": "neo4jcon",
				"description": "",
				"connectionURI": "bolt://localhost:7687",
				"credType": "None",
				"username": "",
				"password": ""
			}
		},
		"databaseName": "neo4j",
		"accessMode": "Write"
	}
}`

func getActivityMetadata() *activity.Metadata {

	if activityMetadata == nil {
		jsonMetadataBytes, err := ioutil.ReadFile("activity.json")
		if err != nil {
			panic("No Json Metadata found for activity.json path")
		}

		activityMetadata = activity.ToMetadata(string(jsonMetadataBytes))
	}

	return activityMetadata
}
func TestMatchQuery(t *testing.T) {
	log.RootLogger().Info("****TEST : Executing start****")
	m := make(map[string]interface{})
	err1 := json.Unmarshal([]byte(settingsRead), &m)
	assert.Nil(t, err1)
	mf := mapper.NewFactory(resolve.GetBasicResolver())

	support.RegisterAlias("connection", "connection", "github.com/project-flogo/datastore-contrib/neo4j/connection")
	fmt.Println("=======Settings========", m["settings"])
	iCtx := test.NewActivityInitContext(m["settings"], mf)
	act, err := New(iCtx)
	assert.Nil(t, err)
	tc := test.NewActivityContext(act.Metadata())
	//tc.SetInput("cypherQuery", "MATCH (n:Movie) RETURN n LIMIT 25")
	tc.SetInput("cypherQuery", "MATCH (n) RETURN n LIMIT 25")
	//tc.SetInput("cypherQuery", "MATCH (p:Person)-[:ACTED_IN]->(n:Movie) RETURN p LIMIT 25")
	_, err = act.Eval(tc)
	// Getting outputs
	testOutput := tc.GetOutput("response")
	jsonOutput, _ := json.Marshal(testOutput)
	log.RootLogger().Infof("jsonOutput is : %s", string(jsonOutput))
	log.RootLogger().Info("****TEST : Executing ends****")
	assert.Nil(t, err)
}

func TestCreateQuery(t *testing.T) {
	log.RootLogger().Info("****TEST : Executing start****")
	m := make(map[string]interface{})
	err1 := json.Unmarshal([]byte(settingsWrite), &m)
	assert.Nil(t, err1)
	mf := mapper.NewFactory(resolve.GetBasicResolver())

	support.RegisterAlias("connection", "connection", "github.com/project-flogo/datastore-contrib/neo4j/connection")
	fmt.Println("=======Settings========", m["settings"])
	iCtx := test.NewActivityInitContext(m["settings"], mf)
	act, err := New(iCtx)
	assert.Nil(t, err)
	tc := test.NewActivityContext(act.Metadata())
	//tc.SetInput("cypherQuery", "MATCH (n:Movie) RETURN n LIMIT 25")
	tc.SetInput("cypherQuery", "CREATE (n:Item { id: $id, name: $name }) RETURN n.id, n.name")
	tc.SetInput("queryParams", map[string]interface{}{"id": 11, "name": "Neel"})
	_, err = act.Eval(tc)
	// Getting outputs
	testOutput := tc.GetOutput("response")
	jsonOutput, _ := json.Marshal(testOutput)
	log.RootLogger().Infof("jsonOutput is : %s", string(jsonOutput))
	log.RootLogger().Info("****TEST : Executing ends****")
	assert.Nil(t, err)
}

func TestUpdateQuery(t *testing.T) {
	log.RootLogger().Info("****TEST : Executing start****")
	m := make(map[string]interface{})
	err1 := json.Unmarshal([]byte(settingsWrite), &m)
	assert.Nil(t, err1)
	mf := mapper.NewFactory(resolve.GetBasicResolver())

	support.RegisterAlias("connection", "connection", "github.com/project-flogo/datastore-contrib/neo4j/connection")
	fmt.Println("=======Settings========", m["settings"])
	iCtx := test.NewActivityInitContext(m["settings"], mf)
	act, err := New(iCtx)
	assert.Nil(t, err)
	tc := test.NewActivityContext(act.Metadata())
	//tc.SetInput("cypherQuery", "MATCH (n:Movie) RETURN n LIMIT 25")
	tc.SetInput("cypherQuery", "MATCH (p:Person {name: 'Tom Cruise'}) SET p.born = 2020	RETURN p")
	//tc.SetInput("queryParams", map[string]interface{}{"id": 11, "name": "Neel"})
	_, err = act.Eval(tc)
	// Getting outputs
	testOutput := tc.GetOutput("response")
	jsonOutput, _ := json.Marshal(testOutput)
	log.RootLogger().Infof("jsonOutput is : %s", string(jsonOutput))
	log.RootLogger().Info("****TEST : Executing ends****")
	assert.Nil(t, err)
}

func TestDeleteQuery(t *testing.T) {
	log.RootLogger().Info("****TEST : Executing start****")
	m := make(map[string]interface{})
	err1 := json.Unmarshal([]byte(settingsWrite), &m)
	assert.Nil(t, err1)
	mf := mapper.NewFactory(resolve.GetBasicResolver())

	support.RegisterAlias("connection", "connection", "github.com/project-flogo/datastore-contrib/neo4j/connection")
	fmt.Println("=======Settings========", m["settings"])
	iCtx := test.NewActivityInitContext(m["settings"], mf)
	act, err := New(iCtx)
	assert.Nil(t, err)
	tc := test.NewActivityContext(act.Metadata())
	//tc.SetInput("cypherQuery", "MATCH (n:Movie) RETURN n LIMIT 25")
	tc.SetInput("cypherQuery", "MATCH (p:Person {name: 'Jack Nicholson'}) DETACH DELETE p")
	//tc.SetInput("queryParams", map[string]interface{}{"id": 11, "name": "Neel"})
	_, err = act.Eval(tc)
	// Getting outputs
	testOutput := tc.GetOutput("response")
	jsonOutput, _ := json.Marshal(testOutput)
	log.RootLogger().Infof("jsonOutput is : %s", string(jsonOutput))
	log.RootLogger().Info("****TEST : Executing ends****")
	assert.Nil(t, err)
}
