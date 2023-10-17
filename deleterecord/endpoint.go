package deleterecord

import (
	"interviewtest/tools"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// MakeDeleteRecordEndpoint function create DELETE endpoint for delete record
func MakeDeleteRecordEndpoint(service Service) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		recID := mux.Vars(request)["id"]

		id, err := strconv.ParseInt(recID, 10, 64)

		if err != nil {
			tools.SetErrResponseWithStatusCode(response, err, http.StatusBadRequest)
			return
		}

		if err := service.Delete(id); err != nil {
			tools.SetErrResponse(response, err)
			return
		}

		response.WriteHeader(http.StatusNoContent)

		log.Debugf("Delete record %s was successful", recID)
	}
}
