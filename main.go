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
package main

import (
	log "aggregator/common/logs"
	http "aggregator/servers/http"
	protobuf "aggregator/servers/protobuf"
	"os"
	"sync"
)

//	@title			Swagger Aggregator API
//	@version		1.0
//	@description	Aggregator Microservice.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/

//	@securityDefinitions.basic	OAuth 2.0

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/

// path used in logs
const pathLOG string = "AGGREGATOR > "

func main() {

	log.Info(pathLOG + "Starting Aggregator [v1.2.6] [2024.09.06] ...")

	// Get ports
	http_port := os.Getenv("HTTP_PORT")
	grpc_port := os.Getenv("GRPC_PORT")

	log.Info(pathLOG + "Using PROMETHEUS_ADDRESS: " + os.Getenv("PROMETHEUS_ADDRESS"))

	var wg sync.WaitGroup

	// Default: HTTP server in port 8080
	if http_port == "" && grpc_port == "" {
		http_port = "8080"
	}

	// Launch HTTP server
	if http_port != "" {
		log.Info("Starting HTTP server...")
		wg.Add(1)
		go http.CreateServer(&wg, "icos", http_port)
	}

	// Launch gRPC server
	if grpc_port != "" {
		log.Info("Starting gRPC server...")
		wg.Add(1)
		go protobuf.CreateServer(&wg, "etim", grpc_port)
	}

	wg.Wait()

}
