/*
Copyright © 2022-2024 EVIDEN

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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

// connectToQuerier example
//
//	@Summary 		get clusters state
//	@Description	get clusters state
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string					true	"Authentication header"
//	@Success		200				{object}	m_icos.Infrastructure	"Ok"
//	@Failure		400				{object}	string					"Bad request"
//	@Router			/ [get]
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
