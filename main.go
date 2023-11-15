package main

import (
	http "icos/server/servers/http"
	//protobuf "icos/server/servers/protobuf"
)

func main() {

	// HTTP server
	http.CreateServer("icos")

	//gRPC server
	//protobuf.CreateServer("cognifog")
}
