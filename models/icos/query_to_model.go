package models

import (
	"fmt"
	"icos/server/querier"
	"os"
)

func TransformQuery() {

	os.Setenv("PROMETHEUS_ADDRESS", "http://query.192.168.137.175.nip.io/") // thanos-query

	q := querier.PromQLQuery{
		Metric: "kube_pod_info{created_by_name='prom-node-exporter'}",
		Params: map[string]string{}}

	// Create clusters and nodes
	var clusters = map[string]ClusterTest{}

	for _, pod := range querier.Query(q.String()) {

		ins := string(pod.Metric["instance"])
		node := string(pod.Metric["node"])

		// Create new cluster if needed
		// Same instances means nodes are in the same cluster
		// New instances are new clusters
		if _, exists := clusters[ins]; !exists {
			var newCluster = ClusterTest{
				Name: ins,
				Node: map[string]NodeTest{},
			}
			clusters[ins] = newCluster
		}

		// Create new node
		newNode := NodeTest{
			Name:              node,
			StaticMetricsTest: StaticMetricsTest{},
		}
		clusters[ins].Node[node] = newNode
	}

	// Add Node Stats
	q = querier.PromQLQuery{
		Metric: "machine_cpu_cores{service='prom-kube-prometheus-kubelet'}",
		Params: map[string]string{}}

	for _, node := range querier.Query(q.String()) {

		nodeName := string(node.Metric["node"])

		for cluster, _ := range clusters {
			if n, exists := clusters[cluster].Node[nodeName]; exists {
				n.StaticMetricsTest.CPUCores = float64(node.Value)
				clusters[cluster].Node[nodeName] = n
			}
		}

	}

	fmt.Println(clusters)

}
