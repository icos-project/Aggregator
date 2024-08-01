/*
Copyright 2023 Bull SAS

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
	"aggregator/querier"
	"strings"
)

/**
 * nodes: Retrieves nterwork interfaces from nodes
 *        It uses OrchInfoNode map struct to identify OCM and Nuvla clusters and nodes.
 * metrics: - node_network_info
 *			- node_network_up
 *			- node_network_address_info
 *			- node_network_speed_bytes
 *
 * 	Infrastrcuture:
 *		Timestamp:
 *		 	...
 *		Agent:Cluster[]
 *			...
 *			Node[]
 *          	...
 *				StaticMetrics:
 *                  ...
 * 				DinamycMetrics:
 *					...
 *             	NetworkInterfaces[]
 *					Interface_name: #DONE
 *					Interface_type: #DONE
 *					Interface_speed: #DONE
 *					Interface_IP: #DONE
 *					Interface_status: #DONE
 *					Interface_subnet_mask: #DONE
 *					Interface_ingress_usage:
 *					Interface_egress_usage:
 *					...
 *				Device[]
 *					...
 */
func processNetworkInterfaces(clusters map[string]Cluster, orchs map[string]OrchInfoNode) {

	q := querier.PromQLQuery{
		Metric: "node_network_info",
		Params: map[string]string{}}

	for _, node := range querier.Query(q.String()) {
		cluster_id := string(node.Metric["k8s_cluster_uid"])
		node_id := string(node.Metric["icos_host_id"])
		node_name := string(node.Metric["icos_host_name"])
		i_name := string(node.Metric["device"])
		i_type := string(node.Metric["duplex"])

		if cluster_id == "" && strings.HasPrefix(node_name, "nuvla") {
			cluster_id = "nuvla"
		}

		if checkClusterNode(cluster_id, node_id, clusters, q.Metric) && cluster_id != "self" {
			newInterface := Interface{
				Name: i_name,
				Type: i_type,
			}

			clusters[cluster_id].Node[node_id].NetworkInterfaces[i_name] = newInterface
		}
	}

	q = querier.PromQLQuery{
		Metric: "node_network_up",
		Params: map[string]string{}}

	for _, node := range querier.Query(q.String()) {
		cluster_id := string(node.Metric["k8s_cluster_uid"])
		node_id := string(node.Metric["icos_host_id"])
		node_name := string(node.Metric["icos_host_name"])
		i_name := string(node.Metric["device"])
		status := int(node.Value)

		status_str := "unknown"
		if status == 1 {
			status_str = "up"
		} else {
			status_str = "down"
		}

		if cluster_id == "" && strings.HasPrefix(node_name, "nuvla") {
			cluster_id = "nuvla"
		}

		if checkClusterNode(cluster_id, node_id, clusters, q.Metric) && cluster_id != "self" {
			iface := clusters[cluster_id].Node[node_id].NetworkInterfaces[i_name]
			iface.Status = status_str
			clusters[cluster_id].Node[node_id].NetworkInterfaces[i_name] = iface
		}
	}

	q = querier.PromQLQuery{
		Metric: "node_network_address_info",
		Params: map[string]string{}}

	for _, node := range querier.Query(q.String()) {
		cluster_id := string(node.Metric["k8s_cluster_uid"])
		node_id := string(node.Metric["icos_host_id"])
		node_name := string(node.Metric["icos_host_name"])
		i_name := string(node.Metric["device"])
		i_address := string(node.Metric["address"])
		i_netmask := string(node.Metric["netmask"])

		if cluster_id == "" && strings.HasPrefix(node_name, "nuvla") {
			cluster_id = "nuvla"
		}

		if checkClusterNode(cluster_id, node_id, clusters, q.Metric) && cluster_id != "self" {
			iface := clusters[cluster_id].Node[node_id].NetworkInterfaces[i_name]
			iface.IP = i_address
			iface.SubnetMask = i_netmask
			clusters[cluster_id].Node[node_id].NetworkInterfaces[i_name] = iface
		}
	}

	q = querier.PromQLQuery{
		Metric: "node_network_speed_bytes",
		Params: map[string]string{}}

	for _, node := range querier.Query(q.String()) {
		cluster_id := string(node.Metric["k8s_cluster_uid"])
		node_id := string(node.Metric["icos_host_id"])
		node_name := string(node.Metric["icos_host_name"])
		i_name := string(node.Metric["device"])
		speed := int64(node.Value)

		if cluster_id == "" && strings.HasPrefix(node_name, "nuvla") {
			cluster_id = "nuvla"
		}

		if checkClusterNode(cluster_id, node_id, clusters, q.Metric) && cluster_id != "self" {
			iface := clusters[cluster_id].Node[node_id].NetworkInterfaces[i_name]
			iface.Speed = speed
			clusters[cluster_id].Node[node_id].NetworkInterfaces[i_name] = iface
		}
	}
}
