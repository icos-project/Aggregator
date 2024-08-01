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
	"fmt"
	"math"
	"strings"
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

func checkCluster(cluster string, clusters map[string]Cluster, metric string) bool {
	_, existsCluster := clusters[cluster]

	if !existsCluster {
		fmt.Println("Unknown cluster ", cluster, ". Error in metric: ", metric)
	}

	return existsCluster
}

func checkClusterNode(cluster string, node string, clusters map[string]Cluster, metric string) bool {
	_, existsCluster := clusters[cluster]
	_, existsNode := clusters[cluster].Node[node]

	if !existsCluster {
		fmt.Println("Unknown cluster ", cluster, ". Error in metric: ", metric)
	} else if !existsNode {
		fmt.Println("Unknown node ", node, " in cluster ", cluster, ". Error in metric: ", metric)
	}

	return existsCluster && existsNode
}

func checkClusterPod(cluster string, pod string, clusters map[string]Cluster, metric string) bool {
	_, existsCluster := clusters[cluster]
	_, existsPod := clusters[cluster].Pod[pod]

	if !existsCluster {
		fmt.Println("Unknown cluster ", cluster, ". Error in metric: ", metric)
	} else if !existsPod {
		fmt.Println("Unknown pod ", pod, " in cluster ", cluster, ". Error in metric: ", metric)
	}

	return existsCluster && existsPod
}

func checkClusterPodContainer(cluster string, pod string, container string, clusters map[string]Cluster, metric string) bool {
	_, existsCluster := clusters[cluster]
	_, existsPod := clusters[cluster].Pod[pod]
	_, existsContainer := clusters[cluster].Pod[pod].Container[container]

	if !existsCluster {
		fmt.Println("Unknown cluster ", cluster, ". Error in metric: ", metric)
	} else if !existsPod {
		fmt.Println("Unknown pod ", pod, " in cluster ", cluster, ". Error in metric: ", metric)
	} else if !existsContainer {
		fmt.Println("Unknown Container ", container, " in pod ", pod, " in cluster ", cluster, ". Error in metric: ", metric)
	}

	return existsCluster && existsPod && existsContainer
}

// NUVLA and OCM clusters / nodes

// check if 'cluster_id' and 'icos_host_name' from 'node_uname_info' query correspond to a Nuvla node / cluster
func isNuvlaCluster(cluster_id string, icos_host_name string, orchs map[string]OrchInfoNode) bool {

	// Example:
	// - cluster_id=""
	// - icos_host_name="icos-uc2-test-001"
	// - OrchInfoNode.icos_host_name="icos-uc2-test-001"
	// ==>  OrchInfoNode[i].IcosHostName == icos_host_name ==> nuvla
	for _, n := range orchs {
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
