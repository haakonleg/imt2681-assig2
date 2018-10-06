package paragliding

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"

	igc "github.com/marni/goigc"
)

func registerTrack(req *Request, db *Database) {
	var request struct {
		URL string `json:"url"`
	}

	// Get the JSON post request
	if err := req.ParseJSONRequest(&request); err != nil {
		req.SendError("Error parsing JSON request", 400)
		return
	}

	// Check that the supplied link is valid
	if valid := ensureIGCLink(request.URL); !valid {
		req.SendError("This is not a valid IGC link", 400)
		return
	}

	// Parse the IGC file
	igc, err := igc.ParseLocation(request.URL)
	if err != nil {
		fmt.Println(err)
		req.SendError("Error parsing IGC track", 400)
		return
	}

	newTrack := createTrack(&igc)
	if err := db.insertObject(newTrack, TRACKS); err != nil {
		req.SendError("Internal database error", 500)
		return
	}

	response := struct {
		ID string `json:"id"`
	}{ID: newTrack.ID.Hex()}
	req.SendJSON(&response)
}

// Ensures that a link points to an IGC resource (but just that it is a valid URL and has an igc extension)
func ensureIGCLink(link string) bool {
	if _, err := url.ParseRequestURI(link); err != nil {
		return false
	}

	ext := strings.ToLower(path.Ext(link))
	if ext != ".igc" {
		return false
	}
	return true
}

// Routes the /track request to handlers
func handleTrackRequest(req *Request, db *Database, path string) {
	// GET /api/track and POST /api/track
	if match, _ := regexp.MatchString("track[/]?$", path); match {
		switch req.r.Method {
		case "GET":
			req.SendText("GET /track")
		case "POST":
			registerTrack(req, db)
		}
		return
	}

	http.NotFound(req.w, req.r)
}
