package updateDocument

import (
	ctx "context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
	connection "github.com/project-flogo/datastore-contrib/mongodb/connection"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var logUpdate = log.ChildLogger(log.RootLogger(), "mongodb-updateDocument")

func init() {
	err := activity.Register(&Activity{}, New)
	if err != nil {
		logUpdate.Errorf("MongoDB Update Document Activity init error : %s ", err.Error())
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
		mdMgr := mcon.GetConnection().(connection.MongoDBManager)
		act := &Activity{mdMgr: mdMgr, operation: settings.Operation, collectionName: settings.CollectionName,
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

// Cleanup of activity
func (a *Activity) Cleanup() error {
	logUpdate.Debugf("cleaning up MongoDB Update Activity")
	ctx, cancel := ctx.WithTimeout(ctx.Background(), 30*time.Second)
	defer cancel()

	if a.mdMgr.IsConnected() {
		return a.mdMgr.Client.Disconnect(ctx)
	}
	return nil
}

// Activity is a stub for your Activity implementation
type Activity struct {
	mdMgr          connection.MongoDBManager
	operation      string
	collectionName string
	database       string
	timeout        int32
}

// Eval implements activity.Activity.Eval
func (a *Activity) Eval(context activity.Context) (done bool, err error) {
	logUpdate.Debugf("Executing  MongoDB Update Document Activity")

	if a.operation == "" {
		return false, activity.NewError("operation is not configured", "MONGODB-1004", nil)
	}
	if a.collectionName == "" {
		return false, activity.NewError("collection Name is not configured", "MONGODB-1004", nil)
	}
	var criteriaJSON, updateJSON string
	qm := make(map[string]interface{})
	um := make(map[string]interface{})
	inputVal := &Input{}
	err = context.GetInputObject(inputVal)
	if err != nil {
		return true, nil
	}

	criteriaJSONBytes, err := json.Marshal(inputVal.Criteria)
	if err != nil {
		return false, fmt.Errorf("Error reading Criteria json %s", err.Error())
	}
	criteriaJSON = string(criteriaJSONBytes)
	logUpdate.Debugf("Update Criteria JSON from input: ", criteriaJSON)
	if criteriaJSON != "" {
		err = json.Unmarshal([]byte(criteriaJSON), &qm)
		if err != nil {
			return false, err
		}
	}
	updateDataBytes, err := json.Marshal(inputVal.UpdateData)
	if err != nil {
		return false, fmt.Errorf("Error reading update JSON data %s", err.Error())
	}
	updateJSON = string(updateDataBytes)
	logUpdate.Debugf("Update Date JSON from input: ", updateJSON)
	if updateJSON != "" {
		err = json.Unmarshal([]byte(updateJSON), &um)
		if err != nil {
			return false, err
		}
	}

	if !a.mdMgr.IsConnected() {
		err := a.mdMgr.Connect()
		if err != nil {
			return false, activity.NewRetriableError(fmt.Sprintf("Failed to ping to server due to error - {%s}", err.Error()), "", nil)
		}
		logUpdate.Debugf("Successful ping to the server")
	}

	db := a.mdMgr.Client.Database(a.database)
	coll := db.Collection(a.collectionName)
	timeout := a.timeout
	if timeout <= 0 {
		timeout = 60 //set a default timeout of 60 seconds if no timeout is specified
	}
	cntx, cancel := ctx.WithTimeout(ctx.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	var result *mongo.UpdateResult
	query := make(map[string]interface{})
	err = json.Unmarshal([]byte(criteriaJSON), &query)
	if err != nil {
		return false, err
	}
	update := make(map[string]interface{})
	err = json.Unmarshal([]byte(updateJSON), &update)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	if a.operation == "Replace One Document" {
		result, err = coll.ReplaceOne(
			cntx, query, update)

	} else if a.operation == "Update One Document" {
		result, err = coll.UpdateOne(
			cntx,
			query,
			bson.D{{"$set", update}},
		)
	} else {
		result, err = coll.UpdateMany(
			cntx,
			query,
			bson.D{{"$set", update}},
		)
	}
	if err != nil {
		return false, err
	}
	context.SetOutput("matchedCount", result.MatchedCount)
	context.SetOutput("updatedCount", result.ModifiedCount)
	logUpdate.Debugf("Execution of MongoDB Update Document Activity completed successfully")
	return true, nil
}
