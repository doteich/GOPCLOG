package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const connectionURL = "mongodb://{{USER}}:{{PASSWORD}}@{{IP}}:{{PORT}}/?directConnection=true"

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

func CreateConnection() {
	ctx = context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionURL))

	MongoClient = client

	if err != nil {
		fmt.Printf("Error while create MongoDB Connection: %v", err)
	}

	// defer func() {
	// 	err := MongoClient.Disconnect(ctx)

	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }()

	err = MongoClient.Ping(ctx, readpref.Primary())

	if err != nil {
		panic(err)
	}

	//go CheckConn()

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

func WriteData(nodeId string, nodeName string, value interface{}, timestamp time.Time, logName string, server string, datatype string) {
	coll := MongoClient.Database("machine-data").Collection("machine1")

	meta := MetaData{NodeId: nodeId, NodeName: nodeName, LogName: logName, Server: server, DataType: datatype}
	newEntry := DBEntry{Meta: meta, Timestamp: timestamp, Value: value}

	result, err := coll.InsertOne(ctx, newEntry)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)
}
