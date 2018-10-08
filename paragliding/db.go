package paragliding

import (
	"context"
	"fmt"
	"log"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
)

// An enum of the database collections
type databaseCollection int

const (
	TRACKS databaseCollection = iota
	WEBHOOKS
)

// Stringer for databaseCollection type
func (dc databaseCollection) String() string {
	switch dc {
	case TRACKS:
		return "tracks"
	case WEBHOOKS:
		return "webhooks"
	}
	return ""
}

// Database contains the mongoDB database context, it also has helper methods for connecting to and querying the database
type Database struct {
	MongoURL string
	DBName   string

	client   *mongo.Client
	database *mongo.Database
}

// CreateConnection creates a connection to the mongoDB server
func (db *Database) createConnection() {
	client, err := mongo.Connect(context.Background(), db.MongoURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	db.client = client
	db.database = db.client.Database(db.DBName)
	db.createTimestampIndex()
}

func (db *Database) InsertObject(collection databaseCollection, object interface{}) (string, error) {
	col := db.database.Collection(collection.String())
	res, err := col.InsertOne(context.Background(), object)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return res.InsertedID.(*bson.Element).Value().ObjectID().Hex(), nil
}

func (db *Database) Find(collection databaseCollection, filter interface{}, opts []findopt.Find) ([]interface{}, error) {
	col := db.database.Collection(collection.String())
	cur, err := col.Find(context.Background(), filter, opts...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer cur.Close(context.Background())

	results := make([]interface{}, 0)
	for cur.Next(context.Background()) {
		switch collection {
		case TRACKS:
			var track Track
			if err := cur.Decode(&track); err != nil {
				return nil, err
			}
			results = append(results, track)
		case WEBHOOKS:
			var webhook Webhook
			if err := cur.Decode(&webhook); err != nil {
				return nil, err
			}
			results = append(results, webhook)
		}
	}

	return results, nil
}

func (db *Database) Update(collection databaseCollection, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	col := db.database.Collection(collection.String())
	ur, err := col.UpdateMany(context.Background(), filter, update)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ur, nil
}

// Creates a descending index on the timestamp field in tracks, to be able to support certain queries and better performance
func (db *Database) createTimestampIndex() {
	indexView := db.database.Collection(TRACKS.String()).Indexes()

	indexModel := mongo.IndexModel{
		Keys: bson.NewDocument(bson.EC.Int32("ts", -1))}

	_, err := indexView.CreateOne(context.Background(), indexModel, nil)
	if err != nil {
		log.Fatal(err)
	}
}
