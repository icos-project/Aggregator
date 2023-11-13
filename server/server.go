package server

import (
	"errors"
	"fmt"
	models "icos/server/models/icos"
	"net/http"
	"os"
)

func CreateServer() {

	// Open server
	http.HandleFunc("/", serveQuery)

	fmt.Println("server listening on :8080")

	err := http.ListenAndServe(":8080", nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}

}

func serveQuery(w http.ResponseWriter, r *http.Request) {

	// Get data in JSON format
	querierData := models.GetInfra()

	// Server response
	w.Write(querierData)
}
