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
	"strconv"
	"strings"
)

/*
*

  - nodes: Retrieves clusters and nodes
    *

  - metrics: - tlum_orch_info

  - - tlum_runtime_info

  - - tlum_host_info
    *

  - Infrastrcuture:

  - Timestamp:

  - ...

  - Agent:Cluster[]

  - ...

  - Node[]

  - ...
*/
func SetNodesV2(clusters map[string]models.Cluster, orchs map[string]models.OrchInfoNode) {

	// CLUSTERS
	// clusters: tlum_orch_info
	logs.GetLogger().Info("\t>> QUERY: tlum_orch_info")
	q := querier.PromQLQuery{
		Metric: "tlum_orch_info",
		Params: map[string]string{}}

	for _, node := range querier.Query(q.String()) {
		orch_id := string(node.Metric["agent_id"])

		var newCluster = models.OrchInfoNode{
			Id:           orch_id,
			Type:         string(node.Metric["type"]),
			Name:         string(node.Metric["agent_name"]), // Nuvla cluster
			Uuid:         string(node.Metric["agent_id"]),
			ClusterId:    string(node.Metric["icos_cluster_id"]), // OCM // icos_cluster_id instead of k8s_cluster_uid
			IcosHostName: string(node.Metric["icos_host_name"]),  // Nuvla
			Engine:       "unknown",                              // default, may be updated
		}
		orchs[orch_id] = newCluster
	}

	// CLUSTERS AND HOSTS
	// clusters: tlum_runtime_info
	logs.GetLogger().Info("\t>> QUERY: tlum_runtime_info")
	q = querier.PromQLQuery{
		Metric: "tlum_runtime_info",
		Params: map[string]string{}}

	hosts := querier.Query(q.String())

	for _, node := range hosts { //querier.Query(q.String()) {
		engine := string(node.Metric["type"])
		icos_cluster_id := string(node.Metric["icos_cluster_id"]) // icos_cluster_id instead of k8s_cluster_uid

		for _, o := range orchs {
			if o.ClusterId == icos_cluster_id {
				// First we get a "copy" of the entry
				if entry, ok := orchs[o.Id]; ok {
					// Then we modify the copy
					entry.Engine = engine

					// Then we reassign map entry
					orchs[o.Id] = entry
				}

				break
			}
		}

	}

	// Clusters and Nodes
	logs.GetLogger().Info("\t>> QUERY: tlum_host_info")
	q = querier.PromQLQuery{
		Metric: "tlum_host_info",
		Params: map[string]string{}}

	for _, node := range querier.Query(q.String()) {
		cluster_id := string(node.Metric["icos_cluster_id"])
		cluster_type := ""
		engine := "UNKNOWN"
		node_id := strings.TrimSpace(string(node.Metric["icos_host_id"]))
		node_name := string(node.Metric["hostname"]) // 'hostname' instead of 'nodename'
		net_host_name := string(node.Metric["net_host_name"])
		icos_host_name := string(node.Metric["icos_host_name"])
		icos_agent_id := string(node.Metric["icos_agent_id"])
		architecture := string(node.Metric["machine"])
		latitude, _ := strconv.ParseFloat(string(node.Metric["latitude"]), 64)
		longitude, _ := strconv.ParseFloat(string(node.Metric["longitude"]), 64)

		if cluster_id != "self" {
			if common.IsNuvlaCluster(icos_host_name, orchs) {
				cluster_id = common.GetNuvlaClusterName(icos_host_name, orchs)
				cluster_type = "nuvla"
				engine = common.GetNuvlaEngine(cluster_id, cluster_id, orchs)
			} else if common.IsOCMCluster(cluster_id, orchs) {
				cluster_type = "ocm"
				engine = common.GetEngine(cluster_id, orchs)
			}

			if cluster_type != "" {
				if _, exists := clusters[cluster_id]; !exists {
					cluster_name := common.GetClusterName(cluster_id, orchs)
					var newCluster = models.Cluster{
						Uuid:        cluster_id,
						Type:        cluster_type,
						Engine:      engine,
						ICOSAgentID: icos_agent_id,
						Node:        map[string]models.Node{},
						Name:        cluster_name,
					}
					clusters[cluster_id] = newCluster
				}

				newNode := models.Node{
					Uuid:         node_id,
					Type:         cluster_type,
					Name:         node_name,
					Engine:       engine,
					NetHostName:  net_host_name,
					IcosHostName: icos_host_name,
					Location: models.Location{
						Latitude:  latitude,
						Longitude: longitude,
					},
					StaticMetrics:     models.StaticMetrics{CPUArchitecture: architecture},
					NetworkInterfaces: map[string]models.Interface{},
					Devices:           map[string]models.Device{},
					Pod:               map[string]models.Pod{},
				}
				clusters[cluster_id].Node[node_id] = newNode
			} else {
				logs.GetLogger().Warn("Unknown cluster type [cluster_id:", cluster_id, "]. Cluster not added to infraestructure.")
			}

		}
	}

	// Rename hosts with the name from metric: tlum_runtime_info and label: node_name
	// nodes: tlum_runtime_info
	logs.GetLogger().Info("\t>> QUERY: tlum_runtime_info")
	q = querier.PromQLQuery{
		Metric: "tlum_runtime_info",
		Params: map[string]string{}}

	for _, node := range querier.Query(q.String()) {
		node_name := string(node.Metric["node_name"])
		node_id := strings.TrimSpace(string(node.Metric["icos_host_id"]))

		for _, c := range clusters {
			for _, n := range c.Node {
				if n.Uuid == node_id {
					n.Name = node_name
					clusters[c.Uuid].Node[n.Uuid] = n
				}
			}
		}
	}

	logs.GetLogger().Debug("\t>> Clusters (Uuid) and nodes (Uuid, IcosHostName) added to infra: ")
	for _, c := range clusters {
		logs.GetLogger().Debug("\t     * " + c.Uuid + " (" + c.Type + ", " + c.Engine + ")")
		for _, n := range clusters[c.Uuid].Node {
			logs.GetLogger().Debug("\t       - " + n.Uuid + " (" + n.IcosHostName + ")")
		}
	}

}
