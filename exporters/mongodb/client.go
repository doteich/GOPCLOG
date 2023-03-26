package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/doteich/OPC-UA-Logger/exporters/logging"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var ctx context.Context
var MongoClient *mongo.Client

type DBEntry struct {
	Meta      MetaData    `bson:"meta"`
	Timestamp time.Time   `bson:"ts"`
	Value     interface{} `bson:"value"`
}

type MetaData struct {
	NodeId   string `bson:"nodeId"`
	NodeName string `bson:"nodeName"`
	LogName  string `bson:"logName"`
	Server   string `bson:"server"`
	DataType string `bson:"dataType"`
}

func CreateConnection(namespace string, username string, password string, connectionString string, connectionType string) {
	ctx = context.Background()

	var connectionURL string

	if connectionType == "srv" {
		connectionURL = "mongodb+srv://" + username + ":" + password + "@" + connectionString
	} else {
		connectionURL = "mongodb://" + username + ":" + password + "@" + connectionString
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionURL))

	if err != nil {
		logging.LogError(err, "Failed to connect to MongoDB", "mongodb")
	}

	MongoClient = client

	err = MongoClient.Ping(ctx, readpref.Primary())

	if err != nil {
		logging.LogError(err, "Failed to ping to MongoDB", "mongodb")
		panic(err)
	}

	opts := options.CreateCollection().SetTimeSeriesOptions(options.TimeSeries().SetGranularity("seconds").SetMetaField("meta").SetTimeField("ts"))

	MongoClient.Database("machine-data").CreateCollection(ctx, namespace, opts)

	logging.LogGeneric("info", "Successfully connected and pinged mongodb: "+connectionString, "mongodb")

}

func CheckConn() {

	ticker := time.NewTicker(10 * time.Second)

	for {

		<-ticker.C

		err := MongoClient.Ping(ctx, readpref.Primary())

		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("PING /n")
	}
}

func WriteData(nodeId string, nodeName string, value interface{}, timestamp time.Time, logName string, server string, datatype string, namespace string) {
	coll := MongoClient.Database("machine-data").Collection(namespace)

	meta := MetaData{NodeId: nodeId, NodeName: nodeName, LogName: logName, Server: server, DataType: datatype}
	newEntry := DBEntry{Meta: meta, Timestamp: timestamp, Value: value}

	_, err := coll.InsertOne(ctx, newEntry)

	if err != nil {
		logging.LogError(err, "Failed to insert document", "mongodb")
	}

}
