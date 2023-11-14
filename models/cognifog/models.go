package models_cognifog

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
	Resources           ComputeResources `json:"resources,omitempty"`
	Available_resources ComputeResources `json:"available_resources,omitempty"`
	Labels              []string         `json:"labels,omitempty"`
}

type ComputeResources struct {
	MilliCPU      int32 `json:"milliCPU,omitempty"`
	MemoryInBytes int64 `json:"memoryInBytes,omitempty"`
}
