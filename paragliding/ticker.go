package paragliding

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/mongodb/mongo-go-driver/bson"

	"github.com/mongodb/mongo-go-driver/mongo/findopt"
)

func latestTimestamp(req *Request, db *Database) {
	// Sort by timestamp in decsending order, and limit to one result
	sort := bson.NewDocument(bson.EC.Int64("ts", -1))
	findopts := []findopt.Find{
		findopt.Limit(1),
		findopt.Sort(sort)}

	tracks, err := db.findTracks(nil, findopts)
	if err != nil {
		req.SendError("Internal database error", http.StatusInternalServerError)
		return
	}
	if len(tracks) < 1 {
		req.SendError("No tracks added yet", http.StatusBadRequest)
		return
	}

	req.SendText(strconv.FormatInt(tracks[0].Ts, 10))
}

func handleTickerRequest(req *Request, db *Database, path string) {
	// This regex matches all paths in a single regex by checking which capture groups have zero length (in other words they were matched)
	if req.r.Method == "GET" {
		if match := regexp.MustCompile("^ticker/?(latest)?([0-9]+)?/?$").FindStringSubmatch(path); match != nil {
			if len(match[2]) != 0 {
				req.SendText("GET /api/ticker/timestamp")
				return
			} else if len(match[1]) != 0 {
				latestTimestamp(req, db)
				return
			} else {
				req.SendText("GET /api/ticker")
				return
			}
		}
	}

	http.NotFound(req.w, req.r)
}
