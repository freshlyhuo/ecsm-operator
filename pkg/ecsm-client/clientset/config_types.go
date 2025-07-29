// file: pkg/ecsm_client/clientset/config_types.go

package clientset

// ConfigItemType 定义了配置项 value 的允许类型。
// 使用自定义类型和常量，可以利用编译器的检查，避免硬编码字符串。
type ConfigItemType string

const (
	// ConfigItemTypeString 表示 value 是一个字符串。
	ConfigItemTypeString ConfigItemType = "string"
	// ConfigItemTypeNumber 表示 value 是一个数字。
	ConfigItemTypeNumber ConfigItemType = "number"
	// ConfigItemTypeJSON 表示 value 是一个 JSON 对象或数组。
	ConfigItemTypeJSON ConfigItemType = "json"
)

// --- Create Request Structures ---

// CreateConfigRequest 定义了创建新配置
type CreateConfigRequest struct {
	Key   string         `json:"key"`
	Type  ConfigItemType `json:"type"`
	Value interface{}    `json:"value"`
}

type ConfigItem struct {
	ID    string         `json:"id"`
	Key   string         `json:"key"`
	Type  ConfigItemType `json:"type"`
	Value interface{}    `json:"value"`
}

// --- List Options and Response Structures ---
// ListServiceOptions 封装了所有可以用于 List 服务的查询参数。
type ListConfigsOptions struct {
	PageNum  int    `json:"pageNum"`  // 必填
	PageSize int    `json:"pageSize"` // 必填
	Key      string `json:"key,omitempty"`
}
