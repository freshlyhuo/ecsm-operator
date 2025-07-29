// file: pkg/ecsm_client/clientset/record_types.go

package clientset

//GetReord 部署记录详情
type RecordGet struct {
	Name        string          `json:"name"`
	Image       ImageInfo       `json:"image"`
	Node        NodeInfo        `json:"node"`
	Action      string          `json:"action"`
	Policy      string          `json:"policy"` // "dynamic" or "static"
	Factor      *int            `json:"factor"` // Factor 是动态部署策略 ("dynamic") 下的容器实例数量。
	Cmd         []string        `json:"cmd"`
	VSOA        *ImageVSOA      `json:"vsoa"`
	Config      *EcsImageConfig `json:"config"` // 假设我们只关心 EcsImageConfig
	CreatedTime string          `json:"createdTime"`
}

type ImageInfo struct {
	Ref         string `json:"ref"`
	Path        string `json:"path"`
	PullPolicy  string `json:"pullPolicy"`
	AutoUpgrade string `json:"autoUpgrade"`
}

// --- List Options and Response Structures ---
// ListRecordOptions 封装了所有可以用于 List 服务的查询参数。
type ListRecordOptions struct {
	PageNum   int    `json:"pageNum"`
	PageSize  int    `json:"pageSize"`
	ServiceID string `json:"serviceId"`
}

//
type RecordList struct {
	Total    int            `json:"total"`
	PageNum  int            `json:"pageNum"`
	PageSize int            `json:"pageSize"`
	Items    []DeployRecord `json:"list"` // 字段名是 "list"
}

//DeployRecord
type DeployRecord struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	CMD         string    `json:"cmd"` //容器启动参数
	CreatedTime string    `json:"createdTime"`
	Image       string    `json:"image"`
	Node        *NodeSpec `json:"node,omitempty"` // <-- 复用共享类型
}
