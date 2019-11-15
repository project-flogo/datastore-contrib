package removeDocument

import (
	ctx "context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var logRemove = log.ChildLogger(log.RootLogger(), "mongodb-removeDocument")

func init() {
	err := activity.Register(&Activity{}, New)
	if err != nil {
		logRemove.Errorf("MongoDB Remove Document Activity init error : %s ", err.Error())
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
		client := mcon.GetConnection().(*mongo.Client)
		act := &Activity{client: client, operation: settings.Operation, collectionName: settings.CollectionName,
			database: settings.Database, timeout: settings.Timeout}
		return act, nil
	}

	return nil, nil
}

var activityMd = activity.ToMetadata(&Input{}, &Output{})

// Metadata implements activity.Activity.Metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

func (a *Activity) Cleanup() error {
	logRemove.Debugf("cleaning up MongoDB Remove Activity")
	ctx, cancel := ctx.WithTimeout(ctx.Background(), 30*time.Second)
	defer cancel()
	return a.client.Disconnect(ctx)
}

// Activity is a stub for your Activity implementation
type Activity struct {
	client         *mongo.Client
	operation      string
	collectionName string
	database       string
	timeout        int32
}

// Eval implements activity.Activity.Eval
func (a *Activity) Eval(context activity.Context) (done bool, err error) {
	logRemove.Debugf("Executing  MongoDB Remove Document Activity")
	if a.operation == "" {
		return false, activity.NewError("Remove Document Operation is not configured", "MONGODB-1004", nil)
	}
	if a.collectionName == "" {
		return false, activity.NewError("Collection Name is not configured for Remove Document activity", "MONGODB-1004", nil)
	}
	m := make(map[string]interface{})
	var result *mongo.DeleteResult
	var error error
	var inputJSON string

	inputVal := &Input{}
	err = context.GetInputObject(inputVal)
	if err != nil {
		return true, nil
	}
	jsonBytes, err := json.Marshal(inputVal.Input)
	if err != nil {
		return false, fmt.Errorf("Error reading input json %s", err.Error())
	}
	inputJSON = string(jsonBytes)
	logRemove.Debugf("InputJSON: ", inputJSON)

	db := a.client.Database(a.database)
	coll := db.Collection(a.collectionName)
	timeout := a.timeout
	if timeout <= 0 {
		timeout = 60 //set a default timeout of 60 seconds if no timeout is specified
	}
	cntx, cancel := ctx.WithTimeout(ctx.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	if inputJSON != "" {
		err = json.Unmarshal(jsonBytes, &m)
		if err != nil {
			return false, err
		}
	}
	if a.operation == "Remove One Document" {
		result, error = coll.DeleteOne(cntx, m)
	} else {
		if inputJSON == "" {
			logRemove.Debugf("Going to delete all the documents from collection %s in MongoDB database %s", a.collectionName, a.database)
			result, error = coll.DeleteMany(cntx, bson.D{})
		} else {
			result, error = coll.DeleteMany(cntx, m)
		}
	}
	if error != nil {
		return false, error
	}
	context.SetOutput("deletedCount", result.DeletedCount)
	return true, nil
}
