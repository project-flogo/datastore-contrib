package eventlistener

import (
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/support/connection"
	mongodb "github.com/project-flogo/datastore-contrib/Mongodb/connector/connection"
)

type Settings struct {
	Connection connection.Manager `md:"mongodbConnection,required"`
}

type HandlerSettings struct {
	Collection   string `md:"collection"`
	ListenInsert bool   `md:"listenInsert,required"`
	ListenUpdate bool   `md:"listenUpdate,required"`
	ListenRemove bool   `md:"listenRemove,required"`
}

type Output struct {
	Output map[string]interface{} `md:"Output"`
}

//FromMap method
func (i *HandlerSettings) FromMap(values map[string]interface{}) error {
	var err error

	i.Collection, err = coerce.ToString(values["collection"])
	if err != nil {
		return err
	}
	i.ListenInsert, err = coerce.ToBool(values["listenInsert"])
	if err != nil {
		return err
	}
	i.ListenUpdate, err = coerce.ToBool(values["listenUpdate"])
	if err != nil {
		return err
	}
	i.ListenRemove, err = coerce.ToBool(values["listenRemove"])
	if err != nil {
		return err
	}

	return nil
}

//ToMap method
func (i *HandlerSettings) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"collection":   i.Collection,
		"listenInsert": i.ListenInsert,
		"listenUpdate": i.ListenUpdate,
		"listenRemove": i.ListenRemove,
	}
}

//FromMap method
func (i *Settings) FromMap(values map[string]interface{}) error {
	var err error

	i.Connection, err = mongodb.GetSharedConfiguration(values["mongodbConnection"])
	if err != nil {
		return err
	}

	return nil
}

//ToMap method
func (i *Settings) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"mongodbConnection": i.Connection,
	}
}
func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"Output": o.Output,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {

	var err error

	o.Output, err = coerce.ToObject(values["Output"])
	if err != nil {
		return err
	}

	return nil
}
