package mongodbtrigger

import (
	"context"
	"encoding/json"

	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/trigger"
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
		mclient := mConn.GetConnection().(*mongo.Client)
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
	stream   *mongo.ChangeStream
}

// Initialize Mongodb Event Listener
func (t *Trigger) Initialize(ctx trigger.InitContext) error {
	t.logger = log.ChildLogger(ctx.Logger(), "mongodb-event-listener")
	t.logger.Infof("============initializing event listener==")
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
		evntLsnr.database = evntLsnr.settings.Database
		evntLsnr.done = make(chan bool)
		t.evntLsnrs = append(t.evntLsnrs, evntLsnr)
		t.logger.Debugf("============collName=== %s", evntLsnr.settings.Collection)
		t.logger.Debugf("========listenInsert=== %b", evntLsnr.settings.ListenInsert)
		t.logger.Debugf("========listenUpdate=== %b", evntLsnr.settings.ListenUpdate)
		t.logger.Debugf("========listenRemove=== %b", evntLsnr.settings.ListenRemove)

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
		// Start polling on a separate Go routine so as to not block engine
		err := evntLsnr.setStream(t.mclient)
		if err != nil {
			return err
		}
		go evntLsnr.listen()
	}
	t.logger.Infof("Trigger - %s  started", t.id)
	return nil
}

func (evntLsnr *EventListener) setStream(mclient *mongo.Client) error {
	pipeline := mongo.Pipeline{}

	if evntLsnr.settings.ListenInsert && !evntLsnr.settings.ListenRemove && !evntLsnr.settings.ListenUpdate {

		pipeline = mongo.Pipeline{bson.D{{"$match",
			bson.D{{"operationType", "insert"}},
		}}}
	} else if evntLsnr.settings.ListenUpdate && !evntLsnr.settings.ListenInsert && !evntLsnr.settings.ListenRemove {

		pipeline = mongo.Pipeline{bson.D{{"$match",
			bson.D{{"operationType", "update"}},
		}}}
	} else if evntLsnr.settings.ListenRemove && !evntLsnr.settings.ListenInsert && !evntLsnr.settings.ListenUpdate {

		pipeline = mongo.Pipeline{bson.D{{"$match",
			bson.D{{"operationType", "delete"}},
		}}}
	} else if evntLsnr.settings.ListenInsert && evntLsnr.settings.ListenUpdate && !evntLsnr.settings.ListenRemove {

		pipeline = mongo.Pipeline{bson.D{{"$match", bson.D{{"$or",
			bson.A{
				bson.D{{"operationType", "insert"}},
				bson.D{{"operationType", "update"}}}}},
		}}}
	} else if evntLsnr.settings.ListenInsert && evntLsnr.settings.ListenRemove && !evntLsnr.settings.ListenUpdate {

		pipeline = mongo.Pipeline{bson.D{{"$match", bson.D{{"$or",
			bson.A{
				bson.D{{"operationType", "insert"}},
				bson.D{{"operationType", "delete"}}}}},
		}}}
	} else if evntLsnr.settings.ListenRemove && evntLsnr.settings.ListenUpdate && !evntLsnr.settings.ListenInsert {

		pipeline = mongo.Pipeline{bson.D{{"$match", bson.D{{"$or",
			bson.A{
				bson.D{{"operationType", "delete"}},
				bson.D{{"operationType", "update"}}}}},
		}}}
	} else {
		pipeline = mongo.Pipeline{}
	}

	db := mclient.Database(evntLsnr.database)
	var stream *mongo.ChangeStream
	var err error
	if evntLsnr.settings.Collection != "" {
		coll := db.Collection(evntLsnr.settings.Collection)

		evntLsnr.logger.Infof("listening on collection:: %s in Database:: %s", evntLsnr.settings.Collection, evntLsnr.database)

		stream, err = coll.Watch(context.Background(), pipeline)
	} else { // listening on database if no collection name specified
		evntLsnr.logger.Infof("listening on all collections of database:: %s", evntLsnr.database)
		stream, err = db.Watch(context.Background(), pipeline)
	}

	if err != nil {
		evntLsnr.logger.Errorf("Failed to open change stream %s", err)
		return err
	}
	evntLsnr.stream = stream

	return nil
}

func (evntLsnr *EventListener) listen() {
	evntLsnr.logger.Infof("============listening====")
	for {
		select {
		case <-evntLsnr.done:
			evntLsnr.logger.Infof("stopped listening...")
			// Exit
			return
		default:
			ok := evntLsnr.stream.Next(context.Background())
			if ok {
				var res bson.D
				err := evntLsnr.stream.Decode(&res)
				if err != nil {
					evntLsnr.logger.Errorf("got error while decoding stream %s", err)
				}
				if len(res) == 0 {
					evntLsnr.logger.Infof("result is empty, was expecting change document")
				}
				stringOp := evntLsnr.stream.Current.String()
				go evntLsnr.process(stringOp)
			} else {
				err := evntLsnr.stream.Err()
				if err != nil {
					//if err is not nil, it means something bad happened, let's finish our func
					evntLsnr.logger.Errorf("Error while listening to the MongoDB event stream %s", err)
					evntLsnr.stream.Close(context.Background())
					evntLsnr.logger.Infof("Stopped Listening...stream closed")
					return
				}
			}
		}
	}

}
func (evntLsnr *EventListener) process(stringOp string) {
	evntLsnr.logger.Infof("started processing record...")
	var outputRoot = map[string]interface{}{}
	var jsonOp map[string]interface{}
	err := json.Unmarshal([]byte(stringOp), &jsonOp)
	if err != nil {
		evntLsnr.logger.Errorf("got error while unmarshelling %s", err)
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
		evntLsnr.logger.Errorf("Failed to process record from collection [%s], due to error - %s", evntLsnr.settings.Collection, err.Error())
	} else {
		// record is successfully processed.
		evntLsnr.logger.Infof("Record from collection [%s] in Database [%s] is successfully processed", evntLsnr.settings.Collection, evntLsnr.database)
	}

}

// Stop implements trigger.Trigger.Stop
func (t *Trigger) Stop() error {
	t.logger.Infof("Stopping Trigger - %s", t.id)
	t.mclient.Disconnect(context.Background())
	t.logger.Infof("Trigger - %s  stopped", t.id)
	return nil
}
