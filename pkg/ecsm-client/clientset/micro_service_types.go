// file: pkg/ecsm_client/clientset/micro_service_types.go

package clientset

// --- List Options and Response Structures ---
// ListMicroServiceOptions 封装了所有可以用于 List 服务的查询参数。
type ListMicroServicesOptions struct {
	PageNum  int    `json:"pageNum"`  // 必填
	PageSize int    `json:"pageSize"` // 必填
	KeyWord  string `json:"keyWord,omitempty"`
	// 注意：API 文档中的 'projectId' 字段名可能会引起混淆，因为它指的是镜像ID，
	// 我们在结构体中用更明确的名字 ImageID。
	ImageID int    `json:"projectId,omitempty"`
	NodeID  string `json:"nodeId,omitempty"`
	Label   string `json:"label,omitempty"`
}

type MicroServiceList struct {
	Total    int                   `json:"total"`
	PageNum  int                   `json:"pageNum"`
	PageSize int                   `json:"pageSize"`
	Items    []MicroServiceListRow `json:"list"` // 字段名是 "list"
}

type MicroServiceListRow struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	ImageName      string `json:"imageName"`
	HealthInstance int    `json:"healthInstance"`
	Instance       int    `json:"instance"`
	LoadBalance    string `json:"loadBalance"` //负载均衡策略,roundRobin：轮询;  masterSlave：主备
}

// --- Get Options and Response Structures ---
type MicroServiceGet struct {
	BoDYnamic         bool                    `json:"boDynamic"`
	ID                string                  `json:"id"`
	Name              string                  `json:"name"`
	ImageName         string                  `json:"imageName"`
	HealthInstance    int                     `json:"healthInstance"`
	Instance          int                     `json:"instance"`
	LoadBalance       string                  `json:"loadBalance"`
	LoadBalanceDetail []LoadBalanceDetailSpec `json:"loadBalanceDetail,omitempty"` //负载均衡策略详情（masterSlave 才有）
}

type LoadBalanceDetailSpec struct {
	Master string `json:"master"`
	TaskID string `json:"id"` //备份节点的 taskId,文档中使用"id"
}

// --- Update Request Structures ---

// UpdateMicroServiceRequest 定义了更新一个服务时，ECSM API 所需的 payload。
type UpdateMicroServiceRequest struct {
	ID                string                  `json:"id"`
	LoadBalance       string                  `json:"loadBalance"`
	LoadBalanceDetail []LoadBalanceDetailSpec `json:"loadBalanceDetail,omitempty"` //负载均衡策略详情（masterSlave 才有）
}

// --- NEW: List MicroService Instances Structures ---

// ListMicroServiceInstancesOptions 封装了查询微服务实例列表的查询参数。
type ListMicroServiceInstancesOptions struct {
	// ID 是微服务的主键，必填
	ID string `json:"id"`
	// PageNum 页码，如果传 -1，则展示全部数据
	PageNum int `json:"pageNum"`
	// PageSize 每页的实例数量
	PageSize int `json:"pageSize"`
	// KeyWord 根据容器名称模糊查询，可选
	KeyWord string `json:"keyWord,omitempty"`
}

// MicroServiceInstanceList 是查询微服务实例列表的响应体结构。
type MicroServiceInstanceList struct {
	Total    int             `json:"total"`
	PageNum  int             `json:"pageNum"`
	PageSize int             `json:"pageSize"`
	Items    []ContainerInfo `json:"list"` // API文档中字段名为 "list"
}
