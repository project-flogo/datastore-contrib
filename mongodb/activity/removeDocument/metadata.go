package removeDocument

import (
	"github.com/project-flogo/core/data/coerce"
)

// Settings structure
type Settings struct {
	Connection     string `md:"connection,required"`                                                   // The MongoDB connection
	Operation      string `md:"operation,required,allowed(Remove One Document,Remove Many Documents)"` // Operation to perform
	CollectionName string `md:"collectionName,required"`                                               // The collection within the MongoDB database to Remove Documents
	Database       string `md:"databaseName,required"`                                                 // MongoDB databse to Remove Documents
	Timeout        int32  `md:"timeout"`                                                               // Timeout in seconds for the activity's operations
}

//Input structure
type Input struct {
	Input interface{} `md:"criteria,required"` // The JSON Request Object that will serve as search parameter for deciding which documents to delete
}

//Output structure
type Output struct {
	DeletedCount int64 `md:"deletedCount"` // A number indicating total Documents that were deleted by this activity

}

//FromMap method
func (i *Input) FromMap(values map[string]interface{}) error {
	i.Input, _ = values["criteria"]
	return nil
}

//ToMap method
func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"criteria": i.Input,
	}
}

//ToMap Output
func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"deletedCount": o.DeletedCount,
	}
}

//FromMap Output
func (o *Output) FromMap(values map[string]interface{}) error {
	var err error
	o.DeletedCount, err = coerce.ToInt64(values["deletedCount"])
	if err != nil {
		return err
	}
	return nil
}
