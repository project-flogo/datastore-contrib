package insertDocument

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
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
		"operation": "Insert One Document",
        "databaseName": "trigger",
        "collectionName": "collection",
		"timeout": 0,
		"continueOnErr": false
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
		"operation": "Insert Many Documents",
        "databaseName": "trigger",
        "collectionName": "collection",
		"timeout": 0,
		"continueOnErr": false
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
		"operation": "Insert Many Documents",
        "databaseName": "trigger",
        "collectionName": "collection",
		"timeout": 0,
		"continueOnErr": true
	}
}`

var activityMetadata *activity.Metadata

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

func Test_InsertOne(t *testing.T) {
	log.RootLogger().Info("****TEST : Executing InsertOne start****")
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
	tc.SetInput("data", map[string]string{"_id": "1", "age": "22", "empid": "21", "name": "22OctoberTest"})
	_, err = act.Eval(tc)
	// Getting outputs
	testOPInsertedIDs := tc.GetOutput("insertedId")
	testOPTotalCount := tc.GetOutput("totalCount").(int)
	testOPSuccessCount := tc.GetOutput("successCount").(int)
	testOPFailureCount := tc.GetOutput("failureCount").(int)
	log.RootLogger().Infof("Insert Document output (InsertedIDs) is : %s", testOPInsertedIDs)
	log.RootLogger().Infof("Insert Document output (Total count) is : %d", testOPTotalCount)
	log.RootLogger().Infof("Insert Document output (Success count) is : %d", testOPSuccessCount)
	log.RootLogger().Infof("Insert Document output (Failur count) is : %d", testOPFailureCount)
	if err != nil {
		fmt.Printf("Error in activity eval : %s ", err.Error())
	}
	assert.Nil(t, err)
}

func Test_InsertMany_DontContinueOnError(t *testing.T) {
	log.RootLogger().Info("****TEST : Executing Test_InsertMany_DontContinueOnError start****")
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
	//Dealing with activity's inputs
	var objects []map[string]string
	val1 := make(map[string]string)
	for i := 0; i < 3; i++ {
		val1 = map[string]string{"_id": strconv.Itoa(i), "age": "22", "empid": "21", "name": "22OctoberTest"}
		objects = append(objects, val1)
	}
	tc.SetInput("data", objects)
	_, err = act.Eval(tc)
	// Getting generated output
	testOPInsertedIDs := tc.GetOutput("insertedId")
	testOPTotalCount := tc.GetOutput("totalCount").(int)
	testOPSuccessCount := tc.GetOutput("successCount").(int)
	testOPFailureCount := tc.GetOutput("failureCount").(int)
	log.RootLogger().Infof("Insert Document output (InsertedIDs) is : %s", testOPInsertedIDs)
	log.RootLogger().Infof("Insert Document output (Total count) is : %d", testOPTotalCount)
	log.RootLogger().Infof("Insert Document output (Success count) is : %d", testOPSuccessCount)
	log.RootLogger().Infof("Insert Document output (Failur count) is : %d", testOPFailureCount)
	if err != nil {
		fmt.Printf("Error in activity eval : %s ", err.Error())
	}
	assert.Nil(t, err)
}

func Test_InsertMany_ContinueOnError(t *testing.T) {
	log.RootLogger().Info("****TEST : Executing Test_InsertMany_ContinueOnError start****")
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
	//Dealing with activity's inputs
	var objects []map[string]string
	val1 := make(map[string]string)
	for i := 10; i < 13; i++ {
		val1 = map[string]string{"_id": strconv.Itoa(i), "age": "22", "empid": "21", "name": "22OctoberTest"}
		objects = append(objects, val1)
	}
	tc.SetInput("data", objects)
	_, err = act.Eval(tc)
	// Getting generated output
	testOPInsertedIDs := tc.GetOutput("insertedId")
	testOPTotalCount := tc.GetOutput("totalCount").(int)
	testOPSuccessCount := tc.GetOutput("successCount").(int)
	testOPFailureCount := tc.GetOutput("failureCount").(int)
	log.RootLogger().Infof("Insert Document output (InsertedIDs) is : %s", testOPInsertedIDs)
	log.RootLogger().Infof("Insert Document output (Total count) is : %d", testOPTotalCount)
	log.RootLogger().Infof("Insert Document output (Success count) is : %d", testOPSuccessCount)
	log.RootLogger().Infof("Insert Document output (Failur count) is : %d", testOPFailureCount)
	if err != nil {
		fmt.Printf("Error in activity eval : %s ", err.Error())
	}
	assert.Nil(t, err)
}
