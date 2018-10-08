package paragliding

import (
	"net/http"
	"regexp"
	"strconv"
)

// Routes the /api/track/ requests to handlers
func handleTrackRequest(req *Request, db *Database, path string) {
	// GET/POST /api/track
	if match, _ := regexp.MatchString("^track/?$", path); match {
		switch req.r.Method {
		case "GET":
			getAllTracks(req, db)
		case "POST":
			registerTrack(req, db)
		}
		return
	}

	// GET /api/track/{id}
	if req.r.Method == "GET" {
		if match := regexp.MustCompile("^track/([a-z0-9]{24})/?$").FindStringSubmatch(path); match != nil {
			getTrack(req, db, match[1])
			return
		}

		// GET track/{id}/{field}
		if match := regexp.MustCompile("^track/([a-z0-9]{24})/(pilot|glider|glider_id|track_length|H_date|track_src_url)/?$").FindStringSubmatch(path); match != nil {
			getTrackField(req, db, match[1], match[2])
			return
		}
	}

	http.NotFound(req.w, req.r)
}

// Routes the /api/ticker/ requests to handlers
func handleTickerRequest(req *Request, db *Database, path string, tickerLimit int64) {
	// This regex matches all paths in a single regex by checking which capture groups a nonzero length (in other words they were matched)
	if req.r.Method == "GET" {
		if match := regexp.MustCompile("^ticker/?(latest)?([0-9]+)?/?$").FindStringSubmatch(path); match != nil {
			if len(match[2]) != 0 {
				timestamp, err := strconv.ParseInt(match[2], 10, 64)
				if err != nil {
					req.SendError("Invalid timestamp", http.StatusBadRequest)
					return
				}

				getTicker(req, db, timestamp, tickerLimit)
				return
			} else if len(match[1]) != 0 {
				latestTimestamp(req, db)
				return
			} else {
				getTicker(req, db, 0, tickerLimit)
				return
			}
		}
	}

	http.NotFound(req.w, req.r)
}

// Routes the /api/webhook/ requests to handlers
func handleWebhookRequest(req *Request, db *Database, path string) {
	// Match all webhook requests in one regex by checking if capture group is non-zero
	if match := regexp.MustCompile("^webhook/new_track/?(/[a-z0-9]{24})?/?$").FindStringSubmatch(path); match != nil {
		// POST /api/webhook/new_track/
		if len(match[1]) == 0 && req.r.Method == "POST" {
			registerWebhook(req, db)
			return
			// GET /api/webhook/new/track/{webhook_id}
		} else if req.r.Method == "GET" {
			getWebhook(req, db)
			return
			// DELETE /api/webhook/new/track/{webhook_id}
		} else if req.r.Method == "DELETE" {
			deleteWebhook(req, db)
			return
		}
	}

	http.NotFound(req.w, req.r)
}
