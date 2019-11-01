package eventlistener

import (
	"context"
	"encoding/json"

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
	handler      trigger.Handler
	collName     string
	listenInsert bool
	listenUpdate bool
	listenRemove bool
	mclient      *mongo.Client
	database     string
	done         chan bool
	logger       log.Logger
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
		evntLsnr.handler = handler
		evntLsnr.logger = t.logger
		evntLsnr.collName = s.Collection
		evntLsnr.listenInsert = s.ListenInsert
		evntLsnr.listenUpdate = s.ListenUpdate
		evntLsnr.listenRemove = s.ListenRemove
		evntLsnr.mclient = t.mclient
		evntLsnr.database = config.Database
		evntLsnr.done = make(chan bool)
		t.evntLsnrs = append(t.evntLsnrs, evntLsnr)
		t.logger.Debugf("============collName=== %s", evntLsnr.collName)
		t.logger.Debugf("========listenInsert=== %b", evntLsnr.listenInsert)

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
		if evntLsnr.listenInsert == true && (evntLsnr.listenRemove == false && evntLsnr.listenUpdate == false) {
			eventOption = 1
		} else if evntLsnr.listenUpdate == true && (evntLsnr.listenInsert == false && evntLsnr.listenRemove == false) {
			eventOption = 2
		} else if evntLsnr.listenRemove == true && (evntLsnr.listenInsert == false && evntLsnr.listenUpdate == false) {
			eventOption = 3
		} else if (evntLsnr.listenInsert == true && evntLsnr.listenUpdate == true) && evntLsnr.listenRemove == false {
			eventOption = 4
		} else if (evntLsnr.listenInsert == true && evntLsnr.listenRemove == true) && evntLsnr.listenUpdate == false {
			eventOption = 5
		} else if (evntLsnr.listenRemove == true && evntLsnr.listenUpdate == true) && evntLsnr.listenInsert == false {
			eventOption = 6
		}
		t.logger.Debugf("====eventOption=== %d", eventOption)
		switch eventOption {
		case 1:
			pipeline = mongo.Pipeline{bson.D{{"$match",
				bson.D{{"operationType", "insert"}},
			}}}
		case 2:
			pipeline = mongo.Pipeline{bson.D{{"$match",
				bson.D{{"operationType", "update"}},
			}}}
		case 3:
			pipeline = mongo.Pipeline{bson.D{{"$match",
				bson.D{{"operationType", "delete"}},
			}}}
		case 4:
			pipeline = mongo.Pipeline{bson.D{{"$match", bson.D{{"$or",
				bson.A{
					bson.D{{"operationType", "insert"}},
					bson.D{{"operationType", "update"}}}}},
			}}}
		case 5:
			pipeline = mongo.Pipeline{bson.D{{"$match", bson.D{{"$or",
				bson.A{
					bson.D{{"operationType", "insert"}},
					bson.D{{"operationType", "delete"}}}}},
			}}}
		case 6:
			pipeline = mongo.Pipeline{bson.D{{"$match", bson.D{{"$or",
				bson.A{
					bson.D{{"operationType", "delete"}},
					bson.D{{"operationType", "update"}}}}},
			}}}
		default:
			pipeline = mongo.Pipeline{}
		}
		db := evntLsnr.mclient.Database(evntLsnr.database)
		coll := db.Collection(evntLsnr.collName)

		var stream *mongo.ChangeStream
		var err error
		if evntLsnr.collName != "" {
			t.logger.Infof("listening on collection:: %s", evntLsnr.collName)
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
		entLsnr.logger.Errorf("Failed to process record from collection [%s], due to error - %s", evntLsnr.collName, err.Error())
	} else {
		// record is successfully processed.
		entLsnr.logger.Infof("Record from collection [%s] is successfully processed", evntLsnr.collName)
	}

}

// Stop implements trigger.Trigger.Stop
func (t *Trigger) Stop() error {
	t.logger.Infof("Stopping Trigger - %s", t.id)
	for _, evntLsnr := range t.evntLsnrs {
		evntLsnr.mclient.Disconnect(context.Background())
		t.logger.Infof("client disconnected")
	}
	t.logger.Infof("Trigger - %s  stopped", t.id)
	return nil
}
