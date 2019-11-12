package queryDocument

import (
	ctx "context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var logquery = log.ChildLogger(log.RootLogger(), "mongodb-querydocument")

func init() {
	err := activity.Register(&Activity{}, New)
	if err != nil {
		logquery.Errorf("MongoDB Query Activity init error : %s ", err.Error())
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
		act := &Activity{client: client, operation: settings.Operation, collectionName: settings.CollectionName, database: settings.Database, timeout: settings.Timeout}
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
	timeout        int32
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

	if inputVal.Input == nil {
		// should we throw an error or warn?
		logquery.Warnf("No criteria specified for Query!")
		return true, nil
	}
	jsonBytes, err := json.Marshal(inputVal.Input)
	if err != nil {
		return false, fmt.Errorf("Error reading input json %s", err.Error())
	}
	inputJSON = string(jsonBytes)
	logquery.Debugf("InputJSON: ", inputJSON)
	if err != nil {
		return false, fmt.Errorf("Error getting mongodb connection %s", err.Error())
	}

	client := a.client
	db := client.Database(a.database)
	timeout := a.timeout
	if timeout == 0 {
		timeout = 60 //set a default timeout of 60 seconds if no timeout is specified
	}
	coll := db.Collection(collectionName)
	cntx, cancel := ctx.WithTimeout(ctx.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
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
			//return false, err
			if strings.Contains(err.Error(), "no documents in result") {
				val["QueryError"] = err.Error()
			} else {
				return false, err
			}
		}
		context.SetOutput("response", val)
	} else {
		var err error
		var cursor *mongo.Cursor
		var i = 0
		if inputJSON == "" {
			cursor, err = coll.Find(cntx, bson.D{})
		} else {
			cursor, err = coll.Find(cntx, m)
		}
		if err != nil {
			return false, err
		}
		for cursor.Next(cntx) {
			i++
			val1 := make(map[string]interface{})
			err := cursor.Decode(&val1)
			if err != nil {
				return false, err
			}
			objects = append(objects, val1)
		}
		if i > 0 {
			context.SetOutput("response", objects)
		} else {
			val["QueryError"] = "mongo: no documents in result"
			context.SetOutput("response", val)
		}

	}
	return true, nil
}
