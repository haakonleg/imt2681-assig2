package router

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type ResponseType int

const (
	JSON ResponseType = iota
	TEXT
)

// Request contains the context of the HTTP request, it also has some helper methods
type Request struct {
	W    http.ResponseWriter
	R    *http.Request
	Vars map[string]interface{}
}

// SetResponseType sets the response type of the HTTP response
func (req *Request) SetResponseType(responseType ResponseType) {
	switch responseType {
	case JSON:
		req.W.Header().Set("Content-Type", "application/json")
	case TEXT:
		req.W.Header().Set("Content-Type", "text/plain")
	}
}

// SendJSON sends a json response, parameter jsonStruct is a struct
// that contains the json fields
func (req *Request) SendJSON(jsonStruct interface{}, statusCode int) {
	req.SetResponseType(JSON)
	res, err := json.Marshal(jsonStruct)
	if err != nil {
		log.Fatal(err)
	}
	req.W.WriteHeader(statusCode)
	req.W.Write(res)
}

// SendText sends a plain text message as response to the request
func (req *Request) SendText(text string, statusCode int) {
	req.SetResponseType(TEXT)
	req.W.WriteHeader(statusCode)
	fmt.Fprint(req.W, text)
}

// ParseJSONRequest parses the JSON contents of a POST request, takes
// a struct as parameter with JSON fields and writes the contents to it
func (req *Request) ParseJSONRequest(jsonStruct interface{}) error {
	body, err := ioutil.ReadAll(req.R.Body)
	defer req.R.Body.Close()

	if err != nil {
		fmt.Println(err)
		return err
	}

	err = json.Unmarshal(body, jsonStruct)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// SendError sends a JSON error message in response to the request in case an error occured
func (req *Request) SendError(message string, statusCode int) {
	err := struct {
		Error string `json:"error"`
	}{Error: message}
	req.SendJSON(&err, statusCode)
}

// Redirect redirects the request to another URL
func (req *Request) Redirect(path string) {
	http.Redirect(req.W, req.R, path, http.StatusPermanentRedirect)
}
