package main

import (
	"context"
	"fmt"
	"interviewtest/appconfiguration"
	"interviewtest/createrecord"
	"interviewtest/deleterecord"
	"interviewtest/editrecord"
	"interviewtest/getrecord"
	"interviewtest/healthcheck"
	"interviewtest/storage"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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

	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", appConf.ServerPort),
		Handler: corsOptions(myRouter),
	}

	idleConnectionsClosed := shutDownServer(&srv)

	log.Infof("Server listening on %s", srv.Addr)

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe Error: %v", err)
	}

	<-idleConnectionsClosed

	log.Println("Server shutdown gracefully")
}

func shutDownServer(srv *http.Server) chan struct{} {
	idleConnectionsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP Server Shutdown Error: %v", err)
		}
		close(idleConnectionsClosed)
	}()

	return idleConnectionsClosed
}

func corsOptions(myRouter *mux.Router) http.Handler {
	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type"})
	allowedMethods := handlers.AllowedMethods([]string{http.MethodGet, http.MethodPost, http.MethodOptions, http.MethodPut, http.MethodPost})

	return handlers.CORS(allowedHeaders, allowedMethods)(myRouter)
}
