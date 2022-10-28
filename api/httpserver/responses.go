package httpserver

import (
	"encoding/json"
	"log"
	"net/http"
)

const contentType = "Content-Type"
const applicationJson = "application/json"

func apiResponse(w http.ResponseWriter, statusCode int, responseBody interface{}) {
	response, err := json.Marshal(responseBody)
	if err != nil {
		log.Fatalf("error marshalling response. error=(%v)", err)
	}

	w.Header().Set(contentType, applicationJson)
	w.WriteHeader(statusCode)
	_, err = w.Write(response)
	if err != nil {
		log.Fatalf("error writing response. error=(%v)", err)
	}
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponses struct {
	Errors []ErrorResponse `json:"errors"`
}

func NewErrorResponses(errors ...ErrorResponse) ErrorResponses {
	var errResponses []ErrorResponse
	errResponses = append(errResponses, errors...)
	return ErrorResponses{errResponses}
}
