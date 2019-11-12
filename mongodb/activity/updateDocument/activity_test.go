package updateDocument

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
	_ "github.com/project-flogo/datastore-contrib/mongodb/connection"
	"github.com/stretchr/testify/assert"
)

var activityMetadata *activity.Metadata

var settingsjson = `{
	"settings": {
		"connection": {
			"id": "e1e890d0-de91-11e9-aef0-13201957902e",
			"name": "mongocon",
			"ref": "github.com/project-flogo/datastore-contrib/mongodb/connection",
			"settings": {
				"name": "mongocon",
				"description": "",
				"connectionURI": "mongodb://admin:admin@10.102.169.188:27017",
				"credType": "None",
				"ssl": false
			}
		},
		"operation": "Update One Document",
        "databaseName": "sample",
        "collectionName": "test",
        "timeout": 0
	}
}`
var settingsjson1 = `{
	"settings": {
		"connection": {
			"id": "e1e890d0-de91-11e9-aef0-13201957902e",
			"name": "mongocon",
			"ref": "github.com/project-flogo/datastore-contrib/mongodb/connection",
			"settings": {
				"name": "mongocon",
				"description": "",
				"connectionURI": "mongodb://admin:admin@10.102.169.188:27017",
				"credType": "None",
				"ssl": false
			}
		},
		"operation": "Update Many Documents",
        "databaseName": "sample",
		"collectionName": "deletetesting",
        "timeout": 0
	}
}`
var settingsjson2 = `{
	"settings": {
		"connection": {
			"id": "e1e890d0-de91-11e9-aef0-13201957902e",
			"name": "mongocon",
			"ref": "github.com/project-flogo/datastore-contrib/mongodb/connection",
			"settings": {
				"name": "mongocon",
				"description": "",
				"connectionURI": "mongodb://admin:admin@10.102.169.188:27017",
				"credType": "None",
				"ssl": false
			}
		},
		"operation": "Replace One Document",
        "databaseName": "sample",
		"collectionName": "test",
        "timeout": 0
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
func Test_UpdateOne(t *testing.T) {
	log.RootLogger().Info("****TEST : Executing Update start****")
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
	//Setting inputs
	tc.SetInput("criteria", map[string]string{"location": "Palo Alto"})
	tc.SetInput("updateData", map[string]string{"name": "Test Last"})

	_, err = act.Eval(tc)
	// Getting outputs
	testOutputMatchedCount := tc.GetOutput("matchedCount").(int64)
	testOutputModifiedCount := tc.GetOutput("updatedCount").(int64)
	log.RootLogger().Infof("Update Document output (Matched count) is : %d", testOutputMatchedCount)
	log.RootLogger().Infof("Update Document output (Updated count) is : %d", testOutputModifiedCount)
	assert.Nil(t, err)
}
func Test_UpdateManyPass(t *testing.T) {
	log.RootLogger().Info("****TEST : Executing Update start****")
	m := make(map[string]interface{})
	err1 := json.Unmarshal([]byte(settingsjson1), &m)
	assert.Nil(t, err1)
	mf := mapper.NewFactory(resolve.GetBasicResolver())
	support.RegisterAlias("connection", "connection", "github.com/project-flogo/datastore-contrib/mongodb/connection")
	fmt.Println("=======Settings========", m["settings"])
	iCtx := test.NewActivityInitContext(m["settings"], mf)
	act, err := New(iCtx)
	assert.Nil(t, err)
	tc := test.NewActivityContext(act.Metadata())
	//Setting inputs
	tc.SetInput("criteria", map[string]string{})
	tc.SetInput("updateData", map[string]string{"name": "Common Name"})

	_, err = act.Eval(tc)
	// Getting outputs
	testOutputMatchedCount := tc.GetOutput("matchedCount").(int64)
	testOutputModifiedCount := tc.GetOutput("updatedCount").(int64)
	log.RootLogger().Infof("Update Document output (Matched count) is : %d", testOutputMatchedCount)
	log.RootLogger().Infof("Update Document output (Updated count) is : %d", testOutputModifiedCount)
	assert.Nil(t, err)
}

func Test_ReplacePass(t *testing.T) {
	log.RootLogger().Info("****TEST : Executing Update start****")
	m := make(map[string]interface{})
	err1 := json.Unmarshal([]byte(settingsjson2), &m)
	assert.Nil(t, err1)
	mf := mapper.NewFactory(resolve.GetBasicResolver())
	support.RegisterAlias("connection", "connection", "github.com/project-flogo/datastore-contrib/mongodb/connection")
	fmt.Println("=======Settings========", m["settings"])
	iCtx := test.NewActivityInitContext(m["settings"], mf)
	act, err := New(iCtx)
	assert.Nil(t, err)
	tc := test.NewActivityContext(act.Metadata())
	//Setting inputs
	tc.SetInput("criteria", map[string]string{"location": "San Jose"})
	tc.SetInput("updateData", map[string]string{"location": "Seoul"})
	_, err = act.Eval(tc)
	// Getting outputs
	testOutputMatchedCount := tc.GetOutput("matchedCount").(int64)
	testOutputModifiedCount := tc.GetOutput("updatedCount").(int64)
	log.RootLogger().Infof("Update Document output (Matched count) is : %d", testOutputMatchedCount)
	log.RootLogger().Infof("Update Document output (Updated count) is : %d", testOutputModifiedCount)
	assert.Nil(t, err)
}
