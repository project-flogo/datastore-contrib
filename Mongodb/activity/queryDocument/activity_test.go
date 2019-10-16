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
		"ref": "git.tibco.com/git/product/ipaas/wi-mongodb.git/src/app/Mongodb/connector/connection",
		"settings": {
			  "Name": "mongocon",
			  "Description": "",
			  "ConnectionURI": "mongodb://admin:admin@10.102.169.188:27017",
			  "Database": "test"
			}
		},
	"operation": "Find One Document",
	"collName": "testcollection"
}
}`
var settingsjson1 = `{
	"settings": {
	  "mongoConnection": {
		"id": "e1e890d0-de91-11e9-aef0-13201957902e",
		"name": "mongocon",
		"ref": "git.tibco.com/git/product/ipaas/wi-mongodb.git/src/app/Mongodb/connector/connection",
		"settings": {
			  "Name": "mongocon",
			  "Description": "",
			  "ConnectionURI": "mongodb://admin:admin@10.102.169.188:27017",
			  "Database": "test"
			}
		},
	"operation": "Find Many Documents",
	"collName": "testcollection"
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

	support.RegisterAlias("connection", "connection", "git.tibco.com/git/product/ipaas/wi-mongodb.git/src/app/Mongodb/connector/connection")
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
	log.RootLogger().Info("****TEST : Executing Create folder test for testing conflict behavior replace ends****")
	assert.Nil(t, err)
}
func Test_FindAll(t *testing.T) {
	log.RootLogger().Info("****TEST : Executing Find One start****")
	m := make(map[string]interface{})
	err1 := json.Unmarshal([]byte(settingsjson1), &m)
	assert.Nil(t, err1)
	mf := mapper.NewFactory(resolve.GetBasicResolver())
	support.RegisterAlias("connection", "connection", "git.tibco.com/git/product/ipaas/wi-mongodb.git/src/app/Mongodb/connector/connection")
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
	log.RootLogger().Info("****TEST : Executing Create folder test for testing conflict behavior replace ends****")
	assert.Nil(t, err)
}
func Test_FindMany(t *testing.T) {
	log.RootLogger().Info("****TEST : Executing Find One start****")
	m := make(map[string]interface{})
	err1 := json.Unmarshal([]byte(settingsjson1), &m)
	assert.Nil(t, err1)
	mf := mapper.NewFactory(resolve.GetBasicResolver())
	support.RegisterAlias("connection", "connection", "git.tibco.com/git/product/ipaas/wi-mongodb.git/src/app/Mongodb/connector/connection")
	iCtx := test.NewActivityInitContext(m["settings"], mf)
	act, err := New(iCtx)
	assert.Nil(t, err)

	tc := test.NewActivityContext(act.Metadata())
	tc.SetInput("input", `{"jsonDocument":[{ "empid" : 1 },{"empid" : 2 }]}`)
	_, err = act.Eval(tc)
	// Getting outputs
	testOutput := tc.GetOutput("Output")
	jsonOutput, _ := json.Marshal(testOutput)
	log.RootLogger().Info("jsonOutput is : %s", string(jsonOutput))
	log.RootLogger().Info("****TEST : Executing Create folder test for testing conflict behavior replace ends****")
	assert.Nil(t, err)
}
