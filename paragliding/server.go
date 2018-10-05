package paragliding

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const rootPath = "/paragliding/"
const apiPath = "/paragliding/api/"

// Routes the request to handlers
func routeRequest(db *Db) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get routes
		route := strings.TrimPrefix(r.URL.Path, apiPath)
		req := Request{w, r}

		// GET /api, this is just the base API path
		if len(route) == 0 && req.GetMethod() == GET {
			ApiInfo(&req)
			return
		}

		http.NotFound(w, r)
	})
}

func StartServer(mongoUrl string, port string) {
	// Connect to database
	db := Db{mongoURL: mongoUrl}
	err := db.CreateConnection()
	if err != nil {
		log.Fatal(err.Error())
	} else {
		fmt.Println("Connected to mongoDB")
	}

	// Handle routes
	http.Handle("/", http.NotFoundHandler())
	http.Handle(rootPath, http.RedirectHandler(apiPath, 301))
	http.Handle(apiPath, routeRequest(&db))

	// Start listen
	fmt.Printf("Server listening on port %s\n", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}
