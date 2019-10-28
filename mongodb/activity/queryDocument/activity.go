package queryDocument

import (
	ctx "context"
	"encoding/json"
	"fmt"
	"time"

	mongocon "github.com/project-flogo/datastore-contrib/mongodb/connector/connection"
	// "github.com/TIBCOSoftware/flogo-lib/core/activity"
	// "github.com/TIBCOSoftware/flogo-lib/core/data"
	// "github.com/TIBCOSoftware/flogo-lib/logger"
	//"github.com/TIBCOSoftware/flogo-lib/core/data"
	//"github.com/TIBCOSoftware/flogo-lib/core/data"
	// Need to remove this after testing with latest API
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/support/log"

	//"github.com/project-flogo/core/data/metadata"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// var log = logger.GetLogger("activity-mongodb")

const (
	connectionProp = "mongoConnection"
	ivOperation    = "operation"
	outputProperty = "Output"
	ivInput        = "input"
	ovCount        = "count"
)

func init() {
	fmt.Println("Entered the init")
	err := activity.Register(&Activity{}, New)
	if err != nil {
		fmt.Println(err)
	}
}

// New functioncommon
func New(ctx1 activity.InitContext) (activity.Activity, error) {
	fmt.Println("Entered the new method ")
	settings := &Settings{}
	//fmt.Println(ctx1.Settings())
	//fmt.Println("Before Getting the config")
	settings.Connection, _ = mongocon.GetSharedConfiguration(ctx1.Settings()["mongoConnection"])
	//fmt.Println(ctx1.Settings()["mongoConnection"])
	mcon, _ := settings.Connection.(*mongocon.MongodbSharedConfigManager)
	//	mCtx, _ := ctx.WithTimeout(ctx.Background(), 10*time.Second)
	client := mcon.GetClient()
	// opts := options.Client()

	// // if s.Username != "" && s.Password != "" {
	// // 	opts = opts.SetAuth(options.Credential{
	// // 		Username: s.Username,
	// // 		Password: s.Password,
	// // 	})
	// // }
	// //fmt.Println("Before connect")
	// //fmt.Println(mcon.GetClientConfiguration().ConnectionURI)
	// client, err := mongo.Connect(mCtx, opts.ApplyURI(mcon.GetClientConfiguration().ConnectionURI))

	// if err != nil {
	// 	fmt.Println("Error creating client")
	// 	return nil, err
	// }
	//fmt.Println(client)
	//db := client.Database(s.Database)
	//collection := db.Collection(s.Collection)

	act := &Activity{client: client, operation: ctx1.Settings()["operation"].(string), collectionName: ctx1.Settings()["collName"].(string), database: mcon.GetClientConfiguration().Database}
	return act, nil
}

// Activity is a stub for your Activity implementation
type Activity struct {
	client         *mongo.Client
	operation      string
	collectionName string
	database       string
}

// NewActivity inserts a new activity
// func NewActivity(metadata *activity.Metadata) activity.Activity {
// 	return &Activity{metadata: metadata}
// }
var activityMd = activity.ToMetadata(&Input{}, &Output{})

// Metadata implements activity.Activity.Metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

//GetComplexValue safely get the object value
// func GetComplexValue(complexObject *data.ComplexObject) interface{} {
// 	if complexObject != nil {
// 		return complexObject.Value
// 	}
// 	return nil
// }

//Cleanup method
func (a *Activity) Cleanup() error {

	log.RootLogger().Tracef("cleaning up MongoDB activity")

	ctx, cancel := ctx.WithTimeout(ctx.Background(), 30*time.Second)
	defer cancel()

	return a.client.Disconnect(ctx)
}

// Eval implements activity.Activity.Eval
func (a *Activity) Eval(context activity.Context) (done bool, err error) {
	//log.Info("Executing  mongodb query Activity")
	fmt.Println("======Starting Eval=======")

	method := a.operation
	if method == "" {
		return false, activity.NewError("operation is not configured", "MONGODB-1004", nil)
	}
	collectionName := a.collectionName
	if collectionName == "" {
		return false, activity.NewError("collection Name is not configured", "MONGODB-1004", nil)
	}
	var inputJSON string

	inputVal := &Input{}
	err = context.GetInputObject(inputVal)
	if err != nil {
		return true, nil
	}

	jsonBytes, err := json.Marshal(inputVal.Input["jsonDocument"])
	if err != nil {
		return false, fmt.Errorf("Error reading input json %s", err.Error())
	}
	inputJSON = string(jsonBytes)
	fmt.Println("inputJSON  : ", inputJSON)
	if err != nil {
		return false, fmt.Errorf("Error getting mongodb connection %s", err.Error())
	}

	connectioInfo := a.database
	client := a.client
	db := client.Database(connectioInfo)

	coll := db.Collection(collectionName)
	cntx, cancel := ctx.WithTimeout(ctx.Background(), 60*time.Second)
	defer cancel()
	resp := make(map[string]interface{})
	val := make(map[string]interface{})
	m := make(map[string]interface{})
	var objects []map[string]interface{}

	if inputJSON != "" {
		err = json.Unmarshal(jsonBytes, &m)
		if err != nil {
			fmt.Println("=======Error Parsing Json=====", err)
			return false, err
		}
	}
	if method == "Find One Document" {
		result := coll.FindOne(cntx, m)
		//	result := coll.FindOne(cntx, bson.D{{keyName.(string), keyValue}})
		err := result.Decode(&val)
		if err != nil {
			return false, err
		}
		//log.Debugf("Get Results $#v", result)
		resp["Response"] = val

	} else {
		var err error
		var cursor *mongo.Cursor
		if inputJSON == "" {
			cursor, err = coll.Find(cntx, bson.D{})
		} else {
			cursor, err = coll.Find(cntx, m)
		}
		if err != nil {
			return false, err
		}
		i := 0
		for cursor.Next(cntx) {
			val1 := make(map[string]interface{})
			err := cursor.Decode(&val1)
			if err != nil {
				return false, err
			}
			//err = bson.Unmarshal(bsonraw, &val1)
			if err != nil {
				return false, err
			}
			objects = append(objects, val1)
			//log.Debugf("Get Results $#v", bsonraw.String())
			i++

		}
		resp["Response"] = objects
	}
	//outputComplex := &data.ComplexObject{Metadata: "", Value: resp} // Need to remove this after testing
	context.SetOutput("Output", resp)

	return true, nil
}
