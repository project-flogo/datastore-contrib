package couchbase

import (
	"github.com/project-flogo/core/data/coerce"
)

type Settings struct {
	Username       string `md:"username"` // Cluster username
	Password       string `md:"password"` // Cluster password
	BucketName     string `md:"bucketName,required"` // The bucket name
	BucketPassword string `md:"bucketPassword"`      // The bucket password
	Server         string `md:"server,required"`     // The Couchbase server (e.g. couchbase://127.0.0.1)
	Method         string `md:"method,required,allowed(Insert,Upsert,Remove,Get)"` // The method type (Insert, Upsert, Remove or Get); (default: *Insert*)
	Expiry         int    `md:"expiry,required"`     // The document expiry (default: 0)
}
type Input struct {
	Key  string `md:"key,required"` // The document key identifier
	Data string `md:"data"`         // The document data (when the method is get this field is ignored)
}

type Output struct {
	Data interface{} `md:"data"`    // The result of the method invocation
}

func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"key":  i.Key,
		"data": i.Data,
	}
}

func (i *Input) FromMap(values map[string]interface{}) error {
	var err error

	i.Key, err = coerce.ToString(values["key"])
	if err != nil {
		return err
	}
	i.Data, err = coerce.ToString(values["data"])
	return err
}

func (o *Output) ToMap() map[string]interface{} {

	return map[string]interface{}{
		"data": o.Data,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {

	o.Data, _ = values["data"]
	return nil
}
