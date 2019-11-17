package insertDocument

import (
	"github.com/project-flogo/core/data/coerce"
)

// Settings structure
type Settings struct {
	Connection     string `md:"connection,required"`                                                   // The MongoDB connection
	Operation      string `md:"operation,required,allowed(Insert One Document,Insert Many Documents)"` // Operation to perform
	CollectionName string `md:"collectionName,required"`                                               // The collection within the MongoDB database to Insert Documents
	Database       string `md:"databaseName,required"`                                                 // MongoDB database to Insert Documents
	Timeout        int32  `md:"timeout"`                                                               // Timeout in seconds for the activity's operations
	ContinueOnErr  bool   `md:"continueOnErr"`                                                         // In case of Insert Many Documents operation, should the activity continue to insertDocument when the previous insertDocument operation failed?
}

//Input structure
type Input struct {
	Data interface{} `md:"data,required"` // The JSON Object that will serve as the input data
}

//Output structure
type Output struct {
	InsertedID   string `md:"insertedId"`   // InsertedId of inserted document. In case of Insert Many Documents, a list of IDs is returned
	TotalCount   int    `md:"totalCount"`   // Applicable for Insert Many Documents only. The total numner of Documents that were attempted to be inserted.
	SuccessCount int    `md:"successCount"` // Applicable for Insert Many Documents only. The total number of successful Document Insertions.
	FailureCount int    `md:"failureCount"` // Applicable for Insert Many Documents only. The total number of Document insertions that failed.
}

//FromMap method
func (i *Input) FromMap(values map[string]interface{}) error {
	i.Data, _ = values["data"]
	return nil
}

//ToMap method
func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"data": i.Data,
	}
}

//ToMap Output
func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"insertedId":   o.InsertedID,
		"totalCount":   o.TotalCount,
		"successCount": o.SuccessCount,
		"failureCount": o.FailureCount,
	}
}

//FromMap Output
func (o *Output) FromMap(values map[string]interface{}) error {
	var err error
	o.InsertedID, err = coerce.ToString(values["insertedId"])
	if err != nil {
		return err
	}
	o.TotalCount, err = coerce.ToInt(values["totalCount"])
	if err != nil {
		return err
	}
	o.SuccessCount, err = coerce.ToInt(values["successCount"])
	if err != nil {
		return err
	}
	o.FailureCount, err = coerce.ToInt(values["failureCount"])
	if err != nil {
		return err
	}
	return nil
}
