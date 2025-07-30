// file: pkg/ecsm_client/clientset/template_types.go

package clientset

// --- Create Request Structures ---

type CreateTemplateRequest struct {
	ImageRefs []string `json:"imageRefs"`
	Path      string   `json:"path"`
}

type CreateDictoryRequest struct {
	DictoryName string `json:"name"`
	DictoryPath string `json:"path"`
}

// --- Create Response Structures ---

type CreateTemplateResponse struct {
	ProvsionTmplList []ProvisonTmplRow `json:"provisionTmplList"`
}

type ProvisonTmplRow struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CreateDictoryResponse struct {
	DictoryID   string `json:"id"`
	DictoryPath string `json:"path"`
}

// --- Move Structures ---

type MoveRequest struct {
	Src string `json:"src"`
	Dst string `json:"dst"`
}

type MoveResponse struct {
	ID string `json:"id"`
}

// --- Update Templates Structures ---

type UpdateTemplatesRequest struct {
	Name      string       `json:"name,omitempty"`      //可选
	Templates TemplateSpec `json:"templates,omitempty"` //可选
	Action    string       `json:"action,omitempty"`    //可选的全局部署行为: "run" 或 "load"
}

type TemplateSpec struct {
	Image   ImageForTmpl `json:"image,omitempty"`
	Node    NodeSpec     `json:"node,omitempty"`
	Factor  *int         `json:"factor,omitempty"`
	Policy  string       `json:"policy,omitempty"` // "dynamic" or "static"
	Prepull *bool        `json:"prepull,omitempty"`
}

type ImageForTmpl struct {
	Ref        string          `json:"ref"`
	Config     *EcsImageConfig `json:"config"` // 假设我们只关心 EcsImageConfig
	VSOA       *ImageVSOA      `json:"vsoa,omitempty"`
	PullPolicy string          `json:"pullPolicy,omitempty"`
}

type UpdateTemplateResult struct {
	ID           string        `json:"id"`
	DeployResult *DeployResult `json:"deployResult"`
}

type UpdateTemplatesBatchRequest struct {
	Templates []TemplateUpdateBatchSpec `json:"templates"`
	Action    string                    `json:"action,omitempty"` // 可选的全局部署行为: "run" 或 "load"
}

// TemplateUpdateBatchSpec 定义了在批量更新中单个模板的更新信息。
type TemplateUpdateBatchSpec struct {
	Name     string `json:"name,omitempty"`   //通过id更新无需填写，通过name更新必填
	ID       string `json:"id,omitempty"`     //通过name更新无需填写，通过id更新必填
	ImageRef string `json:"imageRef"`         // 必填：模板的新镜像引用，格式为 name@tag#os
	Action   string `json:"action,omitempty"` // 可选：针对此模板的部署行为，优先级高于全局 Action
}

// UpdateTemplateBatchResult 代表了从更新模板 API 返回的数组中的单个结果。
type UpdateTemplateBatchResult struct {
	ID           string        `json:"id"`
	Name         string        `json:"name,omitempty"` //通过id更新无此字段，通过name更新有
	Result       bool          `json:"result"`
	Message      string        `json:"message"`
	DeployResult *DeployResult `json:"deployResult"`
}

// DeployResult 包含了当更新模板并触发部署时的结果信息。
type DeployResult struct {
	ProvisionTmpld string `json:"provisionTmpld"`

	// 部署结果: "failed", "created", "updated"
	Result string `json:"result"`

	// 部署成功后创建的服务的 ID。如果未创建服务或失败，则可能为 null。
	ProvisionID *string `json:"provisionId"`

	// 部署成功后返回的相关任务信息。
	Tasks []DeployTask `json:"tasks"`

	// 部署失败时的错误信息。如果成功，则为 null。
	Error *string `json:"error"`
}

type DeployTask struct {
	TaskID string `json:"id"`
}

// --- Get Options and Response Structures ---
// GetTemplateTreeOptions  封装了所有可以用于模板的查询参数。
type GetTemplateTreeOptions struct {
	Path  string `json:"path"`
	Level int    `json:"level,omitempty"`
	Model string `json:"model,omitempty"` //"simple" or	"full"
}

// ProvisionTmplTree 定义了模板树的节点结构。
type ProvisionTmplTree struct {
	Name       string                        `json:"name"`
	RealPath   string                        `json:"realpath"`
	ChildCount int                           `json:"childCount"`
	Data       *ProvisionTmplDetail          `json:"data"`     // 仅在 model=full 时返回
	Children   map[string]*ProvisionTmplTree `json:"children"` // 子模板树
}

// ProvisionTmplDetail 包含了模板的详细信息。
type ProvisionTmplDetail struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Kind        string   `json:"kind"` // "folder" 或 "service"
	Hostname    string   `json:"hostname"`
	Node        NodeSpec `json:"node"`
	CreatedTime string   `json:"createdTime"`
	UpdatedTime string   `json:"updatedTime"`
}

type TemplateGet struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Kind        string       `json:"kind"` // "folder" 或 "service"
	Spec        TemplateSpec `json:"spec"`
	CreatedTime string       `json:"createdTime"`
	UpdatedTime string       `json:"updatedTime"`
}

type SearchTemplateOptions struct {
	Key  string `json:"key,omitempty"`
	Path string `json:"path,omitempty"`
	Kind string `json:"kind,omitempty"`
}

type SearchTemplateResult struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Kind        string   `json:"kind"` // "folder" 或 "service"
	Hostname    string   `json:"hostname"`
	Realpath    string   `json:"realpath"`
	Node        NodeSpec `json:"node"` //文档给出NodeSpec另一种定义，查看响应实例还是保留原定义
	CreatedTime string   `json:"createdTime"`
	UpdatedTime string   `json:"updatedTime"`
}

// --- Delete Options and Response Structures ---
type DeleteTempalteResult struct {
	ID string `json:"id"`
}

type DeleteTempaltesResult struct {
	IDs []string `json:"id"`	//文档给出响应参数为id，响应示例为id数组，此处以响应示例为准
}