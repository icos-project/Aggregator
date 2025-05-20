# Relation between queries and infra structs

## Infra

```yaml
Infrastrcuture
    Timestamp
    Cluster[]
        Node[]
            StaticMetrics
            DinamycMetrics
            NetworkInterfaces[]
            Device[]
            Pod[]
                *ParentNodeId* // Uuid from parent Node
                Container[]
```

## Queries

### Clusters: tlum_orch_info, tlum_runtime_info

info about **clusters**

- _tlum_orch_info_: ocm, nuvla
- _tlum_runtime_info_: engine

----------------------------------------------------------

### Clusters, Nodes: node_uname_info

Info about clusters **nodes**. Clusters and Nodes lists / maps are created after this query

- Cluster[ID] => **node_uname_info**.k8s_cluster_uid

A new Cluster struct is created:

```golang
cluster_id := string(node.Metric["k8s_cluster_uid"])

var newCluster = Cluster{
    Uuid:        cluster_id,
    Type:        cluster_type,
    Engine:      engine,
    ICOSAgentID: icos_agent_id,
    ...
}
clusters[cluster_id] = newCluster
```

- Node[ID] => **node_uname_info**.icos_host_id

A new Node struct is created:

```golang
node_id := string(node.Metric["icos_host_id"])
...
newNode := Node{
    Uuid:         node_id,
    Type:         cluster_type,
    Name:         node_name,
    Engine:       engine,
    NetHostName:  net_host_name,
    IcosHostName: icos_host_name, // matches value of 'icos_host_name' in 'kube_pod_info' query
    K8sNodeUid:   k8s_node_uid,   // matches value of 'k8s_node_uid' in 'kube_pod_info' query
    ...
}
clusters[cluster_id].Node[node_id] = newNode
```

Node UID value (**node_uname_info**.k8s_node_uid) used in other metrics to associate elements (Nodes, Pods, Containers) is stored in **Node[].K8sNodeUid**.

----------------------------------------------------------

### Pods: kube_pod_info, kube_pod_status_phase

#### kube_pod_info

Info about nodes **pods**. Pods list / map is created after this query


- **node_uname_info**.k8s_node_uid == **kube_pod_info**.k8s_node_uid => Pod.ParentNodeId
- Pod.Uid => **kube_pod_info**.uid
- Pod[ID] => **kube_pod_info**.pod

A new Pod struct is created:

```golang
pod_name := string(pod.Metric["pod"])
...
parentNodeId := n.Uuid

var podInfo = Pod{
    Uid:                pod_uid,        // uid from pod. Used to match pods in next queries
    ParentNodeName:     parentNodeName, // hidden value: parent Nodename
    ParentNodeId:       parentNodeId,   // hidden value: edge_harbor_host_id / Uuid from parent Node
    ClusterUid:         cluster_uid,    // hidden value: cluster uid
    IcosHostName:       host_name,      // hidden value
    K8sNodeUid:         k8s_node_uid,   // hidden value
    Name:               pod_name,
    ...
}

...
clusters[cluster_uid].Node[parentNodeId].Pod[pod_name] = podInfo
```


#### kube_pod_status_phase

- **kube_pod_info**.uid == **kube_pod_status_phase**.uid

----------------------------------------------------------

### Containers: kube_pod_container_info, container_cpu_utilization_ratio

#### kube_pod_container_info

Containers list / map is created after this query

info about nodes **containers**

- Connection to Pod: **kube_pod_info**.uid == **kube_pod_container_info**.uid
- Container[name] => **kube_pod_container_info**.container

#### container_cpu_utilization_ratio

- Connection to Pod: **kube_pod_info**.uid == **container_cpu_utilization_ratio**.k8s_pod_uid
- **kube_pod_container_info**.uid == **container_cpu_utilization_ratio**.k8s_pod_uid