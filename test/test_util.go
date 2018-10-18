package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/haakonleg/imt2681-assig2/paragliding"
)

const listenPort = "8080"

func init() {
	startServer()
}

func startServer() {
	go func() {
		app := paragliding.App{
			MongoURL:    "mongodb://testuser1:testpass1@ds223063.mlab.com:23063/imt2681-assig2",
			ListenPort:  listenPort,
			DBName:      "imt2681-assig2",
			TickerLimit: 5}
		app.StartServer()
	}()
	time.Sleep(1000 * time.Millisecond)
}

func sendPostRequest(path string, requestBody interface{}, responseBody interface{}) error {
	reqBytes, _ := json.Marshal(requestBody)
	body := bytes.NewBuffer(reqBytes)

	resp, err := http.Post("http://:"+listenPort+path, "application/json", body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("Got status code" + strconv.Itoa(resp.StatusCode))
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(respBytes, responseBody)
	if err != nil {
		return err
	}
	return nil
}

func sendGetRequest(path string, responseBody interface{}, isJSON bool) error {
	resp, err := http.Get("http://:" + listenPort + path)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(path + " got status code" + strconv.Itoa(resp.StatusCode))
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if isJSON {
		err = json.Unmarshal(respBytes, responseBody)
		if err != nil {
			return err
		}
	} else {
		*responseBody.(*string) = string(respBytes)
	}

	return nil
}

func isInArr(arr []string, elem string) bool {
	for _, e := range arr {
		if e == elem {
			return true
		}
	}
	return false
}
