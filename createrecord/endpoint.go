package createrecord

import (
	"encoding/json"
	"interviewtest/record"
	"interviewtest/tools"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// MakePostCreateRecordEndpoint function create POST endpoint for create record
func MakePostCreateRecordEndpoint(service Service) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		var rec record.Record
		decoder := json.NewDecoder(request.Body)
		decoder.DisallowUnknownFields()

		if err := decoder.Decode(&rec); err != nil {
			tools.SetErrResponseWithStatusCode(response, err, http.StatusBadRequest)
			return
		}

		creatRes, err := service.Create(&rec)

		if err != nil {
			tools.SetErrResponseWithStatusCode(response, err, http.StatusInternalServerError)
			return
		}

		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusCreated)

		if err = json.NewEncoder(response).Encode(creatRes); err != nil {
			tools.SetErrResponseWithStatusCode(response, err, http.StatusInternalServerError)
			return
		}

		log.Debugf("Create record %d was successful", creatRes.RecordID)
	}
}
