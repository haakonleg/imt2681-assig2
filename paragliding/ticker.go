package paragliding

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"

	"github.com/mongodb/mongo-go-driver/mongo/findopt"
)

// Finds the timestamp of the latest added track in the database
func findLatestTimestamp(db *Database) (int64, error) {
	// Sort by timestamp in decsending order, and limit to one result
	findopts := []findopt.Find{
		findopt.Sort(bson.NewDocument(bson.EC.Int64("ts", -1))),
		findopt.Projection(bson.NewDocument(bson.EC.Int64("ts", 1))),
		findopt.Limit(1)}

	tracks, err := db.Find(TRACKS, nil, findopts)
	if err != nil {
		return -1, errors.New("Internal database error")
	}
	if len(tracks) < 1 {
		return -1, errors.New("No tracks added yet")
	}

	return tracks[0].(Track).Ts, nil
}

func latestTimestamp(req *Request, db *Database) {
	ts, err := findLatestTimestamp(db)
	if err != nil {
		req.SendError(err.Error(), http.StatusBadRequest)
		return
	}

	req.SendText(strconv.FormatInt(ts, 10))
}

// GET /api/ticker
func getTicker(req *Request, db *Database, timestampLimit int64, tickerLimit int64) {
	// Measure time
	start := time.Now()

	// The response to send
	var ticker struct {
		TLatest    int64    `json:"t_latest"`
		TStart     int64    `json:"t_start"`
		TStop      int64    `json:"t_stop"`
		Tracks     []string `json:"tracks"`
		Processing int64    `json:"processing"`
	}

	// Get latest timestamp
	latestTs, err := findLatestTimestamp(db)
	if err != nil {
		req.SendError(err.Error(), http.StatusBadRequest)
		return
	}
	// Add latest timestamp to struct
	ticker.TLatest = latestTs

	// Sort by timestamp oldest first, set minimum value for timestamp, limit to 5 results
	findopts := []findopt.Find{
		findopt.Projection(bson.NewDocument(bson.EC.Int64("ts", 1))),
		findopt.Sort(bson.NewDocument(bson.EC.Int64("ts", 1))),
		findopt.Max(bson.NewDocument(bson.EC.Int64("ts", timestampLimit))),
		findopt.Limit(tickerLimit)}

	tracks, err := db.Find(TRACKS, nil, findopts)
	if err != nil {
		req.SendError("Internal database error", http.StatusInternalServerError)
		return
	}
	if len(tracks) < 1 {
		req.SendError("No tracks added yet", http.StatusBadRequest)
		return
	}

	// Add start and stop timestamps and IDs to struct
	ticker.TStart = tracks[0].(Track).Ts
	ticker.TStop = tracks[len(tracks)-1].(Track).Ts
	for _, track := range tracks {
		id := track.(Track).ID.Hex()
		ticker.Tracks = append(ticker.Tracks, id)
	}

	// Calculate time it took
	ticker.Processing = int64(time.Since(start) / time.Millisecond)

	req.SendJSON(&ticker, http.StatusOK)
}
