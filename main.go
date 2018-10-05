package main

import (
	"strings"

	"github.com/haakonleg/imt2681-assig2/paragliding"
)

const dbUser = "testuser1"
const dbPass = "testpass1"
const mongoURL = "mongodb://<dbuser>:<dbpassword>@ds223063.mlab.com:23063/imt2681-assig2"

// Main starts the paragliding server by supplying the mongodb connection URL
func main() {
	url := strings.Replace(mongoURL, "<dbuser>", dbUser, 1)
	url = strings.Replace(url, "<dbpassword>", dbPass, 1)

	paragliding.StartServer(url)
}
