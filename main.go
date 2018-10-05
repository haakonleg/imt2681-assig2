package main

import (
	"os"
	"strings"

	"github.com/haakonleg/imt2681-assig2/paragliding"
)

const dbUser = "testuser1"
const dbPass = "testpass1"
const mongoURL = "mongodb://<dbuser>:<dbpassword>@ds223063.mlab.com:23063/imt2681-assig2"

const defaultPort = "8080"

// Main starts the paragliding server by supplying the mongodb connection URL
func main() {
	url := strings.Replace(mongoURL, "<dbuser>", dbUser, 1)
	url = strings.Replace(url, "<dbpassword>", dbPass, 1)

	// Get port
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = defaultPort
	}

	paragliding.StartServer(url, port)
}
