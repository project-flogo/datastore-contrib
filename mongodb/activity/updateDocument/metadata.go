package updateDocument

import "github.com/project-flogo/core/data/coerce"

// Settings structure
type Settings struct {
	Connection     string `md:"connection,required"`     // The MongoDB connection
	Operation      string `md:"operation,required"`      // Operation to perform: Update One Document, Update Many Documents or Replace One Document
	CollectionName string `md:"collectionName,required"` // The collection within the MongoDB database to Update Documents
	Database       string `md:"databaseName,required"`   // MongoDB databse to Update Documents
	Timeout        int32  `md:"timeout"`                 // Timeout in seconds for the activity's operations
}

//Input structure
type Input struct {
	Criteria   interface{} `md:"criteria,required"`   // The JSON Request Object that will serve as search parameter for deciding which documents to update
	UpdateData interface{} `md:"updateData,required"` // The JSON Request Object that will serve as the update data
}

//Output structure
type Output struct {
	MatchedCount int64 `md:"matchedCount"` // A number indicating total Documents that were matched for update
	UpdatedCount int64 `md:"updatedCount"` // A number indicating total Documents that were updated by this activity

}

//FromMap method
func (i *Input) FromMap(values map[string]interface{}) error {
	i.Criteria, _ = values["criteria"]
	i.UpdateData, _ = values["updateData"]
	return nil
}

//ToMap method
func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"criteria":   i.Criteria,
		"updateData": i.UpdateData,
	}
}

//ToMap Output
func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"matchedCount": o.MatchedCount,
		"updatedCount": o.UpdatedCount,
	}
}

//FromMap Output
func (o *Output) FromMap(values map[string]interface{}) error {
	var err error
	o.MatchedCount, err = coerce.ToInt64(values["matchedCount"])
	if err != nil {
		return err
	}
	o.UpdatedCount, err = coerce.ToInt64(values["updatedCount"])
	if err != nil {
		return err
	}
	return nil
}
