package test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/haakonleg/imt2681-assig2/ticker"
)

func TestGetLatestTimestamp(t *testing.T) {
	t.Parallel()
	fmt.Println("Running test TestGetLatestTimestamp...")

	var response string
	if err := sendGetRequest("/paragliding/api/ticker/latest", &response, false); err != nil {
		t.Fatalf(err.Error())
	}

	ts, err := strconv.ParseInt(response, 10, 64)
	if err != nil {
		t.Fatalf(err.Error())
	}

	tm := time.Unix(ts/int64(1000), 0)
	fmt.Println(tm.String())
}

func TestGetTicker(t *testing.T) {
	t.Parallel()
	fmt.Println("Running test TestGetTicker...")

	var response string
	if err := sendGetRequest("/paragliding/api/ticker/latest", &response, false); err != nil {
		t.Fatalf(err.Error())
	}

	ts, err := strconv.ParseInt(response, 10, 64)
	if err != nil {
		t.Fatalf(err.Error())
	}

	ticker := new(ticker.GetTickerResponse)
	if err := sendGetRequest("/paragliding/api/ticker", ticker, true); err != nil {
		t.Fatalf(err.Error())
	}

	if ticker.TLatest != ts {
		t.Fatalf("Expected: %d. Got: %d", ts, ticker.TLatest)
	}

	tsTest := int64(1539381600000)
	if err := sendGetRequest("/paragliding/api/ticker/1539381600000", ticker, true); err != nil {
		t.Fatalf(err.Error())
	}

	if ticker.TStart <= tsTest {
		t.Fatalf("Expected Ticker.TStart: %d to be higher than parameter timestamp %d", ticker.TStart, tsTest)
	}
}
