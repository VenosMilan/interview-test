package main

import (
	"fmt"
	"interviewtest/appconfiguration"
	"interviewtest/createrecord"
	"interviewtest/deleterecord"
	"interviewtest/editrecord"
	"interviewtest/getrecord"
	"interviewtest/healthcheck"
	"interviewtest/storage"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func main() {
	appConf := appconfiguration.NewAppConfiguration()

	storageService, err := storage.NewService(appConf.BinaryFilePath)

	if err != nil {
		log.Fatal(err)
	}

	defer storageService.Close()

	getRecordService := getrecord.NewService(storageService)
	createRecordService := createrecord.NewService(storageService)
	deleteRecordService := deleterecord.NewService(storageService)
	putRecordService := editrecord.NewService(storageService)

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.Handle("/readyz", healthcheck.MakeGetReadyEndpoint()).Methods(http.MethodGet)
	myRouter.Handle("/records", createrecord.MakePostCreateRecordEndpoint(createRecordService)).Methods(http.MethodPost)
	myRouter.Handle("/records/{id:[0-9]+}", deleterecord.MakeDeleteRecordEndpoint(deleteRecordService)).Methods(http.MethodDelete)
	myRouter.Handle("/records/{id:[0-9]+}", editrecord.MakePutRecordEndpoint(putRecordService)).Methods(http.MethodPut)
	myRouter.Handle("/records/{id:[0-9]+}", getrecord.MakeGetRecordEndpoint(getRecordService)).Methods(http.MethodGet)

	log.Info("Service running")

	if err := http.ListenAndServe(fmt.Sprintf(":%s", appConf.ServerPort), corsOptions(myRouter)); err != nil {
		log.Fatal(errors.WithStack(err))
	}
}

func corsOptions(myRouter *mux.Router) http.Handler {
	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type"})
	allowedMethods := handlers.AllowedMethods([]string{http.MethodGet, http.MethodPost, http.MethodOptions, http.MethodPut, http.MethodPost})

	return handlers.CORS(allowedHeaders, allowedMethods)(myRouter)
}
