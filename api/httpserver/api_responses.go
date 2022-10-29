package httpserver

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
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

func NewErrorResponse(message string) ErrorResponse {
	return ErrorResponse{Message: strings.ToLower(message)}
}
