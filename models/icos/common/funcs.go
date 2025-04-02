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
package common

import (
	"math"
	"strings"

	logs "aggregator/common/logs"
	"aggregator/models/icos/models"
)

func MaxInt64(list []int64) int64 {
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
func IsNuvlaCluster(icos_host_name string, orchs map[string]models.OrchInfoNode) bool {
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

func GetNuvlaClusterName(icos_host_name string, orchs map[string]models.OrchInfoNode) string {

	for _, n := range orchs {
		if n.Type == strings.ToLower("nuvla") && n.IcosHostName == icos_host_name {
			return n.Name
		}
	}

	return ""
}

// check if 'k8s_cluster_uid' from 'node_uname_info' query correspond to an OCM node / cluster
func IsOCMCluster(cluster_ud string, orchs map[string]models.OrchInfoNode) bool {

	for _, n := range orchs {
		if cluster_ud == "" {
			logs.GetLogger().Warn("'k8s_cluster_uid' value is empty [cluster_ud:", cluster_ud, "]")
			return false
		}

		if n.Type == strings.ToLower("ocm") && n.ClusterId == cluster_ud {
			return true
		}
	}

	return false
}

// get Id from nuvla node using the icos_host_name value
func GetNuvlaNodeId(icos_host_name string, orchs map[string]models.OrchInfoNode) string {

	for _, n := range orchs {
		if n.IcosHostName == icos_host_name {
			return n.Id
		}
	}

	return "" // NOT FOUND / already deleted
}

// get engine from node
func GetEngine(cluster_id string, orchs map[string]models.OrchInfoNode) string {
	for _, n := range orchs {
		if n.ClusterId == cluster_id {
			return n.Engine
		}
	}

	return "unknown" // NOT FOUND / already deleted
}

// get ClusterLink
func GetClusterClusterLink(clusterLinks []models.ClusterLinks, icos_cluster_id string, icos_agent_id string) bool {
	if len(clusterLinks) == 0 {
		return false
	}

	for _, c := range clusterLinks {
		if c.ICOSAgentID == icos_agent_id && c.ICOSClusterID == icos_cluster_id {
			return true
		}
	}

	return false // NOT FOUND
}

// get name from cluster
func GetClusterName(icos_cluster_id string, orchs map[string]models.OrchInfoNode) string {
	for _, n := range orchs {
		if n.ClusterId == icos_cluster_id {
			return n.Name
		}
	}

	return "unknown" // NOT FOUND / already deleted
}

// get engine from nuvla node using the icos_host_name or cluster_id value
func GetNuvlaEngine(icos_host_name string, cluster_id string, orchs map[string]models.OrchInfoNode) string {
	for _, n := range orchs {
		if n.IcosHostName == icos_host_name || n.ClusterId == cluster_id {
			return n.Engine
		}
	}

	return "unknown" // NOT FOUND / already deleted
}

// Nodes, pods and containers info

// getNodeInfo
func GetNodeInfo(icos_host_name string, nodes map[string]models.Node) models.Node {
	for _, n := range nodes {

		if strings.EqualFold(n.IcosHostName, icos_host_name) {
			//logs.GetLogger().Debug("Node matches input values: pod_name [", pod_name, "] parentNodeUid [", parentNodeUid, "]. Returning node ...")
			return n
		}
	}

	logs.GetLogger().Warn("Node not found. Input values: 'icos_host_name' [", icos_host_name, "]")
	return models.Node{} // NOT FOUND
}

// Checks

// checkCluster checks if cluster is not "self" and if it exists in clusters list
func CheckCluster(clusterid string, clusters map[string]models.Cluster, metric string) bool {
	if clusterid == "self" {
		return false
	}

	if _, existsCluster := clusters[clusterid]; existsCluster {
		return true
	}

	logs.GetLogger().Warn("Unknown CLUSTER [", clusterid, "] {", metric, "}")
	return false
}

// checkClusterNode
func CheckClusterNode(clusterid string, nodeid string, clusters map[string]models.Cluster, metric string) bool {
	if CheckCluster(clusterid, clusters, metric) {
		if _, existsNode := clusters[clusterid].Node[nodeid]; existsNode {
			return true
		}
		logs.GetLogger().Warn("Unknown NODE [", nodeid, "] in CLUSTER [", clusterid, "] {", metric, "}")
	}

	return false
}

// checkClusterNodePod
func CheckClusterNodePod(clusterid string, nodeid string, podid string, clusters map[string]models.Cluster, metric string) bool {
	if CheckClusterNode(clusterid, nodeid, clusters, metric) {
		if _, existsPod := clusters[clusterid].Node[nodeid].Pod[podid]; existsPod {
			return true
		}
		logs.GetLogger().Warn("Unknown POD [", podid, "] in NODE [", nodeid, "] and CLUSTER [", clusterid, "] {", metric, "}")
	}

	return false
}

// checkClusterNodePodContainer
func CheckClusterNodePodContainer(clusterid string, nodeid string, podid string, containerid string, clusters map[string]models.Cluster, metric string) bool {

	if CheckClusterNodePod(clusterid, nodeid, podid, clusters, metric) {
		if _, existsContainer := clusters[clusterid].Node[nodeid].Pod[podid].Container[containerid]; existsContainer {
			return true
		}
		logs.GetLogger().Warn("Unknown CONTAINER [", containerid, "] in POD [", podid, "] and NODE [", nodeid, "] and CLUSTER [", clusterid, "] {", metric, "}")
	}

	return false
}
