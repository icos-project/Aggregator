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
package models_icos

import (
	logs "aggregator/common/logs"
	"aggregator/querier"
	"encoding/json"
	"math"
	"strconv"
	"strings"
	"time"
)

// path used in logs
const pathLOG string = "ICOS > "

func GetInfra() []byte {

	// Get metrics from Thanos
	logs.GetLogger().Info(pathLOG + "> > > > > > Getting metrics from Thanos ...")
	clusters := queryPrometheus()

	// Convert to JSON
	json, err := json.MarshalIndent(clusters, "", "\t")
	if err != nil {
		logs.GetLogger().Error("Error marshaling models: ", err)
	}

	return json
}

func queryPrometheus() Infrastructure {

	var clusters = map[string]Cluster{}
	var orchs = map[string]OrchInfoNode{}

	// Timestamps
	logs.GetLogger().Debug(pathLOG + "QUERY: timestamp(up)")
	q := querier.PromQLQuery{
		Metric: "timestamp(up)",
		Params: map[string]string{}}

	oldest := math.Inf(1)
	for _, timestamp := range querier.Query(q.String()) {
		if float64(timestamp.Value) < oldest {
			oldest = float64(timestamp.Value)
		}
	}

	// clusters: tlum_orch_info
	logs.GetLogger().Debug(pathLOG + "QUERY: tlum_orch_info")
	q = querier.PromQLQuery{
		Metric: "tlum_orch_info",
		Params: map[string]string{}}

	for _, node := range querier.Query(q.String()) {
		orch_id := string(node.Metric["agent_id"])

		var orchInfo = OrchInfoNode{
			Id:            orch_id,
			Type:          string(node.Metric["type"]),
			Name:          string(node.Metric["agent_name"]), // Nuvla cluster
			Engine:        "unknown",                         // default, may be updated
			Uuid:          string(node.Metric["agent_id"]),
			K8sClusterUid: string(node.Metric["k8s_cluster_uid"]), // OCM
			IcosHostName:  string(node.Metric["icos_host_name"]),  // Nuvla
			K8sNodeName:   string(node.Metric["k8s_node_name"]),
		}
		orchs[orch_id] = orchInfo
	}

	// clusters: tlum_runtime_info
	logs.GetLogger().Debug(pathLOG + "QUERY: tlum_runtime_info")
	q = querier.PromQLQuery{
		Metric: "tlum_runtime_info",
		Params: map[string]string{}}

	for _, node := range querier.Query(q.String()) {
		engine := string(node.Metric["type"])
		icos_host_name := string(node.Metric["icos_host_name"])
		k8s_cluster_uid := string(node.Metric["k8s_cluster_uid"])

		for _, o := range orchs {
			if o.IcosHostName == icos_host_name || o.K8sClusterUid == k8s_cluster_uid {

				// First we get a "copy" of the entry
				if entry, ok := orchs[o.Id]; ok {

					// Then we modify the copy
					entry.Engine = engine

					// Then we reassign map entry
					orchs[o.Id] = entry
				}

			}
		}

	}

	// Clusters and Nodes
	logs.GetLogger().Debug(pathLOG + "QUERY: node_uname_info")
	q = querier.PromQLQuery{
		Metric: "node_uname_info",
		Params: map[string]string{}}

	for _, node := range querier.Query(q.String()) {

		cluster_id := string(node.Metric["k8s_cluster_uid"])
		cluster_type := "unknown"
		node_id := string(node.Metric["icos_host_id"])
		node_name := string(node.Metric["nodename"])
		net_host_name := string(node.Metric["net_host_name"])
		icos_host_name := string(node.Metric["icos_host_name"])
		icos_agent_id := string(node.Metric["icos_agent_id"])
		k8s_node_uid := string(node.Metric["k8s_node_uid"])
		architecture := string(node.Metric["machine"])
		latitude, _ := strconv.ParseFloat(string(node.Metric["icos_alpha_latitude"]), 8)
		longitude, _ := strconv.ParseFloat(string(node.Metric["icos_alpha_longitude"]), 8)

		if cluster_id != "self" {
			if isNuvlaCluster(icos_host_name, orchs) {
				cluster_id = getNuvlaClusterName(icos_host_name, orchs)
				cluster_type = "nuvla"
			} else if isOCMCluster(cluster_id, orchs) {
				cluster_type = "ocm"
			}

			engine := getEngine(icos_host_name, cluster_id, orchs)

			//logs.GetLogger().Debug(pathLOG+"[node_uname_info] [cluster_id:", cluster_id, "] [icos_host_name:", icos_host_name, "] [cluster_type:", cluster_type, "]")

			if cluster_type != "" {
				if _, exists := clusters[cluster_id]; !exists {
					var newCluster = Cluster{
						Uuid:        cluster_id,
						Type:        cluster_type,
						Engine:      engine,
						ICOSAgentID: icos_agent_id,
						Node:        map[string]Node{},
					}
					clusters[cluster_id] = newCluster
				}

				newNode := Node{
					Uuid:         node_id,
					Type:         cluster_type,
					Name:         node_name,
					Engine:       engine,
					NetHostName:  net_host_name,
					IcosHostName: icos_host_name, // matches value of 'icos_host_name' in 'kube_pod_info' query
					K8sNodeUid:   k8s_node_uid,   // matches value of 'k8s_node_uid' in 'kube_pod_info' query
					Location: Location{
						Latitude:  latitude,
						Longitude: longitude,
					},
					StaticMetrics:     StaticMetrics{CPUArchitecture: architecture},
					NetworkInterfaces: map[string]Interface{},
					Devices:           map[string]Device{},
					Pod:               map[string]Pod{},
				}
				clusters[cluster_id].Node[node_id] = newNode
			} else {
				logs.GetLogger().Warn(pathLOG+"Unknown cluster type [cluster_id:", cluster_id, "]. Cluster not added to infraestructure.")
			}

		}
	}

	// Cluster - Node - vulnerabilities and SCA_score
	//   node_uname_info > net_host_name="10.150.0.144", nodename="icosedge"
	//					 ==> Node.NetHostName = net_host_name (NEW), Node.Name = nodename
	//   SCA_score       > agent_hostname="icosedge", agent_ip="10.150.0.144"
	//   vulnerabilities > agent_hostname="icosedge", agent_ip="10.150.0.144"
	//					 ==> Node.NetHostName == agent_ip && Node.Name == agent_hostname

	// Cluster - Node - SCA_score
	logs.GetLogger().Debug(pathLOG + "QUERY: SCA_score")
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
	logs.GetLogger().Debug(pathLOG + "QUERY: vulnerabilities")
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
	logs.GetLogger().Debug(pathLOG + "QUERY: kube_node_status_capacity{resource='cpu'} * on(icos_agen...")
	q = querier.PromQLQuery{
		Metric: "kube_node_status_capacity{resource='cpu'} * on(icos_agent_id, icos_host_name) group_left(icos_host_id) node_uname_info",
		Params: map[string]string{}}

	for _, node := range querier.Query(q.String()) {
		cluster_id := string(node.Metric["k8s_cluster_uid"])
		node_id := string(node.Metric["icos_host_id"])
		cores := int32(node.Value)

		if checkClusterNode(cluster_id, node_id, clusters, q.Metric) {
			n := clusters[cluster_id].Node[node_id]
			n.StaticMetrics.CPUCores = cores
			clusters[cluster_id].Node[node_id] = n
		}
	}

	logs.GetLogger().Debug(pathLOG + "QUERY: node_cpu_frequency_max_hertz")
	q = querier.PromQLQuery{
		Metric: "node_cpu_frequency_max_hertz",
		Params: map[string]string{}}

	freqs := make(map[[2]string][]int64)
	for _, node := range querier.Query(q.String()) {
		cluster_id := string(node.Metric["k8s_cluster_uid"])
		node_id := string(node.Metric["icos_host_id"])
		frequency := int64(node.Value)

		if checkClusterNode(cluster_id, node_id, clusters, q.Metric) {
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
		maxFrequency := maxInt64(f)
		n := clusters[comb[0]].Node[comb[1]]
		n.StaticMetrics.CPUMaxFrequency = maxFrequency
		clusters[comb[0]].Node[comb[1]] = n
	}

	logs.GetLogger().Debug(pathLOG + "QUERY: kube_node_status_capacity{resource='memory'} * on(...")
	q = querier.PromQLQuery{
		Metric: "kube_node_status_capacity{resource='memory'} * on(icos_agent_id, icos_host_name) group_left(icos_host_id) node_uname_info",
		Params: map[string]string{}}

	for _, node := range querier.Query(q.String()) {
		cluster_id := string(node.Metric["k8s_cluster_uid"])
		node_id := string(node.Metric["icos_host_id"])
		ram := int64(node.Value)

		if checkClusterNode(cluster_id, node_id, clusters, q.Metric) {
			n := clusters[cluster_id].Node[node_id]
			n.StaticMetrics.RAMMemory = ram
			clusters[cluster_id].Node[node_id] = n
		}
	}

	// Cluster - Node - DynamicMetrics
	logs.GetLogger().Debug(pathLOG + "QUERY: node_thermal_zone_temp")
	q = querier.PromQLQuery{
		Metric: "node_thermal_zone_temp",
		Params: map[string]string{}}

	for _, node := range querier.Query(q.String()) {
		cluster_id := string(node.Metric["k8s_cluster_uid"])
		node_id := string(node.Metric["icos_host_id"])
		temp := float64(node.Value)

		if checkClusterNode(cluster_id, node_id, clusters, q.Metric) {
			n := clusters[cluster_id].Node[node_id]
			n.DynamicMetrics.CPUTemperature = temp
			clusters[cluster_id].Node[node_id] = n
		}
	}

	logs.GetLogger().Debug(pathLOG + "QUERY: scaph_host_energy_microjoules_total * on...")
	q = querier.PromQLQuery{
		Metric: "scaph_host_energy_microjoules_total * on(icos_agent_id, icos_host_name) group_left(icos_host_id) node_uname_info",
		Params: map[string]string{}}

	for _, node := range querier.Query(q.String()) {
		cluster_id := string(node.Metric["k8s_cluster_uid"])
		node_id := string(node.Metric["icos_host_id"])
		energy := float64(node.Value) / 1000000

		if checkClusterNode(cluster_id, node_id, clusters, q.Metric) {
			n := clusters[cluster_id].Node[node_id]
			n.DynamicMetrics.CPUEnergyConsumption = energy
			clusters[cluster_id].Node[node_id] = n
		}
	}

	logs.GetLogger().Debug(pathLOG + "QUERY: node_memory_MemFree_bytes")
	q = querier.PromQLQuery{
		Metric: "node_memory_MemFree_bytes",
		Params: map[string]string{}}

	for _, node := range querier.Query(q.String()) {
		cluster_id := string(node.Metric["k8s_cluster_uid"])
		node_id := string(node.Metric["icos_host_id"])
		node_name := string(node.Metric["icos_host_name"])
		ram := int64(node.Value)

		if isNuvlaCluster(node_name, orchs) {
			cluster_id = getNuvlaClusterName(node_name, orchs)
		}

		if checkClusterNode(cluster_id, node_id, clusters, q.Metric) {
			n := clusters[cluster_id].Node[node_id]
			n.DynamicMetrics.FreeRAM = ram
			n.DynamicMetrics.UsedRAM = n.StaticMetrics.RAMMemory - ram
			clusters[cluster_id].Node[node_id] = n
		}
	}

	// Get Network interfaces and add them to infra
	logs.GetLogger().Debug(pathLOG + "Get Network interfaces and add them to infra")
	processNetworkInterfaces(clusters, orchs)

	// Cluster - Node - Devices
	logs.GetLogger().Debug(pathLOG + "QUERY: node_mounted")
	q = querier.PromQLQuery{
		Metric: "node_mounted",
		Params: map[string]string{}}

	for _, device := range querier.Query(q.String()) {
		cluster_id := string(device.Metric["k8s_cluster_uid"])
		node_id := string(device.Metric["icos_host_id"])
		node_name := string(device.Metric["icos_host_name"])
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

		if isNuvlaCluster(node_name, orchs) {
			cluster_id = getNuvlaClusterName(node_name, orchs)
		}

		if checkClusterNode(cluster_id, node_id, clusters, q.Metric) {
			newDev := Device{
				Name:   device_name,
				Type:   device_type,
				Status: device_status,
				Path:   device_path,
			}

			clusters[cluster_id].Node[node_id].Devices[device_name] = newDev
		}

	}

	// Cluster - Node - Labels
	logs.GetLogger().Debug(pathLOG + "QUERY: tlum_host_labels")
	q = querier.PromQLQuery{
		Metric: "tlum_host_labels",
		Params: map[string]string{}}

	for _, host_labels := range querier.Query(q.String()) {
		node_id := string(host_labels.Metric["icos_host_id"])
		cluster_id := string(host_labels.Metric["k8s_cluster_uid"])
		node_name := string(host_labels.Metric["icos_host_name"])

		if isNuvlaCluster(node_name, orchs) {
			cluster_id = getNuvlaClusterName(node_name, orchs)
		}

		if checkClusterNode(cluster_id, node_id, clusters, q.Metric) {
			//labels := []Label{}
			labels := make(map[string]string)

			for k, v := range host_labels.Metric {
				//fmt.Println(k, " ", v) // labels: key + value

				if strings.HasPrefix(string(k), "label_") {
					/*newLabel := Label{
						Key:   string(k),
						Value: string(v),
					}

					labels = append(labels, newLabel)*/
					lab := strings.Replace(string(k), "label_", "", 1)
					labels[lab] = string(v)
				}

			}

			n := clusters[cluster_id].Node[node_id]
			n.Labels = labels
			clusters[cluster_id].Node[node_id] = n

		}
	}

	//////////////////////////////////////

	// Cluster - Pods
	var podsList = map[string]Pod{} // list of all pods. Is initialized using the UUID value, not the pod name (used in INFRA)

	// Relationship between Nodes (node_uname_info) and Pods (kube_pod_info)
	// - QUERY: kube_pod_info <---> node_uname_info:
	//   [x] kube_pod_info.k8s_node_uid <---> node_uname_info.k8s_node_uid)
	logs.GetLogger().Debug("> Getting Pods information and adding them to nodes ...")
	logs.GetLogger().Debug("QUERY <pods>: 'kube_pod_info'")
	q = querier.PromQLQuery{
		Metric: "kube_pod_info",
		Params: map[string]string{}}

	for _, pod := range querier.Query(q.String()) {
		pod_uid := string(pod.Metric["uid"])
		parentNodeName := string(pod.Metric["node"])
		cluster_uid := string(pod.Metric["k8s_cluster_uid"])
		k8s_node_uid := string(pod.Metric["k8s_node_uid"])
		host_name := string(pod.Metric["icos_host_name"])
		pod_name := string(pod.Metric["pod"])
		pod_ip := string(pod.Metric["pod_ip"])

		if isNuvlaCluster(host_name, orchs) {
			cluster_uid = getNuvlaClusterName(host_name, orchs)
		}

		n := getNodeInfo(k8s_node_uid, clusters[cluster_uid].Node)

		if n.Uuid != "" {
			parentNodeId := n.Uuid

			var podInfo = Pod{
				Uid:                pod_uid,        // uid from pod. Used to match pods in next queries
				ParentNodeName:     parentNodeName, // hidden value: parent Nodename
				ParentNodeId:       parentNodeId,   // hidden value: edge_harbor_host_id / Uuid from parent Node
				ClusterUid:         cluster_uid,    // hidden value: cluster uid
				IcosHostName:       host_name,      // hidden value
				K8sNodeUid:         k8s_node_uid,   // hidden value
				Name:               pod_name,
				IP:                 pod_ip,
				Status:             "",
				NumberOfContainers: 0,
				NumberOfApps:       0,
				Container:          map[string]Container{},
				Workload:           map[string]Workload{},
			}
			podsList[pod_uid] = podInfo

			if checkCluster(cluster_uid, clusters, q.Metric) {
				clusters[cluster_uid].Node[parentNodeId].Pod[pod_name] = podInfo
			}
		}
	}

	// Cluster - Pod - Status
	// Relationship between Pods (kube_pod_info && kube_pod_status_phase == 1)
	//   [x] kube_pod_info.pod_name <---> kube_pod_status_phase.pod_name
	logs.GetLogger().Debug(pathLOG + "QUERY: kube_pod_status_phase == 1")
	q = querier.PromQLQuery{
		Metric: "kube_pod_status_phase == 1",
		Params: map[string]string{}}

	for _, pod := range querier.Query(q.String()) {
		pod_uid := string(pod.Metric["uid"])
		cluster_id := string(pod.Metric["k8s_cluster_uid"])
		pod_name := string(pod.Metric["pod"])
		status := string(pod.Metric["phase"])
		parentNodeId := podsList[pod_uid].ParentNodeId

		if checkClusterNodePod(cluster_id, parentNodeId, pod_name, clusters, q.Metric) {
			podInfo := podsList[pod_uid]
			podInfo.Status = status
			clusters[cluster_id].Node[parentNodeId].Pod[pod_name] = podInfo
		}
	}

	// tlum_workload_info
	logs.GetLogger().Debug(pathLOG + "QUERY: tlum_workload_info")
	q = querier.PromQLQuery{
		Metric: "tlum_workload_info",
		Params: map[string]string{}}

	for _, w := range querier.Query(q.String()) {
		icos_app_name := string(w.Metric["icos_app_name"])
		icos_app_instance := string(w.Metric["icos_app_instance"])
		icos_app_component := string(w.Metric["icos_app_component"])
		k8s_pod_uid := string(w.Metric["k8s_pod_uid"])
		k8s_cluster_uid := string(w.Metric["k8s_cluster_uid"])
		parentNodeId := ""
		podName := string(w.Metric["name"]) //""

		if _, existsPod := podsList[k8s_pod_uid]; existsPod {
			parentNodeId = podsList[k8s_pod_uid].ParentNodeId
			//podName = podsList[k8s_pod_uid].Name
		}

		if checkClusterNodePod(k8s_cluster_uid, parentNodeId, podName, clusters, q.Metric) && icos_app_instance != "" {
			newWorkload := Workload{
				AppName:      icos_app_name,
				AppInstance:  icos_app_instance,
				AppComponent: icos_app_component,
			}
			clusters[k8s_cluster_uid].Node[parentNodeId].Pod[podName].Workload[icos_app_instance] = newWorkload
		}
	}

	// Cluster - Pod - Containers
	// Relationship between Pods (kube_pod_info && kube_pod_container_info)
	//   [x] kube_pod_info.pod_name <---> kube_pod_container_info.pod_name
	logs.GetLogger().Debug(pathLOG + "QUERY: kube_pod_container_info")
	q = querier.PromQLQuery{
		Metric: "kube_pod_container_info",
		Params: map[string]string{}}

	for _, container := range querier.Query(q.String()) {
		pod_uid := string(container.Metric["uid"])
		cluster_id := string(container.Metric["k8s_cluster_uid"])
		pod_name := string(container.Metric["pod"])
		cont_name := string(container.Metric["container"])
		node := string(container.Metric["icos_host_id"])
		parentNodeId := podsList[pod_uid].ParentNodeId

		if checkClusterNodePod(cluster_id, parentNodeId, pod_name, clusters, q.Metric) {
			newContainer := Container{
				Name: cont_name,
				Node: node,
			}
			clusters[cluster_id].Node[parentNodeId].Pod[pod_name].Container[cont_name] = newContainer
		}
	}

	// Cluster - Pod - Number of containers
	logs.GetLogger().Debug(pathLOG + "Cluster - Pod - Number of containers")
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
	logs.GetLogger().Debug(pathLOG + "QUERY: container_cpu_utilization_ratio")
	q = querier.PromQLQuery{
		Metric: "container_cpu_utilization_ratio",
		Params: map[string]string{}}

	for _, container := range querier.Query(q.String()) {
		node_id := string(container.Metric["k8s_node_uid"])
		pod_uid := string(container.Metric["k8s_pod_uid"])
		cluster_id := string(container.Metric["k8s_cluster_uid"])
		cont_name := string(container.Metric["k8s_container_name"])
		node_name := string(container.Metric["icos_host_name"])
		value := container.Value

		n := getNodeInfo(node_id, clusters[cluster_id].Node)
		pod_name := ""

		if p, existsPod := podsList[pod_uid]; existsPod {
			pod_name = p.Name
		}

		if isNuvlaCluster(node_name, orchs) {
			cluster_id = getNuvlaClusterName(node_name, orchs)
		}

		if checkClusterNodePodContainer(cluster_id, n.Uuid, pod_name, cont_name, clusters, q.Metric) {
			cont := clusters[cluster_id].Node[n.Uuid].Pod[pod_name].Container[cont_name]
			cont.CPUUsage = float64(value)
			clusters[cluster_id].Node[n.Uuid].Pod[pod_name].Container[cont_name] = cont
		}
	}

	// Cluster
	logs.GetLogger().Debug(pathLOG + "QUERY: tlum_ocm_agent_info")
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
	logs.GetLogger().Debug(pathLOG + "Cluster and Node renaming")
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
					id_nuvla_node := getNuvlaNodeId(key_n, orchs)
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
