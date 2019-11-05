package queryDocument

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
	"github.com/stretchr/testify/assert"
)

var activityMetadata *activity.Metadata

var settingsjson = `{
	"settings": {
	  "mongoConnection": {
		"id": "e1e890d0-de91-11e9-aef0-13201957902e",
		"name": "mongocon",
		"ref": "github.com/project-flogo/datastore-contrib/mongodb/connection",
		"settings": {
			  "Name": "mongocon",
			  "Description": "",
			  "ConnectionURI": "mongodb://admin:admin@10.102.169.188:27017",
			  "Database": "test"
			}
		},
	"operation": "Find One Document",
	"collectionName": "testcollection"
}
}`
var settingsjson1 = `{
	"settings": {
	  "mongoConnection": {
		"id": "e1e890d0-de91-11e9-aef0-13201957902e",
		"name": "mongocon",
		"ref": "github.com/project-flogo/datastore-contrib/mongodb/connection",
		"settings": {
			  "Name": "mongocon",
			  "Description": "",
			  "ConnectionURI": "mongodb://admin:admin@10.102.169.188:27017",
			  "Database": "test"
			}
		},
	"operation": "Find Many Documents",
	"collectionName": "testcollection"
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
func Test_FindOne(t *testing.T) {
	log.RootLogger().Info("****TEST : Executing Find One start****")
	m := make(map[string]interface{})
	err1 := json.Unmarshal([]byte(settingsjson), &m)
	assert.Nil(t, err1)
	mf := mapper.NewFactory(resolve.GetBasicResolver())

	support.RegisterAlias("connection", "connection", "github.com/project-flogo/datastore-contrib/mongodb/connection")
	fmt.Println("=======Settings========", m["settings"])
	iCtx := test.NewActivityInitContext(m["settings"], mf)
	act, err := New(iCtx)
	assert.Nil(t, err)
	tc := test.NewActivityContext(act.Metadata())
	tc.SetInput("input", `{"jsonDocument":{ "empid" : 1 }}`)

	_, err = act.Eval(tc)
	// Getting outputs
	testOutput := tc.GetOutput("Output")
	jsonOutput, _ := json.Marshal(testOutput)
	log.RootLogger().Info("jsonOutput is : %s", string(jsonOutput))
	log.RootLogger().Info("****TEST : Executing Find One ends****")
	assert.Nil(t, err)
}
func Test_FindAll(t *testing.T) {
	log.RootLogger().Info("****TEST : Executing Find All start****")
	m := make(map[string]interface{})
	err1 := json.Unmarshal([]byte(settingsjson1), &m)
	assert.Nil(t, err1)
	mf := mapper.NewFactory(resolve.GetBasicResolver())
	support.RegisterAlias("connection", "connection", "github.com/project-flogo/datastore-contrib/mongodb/connection")
	iCtx := test.NewActivityInitContext(m["settings"], mf)
	act, err := New(iCtx)
	assert.Nil(t, err)
	tc := test.NewActivityContext(act.Metadata())
	tc.SetInput("input", `{"jsonDocument":{ "location" : "Hyderabad" }}`)

	_, err = act.Eval(tc)
	// Getting outputs
	testOutput := tc.GetOutput("Output")
	jsonOutput, _ := json.Marshal(testOutput)
	log.RootLogger().Info("jsonOutput is : %s", string(jsonOutput))
	log.RootLogger().Info("****TEST : Executing Find All ends****")
	assert.Nil(t, err)
}
func Test_FindMany(t *testing.T) {
	log.RootLogger().Info("****TEST : Executing Find Many start****")
	m := make(map[string]interface{})
	err1 := json.Unmarshal([]byte(settingsjson1), &m)
	assert.Nil(t, err1)
	mf := mapper.NewFactory(resolve.GetBasicResolver())
	support.RegisterAlias("connection", "connection", "github.com/project-flogo/datastore-contrib/mongodb/connection")
	iCtx := test.NewActivityInitContext(m["settings"], mf)
	act, err := New(iCtx)
	assert.Nil(t, err)

	tc := test.NewActivityContext(act.Metadata())
	tc.SetInput("input", `{"jsonDocument":[{ "location" : "Hyderabad" },{"location" : "Chennai" }]}`)
	_, err = act.Eval(tc)
	// Getting outputs
	testOutput := tc.GetOutput("Output")
	jsonOutput, _ := json.Marshal(testOutput)
	log.RootLogger().Info("jsonOutput is : %s", string(jsonOutput))
	log.RootLogger().Info("****TEST : Executing Find Many ends****")
	assert.Nil(t, err)
}
