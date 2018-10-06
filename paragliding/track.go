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

// GET /api/track
// Returns an array of IDs of all tracks stored in the database
func getAllTracks(req *Request, db *Database) {
	// Get all tracks in database
	tracks, err := db.findTracks(nil)
	if err != nil {
		req.SendError("Internal database error", http.StatusInternalServerError)
		return
	}

	// Retrieve all the ids into a slice
	ids := make([]string, 0, len(tracks))
	for _, track := range tracks {
		ids = append(ids, track.ID.Hex())
	}

	req.SendJSON(&ids, http.StatusOK)
}

// POST /api/track
// Register/upload a track
func registerTrack(req *Request, db *Database) {
	var request struct {
		URL string `json:"url"`
	}

	// Get the JSON post request
	if err := req.ParseJSONRequest(&request); err != nil {
		req.SendError("Error parsing JSON request", http.StatusBadRequest)
		return
	}

	// Check that the supplied link is valid
	if valid := ensureIGCLink(request.URL); !valid {
		req.SendError("This is not a valid IGC link", http.StatusBadRequest)
		return
	}

	// Parse the IGC file
	igc, err := igc.ParseLocation(request.URL)
	if err != nil {
		fmt.Println(err)
		req.SendError("Error parsing IGC track", http.StatusBadRequest)
		return
	}

	// Send response containing the ID to the inserted track
	newTrack := createTrack(&igc)

	id, err := db.insertObject(newTrack, TRACKS)
	if err != nil {
		req.SendError("Internal database error", http.StatusInternalServerError)
		return
	}

	response := struct {
		ID string `json:"id"`
	}{ID: id}
	req.SendJSON(&response, http.StatusOK)
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
			getAllTracks(req, db)
		case "POST":
			registerTrack(req, db)
		}
		return
	}

	http.NotFound(req.w, req.r)
}
