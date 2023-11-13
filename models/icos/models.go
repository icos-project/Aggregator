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
	Node                  map[string]Node       `json:"node,omitempty"`
	Pod                   map[string]Pod        `json:"pod,omitempty"`
	Any                   any                   `json:"any,omitempty"`
}

type Pod struct {
	Name               string               `json:"name,omitempty"`
	Status             string               `json:"status,omitempty"`
	NumberOfContainers int32                `json:"numberOfContainers,omitempty"`
	NumberOfApps       int32                `json:"numberOfApps,omitempty"`
	Container          map[string]Container `json:"container,omitempty"`
}

type Container struct {
	Name            string  `json:"name,omitempty"`
	Node            string  `json:"node,omitempty"`
	Port            []Port  `json:"port,omitempty"`
	ContainerMemory string  `json:"containerMemory,omitempty"`
	CPUUsage        float64 `json:"cpuUsage,omitempty"`
	IP              string  `json:"ip,omitempty"`
}

type Port struct {
	Port string `json:"port,omitempty"`
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
