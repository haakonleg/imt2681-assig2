package main

import (
	"fmt"
	"log"
	"os"

	"github.com/haakonleg/imt2681-assig2/clocktrigger"
	"github.com/haakonleg/imt2681-assig2/mdb"
)

const (
	dbName = "imt2681-assig2"
)

func main() {
	whURL := os.Getenv("CTRIGGER_URL")
	if len(whURL) == 0 {
		log.Fatal("CTRIGGER_URL environment variable is not set (put slack webhook url in here)")
	}

	mongoURL := os.Getenv("PARAGLIDING_MONGO")
	if len(mongoURL) == 0 {
		log.Fatal("PARAGLIDING_MONGO environment variable is not set (put mongodb url in here)")
	}

	// Try connect to mongoDB
	db := &mdb.Database{MongoURL: mongoURL, DBName: dbName}
	db.CreateConnection()
	fmt.Println("Connected to mongoDB")

	clocktrigger.Start(db, whURL)
}
