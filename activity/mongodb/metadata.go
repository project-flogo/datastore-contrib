package mongodb

import (
	"github.com/project-flogo/core/data/coerce"
)

type Settings struct {
	URI        string `md:"uri,required"` // The MongoDB connection URI
	Method     string `md:"method,required,allowed(GET,INSERT,UPDATE,DELETE)"` // The method type
	DbName     string `md:"dbName,required"` // The name of the database
	Collection string `md:"collection, required"` // The collection to work on
  Username   string `md:"username"` // The username of the client
	Password   string `md:"password"` // The password of the client
}

type Input struct {
	KeyName  string      `md:"keyName"`  // The name of the key to use when looking up an object (used in GET, UPDATE and DELETE)
	KeyValue string      `md:"keyValue"` // The value of the key to use when looking up an object (used in GET, UPDATE, and DELETE)
	Data     interface{} `md:"data"`     // The bson document to insert in mongodb
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
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{

		"data":  o.Data,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {
	o.Data, _ = values["data"]
	return nil
}
