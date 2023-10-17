package healthcheck

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

// MakeGetReadyEndpoint function for create GET endpoint,
// Readiness probe, signaling the service is ready to accept requests
func MakeGetReadyEndpoint() http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "text/html")

		_, err := response.Write([]byte("OK"))

		if err != nil {
			log.Printf("Error service is not ready %s", err.Error())
		}

		log.Debug("Service is OK")
	}
}
