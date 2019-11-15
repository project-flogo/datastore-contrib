package queryDocument

// Settings structure
type Settings struct {
	Connection     string `md:"connection,required"`                                               // The MongoDB connection
	Operation      string `md:"operation,required,allowed(Find One Document,Find Many Documents)"` // Operation to perform
	CollectionName string `md:"collectionName,required"`                                           // The collection within the MongoDB database to query
	Database       string `md:"databaseName,required"`                                             // MongoDB databse to query
	Timeout        int32  `md:"timeout"`                                                           // Timeout in seconds for the activity's operations
}

//Input structure
type Input struct {
	Input interface{} `md:"criteria,required"` // The JSON Request Object that will serve as search parameter for the query
}

//Output structure
type Output struct {
	Output interface{} `md:"response"` // The JSON Response for Querying one or more documents from a collection
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
		"response": o.Output,
	}
}

//FromMap Output
func (o *Output) FromMap(values map[string]interface{}) error {
	var err error
	o.Output, _ = (values["response"])
	if err != nil {
		return err
	}

	return nil
}
