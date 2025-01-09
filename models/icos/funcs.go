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
	"math"
	"strings"

	logs "aggregator/common/logs"
)

func maxInt64(list []int64) int64 {
	max := int64(math.Inf(-1))
	for _, n := range list {
		if n > max {
			max = n
		}
	}
	return max
}

// NUVLA and OCM clusters / nodes

// check if 'cluster_id' and 'icos_host_name' from 'node_uname_info' query correspond to a Nuvla node / cluster
// Example:
// - cluster_id=""
// - icos_host_name="icos-uc2-test-001"
// - OrchInfoNode.icos_host_name="icos-uc2-test-001"
// ==>  OrchInfoNode[i].IcosHostName == icos_host_name ==> nuvla
func isNuvlaCluster(icos_host_name string, orchs map[string]OrchInfoNode) bool {
	for _, n := range orchs {
		if icos_host_name == "" {
			//logs.GetLogger().Warn(pathLOG+"'icos_host_name' value is empty [icos_host_name:", icos_host_name, "]")
			return false
		}

		if n.Type == strings.ToLower("nuvla") && n.IcosHostName == icos_host_name {
			return true
		}
	}

	return false
}

func getNuvlaClusterName(icos_host_name string, orchs map[string]OrchInfoNode) string {

	for _, n := range orchs {
		if n.Type == strings.ToLower("nuvla") && n.IcosHostName == icos_host_name {
			return n.Name
		}
	}

	return ""
}

// check if 'k8s_cluster_uid' from 'node_uname_info' query correspond to an OCM node / cluster
func isOCMCluster(k8s_cluster_uid string, orchs map[string]OrchInfoNode) bool {

	for _, n := range orchs {
		if k8s_cluster_uid == "" {
			logs.GetLogger().Warn(pathLOG+"'k8s_cluster_uid' value is empty [k8s_cluster_uid:", k8s_cluster_uid, "]")
			return false
		}

		if n.Type == strings.ToLower("ocm") && n.K8sClusterUid == k8s_cluster_uid {
			return true
		}
	}

	return false
}

// get Id from nuvla node using the icos_host_name value
func getNuvlaNodeId(icos_host_name string, orchs map[string]OrchInfoNode) string {

	for _, n := range orchs {
		if n.IcosHostName == icos_host_name {
			return n.Id
		}
	}

	return "" // NOT FOUND / already deleted
}

// get engine from nuvla node using the icos_host_name or cluster_id value
func getEngine(icos_host_name string, cluster_id string, orchs map[string]OrchInfoNode) string {
	for _, n := range orchs {
		if n.IcosHostName == icos_host_name || n.Id == cluster_id {
			return n.Engine
		}
	}

	return "unknown" // NOT FOUND / already deleted
}

// Nodes, pods and containers info

// getNodeInfo
func getNodeInfo(nodeId string, nodes map[string]Node) Node {
	for _, n := range nodes {
		if strings.EqualFold(n.K8sNodeUid, nodeId) {
			//logs.GetLogger().Debug("Node matches input values: pod_name [", pod_name, "] parentNodeUid [", parentNodeUid, "]. Returning node ...")
			return n
		}
	}

	logs.GetLogger().Warn("Node not found. Input values: parentNodeUid [", nodeId, "]")
	return Node{} // NOT FOUND
}

// Checks

// checkCluster checks if cluster is not "self" and if it exists in clusters list
func checkCluster(clusterid string, clusters map[string]Cluster, metric string) bool {
	if clusterid == "self" {
		return false
	}

	if _, existsCluster := clusters[clusterid]; existsCluster {
		return true
	}

	logs.GetLogger().Warn(pathLOG+"Unknown CLUSTER ", clusterid, " [", metric, "]")
	return false
}

// checkClusterNode
func checkClusterNode(clusterid string, nodeid string, clusters map[string]Cluster, metric string) bool {
	if checkCluster(clusterid, clusters, metric) {
		if _, existsNode := clusters[clusterid].Node[nodeid]; existsNode {
			return true
		}
		logs.GetLogger().Warn(pathLOG+"Unknown NODE ", nodeid, " in CLUSTER ", clusterid, " [", metric, "]")
	}

	return false
}

// checkClusterNodePod
func checkClusterNodePod(clusterid string, nodeid string, podid string, clusters map[string]Cluster, metric string) bool {
	if checkClusterNode(clusterid, nodeid, clusters, metric) {
		if _, existsPod := clusters[clusterid].Node[nodeid].Pod[podid]; existsPod {
			return true
		}
		logs.GetLogger().Warn(pathLOG+"Unknown POD ", podid, " in NODE ", nodeid, " and CLUSTER ", clusterid, " [", metric, "]")
	}

	return false
}

// checkClusterNodePodContainer
func checkClusterNodePodContainer(clusterid string, nodeid string, podid string, containerid string, clusters map[string]Cluster, metric string) bool {

	if checkClusterNodePod(clusterid, nodeid, podid, clusters, metric) {
		if _, existsContainer := clusters[clusterid].Node[nodeid].Pod[podid].Container[containerid]; existsContainer {
			return true
		}
		logs.GetLogger().Warn(pathLOG+"Unknown CONTAINER ", containerid, " in POD ", podid, " and NODE ", nodeid, " and CLUSTER ", clusterid, " [", metric, "]")
	}

	return false
}
