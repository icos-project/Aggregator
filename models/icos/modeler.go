package models

import (
	"encoding/json"
	"fmt"
	"icos/server/querier"
)

func TransformQuery() []byte {

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

	// Convert to JSON
	json, err := json.Marshal(clusters)
	if err != nil {
		fmt.Printf("Error marshaling models: %v\n", err)
	}

	return json
}
