package queryDocument

import (
	ctx "context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
	mongocon "github.com/project-flogo/datastore-contrib/mongodb/connection"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var logquery = log.ChildLogger(log.RootLogger(), "mongodb-querydocument")

func init() {
	err := activity.Register(&Activity{}, New)
	if err != nil {
		fmt.Println(err)
	}
}

// New functioncommon
func New(ctx activity.InitContext) (activity.Activity, error) {
	settings := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), settings, true)
	if err != nil {
		return nil, err
	}
	if settings.Connection != "" {

		mcon, toConnerr := coerce.ToConnection(settings.Connection)
		if toConnerr != nil {
			return nil, toConnerr
		}
		client := mcon.(*mongocon.MongodbSharedConfigManager).GetClient()
		act := &Activity{client: client, operation: settings.Operation, collectionName: settings.CollectionName, database: mcon.(*mongocon.MongodbSharedConfigManager).GetClientConfiguration().Database}
		return act, nil
	}
	return nil, nil
}

// Activity is a stub for your Activity implementation
type Activity struct {
	client         *mongo.Client
	operation      string
	collectionName string
	database       string
}

var activityMd = activity.ToMetadata(&Input{}, &Output{})

// Metadata implements activity.Activity.Metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

//Cleanup method
func (a *Activity) Cleanup() error {

	logquery.Debugf("cleaning up MongoDB activity")

	ctx, cancel := ctx.WithTimeout(ctx.Background(), 30*time.Second)
	defer cancel()

	return a.client.Disconnect(ctx)
}

// Eval implements activity.Activity.Eval
func (a *Activity) Eval(context activity.Context) (done bool, err error) {
	logquery.Debugf("Executing  mongodb query Activity")
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
	logquery.Debugf("InputJSON: ", inputJSON)
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
			return false, err
		}
	}
	if method == "Find One Document" {
		result := coll.FindOne(cntx, m)
		err := result.Decode(&val)
		if err != nil {
			return false, err
		}
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
		for cursor.Next(cntx) {
			val1 := make(map[string]interface{})
			err := cursor.Decode(&val1)
			if err != nil {
				return false, err
			}
			objects = append(objects, val1)
		}
		resp["Response"] = objects
	}
	//outputComplex := &data.ComplexObject{Metadata: "", Value: resp} // Need to remove this after testing
	context.SetOutput("output", resp)

	return true, nil
}
