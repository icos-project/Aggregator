package models

type Controllers []*Controller

type Controller struct {
	Type                  string                `json:"type,omitempty"`
	Name                  string                `json:"name,omitempty"`
	Location              Location              `json:"location,omitempty"`
	ServiceLevelAgreement ServiceLevelAgreement `json:"serviceLevelAgreement,omitempty"`
	API                   API                   `json:"API,omitempty"`
	Any                   any                   `json:"any,omitempty"`
}

type Agents []*Cluster

type Cluster struct {
	Type                  string                `json:"type,omitempty"`
	Name                  string                `json:"name,omitempty"`
	Location              Location              `json:"location,omitempty"`
	ServiceLevelAgreement ServiceLevelAgreement `json:"serviceLevelAgreement,omitempty"`
	API                   API                   `json:"API,omitempty"`
	Node                  Node                  `json:"node,omitempty"`
	Any                   any                   `json:"any,omitempty"`
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
	Type           string         `json:"type,omitempty"`
	Name           string         `json:"name,omitempty"`
	StaticMetrics  StaticMetrics  `json:"staticMetrics,omitempty"`
	DynamicMetrics DynamicMetrics `json:"dynamicMetrics,omitempty"`
}

type StaticMetrics struct { //TODO: separate in []CPU?, []GPU
	CPUCores        float64   `json:"cpuCores,omitempty"`
	CPUMaxFrecuency string    `json:"cpuMaxFrecuency,omitempty"`
	GPUCores        float64   `json:"gpuCores,omitempty"`
	GPUMaxFrecuency string    `json:"gpuMaxFrecuency,omitempty"`
	GPURAMMemory    string    `json:"gpuRAMMemory,omitempty"`
	RAMMemory       string    `json:"RAMMemory,omitempty"`
	Storage         []Storage `json:"storage,omitempty"`
}

type DynamicMetrics struct { //TODO: separate in []CPU?, []GPU
	CPUFrecuency         string         `json:"cpuFrecuency,omitempty"`
	CPUTemperature       float64        `json:"cpuTemperature,omitempty"`
	CPUEnergyConsumption float64        `json:"cpuEnergyConsumption,omitempty"`
	GPUFrecuency         string         `json:"gpuFrecuency,omitempty"`
	GPUTemperature       float64        `json:"gpuTemperature,omitempty"`
	GPUEnergyConsumption float64        `json:"gpuEnergyConsumption,omitempty"`
	RAMUsage             string         `json:"ramUsage,omitempty"`
	NetworkUsage         []NetworkUsage `json:"networkUsage,omitempty"`
}

type Storage struct { //TODO: complete
}

type NetworkUsage struct { //TODO: complete

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
