package server_cognifog

import (
	context "context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	md "icos/server/models/cognifog"
	pb "icos/server/servers/protobuf/cognifog"

	"google.golang.org/grpc"
)

type server_cognifog struct {
	pb.UnimplementedAggregatorServer
}

func (s *server_cognifog) ConnectToQuerier(ctx context.Context, in *pb.Empty) (*pb.InfrastructureModel, error) {
	return md.GetInfra(), nil
}

func CreateServer(project string) {

	port, err := strconv.Atoi(getenv("AGGREGATOR_PORT", "8080"))

	if err != nil {
		log.Fatalf("Invalid port: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	switch project {
	case "cognifog":
		pb.RegisterAggregatorServer(s, &server_cognifog{})
	default:
		fmt.Printf("Project '%s' not found\n", project)
		os.Exit(1)
	}

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
