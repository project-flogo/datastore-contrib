package mongodbtrigger

import (
	"github.com/project-flogo/core/data/coerce"
)

// Settings structure
type Settings struct {
	Connection string `md:"mongodbConnection,required"` // The MongoDB connection
}

// HandlerSettings structure
type HandlerSettings struct {
	Collection   string `md:"collectionName"`        // The collection to listen to for changes. If left blank, listens to all collections in a DB
	ListenInsert bool   `md:"listenInsert,required"` // Should the trigger listen to Insert events?
	ListenUpdate bool   `md:"listenUpdate,required"` // Should the trigger listen to Update events?
	ListenRemove bool   `md:"listenRemove,required"` // Should the trigger listen to Remove events?
}

// Output structure
type Output struct {
	Output map[string]interface{} `md:"output"` //The Output of the trigger
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
