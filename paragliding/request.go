package paragliding

import (
	"encoding/json"
	"log"
	"net/http"
)

type ResponseType int
type HTTPMethod int

const (
	JSON ResponseType = 0
	TEXT ResponseType = 1

	UNKNOWN HTTPMethod = 0
	GET     HTTPMethod = 1
	POST    HTTPMethod = 2
)

type Request struct {
	w http.ResponseWriter
	r *http.Request
}

func (req *Request) SetResponseType(responseType ResponseType) {
	switch responseType {
	case JSON:
		req.w.Header().Set("Content-Type", "application/json")
	case TEXT:
		req.w.Header().Set("Content-Type", "text/plain")
	}
}

func (req *Request) GetMethod() HTTPMethod {
	switch req.r.Method {
	case "GET":
		return GET
	case "POST":
		return POST
	}
	return UNKNOWN
}

// SendJSON sends a json response, parameter jsonStruct is a struct
// that contains the json fields
func (req *Request) SendJSON(jsonStruct interface{}) {
	req.SetResponseType(JSON)
	res, err := json.Marshal(jsonStruct)
	if err != nil {
		log.Fatal(err)
	}
	req.w.Write(res)
}
