package servicediscovery

type subReq struct {
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
		//Port      string   `json:"port" binding:"max=65535,min=1"` // user service port, eg. 8080
		Name      string   `json:"name" binding:"required"`        // eg. DemoService
		Group     string   `json:"group" binding:"required"`
	} `json:"service"`
}

type unsubReq struct {}

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
		//Port      string   `json:"port" binding:"max=65535,min=1"` // user service port, eg. 8080
		Name      string   `json:"name" binding:"required"`        // eg. DemoService
		Group     string   `json:"group" binding:"required"`
	} `json:"service"`
}

type unpubReq struct {
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
		//Port      string   `json:"port" binding:"max=65535,min=1"` // user service port, eg. 8080
		Name      string   `json:"name" binding:"required"`        // eg. DemoService
		Group     string   `json:"group" binding:"required"`
	} `json:"service"`
}

// response struct for all requests
type resp struct {
	Errno  int    `json:"err_no"`
	ErrMsg string `json:"err_msg"`
}
