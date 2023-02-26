package mongodb

import (
	"context"
	"fmt"
	"time"

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

func CreateConnection(namespace string, username string, password string, url string, port int) {
	ctx = context.Background()

	connectionURL := "mongodb://" + username + ":" + password + "@" + url + ":" + fmt.Sprint(port) + "/?directConnection=true"

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionURL))

	MongoClient = client

	if err != nil {
		fmt.Printf("Error while create MongoDB Connection: %v", err)
	}

	err = MongoClient.Ping(ctx, readpref.Primary())

	if err != nil {
		panic(err)
	}

	opts := options.CreateCollection().SetTimeSeriesOptions(options.TimeSeries().SetGranularity("seconds").SetMetaField("meta").SetTimeField("ts"))

	MongoClient.Database("machine-data").CreateCollection(ctx, namespace, opts)

	fmt.Println("Successfully connected and pinged.")
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
		fmt.Println(err)
	}

}
