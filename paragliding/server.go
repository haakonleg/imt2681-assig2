package paragliding

import (
	"fmt"
	"log"
)

// Database
var db Db

func StartServer(mongoUrl string) {
	// Connect to database
	err := db.CreateConnection(mongoUrl)
	if err != nil {
		log.Fatal(err.Error())
	} else {
		fmt.Println("Connected to mongoDB")
	}
}
