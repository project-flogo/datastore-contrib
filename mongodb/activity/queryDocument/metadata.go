package queryDocument

// Settings structure
type Settings struct {
	Connection     string `md:"mongoConnection,required"` // The MongoDB connection
	Operation      string `md:"operation,required"`       // Operation to perform: Find One Document or Find Many Documents
	CollectionName string `md:"collectionName,required"`  // Name of collection to query from
}

//Input structure
type Input struct {
	Input interface{} `md:"jsonDocument,required"` // The JSON Request Object that will serve as search parameter for the query
}

//Output structure
type Output struct {
	Output interface{} `md:"response"` // The JSON Response for Querying one or more documents from a collection
}

//FromMap method
func (i *Input) FromMap(values map[string]interface{}) error {
	i.Input, _ = values["jsonDocument"]
	return nil
}

//ToMap method
func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"jsonDocument": i.Input,
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
