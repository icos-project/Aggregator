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
package models_etim

import (
	"aggregator/querier"
	pb "aggregator/servers/protobuf/etim"
	"fmt"
)

func GetInfra() *pb.InfrastructureModel {

	// Get metrics from Thanos
	infra := queryPrometheus()

	// Convert metrics to server.proto structs
	newInfra := ConvertMetrics(infra)

	return newInfra
}

func queryPrometheus() map[string]Cluster {

	var clusters = map[string]Cluster{}

	// Cluster - Nodes
	q := querier.PromQLQuery{
		Metric: "node_uname_info",
		Params: map[string]string{}}

	for _, node := range querier.Query(q.String()) {
		cluster_id := string(node.Metric["cluster"])
		node_id := string(node.Metric["node"])
		architecture := string(node.Metric["machine"])

		if _, exists := clusters[cluster_id]; !exists {
			var newCluster = Cluster{
				Nodes: map[string]Node{},
			}
			clusters[cluster_id] = newCluster
		}

		newNode := Node{
			Id:                  node_id,
			CPUArchitecture:     architecture,
			Resources:           ComputeResources{},
			Available_resources: ComputeResources{},
		}
		clusters[cluster_id].Nodes[node_id] = newNode
	}

	// Cluster - Node - Resources
	q = querier.PromQLQuery{
		Metric: "node_memory_MemTotal_bytes",
		Params: map[string]string{}}

	for _, result := range querier.Query(q.String()) {

		node_id := string(result.Metric["node"])
		cluster_id := string(result.Metric["cluster"])

		if checkClusterNode(cluster_id, node_id, clusters, q.Metric) {
			node := clusters[cluster_id].Nodes[node_id]
			node.Resources.MemoryInBytes = int64(result.Value)
			clusters[cluster_id].Nodes[node_id] = node
		}
	}

	// Cluster - Node - Available_resources
	q = querier.PromQLQuery{
		Metric: "node_memory_MemAvailable_bytes",
		Params: map[string]string{}}

	for _, result := range querier.Query(q.String()) {

		node_id := string(result.Metric["node"])
		cluster_id := string(result.Metric["cluster"])

		if checkClusterNode(cluster_id, node_id, clusters, q.Metric) {
			node := clusters[cluster_id].Nodes[node_id]
			node.Available_resources.MemoryInBytes = int64(result.Value)
			clusters[cluster_id].Nodes[node_id] = node
		}
	}

	return clusters
}

func ConvertMetrics(clusters map[string]Cluster) *pb.InfrastructureModel {

	// Create proto strcuture
	var infra = &pb.InfrastructureModel{
		Nodes: []*pb.InfrastructureModel_Node{},
	}

	// Fill strcuture with collected data
	for _, cluster := range clusters {
		for _, node := range cluster.Nodes {
			newNode := pb.InfrastructureModel_Node{
				Id:              node.Id,
				NodeType:        pb.InfrastructureModel_Node_NodeType(node.Node_type),
				CpuArchitecture: node.CPUArchitecture,
				Resources: &pb.ComputeResources{
					MilliCPU:      node.Resources.MilliCPU,
					MemoryInBytes: node.Resources.MemoryInBytes,
				},
				AvailableResources: &pb.ComputeResources{
					MilliCPU:      node.Available_resources.MilliCPU,
					MemoryInBytes: node.Available_resources.MemoryInBytes,
				},
				Labels: node.Labels,
			}
			infra.Nodes = append(infra.Nodes, &newNode)
		}
	}

	return infra
}

func checkClusterNode(cluster string, node string, clusters map[string]Cluster, metric string) bool {
	_, existsCluster := clusters[cluster]
	_, existsNode := clusters[cluster].Nodes[node]

	if !existsCluster {
		fmt.Println("Unknown cluster ", cluster, ". Error in metric: ", metric)
	} else if !existsNode {
		fmt.Println("Unknown node ", node, " in cluster ", cluster, ". Error in metric: ", metric)
	}

	return existsCluster && existsNode
}
