package server

import (
	"errors"
	"fmt"
	models "icos/server/models/icos"

	mid "icos/server/middlewares"
	"net/http"
	"os"
)

func CreateServer() {

	// Open server
	http.HandleFunc("/", mid.SetMiddlewareLog(mid.SetMiddlewareJSON(mid.JWTValidation(serveQuery))))

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
	querierData := models.TransformQuery()

	// Server response
	w.Write(querierData)
}
