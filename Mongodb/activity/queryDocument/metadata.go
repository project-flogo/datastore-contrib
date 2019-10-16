package queryDocument

import (
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/support/connection"
	mongocon "github.com/project-flogo/datastore-contrib/Mongodb/connector/connection"
)

// Settings structure
type Settings struct {
	//JSONDocument interface{}        `md:"jsonDocument,required"`
	Connection connection.Manager `md:"mongoConnection,required"`
	Operation  string             `md:"operation,required"`
	CollName   string             `md:"collName,required"`
}

//Input structure
type Input struct {
	Input map[string]interface{} `md:"input,required"`
}

//Output structure
type Output struct {
	Output interface{} `md:"Output"` //

}

//FromMap method
func (i *Settings) FromMap(values map[string]interface{}) error {
	var err error

	i.CollName, err = coerce.ToString(values["collName"])
	if err != nil {
		return err
	}
	i.Operation, err = coerce.ToString(values["operation"])
	if err != nil {
		return err
	}
	i.Connection, err = mongocon.GetSharedConfiguration(values["mongoConnection"])
	if err != nil {
		return err
	}

	return nil
}

//ToMap method
func (i *Settings) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"collName":        i.CollName,
		"operaton":        i.Operation,
		"mongoConnection": i.Connection,
	}
}

//FromMap method
func (i *Input) FromMap(values map[string]interface{}) error {
	var err error

	i.Input, err = coerce.ToObject(values["input"])
	if err != nil {
		return err
	}

	return nil
}

//ToMap method
func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"input": i.Input,
	}
}

//ToMap Output
func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"Output": o.Output,
	}
}

//FromMap Output
func (o *Output) FromMap(values map[string]interface{}) error {

	var err error
	o.Output, _ = (values["Output"])
	if err != nil {
		return err
	}

	return nil
}
