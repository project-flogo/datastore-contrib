package executeCypherQuery

import (
	"github.com/project-flogo/core/data/coerce"
)

// Settings structure
type Settings struct {
	Connection   string `md:"connection,required"` // The Neo4j connection
	DatabaseName string `md:"databaseName,required"`
	AccessMode   string `md:"accessMode,required"`
}

//Input structure
type Input struct {
	CypherQuery string                 `md:"cypherQuery,required"` // The cypher query
	QueryParams map[string]interface{} `md:"queryParams"`
}

//Output structure
type Output struct {
	Output interface{} `md:"response"` // The JSON Response of the query
}

//FromMap method
func (i *Input) FromMap(values map[string]interface{}) error {
	i.CypherQuery, _ = coerce.ToString(values["cypherQuery"])
	i.QueryParams, _ = coerce.ToObject(values["queryParams"])
	return nil
}

//ToMap method
func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"cypherQuery": i.CypherQuery,
		"queryParams": i.QueryParams,
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
