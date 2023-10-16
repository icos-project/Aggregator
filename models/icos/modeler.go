package models

import (
	"encoding/json"
	"fmt"
	"icos/server/querier"
	"os"
)

func TransformQuery() []byte {

	os.Setenv("PROMETHEUS_ADDRESS", "http://query.192.168.137.175.nip.io/") // thanos-query

	q := querier.PromQLQuery{
		Metric: "kube_pod_info{created_by_name='prom-node-exporter'}",
		Params: map[string]string{}}

	// Create clusters and nodes
	var clusters = map[string]Cluster{}

	for _, pod := range querier.Query(q.String()) {

		ins := string(pod.Metric["instance"])
		node := string(pod.Metric["node"])

		// Create new cluster if needed
		// Same instances means nodes are in the same cluster
		// New instances are new clusters
		if _, exists := clusters[ins]; !exists {
			var newCluster = Cluster{
				Name:       ins,
				Node:       map[string]Node{},
				Deployment: map[string]Deployment{},
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

	// Add node stats			TODO: COMPLETE
	q = querier.PromQLQuery{
		Metric: "machine_cpu_cores{service='prom-kube-prometheus-kubelet'}",
		Params: map[string]string{}}

	for _, node := range querier.Query(q.String()) {

		nodeName := string(node.Metric["node"])

		for cluster, _ := range clusters {
			if n, exists := clusters[cluster].Node[nodeName]; exists {
				n.StaticMetrics.CPUCores = float64(node.Value)
				clusters[cluster].Node[nodeName] = n
			}
		}

	}

	// Add deployment		TODO: Is it correct?
	q = querier.PromQLQuery{
		Metric: "kube_pod_container_info",
		Params: map[string]string{}}

	for _, container := range querier.Query(q.String()) {

		cluster := string(container.Metric["instance"])
		dep := string(container.Metric["container"])
		con := string(container.Metric["pod"])

		if _, exists := clusters[cluster].Deployment[dep]; !exists {
			newDeployment := Deployment{
				Container: map[string]Container{},
			}
			clusters[cluster].Deployment[dep] = newDeployment
		}

		newContainer := Container{
			Name: con,
		}
		clusters[cluster].Deployment[dep].Container[con] = newContainer
	}

	// Convert to JSON
	json, err := json.MarshalIndent(clusters, "", "\t")
	if err != nil {
		fmt.Printf("Error marshaling models: %v\n", err)
	}

	return json
}
