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

type Controllers []*Controller

type Controller struct {
	Type                  string                `json:"type,omitempty"`
	Name                  string                `json:"name,omitempty"`
	Location              Location              `json:"location,omitempty"`
	ServiceLevelAgreement ServiceLevelAgreement `json:"serviceLevelAgreement,omitempty"`
	API                   API                   `json:"API,omitempty"`
	Any                   any                   `json:"any,omitempty"`
}

type Infrastructure struct {
	Timestamp  Timestamp             `json:"timestamp,omitempty"`
	Controller map[string]Controller `json:"controller,omitempty"`
	Cluster    map[string]Cluster    `json:"cluster,omitempty"`
}

type Timestamp struct {
	OldestTimestamp float64 `json:"oldestTimestamp,omitempty"`
	TimeSinceOldest float64 `json:"timeSinceOldest,omitempty"`
}

type Agents []*Cluster

type Cluster struct {
	Type                  string                `json:"type,omitempty"`
	Uuid                  string                `json:"uuid,omitempty"`
	Name                  string                `json:"name,omitempty"`
	ServiceLevelAgreement ServiceLevelAgreement `json:"serviceLevelAgreement,omitempty"`
	API                   API                   `json:"API,omitempty"`
	Node                  map[string]Node       `json:"node,omitempty"`
	Pod                   map[string]Pod        `json:"pod,omitempty"`
	Any                   any                   `json:"any,omitempty"`
}

type Pod struct {
	Name               string               `json:"name,omitempty"`
	IP                 string               `json:"ip,omitempty"`
	Status             string               `json:"status,omitempty"`
	NumberOfContainers int32                `json:"numberOfContainers,omitempty"`
	NumberOfApps       int32                `json:"numberOfApps,omitempty"`
	Container          map[string]Container `json:"container,omitempty"`
}

type Container struct {
	Name            string  `json:"name,omitempty"`
	IP              string  `json:"ip,omitempty"`
	Node            string  `json:"node,omitempty"`
	Port            string  `json:"port,omitempty"`
	ContainerMemory string  `json:"containerMemory,omitempty"`
	CPUUsage        float64 `json:"cpuUsage,omitempty"`
}

type Location struct {
	Name      string  `json:"name,omitempty"`
	Continent string  `json:"continent,omitempty"`
	Country   string  `json:"country,omitempty"`
	City      string  `json:"city,omitempty"`
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
}

type Node struct {
	Type              string               `json:"type,omitempty"`
	Uuid              string               `json:"uuid,omitempty"`
	Name              string               `json:"name,omitempty"`
	NetHostName       string               `json:"-"` //`json:"netHostName,omitempty"`
	Location          Location             `json:"location,omitempty"`
	Vulnerabilities   map[string]int32     `json:"vulnerabilities,omitempty"`
	ScaScore          int32                `json:"ScaScore,omitempty"`
	StaticMetrics     StaticMetrics        `json:"staticMetrics,omitempty"`
	DynamicMetrics    DynamicMetrics       `json:"dynamicMetrics,omitempty"`
	NetworkInterfaces map[string]Interface `json:"networkInterfaces,omitempty"`
	Devices           map[string]Device    `json:"devices,omitempty"`
}

type StaticMetrics struct {
	CPUArchitecture string    `json:"cpuArchitecture,omitempty"`
	CPUCores        int32     `json:"cpuCores,omitempty"`
	CPUMaxFrequency int64     `json:"cpuMaxFrequency,omitempty"`
	GPUCores        float64   `json:"gpuCores,omitempty"`
	GPUMaxFrequency string    `json:"gpuMaxFrequency,omitempty"`
	GPURAMMemory    string    `json:"gpuRAMMemory,omitempty"`
	RAMMemory       int64     `json:"RAMMemory,omitempty"`
	Storage         []Storage `json:"storage,omitempty"`
}

type DynamicMetrics struct {
	UpTime               float64          `json:"upTime,omitempty"`
	CPUFrequency         string           `json:"cpuFrequency,omitempty"`
	CPUTemperature       float64          `json:"cpuTemperature,omitempty"`
	CPUEnergyConsumption float64          `json:"cpuEnergyConsumption,omitempty"`
	GPUFrequency         string           `json:"gpuFrequency,omitempty"`
	GPUTemperature       float64          `json:"gpuTemperature,omitempty"`
	GPUEnergyConsumption float64          `json:"gpuEnergyConsumption,omitempty"`
	FreeRAM              int64            `json:"freeRAM,omitempty"`
	Storage              AvailableStorage `json:"availableStorage,omitempty"`
	//NetworkUsage         NetworkUsage     `json:"networkUsage,omitempty"`
}

type Device struct {
	Name   string `json:"name,omitempty"`
	Type   string `json:"type,omitempty"`
	Status string `json:"status,omitempty"`
	Path   string `json:"path,omitempty"`
}

type Storage struct {
	Name     string  `json:"name,omitempty"`
	Type     string  `json:"type,omitempty"`
	Capacity float64 `json:"capacity,omitempty"`
}

type AvailableStorage struct {
	Name string  `json:"name,omitempty"`
	Free float64 `json:"free,omitempty"`
}

type Interface struct {
	Name          string `json:"name,omitempty"`
	Type          string `json:"type,omitempty"`
	Speed         int64  `json:"speed,omitempty"`
	IP            string `json:"ip,omitempty"`
	Status        string `json:"status,omitempty"`
	SubnetMask    string `json:"subnetMask,omitempty"`
	IngressUssage string `json:"ingressUssage,omitempty"`
	EngressUssage string `json:"engressUssage,omitempty"`
}

type API struct {
	CommunicationProtocol string `json:"commProtocol,omitempty"`
	ProtocolVersion       string `json:"protocolVersion,omitempty"`
	DataFormat            string `json:"dataFormat,omitempty"`
	Authentication        string `json:"authentication,omitempty"`
	Authorization         string `json:"authorization,omitempty"`
}

type ServiceLevelAgreement struct { //TODO: complete
	Name string `json:"name,omitempty"`
}

type OrchInfoNode struct {
	Id            string `json:"id,omitempty"`
	Type          string `json:"type,omitempty"` // "ocm" or "nuvla"
	Name          string `json:"agent_name,omitempty"`
	Uuid          string `json:"agent_id,omitempty"`
	K8sClusterUid string `json:"k8s_cluster_uid,omitempty"` // OCM
	IcosHostName  string `json:"icos_host_name,omitempty"`  // Nuvla
	K8sNodeName   string `json:"k8s_node_name,omitempty"`
}
