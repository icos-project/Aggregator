package main

import (
	http "aggregator/servers/http"
	protobuf "aggregator/servers/protobuf"
	"fmt"
	"os"
	"sync"
)

func main() {

	// Get ports
	http_port := os.Getenv("HTTP_PORT")
	grpc_port := os.Getenv("GRPC_PORT")

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
		go protobuf.CreateServer(&wg, "cognifog", grpc_port)
	}

	wg.Wait()

}
