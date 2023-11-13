package models_cognifog

import (
	"icos/server/querier"
	pb "icos/server/servers/protobuf/cognifog"
)

func GetInfra() *pb.InfrastructureModel {

	// Get metrics from Thanos
	//clusters := queryThanos()
	//println(clusters)

	// Convert metrics to server.proto structs
	newClusters := ConvertMetrics( /*clusters*/ )

	return newClusters
}

func queryThanos() map[string]Cluster {

	// Create clusters and nodes
	q := querier.PromQLQuery{
		Metric: "kube_pod_info{created_by_name='prom-node-exporter'}",
		Params: map[string]string{}}

	var clusters = map[string]Cluster{}

	for _, pod := range querier.Query(q.String()) {

		ins := string(pod.Metric["instance"])
		node := string(pod.Metric["node"])

		// Create new cluster if needed
		// Same instances means nodes are in the same cluster
		// New instances are new clusters
		if _, exists := clusters[ins]; !exists {
			var newCluster = Cluster{
				Name: ins,
				Node: map[string]Node{},
			}
			clusters[ins] = newCluster
		}

		// Create new node
		newNode := Node{
			Name:          node,
			StaticMetrics: StaticMetrics{},
		}
		clusters[ins].Node[node] = newNode
	}

	// Add node stats
	q = querier.PromQLQuery{
		Metric: "machine_cpu_cores{service='prom-kube-prometheus-kubelet'}",
		Params: map[string]string{}}

	for _, node := range querier.Query(q.String()) {

		nodeName := string(node.Metric["node"])

		for cluster := range clusters {
			if n, exists := clusters[cluster].Node[nodeName]; exists {
				n.StaticMetrics.CPUCores = float64(node.Value)
				clusters[cluster].Node[nodeName] = n
			}
		}

	}

	return clusters
}

func ConvertMetrics( /*clusters map[string]Cluster*/ ) *pb.InfrastructureModel {

	var model = &pb.InfrastructureModel{
		Nodes: []*pb.InfrastructureModel_Node{
			{Id: "id1", NodeType: 1},
			{Id: "id2", NodeType: 0},
		},
	}

	return model
}
