package getrecord

import (
	"encoding/json"
	"interviewtest/tools"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// MakeGetRecordEndpoint function create GET endpoint for get record
func MakeGetRecordEndpoint(service Service) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		recID := mux.Vars(request)["id"]

		id, err := strconv.ParseInt(recID, 10, 64)

		if err != nil {
			tools.SetErrResponseWithStatusCode(response, err, http.StatusNotFound)
			return
		}

		rec, err := service.GetRecord(id)

		if err != nil {
			tools.SetErrResponse(response, err)
			return
		}

		response.Header().Set("Content-Type", "application/json")

		if err = json.NewEncoder(response).Encode(rec); err != nil {
			tools.SetErrResponseWithStatusCode(response, err, http.StatusInternalServerError)
			return
		}

		log.Debugf("Get record by id %s was successful", recID)
	}
}
