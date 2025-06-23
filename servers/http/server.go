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
	"net/http"
	"os"
	"strconv"
	"sync"

	logs "aggregator/common/logs"
	m_icos "aggregator/models/icos"
	_ "aggregator/models/icos/models"
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
		logs.GetLogger().Fatal("Invalid port: ", err)
	}

	logs.GetLogger().Info("server listening on: " + port)
	err = http.ListenAndServe((":" + port), nil)

	if errors.Is(err, http.ErrServerClosed) {
		logs.GetLogger().Warn("server closed\n")
	} else if err != nil {
		logs.GetLogger().Error("error starting server: ", err)
		os.Exit(1)
	}

}

// connectToQuerier example
//
//	@Summary 		get clusters state
//	@Description	get clusters state
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string				true	"Authentication header"
//	@Success		200			{object}	models.Infrastructure	"Ok"
//	@Failure		400			{object}	string				"Bad request"
//	@Router			/ [get]
func connectToQuerier(project string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var querierData []byte

		// Get data in JSON format
		switch project {
		case "icos":
			querierData = m_icos.GetInfra()
		default:
			logs.GetLogger().Info("Project" + project + "not found")
			os.Exit(1)
		}

		// Server response
		w.Write(querierData)

	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Aggregator working properly!")
}
