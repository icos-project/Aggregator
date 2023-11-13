package models_icos

import (
	"encoding/json"
	"fmt"
	"icos/server/querier"
)

func GetInfra() []byte {

	// Get metrics from Thanos
	clusters := queryThanos()

	// Convert to JSON
	json, err := json.MarshalIndent(clusters, "", "\t")
	if err != nil {
		fmt.Printf("Error marshaling models: %v\n", err)
	}

	return json
}

func queryThanos() map[string]Cluster {

	// Create clusters and nodes
	q := querier.PromQLQuery{
		Metric: "kube_node_info",
		Params: map[string]string{}}

	var clusters = map[string]Cluster{}

	for _, node := range querier.Query(q.String()) {

		cluster_id := string(node.Metric["icos_agent_cluster_id"])
		node_id := string(node.Metric["system_uuid"])
		node_name := string(node.Metric["node"])

		// Create new cluster if needed
		// Same instances means nodes are in the same cluster
		// New instances are new clusters
		if _, exists := clusters[cluster_id]; !exists {
			var newCluster = Cluster{
				Name: cluster_id,
				Node: map[string]Node{},
				Pod:  map[string]Pod{},
			}
			clusters[cluster_id] = newCluster
		}

		// Create new node
		newNode := Node{
			Name: node_name,
			//StaticMetrics: StaticMetrics{},
		}
		clusters[cluster_id].Node[node_id] = newNode
	}

	// Add deployment
	q = querier.PromQLQuery{
		Metric: "kube_pod_container_info",
		Params: map[string]string{}}

	for _, container := range querier.Query(q.String()) {

		cluster := string(container.Metric["icos_agent_cluster_id"])
		pod := string(container.Metric["k8s_pod_uid"])
		cont_id := string(container.Metric["uid"])
		cont_name := string(container.Metric["container"])

		if _, exists := clusters[cluster].Pod[pod]; !exists {
			newPod := Pod{
				Container: map[string]Container{},
			}
			clusters[cluster].Pod[pod] = newPod
		}

		newContainer := Container{
			Name: cont_name,
		}
		clusters[cluster].Pod[pod].Container[cont_id] = newContainer
	}

	return clusters
}
