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
package server_protobuff

import (
	context "context"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"

	log "aggregator/common/logs"
	md "aggregator/models/cognifog"
	pb "aggregator/servers/protobuf/cognifog"

	"google.golang.org/grpc"
)

type server_cognifog struct {
	pb.UnimplementedAggregatorServer
}

func (s *server_cognifog) ConnectToQuerier(ctx context.Context, in *pb.Empty) (*pb.InfrastructureModel, error) {
	return md.GetInfra(), nil
}

func CreateServer(wg *sync.WaitGroup, project string, port string) {

	defer wg.Done()

	p, err := strconv.Atoi(port)

	if err != nil {
		log.Fatal("Invalid port: ", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", p))

	if err != nil {
		log.Fatal("failed to listen: ", err)
	}

	s := grpc.NewServer()

	switch project {
	case "cognifog":
		pb.RegisterAggregatorServer(s, &server_cognifog{})
	default:
		log.Info("Project" + project + "not found")
		os.Exit(1)
	}

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Error("failed to serve: ", err)
	}

}
