package queryDocument

import (
	"github.com/project-flogo/core/data/coerce"
)

// Settings structure
type Settings struct {
	Connection     string `md:"mongoConnection,required"`
	Operation      string `md:"operation,required"`
	CollectionName string `md:"collectionName,required"`
}

//Input structure
type Input struct {
	Input map[string]interface{} `md:"input,required"`
}

//Output structure
type Output struct {
	Output interface{} `md:"output"` //

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
		"output": o.Output,
	}
}

//FromMap Output
func (o *Output) FromMap(values map[string]interface{}) error {
	var err error
	o.Output, _ = (values["output"])
	if err != nil {
		return err
	}

	return nil
}
