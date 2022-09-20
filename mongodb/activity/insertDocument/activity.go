package insertDocument

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

	connection "github.com/project-flogo/datastore-contrib/mongodb/connection"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var logInsert = log.ChildLogger(log.RootLogger(), "mongodb-insertDocument")

func init() {
	err := activity.Register(&Activity{}, New)
	if err != nil {
		logInsert.Errorf("MongoDB Insert Document Activity init error : %s ", err.Error())
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
			database: settings.Database, timeout: settings.Timeout, contErr: settings.ContinueOnErr}
		return act, nil
	}

	return nil, nil
}

// Activity is a stub for your Activity implementation
type Activity struct {
	mdMgr          connection.MongoDBManager
	operation      string
	collectionName string
	database       string
	timeout        int32
	contErr        bool
}

var activityMd = activity.ToMetadata(&Input{}, &Output{})

// Metadata implements activity.Activity.Metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

//Cleanup method
func (a *Activity) Cleanup() error {
	logInsert.Debugf("cleaning up MongoDB Insert Activity")
	ctx, cancel := ctx.WithTimeout(ctx.Background(), 30*time.Second)
	defer cancel()
	if a.mdMgr.IsConnected() {
		return a.mdMgr.Client.Disconnect(ctx)
	}
	return nil

}

// Eval implements activity.Activity.Eval
func (a *Activity) Eval(context activity.Context) (done bool, err error) {

	if a.operation == "" {
		return false, activity.NewError("operation is not configured", "MONGODB-1003", nil)
	}
	if a.collectionName == "" {
		return false, activity.NewError("collection Name is not configured", "MONGODB-1004", nil)
	}
	var inputJSON string
	inputVal := &Input{}
	err = context.GetInputObject(inputVal)
	if err != nil {
		return true, nil
	}

	dataJSONByte, err := json.Marshal(inputVal.Data)
	if err != nil {
		return false, fmt.Errorf("Error reading Data json %s", err.Error())
	}
	inputJSON = string(dataJSONByte)
	logInsert.Debugf("Insert Data JSON from input: ", inputJSON)
	if inputJSON == "" || inputJSON == "null" {
		return false, fmt.Errorf("Input Data cannot be null ")
	}

	if !a.mdMgr.IsConnected() {
		err := a.mdMgr.Connect()
		if err != nil {
			logInsert.Errorf("===ping error===")
			return false, activity.NewRetriableError(fmt.Sprintf("Failed to ping to server due to error - {%s}", err.Error()), "", nil)
		}
		logInsert.Debugf("===Ping success===")
	}

	db := a.mdMgr.Client.Database(a.database)
	coll := db.Collection(a.collectionName)
	timeout := a.timeout
	if timeout <= 0 {
		timeout = 60 //set a default timeout of 60 seconds if no timeout is specified
	}
	cntx, cancel := ctx.WithTimeout(ctx.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	failedNum := 0
	succeedNum := 0
	totalCount := 0
	logInsert.Debugf("Insert Operation to perform is [%s] ", a.operation)
	if a.operation == "Insert One Document" {
		m := make(map[string]interface{})
		err = json.Unmarshal([]byte(inputJSON), &m)
		if err != nil {
			return false, err
		}
		result, err := coll.InsertOne(cntx, m)
		if err != nil {
			return false, err
		}
		context.SetOutput("insertedId", result.InsertedID)
	} else {
		insertOptions := *options.InsertMany()
		// SetOrdered configures the ordered option. If true, when a write fails, the function will return without attempting
		// remaining writes. Defaults to true.
		if a.contErr {
			logInsert.Debug("Continue on Error is set to true")
			insertOptions.SetOrdered(false)
		} else {
			logInsert.Debug("Continue on Error is set to false")
		}

		var m []interface{}
		err = json.Unmarshal([]byte(inputJSON), &m)
		if err != nil {
			return false, err
		}

		result, err := coll.InsertMany(cntx, m, &insertOptions)
		totalCount = len(m)
		if err != nil {
			msg := err.Error()
			if a.contErr {
				if strings.Contains(msg, "E1100") {
					logInsert.Warnf(
						"Duplicate Key Error occurred while performing Inserting Many Documents."+
							" Since Continue on Error is set to True, error is being handled gracefully."+
							" Original Error for reference : [%s] ", msg)
					for strings.Contains(msg, "E1100") {
						failedNum = failedNum + 1
						index := strings.Index(msg, "E11000") + 6
						msg = msg[index:len(msg)]
					}

					succeedNum = totalCount - failedNum
				} else {
					logInsert.Errorf("Unexpected Error while performing Insert Many Documents Operation : [%s] ", msg)
					return false, err
				}
			} else {
				/* TODO: Handle this scenario better
				Since API does not properly return list of InsertedIDs in case of continue on Error = false
				and there was an error in one of the inserts, there is no good way to handle this case.
				Users will see incorrect data in the output.
				*/
				logInsert.Warnf(
					"Error encountered during Insert Many Documents with Continue on Error set to false."+
						"Users should check the mongoDB backend for the result of this particular execution for clarity."+
						"Original error for reference : [%s] ", msg)
			}
		} else {
			succeedNum = totalCount
		}
		context.SetOutput("totalCount", totalCount)
		context.SetOutput("successCount", succeedNum)
		context.SetOutput("failureCount", failedNum)
		context.SetOutput("insertedId", result.InsertedIDs)
	}
	logInsert.Debugf("Execution of MongoDB Insert Document Activity completed successfully")
	return true, nil
}
