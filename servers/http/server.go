package server_http

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	m_icos "aggregator/models/icos"
	responses "aggregator/servers/http/responses"
	mid "aggregator/servers/middlewares"
)

func CreateServer(wg *sync.WaitGroup, project string, port string) {

	defer wg.Done()

	// Routes
	key := os.Getenv("KEY")
	if key != "" {
		mid.SetPublicKey(key)
		http.HandleFunc("/", mid.SetMiddlewareLog(mid.SetMiddlewareJSON(mid.JWTValidation(connectToQuerier(project)))))
	} else {
		http.HandleFunc("/", connectToQuerier(project))
	}

	http.HandleFunc("/healthz", healthCheck)

	_, err := strconv.Atoi(port)
	if err != nil {
		log.Fatalf("Invalid port: %v", err)
	}

	fmt.Printf("server listening on :%s\n", port)
	err = http.ListenAndServe((":" + port), nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}

}

func connectToQuerier(project string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var querierData []byte

		// Get data in JSON format
		switch project {
		case "icos":
			querierData = m_icos.GetInfra()
		default:
			fmt.Printf("Project '%s' not found\n", project)
			os.Exit(1)
		}

		// Server response
		w.Write(querierData)

	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Aggregator working properly!")
}
