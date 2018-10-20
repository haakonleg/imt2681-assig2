package clocktrigger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/haakonleg/imt2681-assig2/mdb"
	"github.com/haakonleg/imt2681-assig2/util"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
)

const (
	intervalMin = 1
)

type SlackWebhook struct {
	Username string `json:"username"`
	IconURL  string `json:"icon_url"`
	Text     string `json:"text"`
}

func Start(db *mdb.Database, whURL string) {
	// Infinite loop
	for {
		time.Sleep(intervalMin * time.Minute)
		fmt.Println("Checking for new tracks...")
		checkNewTracks(db, whURL)
	}
}

// Check whether new tracks have been added since last check
func checkNewTracks(db *mdb.Database, whURL string) {
	// Check if tracks have been added in the last 10 minutes
	findopts := []findopt.Find{
		findopt.Projection(bson.NewDocument(bson.EC.Int64("ts", 1)))}
	filter := bson.NewDocument(
		bson.EC.SubDocumentFromElements("ts",
			bson.EC.Int64("$gt", util.NowMilli()-intervalMin*60000)))

	tracks := make([]*mdb.Track, 0)
	if err := db.Find(mdb.TRACKS, filter, findopts, &tracks); err != nil {
		fmt.Println(err)
		return
	}

	if len(tracks) > 0 {
		slackWebhook(tracks, whURL)
	}
}

func slackWebhook(tracks []*mdb.Track, whURL string) {
	// Build request
	buf := new(bytes.Buffer)

	buf.WriteString("New tracks: " + tracks[0].ID.Hex())
	for i := 1; i < len(tracks); i++ {
		buf.WriteString(", ")
		buf.WriteString(tracks[i].ID.Hex())
	}

	content := &SlackWebhook{
		Username: "paragliding_bot",
		IconURL:  "https://hakkon.me/images/pepe.png",
		Text:     buf.String()}

	request, _ := json.Marshal(content)
	buf.Reset()
	buf.Write(request)

	resp, err := http.Post(whURL, "application/json", buf)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp.StatusCode)
}
