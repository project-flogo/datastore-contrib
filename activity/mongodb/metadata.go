package mongodb

import (
	"github.com/project-flogo/core/data/coerce"
)

type Settings struct {
	URI        string `md:"uri,required"`
	Method     string `md:"method,required,allowed(GET,INSERT,UPDATE,DELETE)"`
	DbName     string `md:"dbName,required"`
	Collection string `md:"collection"`
}

type Input struct {
	KeyName  string      `md:"keyName"`
	KeyValue string      `md:"keyValue"`
	Data     interface{} `md:"data"`
}

func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"keyName":  i.KeyName,
		"keyValue": i.KeyValue,
		"data":     i.Data,
	}
}

func (i *Input) FromMap(values map[string]interface{}) error {
	var err error

	i.KeyName, err = coerce.ToString(values["keyName"])
	if err != nil {
		return err
	}
	i.KeyValue, err = coerce.ToString(values["keyValue"])
	if err != nil {
		return err
	}

	i.Data, _ = values["data"]

	return nil
}

type Output struct {
	Data  interface{} `md:"data"`
	Count int32       `md:"count"`
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{

		"data":  o.Data,
		"count": o.Count,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {
	var err error
	o.Count, err = coerce.ToInt32(values["count"])
	if err != nil {
		return err
	}

	o.Data, _ = values["data"]
	return nil
}
