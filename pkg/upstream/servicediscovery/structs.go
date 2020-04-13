package servicediscovery

type subReq struct {
	Interface    string `json:"interface" binding:"required"` // eg. com.mosn.service
	Method       string `json:"method" binding:"required"` // eg. GetUser
	RegistryAddr string `json:"registry_addr" binding:"required"`  // eg. 127.0.0.1:2181
	RegistryType string `json:"registry_type" binding:"eq=zookeeper"` // zookeeper，etcd，k8s..., currently only zookeeper is supported
}

type unsubReq struct {
	Interface    string `json:"interface" binding:"required"` // eg. com.mosn.service
	Method       string `json:"method" binding:"required"` // eg. GetUser
	RegistryAddr string `json:"registry_addr" binding:"required"`  // eg. 127.0.0.1:2181
	RegistryType string `json:"registry_type" binding:"eq=zookeeper"` // zookeeper，etcd，k8s..., currently only zookeeper is supported
}

type pubReq struct {
	Interface    string `json:"interface" binding:"required"` // eg. com.mosn.service
	Method       string `json:"method" binding:"required"` // eg. GetUser
	RegistryAddr string `json:"registry_addr" binding:"required"`  // eg. 127.0.0.1:2181
	RegistryType string `json:"registry_type" binding:"eq=zookeeper"` // zookeeper，etcd，k8s..., currently only zookeeper is supported
}

type unpubReq struct {
	Interface    string `json:"interface" binding:"required"` // eg. com.mosn.service
	Method       string `json:"method" binding:"required"` // eg. GetUser
	RegistryAddr string `json:"registry_addr" binding:"required"`  // eg. 127.0.0.1:2181
	RegistryType string `json:"registry_type" binding:"eq:zookeeper"` // zookeeper，etcd，k8s..., currently only zookeeper is supported
}
