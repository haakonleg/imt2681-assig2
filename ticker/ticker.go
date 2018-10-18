package ticker

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/haakonleg/imt2681-assig2/mdb"
	"github.com/haakonleg/imt2681-assig2/router"
	"github.com/mongodb/mongo-go-driver/bson"

	"github.com/mongodb/mongo-go-driver/mongo/findopt"
)

type GetTickerResponse struct {
	TLatest    int64    `json:"t_latest"`
	TStart     int64    `json:"t_start"`
	TStop      int64    `json:"t_stop"`
	Tracks     []string `json:"tracks"`
	Processing int64    `json:"processing"`
}

type TickerHandler struct {
	tickerLimit int64
	db          *mdb.Database
}

func NewTickerHandler(tickerLimit int64, db *mdb.Database) *TickerHandler {
	return &TickerHandler{
		tickerLimit: tickerLimit,
		db:          db}
}

// Finds the timestamp of the latest added track in the database
func (th *TickerHandler) findLatestTimestamp() (int64, error) {
	// Sort by timestamp in decsending order, and limit to one result
	findopts := []findopt.Find{
		findopt.Sort(bson.NewDocument(bson.EC.Int64("ts", -1))),
		findopt.Projection(bson.NewDocument(bson.EC.Int64("ts", 1))),
		findopt.Limit(1)}

	tracks := make([]*mdb.Track, 0)
	if err := th.db.Find(mdb.TRACKS, nil, findopts, &tracks); err != nil {
		return -1, errors.New("Internal database error")
	}
	if len(tracks) < 1 {
		return -1, errors.New("No tracks added yet")
	}

	return tracks[0].Ts, nil
}

func (th *TickerHandler) GetLatestTimestamp(req *router.Request) {
	ts, err := th.findLatestTimestamp()
	if err != nil {
		req.SendError(err.Error(), http.StatusBadRequest)
		return
	}

	req.SendText(strconv.FormatInt(ts, 10), http.StatusOK)
}

func ValidateTimestamp(variable string) (bool, interface{}) {
	ts, err := strconv.ParseInt(variable, 10, 64)
	if err != nil {
		return false, nil
	}
	return true, ts
}

// GET /api/ticker
func (th *TickerHandler) GetTicker(req *router.Request) {
	// Check if there is a timestamp limit specified in the request
	timestampLimit, ok := req.Vars["timestamp"]
	if !ok {
		timestampLimit = int64(0)
	}

	// Measure time
	start := time.Now()

	// The response to send
	ticker := new(GetTickerResponse)

	// Get latest timestamp
	latestTs, err := th.findLatestTimestamp()
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
		findopt.Max(bson.NewDocument(bson.EC.Int64("ts", timestampLimit.(int64)))),
		findopt.Limit(th.tickerLimit)}

	tracks := make([]*mdb.Track, 0)
	if err := th.db.Find(mdb.TRACKS, nil, findopts, &tracks); err != nil {
		req.SendError("Internal database error", http.StatusInternalServerError)
		return
	}
	if len(tracks) < 1 {
		req.SendError("No more tracks", http.StatusBadRequest)
		return
	}

	// Add start and stop timestamps and IDs to struct
	ticker.TStart = tracks[0].Ts
	ticker.TStop = tracks[len(tracks)-1].Ts
	for _, tr := range tracks {
		id := tr.ID.Hex()
		ticker.Tracks = append(ticker.Tracks, id)
	}

	// Calculate time it took
	ticker.Processing = int64(time.Since(start) / time.Millisecond)

	req.SendJSON(ticker, http.StatusOK)
}
