package mongodb

import (
	"context"
	"strings"
	"time"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	methodUnknown int8 = iota
	methodGet
	methodInsert
	methodDelete
	methodUpdate
)

func init() {
	_ = activity.Register(&Activity{}, New)
}

func New(ctx activity.InitContext) (activity.Activity, error) {

	s := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), s, true)
	if err != nil {
		return nil, err
	}

  mCtx, _ := context.WithTimeout(context.Background(), 10*time.Second)
  
	opts := options.Client()

	if s.Username != "" && s.Password != "" {
		opts = opts.SetAuth(options.Credential{
			Username: s.Username,
			Password: s.Password,
		})
	}
	client, err := mongo.Connect(mCtx, opts.ApplyURI(s.URI))

	if err != nil {
		return nil, err
	}

	db := client.Database(s.DbName)
	collection := db.Collection(s.Collection)

	act := &Activity{client: client, collection: collection}

	switch strings.ToLower(s.Method) {
	case "get":
		act.method = methodGet
	case "insert":
		act.method = methodInsert
	case "delete":
		act.method = methodDelete
	case "update":
		act.method = methodUpdate
	}

	return act, nil
}

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

type Activity struct {
	client     *mongo.Client
	collection *mongo.Collection
	method     int8
}

// Metadata returns the activity's metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

func (a *Activity) Cleanup() error {

	log.RootLogger().Tracef("cleaning up MongoDB activity")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return a.client.Disconnect(ctx)
}

// Eval implements api.Activity.Eval - Logs the Message
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {
	logger := ctx.Logger()
	input := &Input{}
	err = ctx.GetInputObject(input)
	if err != nil {
		return true, nil
	}
	output := &Output{}

	//todo consider making timeout a setting
	bCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	switch a.method {

	case methodGet:
		result := a.collection.FindOne(bCtx, bson.M{input.KeyName: input.KeyValue})
		err = result.Decode(&output.Data)
		if err != nil {
			return true, err
		}
		logger.Tracef("Get Result: %v", output.Data)

	case methodInsert:

		if input.Data == nil && input.KeyValue == "" {
			// should we throw an error or warn?
			ctx.Logger().Warnf("Nothing to insert")
			return true, nil
		}

		var result *mongo.InsertOneResult

		if input.Data != nil {
			result, err = a.collection.InsertOne(bCtx, bson.D{input.Data.(bson.E)})

		} else {
			result, err = a.collection.InsertOne(bCtx, bson.M{input.KeyName: input.KeyValue})
		}

		if err != nil {
			return true, err
		}

		logger.Tracef("Inserted ID: %v", result.InsertedID)
		output.Data = result.InsertedID

	case methodDelete:
		result, err := a.collection.DeleteOne(bCtx, bson.M{input.KeyName: input.KeyValue}, nil)
		if err != nil {
			return true, err
		}

		logger.Tracef("Delete Count: %d", result.DeletedCount)
		output.Data = result.DeletedCount

	case methodUpdate:
		result, err := a.collection.UpdateOne(bCtx, bson.M{input.KeyName: input.KeyValue}, bson.M{"$set": input.Data})
		if err != nil {
			return true, err
		}

		resultObj := map[string]interface{}{"MatchedCount": result.MatchedCount, "ModifiedCount": result.ModifiedCount,
			"UpsertedCount": result.UpsertedCount, "UpsertedID": result.UpsertedID}

		logger.Tracef("Update Result: %v", resultObj)

		output.Data = resultObj
	}

	err = ctx.SetOutputObject(output)
	if err != nil {
		return false, err
	}
	return true, nil
}
