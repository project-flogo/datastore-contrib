package eventlistener

import (
	"github.com/project-flogo/core/data/coerce"
)

// Settings structure
type Settings struct {
	Connection string `md:"mongodbConnection,required"`
}

// HandlerSettings structure
type HandlerSettings struct {
	Collection   string `md:"collection"`
	ListenInsert bool   `md:"listenInsert,required"`
	ListenUpdate bool   `md:"listenUpdate,required"`
	ListenRemove bool   `md:"listenRemove,required"`
}

// Output structure
type Output struct {
	Output map[string]interface{} `md:"Output"`
}

// FromMap method for HandlerSettings
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

//ToMap method for HandlerSettings
func (i *HandlerSettings) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"collection":   i.Collection,
		"listenInsert": i.ListenInsert,
		"listenUpdate": i.ListenUpdate,
		"listenRemove": i.ListenRemove,
	}
}

//FromMap method for Settings
func (i *Settings) FromMap(values map[string]interface{}) error {
	var err error

	i.Connection, err = coerce.ToString(values["mongodbConnection"])
	if err != nil {
		return err
	}

	return nil
}

// ToMap method for Settings
func (i *Settings) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"mongodbConnection": i.Connection,
	}
}

// ToMap method for Output
func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"Output": o.Output,
	}
}

// FromMap method for Output
func (o *Output) FromMap(values map[string]interface{}) error {

	var err error

	o.Output, err = coerce.ToObject(values["Output"])
	if err != nil {
		return err
	}

	return nil
}
