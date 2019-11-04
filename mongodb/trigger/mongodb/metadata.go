package mongodbtrigger

import (
	"github.com/project-flogo/core/data/coerce"
)

// Settings structure
type Settings struct {
	Connection string `md:"mongodbConnection,required"`
}

// HandlerSettings structure
type HandlerSettings struct {
	Collection   string `md:"collectionName"`
	ListenInsert bool   `md:"listenInsert,required"`
	ListenUpdate bool   `md:"listenUpdate,required"`
	ListenRemove bool   `md:"listenRemove,required"`
}

// Output structure
type Output struct {
	Output map[string]interface{} `md:"output"`
}

// ToMap method for Output
func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"output": o.Output,
	}
}

// FromMap method for Output
func (o *Output) FromMap(values map[string]interface{}) error {

	var err error

	o.Output, err = coerce.ToObject(values["output"])
	if err != nil {
		return err
	}

	return nil
}
