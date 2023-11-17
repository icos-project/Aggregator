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

	var clusters = map[string]Cluster{}

	// Clusters and Nodes
	q := querier.PromQLQuery{
		Metric: "kube_node_info",
		Params: map[string]string{}}

	for _, node := range querier.Query(q.String()) {

		cluster_id := string(node.Metric["icos_agent_cluster_id"])
		node_name := string(node.Metric["node"])

		if _, exists := clusters[cluster_id]; !exists {
			var newCluster = Cluster{
				Name: cluster_id,
				Node: map[string]Node{},
				Pod:  map[string]Pod{},
			}
			clusters[cluster_id] = newCluster
		}

		newNode := Node{
			Name: node_name,
		}
		clusters[cluster_id].Node[node_name] = newNode
	}

	// Pods
	q = querier.PromQLQuery{
		Metric: "kube_pod_info",
		Params: map[string]string{}}

	for _, pod := range querier.Query(q.String()) {

		cluster_id := string(pod.Metric["icos_agent_cluster_id"])
		pod_name := string(pod.Metric["pod"])

		newPod := Pod{
			Name:      pod_name,
			Container: map[string]Container{},
		}
		clusters[cluster_id].Pod[pod_name] = newPod
	}

	// Pod - Status
	q = querier.PromQLQuery{
		Metric: "kube_pod_status_phase == 1",
		Params: map[string]string{}}

	for _, pod := range querier.Query(q.String()) {

		cluster_id := string(pod.Metric["icos_agent_cluster_id"])
		pod_name := string(pod.Metric["pod"])
		status := string(pod.Metric["phase"])

		pod := clusters[cluster_id].Pod[pod_name]
		pod.Status = status
		clusters[cluster_id].Pod[pod_name] = pod
	}

	// Containers
	q = querier.PromQLQuery{
		Metric: "kube_pod_container_info",
		Params: map[string]string{}}

	for _, container := range querier.Query(q.String()) {

		cluster_id := string(container.Metric["icos_agent_cluster_id"])
		pod_name := string(container.Metric["pod"])
		cont_name := string(container.Metric["container"])
		node := string(container.Metric["k8s_node_name"])

		newContainer := Container{
			Name: cont_name,
			Node: node,
		}
		clusters[cluster_id].Pod[pod_name].Container[cont_name] = newContainer
	}

	// Pod - Number of containers
	for cluster_id := range clusters {
		for pod_id := range clusters[cluster_id].Pod {
			pod := clusters[cluster_id].Pod[pod_id]
			pod.NumberOfContainers = int32(len(pod.Container))
			clusters[cluster_id].Pod[pod_id] = pod
		}
	}

	// Container - CPU Usage
	q = querier.PromQLQuery{
		Metric: "container_cpu_utilization_ratio",
		Params: map[string]string{}}

	for _, container := range querier.Query(q.String()) {

		cluster_id := string(container.Metric["icos_agent_cluster_id"])
		pod_name := string(container.Metric["k8s_pod_name"])
		cont_name := string(container.Metric["k8s_container_name"])
		value := container.Value

		if _, exists := clusters[cluster_id].Pod[pod_name].Container[cont_name]; exists {
			cont := clusters[cluster_id].Pod[pod_name].Container[cont_name]
			cont.CPUUsage = float64(value)
			clusters[cluster_id].Pod[pod_name].Container[cont_name] = cont
		}
	}

	return clusters
}
