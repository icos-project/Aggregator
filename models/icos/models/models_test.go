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
package models

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestController(t *testing.T) {

	want := []byte(`{"type":"MetaOrchestrator","name":"ICOS1","location":{"name":"BCN"},"serviceLevelAgreement":{},"API":{}}`)

	c := Controller{Type: "MetaOrchestrator",
		Name:     "ICOS1",
		Location: Location{Name: "BCN"},
	}

	got, err := json.Marshal(c)

	if err != nil {
		t.Error(err)
	} else {
		t.Log(string(got))
		t.Log(string(want))
	}

	if bytes.Compare(want, got) != 0 {
		t.Errorf("Controller model error")
	}

}

func TestCluster(t *testing.T) {
	want := []byte(`{"type":"Kubernetes","uuid":"123","name":"TestCluster","engine":"k8s","icosAgentID":"agent1","clusterLink":true,"serviceLevelAgreement":{},"API":{}}`)

	c := Cluster{
		Type:        "Kubernetes",
		Uuid:        "123",
		Name:        "TestCluster",
		Engine:      "k8s",
		ICOSAgentID: "agent1",
		ClusterLink: true,
	}

	got, err := json.Marshal(c)
	if err != nil {
		t.Error(err)
	}

	if bytes.Compare(want, got) != 0 {
		t.Errorf("Cluster model error\nwant: %s\ngot: %s", string(want), string(got))
	}
}

func TestNode(t *testing.T) {
	want := []byte(`{"type":"Worker","uuid":"node1","name":"worker1","engine":"k8s","location":{"name":"BCN"},"vulnerabilities":{"critical":2},"ScaScore":85,"staticMetrics":{"cpuCores":4,"RAMMemory":8589934592},"dynamicMetrics":{"upTime":3600,"freeRAM":4294967296,"availableStorage":{}}}`)

	n := Node{
		Type:   "Worker",
		Uuid:   "node1",
		Name:   "worker1",
		Engine: "k8s",
		Location: Location{
			Name: "BCN",
		},
		Vulnerabilities: map[string]int32{
			"critical": 2,
		},
		ScaScore: 85,
		StaticMetrics: StaticMetrics{
			CPUCores:  4,
			RAMMemory: 8 * 1024 * 1024 * 1024, // 8GB
		},
		DynamicMetrics: DynamicMetrics{
			UpTime:  3600,
			FreeRAM: 4 * 1024 * 1024 * 1024, // 4GB
		},
	}

	got, err := json.Marshal(n)
	if err != nil {
		t.Error(err)
	}

	if bytes.Compare(want, got) != 0 {
		t.Errorf("Node model error\nwant: %s\ngot: %s", string(want), string(got))
	}
}

func TestPod(t *testing.T) {
	want := []byte(`{"name":"test-pod","ip":"10.0.0.1","status":"Running","numberOfContainers":2,"numberOfApps":1,"container":{"container1":{"name":"container1","ip":"10.0.0.2","node":"node1","port":"8080","containerMemory":"256Mi","cpuUsage":0.5}}}`)

	p := Pod{
		Name:               "test-pod",
		IP:                 "10.0.0.1",
		Status:             "Running",
		NumberOfContainers: 2,
		NumberOfApps:       1,
		Container: map[string]Container{
			"container1": {
				Name:            "container1",
				IP:              "10.0.0.2",
				Node:            "node1",
				Port:            "8080",
				ContainerMemory: "256Mi",
				CPUUsage:        0.5,
			},
		},
	}

	got, err := json.Marshal(p)
	if err != nil {
		t.Error(err)
	}

	if bytes.Compare(want, got) != 0 {
		t.Errorf("Pod model error\nwant: %s\ngot: %s", string(want), string(got))
	}
}

func TestLocation(t *testing.T) {
	want := []byte(`{"name":"Barcelona","continent":"Europe","country":"Spain","city":"Barcelona","latitude":41.3851,"longitude":2.1734}`)

	l := Location{
		Name:      "Barcelona",
		Continent: "Europe",
		Country:   "Spain",
		City:      "Barcelona",
		Latitude:  41.3851,
		Longitude: 2.1734,
	}

	got, err := json.Marshal(l)
	if err != nil {
		t.Error(err)
	}

	if bytes.Compare(want, got) != 0 {
		t.Errorf("Location model error\nwant: %s\ngot: %s", string(want), string(got))
	}
}

func TestInfrastructure(t *testing.T) {
	want := []byte(`{"timestamp":{"oldestTimestamp":1617234567.89,"timeSinceOldest":3600},"controller":{"ctrl1":{"type":"MetaOrchestrator","name":"ICOS1","location":{},"serviceLevelAgreement":{},"API":{}}},"cluster":{"cluster1":{"type":"Kubernetes","name":"TestCluster","clusterLink":false,"serviceLevelAgreement":{},"API":{}}}}`)

	i := Infrastructure{
		Timestamp: Timestamp{
			OldestTimestamp: 1617234567.89,
			TimeSinceOldest: 3600,
		},
		Controller: map[string]Controller{
			"ctrl1": {
				Type: "MetaOrchestrator",
				Name: "ICOS1",
			},
		},
		Cluster: map[string]Cluster{
			"cluster1": {
				Type: "Kubernetes",
				Name: "TestCluster",
			},
		},
	}

	got, err := json.Marshal(i)
	if err != nil {
		t.Error(err)
	}

	if bytes.Compare(want, got) != 0 {
		t.Errorf("Infrastructure model error\nwant: %s\ngot: %s", string(want), string(got))
	}
}
