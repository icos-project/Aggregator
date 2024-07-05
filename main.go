/*
Copyright 2023 Bull SAS

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
	http "aggregator/servers/http"
	protobuf "aggregator/servers/protobuf"
	"fmt"
	"os"
	"sync"
)

func main() {

	fmt.Println("Starting Aggregator [v1.0] ...")

	// Get ports
	http_port := os.Getenv("HTTP_PORT")
	grpc_port := os.Getenv("GRPC_PORT")

	fmt.Println("Using PROMETHEUS_ADDRESS: " + os.Getenv("PROMETHEUS_ADDRESS"))

	var wg sync.WaitGroup

	// Default: HTTP server in port 8080
	if http_port == "" && grpc_port == "" {
		http_port = "8080"
	}

	// Launch HTTP server
	if http_port != "" {
		fmt.Println("Starting HTTP server...")
		wg.Add(1)
		go http.CreateServer(&wg, "icos", http_port)
	}

	// Launch gRPC server
	if grpc_port != "" {
		fmt.Println("Starting gRPC server...")
		wg.Add(1)
		go protobuf.CreateServer(&wg, "etim", grpc_port)
	}

	wg.Wait()

}
