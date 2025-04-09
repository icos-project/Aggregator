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
	"aggregator/models/icos/models"
	"encoding/json"
	"os"
	"testing"
)

func TestQueryPrometheus(t *testing.T) {

	os.Setenv("PROMETHEUS_ADDRESS", "http://query.192.168.137.175.nip.io/") // thanos-query

	infra := queryPrometheus()

	// Check if we have clusters
	if len(infra.Cluster) == 0 {
		t.Error("Expected at least one cluster, got none")
	}

	// Check cluster properties
	for id, cluster := range infra.Cluster {
		if id == "" {
			t.Error("Cluster ID should not be empty")
		}
		if cluster.Type == "" {
			t.Errorf("Cluster %s: Type should not be empty", id)
		}
		if cluster.Uuid == "" {
			t.Errorf("Cluster %s: UUID should not be empty", id)
		}
		if cluster.Name == "" {
			t.Errorf("Cluster %s: Name should not be empty", id)
		}
		if cluster.Engine == "" {
			t.Errorf("Cluster %s: Engine should not be empty", id)
		}

		// Check if cluster has nodes
		if len(cluster.Node) == 0 {
			t.Errorf("Cluster %s: Expected at least one node, got none", id)
		}

		// Check node properties
		for nodeId, node := range cluster.Node {
			if nodeId == "" {
				t.Errorf("Cluster %s: Node ID should not be empty", id)
			}
			checkNodeProperties(t, id, nodeId, node)
		}
	}

	// Check timestamp
	if infra.Timestamp.OldestTimestamp <= 0 {
		t.Error("Expected OldestTimestamp to be positive")
	}
	if infra.Timestamp.TimeSinceOldest <= 0 {
		t.Error("Expected TimeSinceOldest to be positive")
	}
}

func checkNodeProperties(t *testing.T, clusterId, nodeId string, node models.Node) {
	if node.Uuid == "" {
		t.Errorf("Cluster %s, Node %s: UUID should not be empty", clusterId, nodeId)
	}
	if node.Name == "" {
		t.Errorf("Cluster %s, Node %s: Name should not be empty", clusterId, nodeId)
	}

	// Check static metrics
	if node.StaticMetrics.CPUCores <= 0 {
		t.Errorf("Cluster %s, Node %s: CPUCores should be positive", clusterId, nodeId)
	}
	if node.StaticMetrics.RAMMemory <= 0 {
		t.Errorf("Cluster %s, Node %s: RAMMemory should be positive", clusterId, nodeId)
	}

	// Check dynamic metrics exist
	if node.DynamicMetrics == (models.DynamicMetrics{}) {
		t.Errorf("Cluster %s, Node %s: DynamicMetrics should not be empty", clusterId, nodeId)
	}

	// Check vulnerabilities map exists
	if node.Vulnerabilities == nil {
		t.Errorf("Cluster %s, Node %s: Vulnerabilities should not be nil", clusterId, nodeId)
	}

	// Check devices map exists
	if node.Devices == nil {
		t.Errorf("Cluster %s, Node %s: Devices should not be nil", clusterId, nodeId)
	}

	// Check pods map exists
	if node.Pod == nil {
		t.Errorf("Cluster %s, Node %s: Pod map should not be nil", clusterId, nodeId)
	}
}

func TestGetInfra(t *testing.T) {
	jsonBytes := GetInfra()

	// Check if the returned value is not empty
	if len(jsonBytes) == 0 {
		t.Error("GetInfra returned empty JSON")
	}

	// Verify it's valid JSON by attempting to unmarshal it
	var infra models.Infrastructure
	err := json.Unmarshal(jsonBytes, &infra)
	if err != nil {
		t.Errorf("GetInfra returned invalid JSON: %v", err)
	}

	// Check if the unmarshaled data contains the expected structure
	if len(infra.Cluster) == 0 {
		t.Error("Unmarshaled JSON contains no clusters")
	}
}
