package test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/haakonleg/imt2681-assig2/mdb"
	"github.com/haakonleg/imt2681-assig2/track"
)

func TestPostTrack(t *testing.T) {
	t.Parallel()
	fmt.Println("Running test TestPostTrack...")

	// Register two tracks and retrieve IDs
	var ids [2]string
	res, err := postTrack("http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc")
	if err != nil {
		t.Fatalf(err.Error())
	}
	ids[0] = res.ID
	res, err = postTrack("http://skypolaris.org/wp-content/uploads/IGS%20Files/Jarez%20to%20Senegal.igc")
	if err != nil {
		t.Fatalf(err.Error())
	}
	ids[1] = res.ID

	// Get all track IDs
	trackIDs, err := getAllTracks()
	if err != nil {
		t.Fatalf(err.Error())
	}

	// The array should contain the two last track IDs registered
	if !isInArr(trackIDs, ids[0]) {
		t.Fatalf("Expected %s to be in array", ids[0])
	}
	if !isInArr(trackIDs, ids[1]) {
		t.Fatalf("Expected %s to be in array", ids[1])
	}
}

func TestGetTrack(t *testing.T) {
	t.Parallel()
	fmt.Println("Running test TestGetTrack...")

	// Register a track
	res, err := postTrack("http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc")
	if err != nil {
		t.Fatalf(err.Error())
	}

	id := res.ID

	track, err := getTrack(id)

	// Check that the track matches the expected data
	expect := &mdb.Track{
		HDate:       "2016-02-19 00:00:00 +0000 UTC",
		Pilot:       "Miguel Angel Gordillo",
		Glider:      "RV8",
		GliderID:    "EC-XLL",
		TrackLength: "443.26km",
		TrackSrcURL: "http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc"}

	if !reflect.DeepEqual(track, expect) {
		t.Fatalf("Expected %v. Got: %v", expect, track)
	}
}

func TestGetTrackField(t *testing.T) {
	t.Parallel()
	fmt.Println("Running test TestGetTrackField...")

	// Register a track
	res, err := postTrack("http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc")
	if err != nil {
		t.Fatalf(err.Error())
	}
	id := res.ID

	testCases := [][]string{
		{"pilot", "Miguel Angel Gordillo"}, {"glider", "RV8"},
		{"glider_id", "EC-XLL"}, {"track_length", "443.26km"},
		{"H_date", "2016-02-19 00:00:00 +0000 UTC"}, {"track_src_url", "http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc"}}

	for _, testCase := range testCases {
		res, err := getTrackField(id, testCase[0])
		if err != nil {
			t.Fatalf(err.Error())
		}
		if res != testCase[1] {
			t.Fatalf("Expected: %s. Got: %s", testCase[1], res)
		}
	}
}

func postTrack(url string) (*track.PostTrackResponse, error) {
	response := new(track.PostTrackResponse)
	if err := sendPostRequest("/paragliding/api/track", &track.PostTrackRequest{URL: url}, response); err != nil {
		return nil, err
	}
	return response, nil
}

func getTrack(id string) (*mdb.Track, error) {
	response := new(mdb.Track)
	if err := sendGetRequest("/paragliding/api/track/"+id, response, true); err != nil {
		return nil, err
	}
	return response, nil
}

func getAllTracks() ([]string, error) {
	var response []string
	if err := sendGetRequest("/paragliding/api/track", &response, true); err != nil {
		return nil, err
	}
	return response, nil
}

func getTrackField(id string, field string) (string, error) {
	var response string
	if err := sendGetRequest("/paragliding/api/track/"+id+"/"+field, &response, false); err != nil {
		return "", err
	}
	return response, nil
}
