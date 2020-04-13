package servicediscovery

type subReq struct {
	Services []string `json:"service_names"`
}

type unsubReq struct {
	Interface    string `json:"interface" binding:"required"`
	Method       string `json:"method" binding:"required"`
	RegistryAddr string `json:"registry_addr" binding:"required"`
	RegistryType string `json:"registry_type" binding:"eq=zookeeper"` // zookeeper，etcd，k8s..., currently only zookeeper is supported
}

type pubReq struct {
	Interface    string `json:"interface" binding:"required"`
	Method       string `json:"method" binding:"required"`
	RegistryAddr string `json:"registry_addr" binding:"required"`
	RegistryType string `json:"registry_type" binding:"eq=zookeeper"` // zookeeper，etcd，k8s..., currently only zookeeper is supported
}

type unpubReq struct {
	Interface    string `json:"interface"`
	Method       string `json:"methods"`
	RegistryAddr string `json:"registry_addr"`
	RegistryType string `json:"registry_type" binding:"eq:zookeeper"` // zookeeper，etcd，k8s..., currently only zookeeper is supported
}
