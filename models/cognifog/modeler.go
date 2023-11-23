package models_cognifog

import (
	"aggregator/querier"
	pb "aggregator/servers/protobuf/cognifog"
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

	// Create clusters and nodes
	q := querier.PromQLQuery{
		Metric: "kube_node_info",
		Params: map[string]string{}}

	for _, node := range querier.Query(q.String()) {

		cluster_id := string(node.Metric["icos_agent_cluster_id"])
		node_name := string(node.Metric["k8s_node_name"])

		if _, exists := clusters[cluster_id]; !exists {
			var newCluster = Cluster{
				Nodes: map[string]Node{},
			}
			clusters[cluster_id] = newCluster
		}

		// Create new node
		newNode := Node{
			Id:                  node_name,
			Resources:           ComputeResources{},
			Available_resources: ComputeResources{},
		}
		clusters[cluster_id].Nodes[node_name] = newNode
	}

	// Get total memory
	q = querier.PromQLQuery{
		Metric: "node_memory_MemTotal_bytes",
		Params: map[string]string{}}

	for _, result := range querier.Query(q.String()) {

		node_name := string(result.Metric["k8s_node_name"])
		cluster_name := string(result.Metric["icos_agent_cluster_id"])

		node := clusters[cluster_name].Nodes[node_name]
		node.Resources.MemoryInBytes = int64(result.Value)
		clusters[cluster_name].Nodes[node_name] = node
	}

	// Get available memory
	q = querier.PromQLQuery{
		Metric: "node_memory_MemAvailable_bytes",
		Params: map[string]string{}}

	for _, result := range querier.Query(q.String()) {

		node_name := string(result.Metric["k8s_node_name"])
		cluster_name := string(result.Metric["icos_agent_cluster_id"])

		node := clusters[cluster_name].Nodes[node_name]
		node.Available_resources.MemoryInBytes = int64(result.Value)
		clusters[cluster_name].Nodes[node_name] = node
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
				Id:       node.Id,
				NodeType: pb.InfrastructureModel_Node_NodeType(node.Node_type),
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
