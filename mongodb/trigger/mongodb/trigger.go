package mongodbtrigger

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/trigger"
	mongodb "github.com/project-flogo/datastore-contrib/mongodb/connector/connection"
	"go.mongodb.org/mongo-driver/bson"
	mongo "go.mongodb.org/mongo-driver/mongo"
)

var triggerMd = trigger.NewMetadata(&Settings{}, &HandlerSettings{}, &Output{})

func init() {
	_ = trigger.Register(&Trigger{}, &TriggerFactory{})
}

// TriggerFactory My Trigger factory
type TriggerFactory struct {
	metadata *trigger.Metadata
}

//NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &TriggerFactory{metadata: md}
}

//New Creates a new trigger instance for a given id
func (t *TriggerFactory) New(config *trigger.Config) (trigger.Trigger, error) {
	settings := &Settings{}
	err := metadata.MapToStruct(config.Settings, settings, true)
	if err != nil {
		return nil, err
	}
	if settings.Connection != "" {
		mConn, toConnerr := coerce.ToConnection(settings.Connection)
		if toConnerr != nil {
			return nil, toConnerr
		}
		mclient := mConn.(*mongodb.MongodbSharedConfigManager).GetClient()
		return &Trigger{metadata: t.metadata, settings: settings, id: config.Id, mclient: mclient}, nil
	}
	return nil, nil
}

// Metadata implements trigger.Factory.Metadata
func (*TriggerFactory) Metadata() *trigger.Metadata {
	return triggerMd
}

// Trigger is a stub for your Trigger implementation
type Trigger struct {
	metadata  *trigger.Metadata
	settings  *Settings
	evntLsnrs []*EventListener
	mclient   *mongo.Client
	logger    log.Logger
	id        string
}

// EventListener is structure of a single EventListener
type EventListener struct {
	handler  trigger.Handler
	settings *HandlerSettings
	database string
	done     chan bool
	logger   log.Logger
}

// Initialize Mongodb Event Listener
func (t *Trigger) Initialize(ctx trigger.InitContext) error {
	t.logger = log.ChildLogger(ctx.Logger(), "mongodb-event-listener")
	t.logger.Infof("============initializing event listener==")
	mConn, toConnerr := coerce.ToConnection(t.settings.Connection)
	if toConnerr != nil {
		return nil
	}
	config := mConn.(*mongodb.MongodbSharedConfigManager).GetClientConfiguration()
	for _, handler := range ctx.GetHandlers() {
		s := &HandlerSettings{}
		err := metadata.MapToStruct(handler.Settings(), s, true)
		if err != nil {
			return err
		}
		evntLsnr := &EventListener{}
		evntLsnr.settings = s
		evntLsnr.handler = handler
		evntLsnr.logger = t.logger
		evntLsnr.database = config.Database
		evntLsnr.done = make(chan bool)
		t.evntLsnrs = append(t.evntLsnrs, evntLsnr)
		t.logger.Debugf("============collName=== %s", evntLsnr.settings.Collection)
		t.logger.Debugf("========listenInsert=== %b", evntLsnr.settings.ListenInsert)

	}

	return nil
}

// Metadata implements trigger.Trigger.Metadata
func (t *Trigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Start implements trigger.Trigger.Start
func (t *Trigger) Start() error {
	t.logger.Infof("Starting Trigger - %s", t.id)
	for _, evntLsnr := range t.evntLsnrs {
		pipeline := mongo.Pipeline{}
		eventOption := 0
		if evntLsnr.settings.ListenInsert && !evntLsnr.settings.ListenRemove && !evntLsnr.settings.ListenUpdate {
			eventOption = 1
			pipeline = mongo.Pipeline{bson.D{{"$match",
				bson.D{{"operationType", "insert"}},
			}}}
		} else if evntLsnr.settings.ListenUpdate && !evntLsnr.settings.ListenInsert && !evntLsnr.settings.ListenRemove {
			eventOption = 2
			pipeline = mongo.Pipeline{bson.D{{"$match",
				bson.D{{"operationType", "update"}},
			}}}
		} else if evntLsnr.settings.ListenRemove && !evntLsnr.settings.ListenInsert && !evntLsnr.settings.ListenUpdate {
			eventOption = 3
			pipeline = mongo.Pipeline{bson.D{{"$match",
				bson.D{{"operationType", "delete"}},
			}}}
		} else if evntLsnr.settings.ListenInsert && evntLsnr.settings.ListenUpdate && !evntLsnr.settings.ListenRemove {
			eventOption = 4
			pipeline = mongo.Pipeline{bson.D{{"$match", bson.D{{"$or",
				bson.A{
					bson.D{{"operationType", "insert"}},
					bson.D{{"operationType", "update"}}}}},
			}}}
		} else if evntLsnr.settings.ListenInsert && evntLsnr.settings.ListenRemove && !evntLsnr.settings.ListenUpdate {
			eventOption = 5
			pipeline = mongo.Pipeline{bson.D{{"$match", bson.D{{"$or",
				bson.A{
					bson.D{{"operationType", "insert"}},
					bson.D{{"operationType", "delete"}}}}},
			}}}
		} else if evntLsnr.settings.ListenRemove && evntLsnr.settings.ListenUpdate && !evntLsnr.settings.ListenInsert {
			eventOption = 6
			pipeline = mongo.Pipeline{bson.D{{"$match", bson.D{{"$or",
				bson.A{
					bson.D{{"operationType", "delete"}},
					bson.D{{"operationType", "update"}}}}},
			}}}
		} else {
			pipeline = mongo.Pipeline{}
		}
		fmt.Println("====eventOption=== %d", eventOption)
		t.logger.Debugf("====eventOption=== %d", eventOption)
		db := t.mclient.Database(evntLsnr.database)
		var stream *mongo.ChangeStream
		var err error
		if evntLsnr.settings.Collection != "" {
			coll := db.Collection(evntLsnr.settings.Collection)
			t.logger.Infof("listening on collection:: %s", evntLsnr.settings.Collection)
			stream, err = coll.Watch(context.Background(), pipeline)
		} else { // listening on database if no collection name specified
			t.logger.Infof("listening on all collections of database:: %s", evntLsnr.database)
			stream, err = db.Watch(context.Background(), pipeline)
		}

		if err != nil {
			t.logger.Errorf("Failed to open change stream %s", err)
			return err
		}
		// Start polling on a separate Go routine so as to not block engine
		go evntLsnr.listen(stream)
	}
	t.logger.Infof("Trigger - %s  started", t.id)
	return nil
}

func (entLsnr *EventListener) listen(stream *mongo.ChangeStream) {
	entLsnr.logger.Infof("============listening====")
	for {
		select {
		case <-entLsnr.done:
			entLsnr.logger.Infof("stopped listening...")
			// Exit
			return
		default:
			//	entLsnr.logger.Infof("done in listening %b", entLsnr.done)
			ok := stream.Next(context.Background())
			if ok {
				var res bson.D
				err := stream.Decode(&res)
				if err != nil {
					entLsnr.logger.Errorf("got error while decoding stream %s", err)
				}
				if len(res) == 0 {
					entLsnr.logger.Infof("result is empty, was expecting change document")
				}
				stringOp := stream.Current.String()
				go entLsnr.process(stringOp, entLsnr)
			} else {
				err := stream.Err()
				if err != nil {
					//if err is not nil, it means something bad happened, let's finish our func
					stream.Close(context.Background())
					entLsnr.logger.Infof("stopped listening...stream closed")
					return
				}
			}
		}
	}

}
func (entLsnr *EventListener) process(stringOp string, evntLsnr *EventListener) {
	entLsnr.logger.Infof("started processing record...")
	var outputRoot = map[string]interface{}{}
	var jsonOp map[string]interface{}
	err := json.Unmarshal([]byte(stringOp), &jsonOp)
	if err != nil {
		entLsnr.logger.Errorf("got error while unmarshelling %s", err)
		return
	}
	outputRoot["NameSpace"] = jsonOp["ns"]
	outputRoot["OperationType"] = jsonOp["operationType"]
	outputRoot["ResultDocument"] = jsonOp["fullDocument"]
	if jsonOp["operationType"] == "delete" {
		outputRoot["ResultDocument"] = jsonOp["documentKey"]
	}
	outputData := &Output{}
	outputData.Output = outputRoot
	_, err = evntLsnr.handler.Handle(context.Background(), outputData)
	if err != nil {
		entLsnr.logger.Errorf("Failed to process record from collection [%s], due to error - %s", evntLsnr.settings.Collection, err.Error())
	} else {
		// record is successfully processed.
		entLsnr.logger.Infof("Record from collection [%s] is successfully processed", evntLsnr.settings.Collection)
	}

}

// Stop implements trigger.Trigger.Stop
func (t *Trigger) Stop() error {
	t.logger.Infof("Stopping Trigger - %s", t.id)
	t.mclient.Disconnect(context.Background())
	t.logger.Infof("Trigger - %s  stopped", t.id)
	return nil
}
