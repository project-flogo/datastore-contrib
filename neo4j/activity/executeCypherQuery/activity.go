package executeCypherQuery

import (
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
)

var logquery = log.ChildLogger(log.RootLogger(), "neo4j-executecyppherquery")

func init() {
	err := activity.Register(&Activity{}, New)
	if err != nil {
		logquery.Errorf("Neo4j Execute Query Activity init error : %s ", err.Error())
	}
}

// New functioncommon
func New(ctx activity.InitContext) (activity.Activity, error) {
	settings := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), settings, true)
	if err != nil {
		return nil, err
	}
	if settings.Connection != "" {

		neo4jcon, toConnerr := coerce.ToConnection(settings.Connection)
		if toConnerr != nil {
			return nil, toConnerr
		}
		driver := neo4jcon.GetConnection().(neo4j.Driver)
		accessMode := neo4j.AccessModeRead
		if settings.AccessMode != "Read" {
			accessMode = neo4j.AccessModeWrite
		}
		act := &Activity{driver: driver, accessMode: accessMode, databaseName: settings.DatabaseName}
		return act, nil
	}
	return nil, nil
}

// Activity is a stub for your Activity implementation
type Activity struct {
	driver       neo4j.Driver
	accessMode   neo4j.AccessMode
	databaseName string
}

var activityMd = activity.ToMetadata(&Input{}, &Output{})

// Metadata implements activity.Activity.Metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

//Cleanup method
func (a *Activity) Cleanup() error {
	logquery.Debugf("cleaning up Neo4j activity")
	return nil
}

type NodeOutput struct {
	Id     int64
	Labels []string
	Props  map[string]interface{}
}

// Eval implements activity.Activity.Eval
func (a *Activity) Eval(context activity.Context) (done bool, err error) {
	logquery.Debugf("Executing  neo4j cypher query Activity")

	input := &Input{}
	err = context.GetInputObject(input)
	if err != nil {
		return true, nil
	}

	sessionConfig := neo4j.SessionConfig{AccessMode: a.accessMode, DatabaseName: a.databaseName}
	session, err := a.driver.NewSession(sessionConfig)
	if err != nil {
		logquery.Errorf("===session error==", err)
		return false, err
	}

	result, err := session.Run(input.CypherQuery, input.QueryParams)
	if err != nil {
		return false, err
	}

	//nodeList := []NodeOutput{}
	nodeList := []interface{}{}
	for result.Next() {
		keys := result.Record().Keys()
		for i, _ := range keys {
			record := result.Record().GetByIndex(i)
			switch record.(type) {
			case neo4j.Node:
				node := record.(neo4j.Node)
				nodeOutput := NodeOutput{Id: node.Id(),
					Labels: node.Labels(),
					Props:  node.Props(),
				}
				nodeList = append(nodeList, nodeOutput)
			case string:
				node := record.(string)
				nodeList = append(nodeList, node)
			case int64:
				node := record.(int64)
				nodeList = append(nodeList, node)
			case float64:
				node := record.(float64)
				nodeList = append(nodeList, node)
			}
		}
	}
	context.SetOutput("response", nodeList)
	session.Close()

	return true, nil
}
