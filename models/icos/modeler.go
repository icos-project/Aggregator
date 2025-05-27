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
package modeler

import (
	logs "aggregator/common/logs"
	"aggregator/models/icos/common"
	"aggregator/models/icos/infra"
	"aggregator/models/icos/models"
	"aggregator/querier"
	"encoding/json"
	"math"
	"strings"
	"time"
)

func GetInfra() []byte {

	// Get metrics from Thanos
	logs.GetLogger().Info("----------------------------------------------------")
	logs.GetLogger().Info("> Aggregator called: Getting metrics from Thanos ...")
	clusters := queryPrometheus()

	// Convert to JSON
	json, err := json.MarshalIndent(clusters, "", "\t")
	if err != nil {
		logs.GetLogger().Error("Error marshaling models: ", err)
	}

	return json
}

func queryPrometheus() models.Infrastructure {

	var clusters = map[string]models.Cluster{}
	var orchs = map[string]models.OrchInfoNode{}

	// Timestamps
	logs.GetLogger().Info("\t>> QUERY: timestamp(up)")
	q := querier.PromQLQuery{
		Metric: "timestamp(up)",
		Params: map[string]string{}}

	oldest := math.Inf(1)
	for _, timestamp := range querier.Query(q.String()) {
		if float64(timestamp.Value) < oldest {
			oldest = float64(timestamp.Value)
		}
	}

	// Get Clusters and Nodes and add them to infra
	logs.GetLogger().Info("\t>> Getting Clusters and Nodes and adding them to infra ...")
	infra.SetNodes(clusters, orchs)

	// Cluster - Node - vulnerabilities and SCA_score
	//   node_uname_info > net_host_name="10.150.0.144", nodename="icosedge"
	//					 ==> Node.NetHostName = net_host_name (NEW), Node.Name = nodename
	//   SCA_score       > agent_hostname="icosedge", agent_ip="10.150.0.144"
	//   vulnerabilities > agent_hostname="icosedge", agent_ip="10.150.0.144"
	//					 ==> Node.NetHostName == agent_ip && Node.Name == agent_hostname

	// Cluster - Node - SCA_score
	logs.GetLogger().Info("\t>> QUERY: SCA_score")
	q = querier.PromQLQuery{
		Metric: "SCA_score",
		Params: map[string]string{}}

	for _, res := range querier.Query(q.String()) {
		agent_hostname := string(res.Metric["agent_hostname"])
		agent_ip := string(res.Metric["agent_ip"])
		scaScoreValue := int32(res.Value)

		for _, c := range clusters {
			for _, n := range c.Node {
				if n.NetHostName == agent_ip && n.Name == agent_hostname {
					n := clusters[c.Uuid].Node[n.Uuid]
					n.ScaScore = scaScoreValue
					clusters[c.Uuid].Node[n.Uuid] = n
				}
			}
		}
	}

	// Cluster - Node - vulnerabilities
	logs.GetLogger().Info("\t>> QUERY: vulnerabilities")
	q = querier.PromQLQuery{
		Metric: "vulnerabilities",
		Params: map[string]string{}}

	for _, res := range querier.Query(q.String()) {
		agent_hostname := string(res.Metric["agent_hostname"])
		agent_ip := string(res.Metric["agent_ip"])
		vulnerabilitySeverity := string(res.Metric["severity"])
		vulnerabilityValue := int32(res.Value)

		for _, c := range clusters {
			for _, n := range c.Node {
				if n.NetHostName == agent_ip && n.Name == agent_hostname {
					n := clusters[c.Uuid].Node[n.Uuid]

					if len(n.Vulnerabilities) == 0 {
						n.Vulnerabilities = make(map[string]int32, 10)
					}

					n.Vulnerabilities[vulnerabilitySeverity] = vulnerabilityValue
					clusters[c.Uuid].Node[n.Uuid] = n
				}
			}
		}
	}

	// Cluster - Node - StaticMetrics
	logs.GetLogger().Info("\t>> QUERY: kube_node_status_capacity{resource='cpu'} * on(icos_agen...")
	q = querier.PromQLQuery{
		Metric: "kube_node_status_capacity{resource='cpu'} * on(icos_agent_id, icos_host_name) group_left(icos_host_id) node_uname_info",
		Params: map[string]string{}}

	for _, node := range querier.Query(q.String()) {
		cluster_id := string(node.Metric["k8s_cluster_uid"])
		node_id := strings.TrimSpace(string(node.Metric["icos_host_id"]))
		cores := int32(node.Value)

		if common.CheckClusterNode(cluster_id, node_id, clusters, q.Metric) {
			n := clusters[cluster_id].Node[node_id]
			n.StaticMetrics.CPUCores = cores
			clusters[cluster_id].Node[node_id] = n
		}
	}

	logs.GetLogger().Info("\t>> QUERY: node_cpu_frequency_max_hertz")
	q = querier.PromQLQuery{
		Metric: "node_cpu_frequency_max_hertz",
		Params: map[string]string{}}

	freqs := make(map[[2]string][]int64)
	for _, node := range querier.Query(q.String()) {
		cluster_id := string(node.Metric["k8s_cluster_uid"])
		node_id := strings.TrimSpace(string(node.Metric["icos_host_id"]))
		frequency := int64(node.Value)

		if common.CheckClusterNode(cluster_id, node_id, clusters, q.Metric) {
			key := [2]string{cluster_id, node_id}

			if _, exists := freqs[key]; !exists {
				freqs[key] = []int64{frequency}
			} else {
				l := freqs[key]
				l = append(l, frequency)
				freqs[key] = l
			}
		}
	}
	for comb, f := range freqs {
		maxFrequency := common.MaxInt64(f)
		n := clusters[comb[0]].Node[comb[1]]
		n.StaticMetrics.CPUMaxFrequency = maxFrequency
		clusters[comb[0]].Node[comb[1]] = n
	}

	logs.GetLogger().Info("\t>> QUERY: kube_node_status_capacity{resource='memory'} * on(...")
	q = querier.PromQLQuery{
		Metric: "kube_node_status_capacity{resource='memory'} * on(icos_agent_id, icos_host_name) group_left(icos_host_id) node_uname_info",
		Params: map[string]string{}}

	for _, node := range querier.Query(q.String()) {
		cluster_id := string(node.Metric["k8s_cluster_uid"])
		node_id := strings.TrimSpace(string(node.Metric["icos_host_id"]))
		ram := int64(node.Value)

		if common.CheckClusterNode(cluster_id, node_id, clusters, q.Metric) {
			n := clusters[cluster_id].Node[node_id]
			n.StaticMetrics.RAMMemory = ram
			clusters[cluster_id].Node[node_id] = n
		}
	}

	// Cluster - Node - DynamicMetrics
	logs.GetLogger().Info("\t>> QUERY: node_thermal_zone_temp")
	q = querier.PromQLQuery{
		Metric: "node_thermal_zone_temp",
		Params: map[string]string{}}

	for _, node := range querier.Query(q.String()) {
		cluster_id := string(node.Metric["k8s_cluster_uid"])
		node_id := strings.TrimSpace(string(node.Metric["icos_host_id"]))
		temp := float64(node.Value)

		if common.CheckClusterNode(cluster_id, node_id, clusters, q.Metric) {
			n := clusters[cluster_id].Node[node_id]
			n.DynamicMetrics.CPUTemperature = temp
			clusters[cluster_id].Node[node_id] = n
		}
	}

	logs.GetLogger().Info("\t>> QUERY: scaph_host_energy_microjoules_total * on...")
	q = querier.PromQLQuery{
		Metric: "scaph_host_energy_microjoules_total * on(icos_agent_id, icos_host_name) group_left(icos_host_id) node_uname_info",
		Params: map[string]string{}}

	for _, node := range querier.Query(q.String()) {
		cluster_id := string(node.Metric["k8s_cluster_uid"])
		node_id := strings.TrimSpace(string(node.Metric["icos_host_id"]))
		energy := float64(node.Value) / 1000000

		if common.CheckClusterNode(cluster_id, node_id, clusters, q.Metric) {
			n := clusters[cluster_id].Node[node_id]
			n.DynamicMetrics.CPUEnergyConsumption = energy
			clusters[cluster_id].Node[node_id] = n
		}
	}

	// Cluster - Node - CPU Frequency
	logs.GetLogger().Info("\t>> QUERY: avg(node_cpu_scaling_frequency_hertz) by (icos_host_id, icos_cluster_id)")
	q = querier.PromQLQuery{
		Metric: "avg(node_cpu_scaling_frequency_hertz) by (icos_host_id, icos_cluster_id)",
		Params: map[string]string{},
	}
	for _, res := range querier.Query(q.String()) {
		cluster_id := string(res.Metric["icos_cluster_id"])
		node_id := strings.TrimSpace(string(res.Metric["icos_host_id"]))
		freq := int64(res.Value)
		if common.CheckClusterNode(cluster_id, node_id, clusters, q.Metric) {
			n := clusters[cluster_id].Node[node_id]
			n.DynamicMetrics.CPUFrequency = freq
			clusters[cluster_id].Node[node_id] = n
		}
	}

	logs.GetLogger().Info("\t>> QUERY: node_memory_MemFree_bytes")
	q = querier.PromQLQuery{
		Metric: "node_memory_MemFree_bytes",
		Params: map[string]string{}}

	for _, node := range querier.Query(q.String()) {
		cluster_id := string(node.Metric["k8s_cluster_uid"])
		node_id := strings.TrimSpace(string(node.Metric["icos_host_id"]))
		ram := int64(node.Value)

		if common.CheckClusterNode(cluster_id, node_id, clusters, q.Metric) {
			n := clusters[cluster_id].Node[node_id]
			n.DynamicMetrics.FreeRAM = ram
			n.DynamicMetrics.UsedRAM = n.StaticMetrics.RAMMemory - ram
			clusters[cluster_id].Node[node_id] = n
		}
	}

	// Get Network interfaces and add them to infra
	logs.GetLogger().Info("\t>> Getting Network interfaces and adding them to infra ...")
	logs.GetLogger().Info("\t>> QUERIES: node_network_info, node_network_up, node_network_address_info, node_network_speed_bytes)")
	infra.SetNetworkInterfaces(clusters, orchs)

	// Cluster - Node - Devices
	logs.GetLogger().Info("\t>> QUERY: node_mounted")
	q = querier.PromQLQuery{
		Metric: "node_mounted",
		Params: map[string]string{}}

	for _, device := range querier.Query(q.String()) {
		cluster_id := string(device.Metric["k8s_cluster_uid"])
		node_id := strings.TrimSpace(string(device.Metric["icos_host_id"]))
		device_name := string(device.Metric["device"])
		device_type := strings.Split(device_name, "_")[0]
		device_status_n := int8(device.Value)
		device_path := string(device.Metric["resource_path"])

		device_status := "unknown"
		switch device_status_n {
		case -1:
			device_status = "detached"
		case 0:
			device_status = "busy"
		case 1:
			device_status = "available"
		}

		if common.CheckClusterNode(cluster_id, node_id, clusters, q.Metric) {
			newDev := models.Device{
				Name:   device_name,
				Type:   device_type,
				Status: device_status,
				Path:   device_path,
			}

			clusters[cluster_id].Node[node_id].Devices[device_name] = newDev
		}

	}

	// Cluster - Node - Labels
	logs.GetLogger().Info("\t>> QUERY: tlum_host_labels")
	q = querier.PromQLQuery{
		Metric: "tlum_host_labels",
		Params: map[string]string{}}

	for _, host_labels := range querier.Query(q.String()) {
		node_id := strings.TrimSpace(string(host_labels.Metric["icos_host_id"]))
		cluster_id := string(host_labels.Metric["icos_cluster_id"])

		if common.CheckClusterNode(cluster_id, node_id, clusters, q.Metric) {
			labels := make(map[string]string)

			for k, v := range host_labels.Metric {
				if strings.HasPrefix(string(k), "label_") {
					lab := strings.Replace(string(k), "label_", "", 1)
					labels[lab] = string(v)
				}
			}

			n := clusters[cluster_id].Node[node_id]
			n.Labels = labels
			clusters[cluster_id].Node[node_id] = n

		}
	}

	// Cluster - Pods
	// Get Pods and add them to infra
	logs.GetLogger().Info("\t>> Getting Pods information and adding them to nodes ...")
	podsList := infra.SetPodsV2(clusters, orchs)

	// Cluster - Pod - Containers
	// Relationship between Pods (kube_pod_info && kube_pod_container_info)
	//   [x] kube_pod_info.pod_name <---> kube_pod_container_info.pod_name
	logs.GetLogger().Info("\t>> QUERY: kube_pod_container_info")
	q = querier.PromQLQuery{
		Metric: "kube_pod_container_info",
		Params: map[string]string{}}

	for _, container := range querier.Query(q.String()) {
		pod_uid := string(container.Metric["uid"])
		cluster_id := string(container.Metric["k8s_cluster_uid"])
		cont_name := string(container.Metric["container"])
		node := strings.TrimSpace(string(container.Metric["icos_host_id"]))
		parentNodeId := podsList[pod_uid].ParentNodeId

		if parentNodeId != "" && common.CheckClusterNodePod(cluster_id, parentNodeId, pod_uid, clusters, q.Metric) {
			newContainer := models.Container{
				Name: cont_name,
				Node: node,
			}
			clusters[cluster_id].Node[parentNodeId].Pod[pod_uid].Container[cont_name] = newContainer
		}
	}

	// Cluster - Pod - Number of containers
	logs.GetLogger().Info("\t>> Cluster - Pod - Setting number of containers ...")
	for cluster_id := range clusters {
		for node_id := range clusters[cluster_id].Node {
			for pod_id := range clusters[cluster_id].Node[node_id].Pod {
				pod := clusters[cluster_id].Node[node_id].Pod[pod_id]
				pod.NumberOfContainers = int32(len(pod.Container))
				clusters[cluster_id].Node[node_id].Pod[pod_id] = pod
			}
		}
	}

	// Cluster - Pod - Container - CPU Usage
	logs.GetLogger().Info("\t>> QUERY: container_cpu_utilization_ratio")
	q = querier.PromQLQuery{
		Metric: "container_cpu_utilization_ratio",
		Params: map[string]string{}}

	for _, container := range querier.Query(q.String()) {
		pod_uid := string(container.Metric["k8s_pod_uid"])
		cluster_id := string(container.Metric["icos_cluster_id"])
		cont_name := string(container.Metric["k8s_container_name"])
		value := container.Value

		if _, existsPod := podsList[pod_uid]; existsPod {
			// icos_host_name is "" ==> not nuvla

			id_node := podsList[pod_uid].ParentNodeId

			if common.CheckClusterNodePodContainer(cluster_id, id_node, pod_uid, cont_name, clusters, q.Metric) {
				cont := clusters[cluster_id].Node[id_node].Pod[pod_uid].Container[cont_name]
				cont.CPUUsage = float64(value)
				clusters[cluster_id].Node[id_node].Pod[pod_uid].Container[cont_name] = cont
			}
		}
	}

	// Cluster
	logs.GetLogger().Info("\t>> QUERY: tlum_ocm_agent_info")
	q = querier.PromQLQuery{
		Metric: "tlum_ocm_agent_info",
		Params: map[string]string{}}

	var clusterIdName = map[string]string{}

	for _, cluster := range querier.Query(q.String()) {
		cluster_id := string(cluster.Metric["k8s_cluster_uid"])
		cluster_name := string(cluster.Metric["name"])

		clusterIdName[cluster_id] = cluster_name
	}

	// Cluster and Node renaming
	logs.GetLogger().Info("\t>> Cluster and Node renaming")
	for key_c, cluster := range clusters {

		// Node
		for key_n, node := range cluster.Node {
			if node.Name != key_n {
				cluster.Node[node.Name] = node
				delete(cluster.Node, key_n)
			}
		}

		// Nuvla nodes (after node renaming)
		if cluster.Type == "nuvla" {
			for key_n, node := range cluster.Node {
				if node.Type == "nuvla" {
					// replace in nodes:
					//   clusters[cluster_id].Node[icos_host_name] ==> key_n (icos_host_name)
					//
					// with:
					//   clusters[cluster_id].Node[id] ==> key_n (id from 'nuvla_device_info' query)
					id_nuvla_node := common.GetNuvlaNodeId(key_n, orchs)
					if len(id_nuvla_node) > 0 {
						cluster.Node[id_nuvla_node] = node
						delete(cluster.Node, key_n)
					}
				}
			}
		}

		// Cluster
		_, existsCluster := clusterIdName[key_c]
		if existsCluster {
			if clusterIdName[key_c] != key_c {
				cluster.Name = clusterIdName[key_c]
				clusters[clusterIdName[key_c]] = cluster
				delete(clusters, key_c)
			}
		}
	}

	// Timestamps
	loc := time.FixedZone("Local", 0)
	tp := time.Date(1970, 1, 1, 0, 0, 0, 0, loc)
	ts := time.Since(tp).Seconds()

	var time = models.Timestamp{
		OldestTimestamp: oldest,
		TimeSinceOldest: ts - oldest,
	}

	var infra = models.Infrastructure{
		Timestamp: time,
		Cluster:   clusters,
	}

	return infra
}
