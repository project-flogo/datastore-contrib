package couchbase

import (
	"github.com/project-flogo/core/data/coerce"
)

type Settings struct {
	Username       string `md:"username"`
	Password       string `md:"password"`
	BucketName     string `md:"bucketName,required"`
	BucketPassword string `md:"bucketPassword"`
	Server         string `md:"server,required"`
}
type Input struct {
	Key    string `md:"key,required"`
	Data   string `md:"data"`
	Method string `md:"method,required,allowed(Insert,Upsert,Remove,Get)"`
	Expiry int32  `md:"expiry,required"`
}

type Output struct {
	Data interface{} `md:"data"`
}

func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"key":    i.Key,
		"data":   i.Data,
		"method": i.Method,
		"expiry": i.Expiry,
	}
}

func (i *Input) FromMap(values map[string]interface{}) error {
	var err error

	i.Key, err = coerce.ToString(values["key"])
	if err != nil {
		return err
	}
	i.Data, err = coerce.ToString(values["data"])
	if err != nil {
		return err
	}
	i.Method, err = coerce.ToString(values["method"])
	if err != nil {
		return err
	}
	i.Expiry, err = coerce.ToInt32(values["expiry"])
	if err != nil {
		return err
	}

	return nil
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
