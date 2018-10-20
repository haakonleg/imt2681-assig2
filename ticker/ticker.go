package ticker

import (
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
func findLatestTimestamp(db *mdb.Database) (int64, *router.Error) {
	// Sort by timestamp in decsending order, and limit to one result
	findopts := []findopt.Find{
		findopt.Sort(bson.NewDocument(bson.EC.Int64("ts", -1))),
		findopt.Projection(bson.NewDocument(bson.EC.Int64("ts", 1))),
		findopt.Limit(1)}

	tracks := make([]*mdb.Track, 0)
	if err := db.Find(mdb.TRACKS, nil, findopts, &tracks); err != nil {
		return -1, &router.Error{StatusCode: http.StatusInternalServerError, Message: "Internal database error"}
	}
	if len(tracks) < 1 {
		return -1, &router.Error{StatusCode: http.StatusBadRequest, Message: "No tracks added yet"}
	}

	return tracks[0].Ts, nil
}

func (th *TickerHandler) GetLatestTimestamp(req *router.Request) {
	ts, err := findLatestTimestamp(th.db)
	if err != nil {
		req.SendError(err)
		return
	}

	req.SendText(strconv.FormatInt(ts, 10), http.StatusOK)
}

func MakeTicker(db *mdb.Database, tickerLimit, timestampLimit int64) (*GetTickerResponse, *router.Error) {
	ticker := new(GetTickerResponse)

	// Measure time
	start := time.Now()

	// Get latest timestamp
	latestTs, err := findLatestTimestamp(db)
	if err != nil {
		return nil, err
	}

	// Add latest timestamp to struct
	ticker.TLatest = latestTs

	// Sort by timestamp oldest first, set minimum value for timestamp, limit results to tickerLimit, if it is over 0
	findopts := []findopt.Find{
		findopt.Projection(bson.NewDocument(bson.EC.Int64("ts", 1))),
		findopt.Sort(bson.NewDocument(bson.EC.Int64("ts", 1))),
		findopt.Max(bson.NewDocument(bson.EC.Int64("ts", timestampLimit)))}

	if tickerLimit <= 0 {
		findopts = append(findopts, findopt.Limit(tickerLimit))
	}

	// Retrieve the track timestamps from DB
	tracks := make([]*mdb.Track, 0)
	if err := db.Find(mdb.TRACKS, nil, findopts, &tracks); err != nil {
		return nil, &router.Error{StatusCode: http.StatusInternalServerError, Message: "Internal database error"}
	}
	if len(tracks) < 1 {
		return nil, &router.Error{StatusCode: http.StatusBadRequest, Message: "No more tracks"}
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

	return ticker, nil
}

// GET /api/ticker
func (th *TickerHandler) GetTicker(req *router.Request) {
	// Check if there is a timestamp limit specified in the request
	timestampLimit, ok := req.Vars["timestamp"].(int64)
	if !ok {
		timestampLimit = 0
	}

	ticker, err := MakeTicker(th.db, th.tickerLimit, timestampLimit)
	if err != nil {
		req.SendError(err)
	}

	req.SendJSON(ticker, http.StatusOK)
}
