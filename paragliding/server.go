package paragliding

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
)

const (
	rootPath = "/paragliding/"
	apiPath  = "/paragliding/api/"
)

// App must be instantiated with the url to the mongodb database, database name and the port for the API to listen on
type App struct {
	MongoURL   string
	DBName     string
	ListenPort string

	db Database
}

// Route the API request to handlers
func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := Request{w, r}
	path := strings.TrimPrefix(r.URL.Path, apiPath)

	// GET /api
	if len(path) == 0 && req.r.Method == "GET" {
		ApiInfo(&req)
		return
	}

	// Get path
	match := regexp.MustCompile("^([^/]+)").FindStringSubmatch(path)
	switch match[0] {
	case "track":
		handleTrackRequest(&req, &app.db, path)
	case "ticker":
	case "webhook":
	case "admin":
	default:
		http.NotFound(req.w, req.r)
	}
}

// StartServer starts listening and serving the API server
func (app *App) StartServer() {
	// Try connect to mongoDB
	app.db = Database{MongoURL: app.MongoURL, DBName: app.DBName}
	if err := app.db.createConnection(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Connected to mongoDB")
	}

	// Configure redirect and 404 not found handler, and direct requests to the API path to the handler
	http.Handle("/", http.NotFoundHandler())
	http.Handle(rootPath, http.RedirectHandler(apiPath, 301))
	http.Handle(apiPath, app)

	// Start listen
	fmt.Printf("Server listening on port %s\n", app.ListenPort)
	if err := http.ListenAndServe(":"+app.ListenPort, nil); err != nil {
		log.Fatal(err.Error())
	}
}