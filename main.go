package main

import (
	aggregator "icos/server/aggregator"

	"errors"
	"fmt"
	"net/http"
	"os"
)

func main() {

	http.HandleFunc("/", aggregator.ServeQuery)

	fmt.Println("server listening on :8080")

	err := http.ListenAndServe(":8080", nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
