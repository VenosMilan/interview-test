package tools

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

// RecordNotFound general error message for non-existent record
var RecordNotFound = errors.New("record not found")

// ErrorResponse structure for error response
type ErrorResponse struct {
	ErrText string `json:"errText"`
}

// SetErrResponse function sets http status code and error text into response
// for non-existent record return 404, in other cases return 500
func SetErrResponse(response http.ResponseWriter, err error) {
	if err != nil && errors.Is(err, RecordNotFound) {
		SetErrResponseWithStatusCode(response, err, http.StatusNotFound)
		return
	}

	SetErrResponseWithStatusCode(response, err, http.StatusInternalServerError)
}

// SetErrResponseWithStatusCode function sets http status code and error text into response
// function sets error text into response as json
func SetErrResponseWithStatusCode(response http.ResponseWriter, err interface{}, statusCode int) {
	response.WriteHeader(statusCode)

	if err != nil {
		log.Printf("Err: %s", err.(error).Error())

		_ = json.NewEncoder(response).Encode(ErrorResponse{ErrText: err.(error).Error()})
	}
}
