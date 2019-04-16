package mongodb

import (
	"context"
	"time"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	methodGet    = "GET"
	methodInsert = "INSERT"
	methodDelete = "DELETE"
	methodUpdate = "UPDATE"
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
	act := &Activity{settings: s}
	return act, nil
}

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

type Activity struct {
	settings *Settings
}

// Metadata returns the activity's metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

// Eval implements api.Activity.Eval - Logs the Message
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {
	logger := ctx.Logger()
	bctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(bctx, options.Client().ApplyURI(a.settings.URI))

	err = client.Connect(bctx)
	input := &Input{}
	err = ctx.GetInputObject(input)
	if err != nil {
		return true, nil
	}

	collection := client.Database(a.settings.DbName).Collection(a.settings.Collection)

	//res, err := collection.InsertOne(bctx, bson.A{"bar", "world", 3.14159, bson.D{{"qux", 12345}}})
	output := &Output{}
	switch a.settings.Method {

	case methodGet:
		result := collection.FindOne(bctx, bson.M{input.KeyName: input.KeyValue})

		logger.Debugf("Result...", result)
		result.Decode(&output.Data)

	case methodInsert:

		if input.Data != nil {
			res, err := collection.InsertOne(bctx, bson.D{input.Data.(bson.E)})

			if err != nil {
				return true, err
			}
			output.Data = res
			break
		}
		res, err := collection.InsertOne(bctx, bson.M{input.KeyName: input.KeyValue})
		if err != nil {
			return true, err
		}
		output.Data = res

	case methodDelete:
		result, err := collection.DeleteOne(bctx, bson.M{input.KeyName: input.KeyValue}, nil)

		if err != nil {
			return true, err
		}

		logger.Debugf("Result...", result)

		output.Data = result

	case methodUpdate:
		result, err := collection.UpdateOne(bctx, bson.M{input.KeyName: input.KeyValue}, bson.M{"$set": input.Data})

		if err != nil {
			return true, err
		}

		logger.Debugf("Result...", result)

		output.Data = result
	}

	return true, nil
}
