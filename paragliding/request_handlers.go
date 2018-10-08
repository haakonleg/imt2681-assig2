package paragliding

import (
	"net/http"
	"regexp"
	"strconv"
)

// matchers
var reTrack = regexp.MustCompile("^track/?(/[a-z0-9]{24})?/?(/(pilot|glider|glider_id|track_length|H_date|track_src_url))?/?$")
var reTicker = regexp.MustCompile("^ticker/?(/latest)?/?(/[0-9]+)?/?$")
var reWebhook = regexp.MustCompile("^webhook/new_track/?(/[a-z0-9]{24})?/?$")

// Routes the /api/track/ requests to handlers
func handleTrackRequest(req *Request, db *Database, path string) {
	if match := reTrack.FindStringSubmatch(path); match != nil {
		if len(match[1]) == 0 && len(match[2]) == 0 {
			switch req.method {
			case GET:
				getAllTracks(req, db)
			case POST:
				registerTrack(req, db)
			}
			return
		} else if len(match[1]) != 0 && len(match[2]) == 0 && req.method == GET {
			getTrack(req, db, match[1][1:])
			return
		} else if len(match[1]) != 0 && len(match[2]) != 0 && req.method == GET {
			getTrackField(req, db, match[1][1:], match[2][1:])
			return
		}
	}

	http.NotFound(req.w, req.r)
}

// Routes the /api/ticker/ requests to handlers
func handleTickerRequest(req *Request, db *Database, path string, tickerLimit int64) {
	// This regex matches all paths in a single regex by checking which capture groups a nonzero length (in other words they were matched)
	if req.method == GET {
		if match := reTicker.FindStringSubmatch(path); match != nil {
			if len(match[1]) == 0 && len(match[2]) == 0 {
				getTicker(req, db, 0, tickerLimit)
				return
			} else if len(match[1]) != 0 && len(match[2]) == 0 {
				latestTimestamp(req, db)
				return
			} else if len(match[1]) == 0 && len(match[2]) != 0 {
				timestamp, err := strconv.ParseInt(match[2][1:], 10, 64)
				if err != nil {
					req.SendError("Invalid timestamp", http.StatusBadRequest)
					return
				}
				getTicker(req, db, timestamp, tickerLimit)
				return
			}
		}
	}

	http.NotFound(req.w, req.r)
}

// Routes the /api/webhook/ requests to handlers
func handleWebhookRequest(req *Request, db *Database, path string) {
	// Match all webhook requests in one regex by checking if capture group is non-zero
	if match := reWebhook.FindStringSubmatch(path); match != nil {
		// POST /api/webhook/new_track/
		if len(match[1]) == 0 && req.method == POST {
			registerWebhook(req, db)
			return
			// GET /api/webhook/new/track/{webhook_id}
		} else if req.method == GET {
			getWebhook(req, db)
			return
			// DELETE /api/webhook/new/track/{webhook_id}
		} else if req.method == DELETE {
			deleteWebhook(req, db)
			return
		}
	}

	http.NotFound(req.w, req.r)
}
