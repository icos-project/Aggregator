package server

import (
	"errors"
	"fmt"
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

	// Data received from models.go
	querierData := []byte(`{
				"Cluster_type": "1",
				"Cluster_name": "B",
				"Location_zone": "Madrid",
				"ServiceLevelAgreement": "C",
				"API": "D",
				"Node": [{"Node_name": "alpa", "Node_type": "1"}, 
						{"Node_name": "beta", "Node_type": "2"}]
				}`)

	w.Write(querierData)
}
