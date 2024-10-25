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
package models_etim

type NodeType int

const (
	COMPUTE NodeType = iota
	NETWORK
)

type Cluster struct {
	Nodes map[string]Node `json:"nodes,omitempty"`
}

type Node struct {
	Id                  string           `json:"id,omitempty"`
	Node_type           NodeType         `json:"node_type,omitempty"`
	CPUArchitecture     string           `json:"cpuArchitecture,omitempty"`
	Resources           ComputeResources `json:"resources,omitempty"`
	Available_resources ComputeResources `json:"available_resources,omitempty"`
	Labels              []string         `json:"labels,omitempty"`
}

type ComputeResources struct {
	MilliCPU      int32 `json:"milliCPU,omitempty"`
	MemoryInBytes int64 `json:"memoryInBytes,omitempty"`
}
