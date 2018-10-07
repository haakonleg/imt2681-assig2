package paragliding

import (
	"context"
	"fmt"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
)

// Name of the database and the collections we have
type databaseCollection int

const (
	TRACKS databaseCollection = iota
)

// Stringer function for database collections
func (dc databaseCollection) String() string {
	switch dc {
	case TRACKS:
		return "tracks"
	default:
		return ""
	}
}

// Database contains the mongoDB database context, it also has helper methods for connecting to and querying the database
type Database struct {
	MongoURL string
	DBName   string

	client   *mongo.Client
	database *mongo.Database
}

// CreateConnection creates a connection to the mongoDB server
func (db *Database) createConnection() error {
	client, err := mongo.Connect(context.Background(), db.MongoURL, nil)
	if err != nil {
		return err
	}
	db.client = client
	db.database = db.client.Database(db.DBName)
	return nil
}

func (db *Database) insertObject(object interface{}, col databaseCollection) (string, error) {
	collection := db.database.Collection(col.String())
	res, err := collection.InsertOne(context.Background(), object)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return res.InsertedID.(*bson.Element).Value().ObjectID().Hex(), nil
}

func (db *Database) findTracks(filter interface{}, opts []findopt.Find) ([]Track, error) {
	collection := db.database.Collection(TRACKS.String())
	cur, err := collection.Find(context.Background(), filter, opts...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer cur.Close(context.Background())

	tracks := make([]Track, 0)
	for cur.Next(context.Background()) {
		var elem Track
		if err := cur.Decode(&elem); err != nil {
			fmt.Println(err)
			return nil, err
		}
		tracks = append(tracks, elem)
	}

	return tracks, nil
}
