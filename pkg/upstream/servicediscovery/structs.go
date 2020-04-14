package servicediscovery

type subReq struct {
	Interface    string `json:"interface" binding:"required"`         // eg. com.mosn.service
	Method       string `json:"method" binding:"required"`            // eg. GetUser
	RegistryAddr string `json:"registry_addr" binding:"required"`     // eg. 127.0.0.1:2181
	RegistryType string `json:"registry_type" binding:"eq=zookeeper"` // zookeeper，etcd，k8s..., currently only zookeeper is supported
}

type unsubReq struct {
	Interface    string `json:"interface" binding:"required"`         // eg. com.mosn.service
	Method       string `json:"method" binding:"required"`            // eg. GetUser
	RegistryAddr string `json:"registry_addr" binding:"required"`     // eg. 127.0.0.1:2181
	RegistryType string `json:"registry_type" binding:"eq=zookeeper"` // zookeeper，etcd，k8s..., currently only zookeeper is supported
}

type pubReq struct {
	Registry struct {
		UserName string `json:"username"`
		Password string `json:"password"`
		Type     string `json:"type" binding:"eq=zookeeper"`  // only zookeeper is supported currently
		Addr     string `json:"addr" binding:"hostname_port"` // eg. 127.0.0.1:2181
		//Timeout  string `json:"timeout" binding:"required"`   // 5s
	} `json:"registry"`
	Service struct {
		Interface string   `json:"interface" binding:"required"`   // eg. com.mosn.service
		Methods   []string `json:"methods" binding:"required"`     // eg. GetUser,GetProfile,UpdateName
		Port      string   `json:"port" binding:"max=65535,min=1"` // user service port, eg. 8080
		Name      string   `json:"name" binding:"required"`        // eg. DemoService
		Group     string   `json:"group" binding:"required"`
	} `json:"service"`
}

type unpubReq struct {
	Interface    string `json:"interface" binding:"required"`         // eg. com.mosn.service
	Method       string `json:"method" binding:"required"`            // eg. GetUser
	RegistryAddr string `json:"registry_addr" binding:"required"`     // eg. 127.0.0.1:2181
	RegistryType string `json:"registry_type" binding:"eq:zookeeper"` // zookeeper，etcd，k8s..., currently only zookeeper is supported
}
