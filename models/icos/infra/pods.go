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
package infra

import (
	"aggregator/common/logs"
	"aggregator/models/icos/common"
	"aggregator/models/icos/models"
	"aggregator/querier"
	"strings"
)

/**
 * nodes: Retrieves pods from nodes
 *
 * metrics: - kube_pod_info
 *			- kube_pod_status_phase == 1
 *			- tlum_workload_info
 *
 * 	Infrastrcuture:
 *		Timestamp:
 *		 	...
 *		Agent:Cluster[]
 *			...
 *			Node[]
 *				...
 *				Pods[]
 *					...
 */
func SetPods(clusters map[string]models.Cluster, orchs map[string]models.OrchInfoNode) map[string]models.Pod {
	// Cluster - Pods
	var podsList = map[string]models.Pod{} // list of all pods. Is initialized using the UUID value, not the pod name (used in INFRA)

	// Relationship between Nodes ('node_uname_info') and Pods ('kube_pod_info')
	// Node Identifier ('Node[id]') in 'clusters[id].Node[id].Pod[id]' is the result of: 'node_uname_info'.icos_host_id
	// [x] 'kube_pod_info'.icos_host_name == 'node_uname_info'.icos_host_name
	// [x] 'kube_pod_info'.icos_cluster_id == 'node_uname_info'.icos_cluster_id
	logs.GetLogger().Info("\t>> QUERY <pods>: 'kube_pod_info'")
	q := querier.PromQLQuery{
		Metric: "kube_pod_info",
		Params: map[string]string{}}

	for _, pod := range querier.Query(q.String()) {
		pod_uid := string(pod.Metric["uid"])
		icos_cluster_id := string(pod.Metric["icos_cluster_id"]) // icos_cluster_id instead of k8s_cluster_uid
		icos_host_name := string(pod.Metric["icos_host_name"])
		pod_name := string(pod.Metric["pod"])
		pod_ip := string(pod.Metric["pod_ip"])

		n := common.GetNodeInfo(icos_host_name, clusters[icos_cluster_id].Node)

		if n.Uuid != "" {
			parentNodeId := n.Uuid

			var podInfo = models.Pod{
				Uid:                pod_uid,         // uid from pod. Used to match pods in next queries
				ParentNodeId:       parentNodeId,    // hidden value: edge_harbor_host_id / Uuid from parent Node
				ClusterUid:         icos_cluster_id, // hidden value: cluster uid
				IcosHostName:       icos_host_name,  // hidden value
				Name:               pod_name,
				IP:                 pod_ip,
				Status:             "",
				NumberOfContainers: 0,
				NumberOfApps:       0,
				Container:          map[string]models.Container{},
				Workload:           map[string]models.Workload{},
			}
			podsList[pod_uid] = podInfo

			if common.CheckCluster(icos_cluster_id, clusters, q.Metric) {
				clusters[icos_cluster_id].Node[parentNodeId].Pod[pod_uid] = podInfo //Pod[pod_name] = podInfo
			}
		}
	}

	logs.GetLogger().Debug("\t>> Clusters (Uuid), Nodes (Uuid, IcosHostName) and Pods added to infra: ")
	for _, c := range clusters {
		logs.GetLogger().Debug("\t     * " + c.Uuid)
		for _, n := range clusters[c.Uuid].Node {
			logs.GetLogger().Debug("\t       - " + n.Uuid + " (" + n.IcosHostName + ")")
			for _, p := range clusters[c.Uuid].Node[n.Uuid].Pod {
				logs.GetLogger().Debug("\t         + " + p.Uid + " (" + p.Name + ")")
			}
		}
	}

	// Cluster - Pod - Status
	// Relationship between Pods ('kube_pod_info' & 'kube_pod_status_phase == 1')
	// [x] 'kube_pod_info'.uid <---> 'kube_pod_status_phase'.uid
	logs.GetLogger().Info("\t>> QUERY: kube_pod_status_phase == 1")
	q = querier.PromQLQuery{
		Metric: "kube_pod_status_phase == 1",
		Params: map[string]string{}}

	for _, pod := range querier.Query(q.String()) {
		pod_uid := string(pod.Metric["uid"])
		icos_cluster_id := string(pod.Metric["icos_cluster_id"]) // icos_cluster_id instead of k8s_cluster_uid
		//pod_name := string(pod.Metric["pod"])
		status := string(pod.Metric["phase"])
		parentNodeId := podsList[pod_uid].ParentNodeId

		if common.CheckClusterNodePod(icos_cluster_id, parentNodeId, pod_uid, clusters, q.Metric) {
			podInfo := podsList[pod_uid]
			podInfo.Status = status
			clusters[icos_cluster_id].Node[parentNodeId].Pod[pod_uid] = podInfo //Pod[pod_name] = podInfo
		}
	}

	// tlum_workload_info
	// Relationship between Pods ('kube_pod_info') and Workloads ('tlum_workload_info'):
	// [x] 'kube_pod_info'.uid == 'tlum_workload_info'.id
	// 	==> pod_uid := string(pod.Metric["uid"])
	//  ==> podsList[pod_uid]
	logs.GetLogger().Info("\t>> QUERY: tlum_workload_info")
	q = querier.PromQLQuery{
		Metric: "tlum_workload_info",
		Params: map[string]string{}}

	for _, w := range querier.Query(q.String()) {
		icos_app_name := string(w.Metric["icos_app_name"])
		icos_app_instance := string(w.Metric["icos_app_instance"])
		icos_app_component := string(w.Metric["icos_app_component"])
		pod_uid := string(w.Metric["id"])                                         // id instead of k8s_pod_uid
		icos_cluster_id := strings.TrimSpace(string(w.Metric["icos_cluster_id"])) // icos_cluster_id instead of k8s_cluster_uid
		parentNodeId := ""
		podName := string(w.Metric["name"]) //""

		logs.GetLogger().Debug("QUERY: tlum_workload_info: " + podName)

		if _, existsPod := podsList[pod_uid]; existsPod {
			parentNodeId = podsList[pod_uid].ParentNodeId
			//podName = podsList[pod_uid].Name

			logs.GetLogger().Debug("QUERY: tlum_workload_info: POD EXISTS!! [parentNodeId: " + parentNodeId + "]")
		} else {
			logs.GetLogger().Warn("QUERY: tlum_workload_info: POD DOESNT EXIST!!")
		}

		if common.CheckClusterNodePod(icos_cluster_id, parentNodeId, pod_uid, clusters, q.Metric) && icos_app_instance != "" {
			newWorkload := models.Workload{
				AppName:      icos_app_name,
				AppInstance:  icos_app_instance,
				AppComponent: icos_app_component,
			}
			clusters[icos_cluster_id].Node[parentNodeId].Pod[pod_uid].Workload[icos_app_instance] = newWorkload

			logs.GetLogger().Debug("QUERY: tlum_workload_info: checkClusterNodePod returned TRUE!!")
		} else {
			logs.GetLogger().Warn("QUERY: tlum_workload_info: checkClusterNodePod returned FALSE!!")
		}
	}

	return podsList
}

/**
 * nodes: Retrieves pods from nodes
 *
 * metrics: - tlum_workload_info
 *			- kube_pod_info
 *			- kube_pod_status_phase == 1
 *
 * 	Infrastrcuture:
 *		Timestamp:
 *		 	...
 *		Agent:Cluster[]
 *			...
 *			Node[]
 *					...
 *				Pods[]
 *					...
 */
func SetPodsV2(clusters map[string]models.Cluster, orchs map[string]models.OrchInfoNode) map[string]models.Pod {
	// Cluster - Pods
	var podsList = map[string]models.Pod{} // list of all pods. Is initialized using the UUID value, not the pod name (used in INFRA)

	// List of all Pods
	logs.GetLogger().Info("\t>> QUERY <pods>: 'tlum_workload_info'")
	q := querier.PromQLQuery{
		Metric: "tlum_workload_info",
		Params: map[string]string{}}

	for _, pod := range querier.Query(q.String()) {
		pod_uid := strings.TrimSpace(string(pod.Metric["id"]))
		icos_cluster_id := strings.TrimSpace(string(pod.Metric["icos_cluster_id"]))
		icos_node_id := strings.TrimSpace(string(pod.Metric["icos_host_id"]))
		icos_host_name := string(pod.Metric["icos_host_name"])
		pod_name := string(pod.Metric["name"])

		if common.CheckClusterNode(icos_cluster_id, icos_node_id, clusters, q.Metric) {
			var p = models.Pod{
				Uid:                pod_uid,
				ParentNodeId:       icos_node_id, // Node id
				ClusterUid:         icos_cluster_id,
				IcosHostName:       icos_host_name, // Node name
				Name:               pod_name,
				IP:                 "",
				Status:             "",
				NumberOfContainers: 0,
				NumberOfApps:       0,
				Container:          map[string]models.Container{},
				Workload:           map[string]models.Workload{},
			}

			podsList[pod_uid] = p
			clusters[icos_cluster_id].Node[icos_node_id].Pod[pod_uid] = p

			// Workload
			icos_app_name := string(pod.Metric["icos_app_name"])
			icos_app_instance := string(pod.Metric["icos_app_instance"])
			icos_app_component := string(pod.Metric["icos_app_component"])

			if icos_app_instance != "" {
				w := models.Workload{
					AppName:      icos_app_name,
					AppInstance:  icos_app_instance,
					AppComponent: icos_app_component,
				}
				clusters[icos_cluster_id].Node[icos_node_id].Pod[pod_uid].Workload[icos_app_instance] = w
				podsList[pod_uid].Workload[icos_app_instance] = w
			}
		}
	}

	// Cluster - Pod - IP
	logs.GetLogger().Info("\t>> QUERY <pods>: 'kube_pod_info'")
	q = querier.PromQLQuery{
		Metric: "kube_pod_info",
		Params: map[string]string{}}

	for _, pod := range querier.Query(q.String()) {
		pod_uid := string(pod.Metric["uid"])
		icos_cluster_id := string(pod.Metric["icos_cluster_id"])
		pod_ip := string(pod.Metric["pod_ip"])
		parentNodeId := podsList[pod_uid].ParentNodeId

		if parentNodeId != "" && common.CheckClusterNodePod(icos_cluster_id, parentNodeId, pod_uid, clusters, q.Metric) {
			p := clusters[icos_cluster_id].Node[parentNodeId].Pod[pod_uid]
			p.IP = pod_ip
			clusters[icos_cluster_id].Node[parentNodeId].Pod[pod_uid] = p
			podsList[pod_uid] = p
		}
	}

	// Cluster - Pod - Status
	logs.GetLogger().Info("\t>> QUERY: kube_pod_status_phase == 1")
	q = querier.PromQLQuery{
		Metric: "kube_pod_status_phase == 1",
		Params: map[string]string{}}

	for _, pod := range querier.Query(q.String()) {
		pod_uid := string(pod.Metric["uid"])
		icos_cluster_id := string(pod.Metric["icos_cluster_id"])
		status := string(pod.Metric["phase"])
		parentNodeId := podsList[pod_uid].ParentNodeId

		if parentNodeId != "" && common.CheckClusterNodePod(icos_cluster_id, parentNodeId, pod_uid, clusters, q.Metric) {
			p := clusters[icos_cluster_id].Node[parentNodeId].Pod[pod_uid]
			p.Status = status
			clusters[icos_cluster_id].Node[parentNodeId].Pod[pod_uid] = p
			podsList[pod_uid] = p
		}
	}

	// Debug: Print all cluster, nodes and pods
	logs.GetLogger().Debug("\t>> Clusters (Uuid), Nodes (Uuid, IcosHostName) and Pods added to infra: ")
	for _, c := range clusters {
		logs.GetLogger().Debug("\t     * " + c.Uuid)
		for _, n := range clusters[c.Uuid].Node {
			logs.GetLogger().Debug("\t       - " + n.Uuid + " (" + n.IcosHostName + ")")
			for _, p := range clusters[c.Uuid].Node[n.Uuid].Pod {
				logs.GetLogger().Debug("\t         + " + p.Uid + " (" + p.Name + ")")
			}
		}
	}

	return podsList
}
