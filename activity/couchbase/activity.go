package couchbase

import (
	"strings"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
	"gopkg.in/couchbase/gocb.v1"
)

const (
	methodUnknown int8 = iota
	methodGet
	methodInsert
	methodUpsert
	methodRemove
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

	cluster, connectError := gocb.Connect(s.Server)
	if connectError != nil {
		ctx.Logger().Errorf("Connection error: %v", connectError)
		return nil, connectError
	}

	err = cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: s.Username,
		Password: s.Password,
	})
	if err != nil {
		return nil, err
	}

	bucket, openBucketError := cluster.OpenBucket(s.BucketName, s.BucketPassword)
	if openBucketError != nil {
		ctx.Logger().Errorf("Error while opening the bucked with the specified credentials: %v", openBucketError)
		return nil, openBucketError
	}

	act := &Activity{bucket: bucket}

	switch strings.ToLower(s.Method) {
	case "get":
		act.method = methodGet
	case "insert":
		act.method = methodInsert
	case "upsert":
		act.method = methodUpsert
	case "remove":
		act.method = methodRemove
	}

	return act, nil
}

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

type Activity struct {
	bucket *gocb.Bucket
	method int8
	expiry int
}

// Metadata returns the activity's metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

func (a *Activity) Cleanup() error {

	log.RootLogger().Tracef("cleaning up Couchbase activity")
	return a.bucket.Close()
}

// Eval implements api.Activity.Eval - Logs the Message
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {
	logger := ctx.Logger()

	input := &Input{}

	err = ctx.GetInputObject(input)
	if err != nil {
		return false, err
	}

	output := &Output{}

	switch a.method {
	case methodInsert:
		cas, methodError := a.bucket.Insert(input.Key, input.Data, uint32(a.expiry))
		if methodError != nil {
			logger.Errorf("Insert error: %v", methodError)
			return false, methodError
		}
		output.Data = cas
	case methodUpsert:
		cas, methodError := a.bucket.Upsert(input.Key, input.Data, uint32(a.expiry))
		if methodError != nil {
			logger.Errorf("Upsert error: %v", methodError)
			return false, methodError
		}
		output.Data = cas
	case methodRemove:
		cas, methodError := a.bucket.Remove(input.Key, 0)
		if methodError != nil {
			logger.Errorf("Remove error: %v", methodError)
			return false, methodError
		}
		output.Data = cas
	case methodGet:
		var document interface{}
		_, methodError := a.bucket.Get(input.Key, &document)
		if methodError != nil {
			logger.Errorf("Get error: %v", methodError)
			return false, methodError
		}
		output.Data = document
	}

	err = ctx.SetOutputObject(output)
	if err != nil {
		return false, err
	}
	return true, nil
}
