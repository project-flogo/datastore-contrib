package mongodb

import (
	"github.com/project-flogo/core/data/coerce"
)

type Settings struct {
	URI string `md:"uri,required"`
}

type Input struct {
	DbName     string      `md:"dbname,required"`
	Collection string      `md:"collection"`
	Method     string      `md:"method"`
	KeyName    string      `md:"keyname"`
	KeyValue   string      `md:"keyvalue"`
	Data       interface{} `md:"data"`
}

func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"dbname":     i.DbName,
		"collection": i.Collection,
		"method":     i.Method,
		"keyname":    i.KeyName,
		"keyvalue":   i.KeyValue,
		"data":       i.Data,
	}
}

func (i *Input) FromMap(values map[string]interface{}) error {
	var err error

	i.DbName, err = coerce.ToString(values["dbname"])
	if err != nil {
		return err
	}
	i.Collection, err = coerce.ToString(values["collection"])
	if err != nil {
		return err
	}
	i.Method, err = coerce.ToString(values["method"])
	if err != nil {
		return err
	}
	i.KeyName, err = coerce.ToString(values["keyname"])
	if err != nil {
		return err
	}
	i.KeyValue, err = coerce.ToString(values["keyvalue"])
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
