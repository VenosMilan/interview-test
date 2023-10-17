package editrecord

import (
	"encoding/json"
	"interviewtest/record"
	"interviewtest/tools"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// MakePutRecordEndpoint function create PUT endpoint for record modification
func MakePutRecordEndpoint(service Service) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		var rec record.Record
		recID := mux.Vars(request)["id"]

		id, err := strconv.ParseInt(recID, 10, 64)

		if err != nil {
			tools.SetErrResponseWithStatusCode(response, err, http.StatusBadRequest)
			return
		}

		decoder := json.NewDecoder(request.Body)
		decoder.DisallowUnknownFields()

		if err := decoder.Decode(&rec); err != nil {
			tools.SetErrResponseWithStatusCode(response, err, http.StatusBadRequest)
			return
		}

		if err := service.Edit(id, &rec); err != nil {
			tools.SetErrResponse(response, err)
			return
		}

		log.Debugf("Edit record %s was successful", recID)
	}
}
