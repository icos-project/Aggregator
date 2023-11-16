package server_icos

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	mid "icos/server/middlewares"
	m_icos "icos/server/models/icos"
	responses "icos/server/responses"
)

func CreateServer(project string) {

	// Routes
	http.HandleFunc("/", mid.SetMiddlewareLog(mid.SetMiddlewareJSON(mid.JWTValidation(connectToQuerier(project)))))
	http.HandleFunc("/healthz", healthCheck)

	// Open server
	port := getenv("AGGREGATOR_PORT", "8080")
	fmt.Printf("server listening on :%s", port)
	err := http.ListenAndServe((":" + port), nil)

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

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
