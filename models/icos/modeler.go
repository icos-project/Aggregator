package models_icos

import (
	"aggregator/querier"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"time"
)

func GetInfra() []byte {

	// Get metrics from Thanos
	clusters := queryPrometheus()

	// Convert to JSON
	json, err := json.MarshalIndent(clusters, "", "\t")
	if err != nil {
		fmt.Printf("Error marshaling models: %v\n", err)
	}

	return json
}

func queryPrometheus() Infrastructure {

	// Timestamps
	q := querier.PromQLQuery{
		Metric: "timestamp(up)",
		Params: map[string]string{}}

	oldest := math.Inf(1)
	for _, timestamp := range querier.Query(q.String()) {
		if float64(timestamp.Value) < oldest {
			oldest = float64(timestamp.Value)
		}
	}

	var clusters = map[string]Cluster{}

	// Clusters and Nodes
	q = querier.PromQLQuery{
		Metric: "kube_node_info",
		Params: map[string]string{}}

	for _, node := range querier.Query(q.String()) {

		cluster_id := string(node.Metric["icos_agent_cluster_id"])
		node_name := string(node.Metric["node"])

		if cluster_id != "self" {
			if _, exists := clusters[cluster_id]; !exists {
				var newCluster = Cluster{
					Name: cluster_id,
					Node: map[string]Node{},
					Pod:  map[string]Pod{},
				}
				clusters[cluster_id] = newCluster
			}

			newNode := Node{
				Name:    node_name,
				Devices: map[string]Device{},
			}
			clusters[cluster_id].Node[node_name] = newNode
		}
	}

	// Cluster - Node - StaticMetrics
	q = querier.PromQLQuery{
		Metric: "kube_node_status_capacity{resource='cpu'}",
		Params: map[string]string{}}

	for _, node := range querier.Query(q.String()) {
		cluster_id := string(node.Metric["icos_agent_cluster_id"])
		node_name := string(node.Metric["node"])
		cores := int32(node.Value)

		if cluster_id != "self" {
			n := clusters[cluster_id].Node[node_name]
			n.StaticMetrics.CPUCores = cores
			clusters[cluster_id].Node[node_name] = n
		}
	}

	q = querier.PromQLQuery{
		Metric: "kube_node_status_capacity{resource='memory'}",
		Params: map[string]string{}}

	for _, node := range querier.Query(q.String()) {
		cluster_id := string(node.Metric["icos_agent_cluster_id"])
		node_name := string(node.Metric["node"])
		ram := int64(node.Value)

		if cluster_id != "self" {
			n := clusters[cluster_id].Node[node_name]
			n.StaticMetrics.RAMMemory = ram
			clusters[cluster_id].Node[node_name] = n
		}
	}

	// Cluster - Node - Devices
	q = querier.PromQLQuery{
		Metric: "node_mounted",
		Params: map[string]string{}}

	for _, node := range querier.Query(q.String()) {
		cluster_id := string(node.Metric["icos_agent_cluster_id"])
		node_name := string(node.Metric["node_name"])
		device_name := string(node.Metric["device"])
		device_type := strings.Split(device_name, "_")[0]

		newDev := Device{
			Name: device_name,
			Type: device_type,
		}

		clusters[cluster_id].Node[node_name].Devices[device_name] = newDev

	}

	// Cluster - Pods
	q = querier.PromQLQuery{
		Metric: "kube_pod_info",
		Params: map[string]string{}}

	for _, pod := range querier.Query(q.String()) {
		cluster_id := string(pod.Metric["icos_agent_cluster_id"])
		pod_name := string(pod.Metric["pod"])
		pod_ip := string(pod.Metric["pod_ip"])

		if cluster_id != "self" {
			newPod := Pod{
				Name:      pod_name,
				Container: map[string]Container{},
			}
			if pod_ip != "" {
				newPod.IP = pod_ip
			}
			clusters[cluster_id].Pod[pod_name] = newPod
		}
	}

	// Cluster - Pod - Status
	q = querier.PromQLQuery{
		Metric: "kube_pod_status_phase == 1",
		Params: map[string]string{}}

	for _, pod := range querier.Query(q.String()) {

		cluster_id := string(pod.Metric["icos_agent_cluster_id"])
		pod_name := string(pod.Metric["pod"])
		status := string(pod.Metric["phase"])

		if cluster_id != "self" {
			pod := clusters[cluster_id].Pod[pod_name]
			pod.Status = status
			clusters[cluster_id].Pod[pod_name] = pod
		}
	}

	// Cluster - Pod - Containers
	q = querier.PromQLQuery{
		Metric: "kube_pod_container_info",
		Params: map[string]string{}}

	for _, container := range querier.Query(q.String()) {

		cluster_id := string(container.Metric["icos_agent_cluster_id"])
		pod_name := string(container.Metric["pod"])
		cont_name := string(container.Metric["container"])
		node := string(container.Metric["k8s_node_name"])

		if cluster_id != "self" {
			newContainer := Container{
				Name: cont_name,
				Node: node,
			}
			clusters[cluster_id].Pod[pod_name].Container[cont_name] = newContainer
		}
	}

	// Cluster - Pod - Number of containers
	for cluster_id := range clusters {
		for pod_id := range clusters[cluster_id].Pod {
			pod := clusters[cluster_id].Pod[pod_id]
			pod.NumberOfContainers = int32(len(pod.Container))
			clusters[cluster_id].Pod[pod_id] = pod
		}
	}

	// Cluster - Pod - Container - CPU Usage
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

	// Timestamps
	loc := time.FixedZone("Local", 0)
	tp := time.Date(1970, 1, 1, 0, 0, 0, 0, loc)
	ts := time.Since(tp).Seconds()

	var time = Timestamp{
		OldestTimestamp: oldest,
		TimeSinceOldest: ts - oldest,
	}

	var infra = Infrastructure{
		Timestamp: time,
		Cluster:   clusters,
	}

	return infra
}
