package main

import (
	"context"
	"log"
	"os"
	"time"

	pb "aggregator/servers/protobuf/cognifog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	addr := "localhost:" + getenv("GRPC_PORT", "8181") // GRPC_PORT is defined in other instance so here is currently empty

	// Set up a connection to the server.
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewAggregatorClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	r, err := c.ConnectToQuerier(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("could not get metrics: %v", err)
	}
	log.Println(r)

}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
