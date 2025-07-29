package clientset

import (
	"context"
	"fmt"
	"strconv"

	"github.com/fx147/ecsm-operator/pkg/ecsm-client/rest"
)

type ServiceGetter interface {
	Services() ServiceInterface
}

// ServiceInterface 提供了所有操作 Service 核心资源的方法。
type ServiceInterface interface {
	// --- 核心 CRUD 操作 ---

	// Create 创建一个新的服务。
	Create(ctx context.Context, service *CreateServiceRequest) (*ServiceCreateResponse, error)

	// Get 根据服务 ID 获取一个服务的详细信息。
	Get(ctx context.Context, serviceID string) (*ServiceGet, error)

	// List 列出所有服务，支持通过 Options 进行过滤。
	List(ctx context.Context, opts ListServicesOptions) (*ServiceList, error)

	ListAll(ctx context.Context, opts ListServicesOptions) ([]ProvisionListRow, error)

	// Update 修改一个已存在的服务。
	Update(ctx context.Context, serviceID string, service *UpdateServiceRequest) (*ServiceCreateResponse, error)

	// Delete 根据服务 ID 删除一个服务。
	Delete(ctx context.Context, serviceID string) (*ServiceDeleteResponse, error)

	// CreateByPath 根据资源模板路径批量创建服务。
	CreateByPath(ctx context.Context, opts CreateByPathOptions) ([]ServiceCreateResponse, error)

	// DeleteByPath 根据资源模板路径批量删除服务。
	DeleteByPath(ctx context.Context, path string) ([]DeleteByPathResult, error)

	// ControlByID 根据 ID 批量操作服务。
	ControlByID(ctx context.Context, serviceIDs []string, action string) (*ControlServicesResponse, error)

	// ControlByLabel 根据服务路径标签批量操作服务。
	ControlByLabel(ctx context.Context, path string, action string) (*ControlServicesResponse, error)

	// --- 特殊操作 (Actions) ---

	// Redeploy 触发一次服务的重新部署。
	Redeploy(ctx context.Context, serviceID string) error

	// ValidateName 校验服务名称是否合法或可用。
	ValidateName(ctx context.Context, req *ValidateNameOptions) (*ValidationResult, error)

	// RollBack 回滚服务到指定的部署记录。
	RollBack(ctx context.Context, req *RollBackRequest) (*Transaction, error)

	// --- 状态与统计 ---

	// GetStatistics 获取服务的统计信息。
	GetStatistics(ctx context.Context) (*ServiceStatistics, error)
}

type serviceClient struct {
	restClient rest.Interface
}

func newServices(restClient rest.Interface) *serviceClient {
	return &serviceClient{restClient: restClient}
}

// Create 实现了 ServiceInterface 的 Create 方法
func (c *serviceClient) Create(ctx context.Context, service *CreateServiceRequest) (*ServiceCreateResponse, error) {
	result := &ServiceCreateResponse{}

	// 开始构建请求
	err := c.restClient.Post().
		Resource("service").
		Body(service).
		Do(ctx).
		Into(result)

	return result, err
}

func (c *serviceClient) Update(ctx context.Context, serviceID string, service *UpdateServiceRequest) (*ServiceCreateResponse, error) {
	// 业务逻辑：确保传入的 serviceID 与 body 中的 ID 一致
	if serviceID != service.ID {
		return nil, fmt.Errorf("serviceID in path (%s) does not match serviceID in body (%s)", serviceID, service.ID)
	}

	result := &ServiceCreateResponse{}

	// 开始构建请求
	err := c.restClient.Put().
		Resource("service").
		Body(service).
		Do(ctx).
		Into(result)

	return result, err
}

func (c *serviceClient) Delete(ctx context.Context, serviceID string) (*ServiceDeleteResponse, error) {
	result := &ServiceDeleteResponse{}

	// 构建请求
	err := c.restClient.Delete().
		Resource("service").
		Name(serviceID).
		Do(ctx).
		Into(result)

	return result, err
}

// DeleteByPath 实现了 ServiceInterface 的 DeleteByPath 方法。
// 它根据指定的父路径，批量删除其下的所有服务。
func (c *serviceClient) DeleteByPath(ctx context.Context, path string) ([]DeleteByPathResult, error) {
	// 准备请求体
	requestBody := struct {
		Path string `json:"path"`
	}{
		Path: path,
	}
 
	var result []DeleteByPathResult
 
	// 开始构建请求
	// Verb: DELETE
	// URL: /service/path
	// Body: requestBody (会被序列化成 JSON)
	err := c.restClient.Delete().
		Resource("service/path").
		Body(requestBody).
		Do(ctx).
		Into(&result)
 
	return result, err
}

func (c *serviceClient) CreateByPath(ctx context.Context, opts CreateByPathOptions) ([]ServiceCreateResponse, error) {
	var allResponse []ServiceCreateResponse

	// 从 opts 中分离出用于 path 的参数，并进行校验
	action := opts.Action
	if action != "run" && action != "load" {
		return nil, fmt.Errorf("invalid action: '%s', must be 'run' or 'load'", action)
	}

	err := c.restClient.Post().
		Resource("service").
		Subresource(action).
		Subresource("templates-path-label").
		Body(opts).
		Do(ctx).
		Into(&allResponse)

	return allResponse, err
}

func (c *serviceClient) Get(ctx context.Context, serviceID string) (*ServiceGet, error) {
	result := &ServiceGet{}

	// 开始构建请求
	err := c.restClient.Get().
		Resource("service").
		Name(serviceID).
		Do(ctx).
		Into(result)

	return result, err
}

func (c *serviceClient) ControlByID(ctx context.Context, serviceIDs []string, action string) (*ControlServicesResponse, error) {
	// 验证操作类型是否有效
	validActions := map[string]bool{
		"start":   true,
		"stop":    true,
		"restart": true,
		"pause":   true,
		"unpause": true,
		"destroy": true,
	}
	if !validActions[action] {
		return nil, fmt.Errorf("invalid action: '%s', must be one of [start, stop, restart, pause, unpause, destroy]", action)
	}

	// 准备请求体
	requestBody := &ControlServicesResponse{
		IDs: serviceIDs,
	}

	result := &ControlServicesResponse{}

	err := c.restClient.Post().
		Resource("service").
		Subresource(action).
		Subresource("ids").
		Body(requestBody).
		Do(ctx).
		Into(result)

	return result, err
}

func (c *serviceClient) ControlByLabel(ctx context.Context, path string, action string) (*ControlServicesResponse, error) {
	// 验证操作类型是否有效
	validActions := map[string]bool{
		"start":   true,
		"stop":    true,
		"restart": true,
		"pause":   true,
		"unpause": true,
		"destroy": true,
	}
	if !validActions[action] {
		return nil, fmt.Errorf("invalid action: '%s', must be one of [start, stop, restart, pause, unpause, destroy]", action)
	}

	// 准备请求体
	requestBody := struct {
		Path string `json:"path"`
	}{
		Path: path,
	}

	result := &ControlServicesResponse{}

	err := c.restClient.Post().
		Resource("service").
		Subresource(action).
		Subresource("path-label").
		Body(requestBody).
		Do(ctx).
		Into(result)

	return result, err
}

// List 实现了 ServiceInterface 的 List 方法。
func (c *serviceClient) List(ctx context.Context, opts ListServicesOptions) (*ServiceList, error) {
	result := &ServiceList{}

	// 开始构建请求
	req := c.restClient.Get().Resource("service")

	// 将 Options 结构体翻译成 URL Query 参数
	req.Param("pageNum", strconv.Itoa(opts.PageNum))
	req.Param("pageSize", strconv.Itoa(opts.PageSize))
	if opts.Name != "" {
		req.Param("name", opts.Name)
	}
	if opts.ImageID != "" {
		// 注意：我们在 Go 结构体中叫 ImageID，但 API 参数是 id
		req.Param("id", opts.ImageID)
	}
	if opts.NodeID != "" {
		req.Param("nodeId", opts.NodeID)
	}
	if opts.Label != "" {
		req.Param("label", opts.Label)
	}

	// 执行请求并解码结果
	err := req.Do(ctx).Into(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *serviceClient) ListAll(ctx context.Context, opts ListServicesOptions) ([]ProvisionListRow, error) {
	var allItems []ProvisionListRow
	opts.PageNum = 1
	if opts.PageSize == 0 {
		opts.PageSize = 100
	}

	for {
		list, err := c.List(ctx, opts)
		if err != nil {
			return nil, err
		}

		if len(list.Items) == 0 {
			break
		}

		allItems = append(allItems, list.Items...)

		if len(allItems) >= list.Total {
			break
		}

		opts.PageNum++
	}
	return allItems, nil
}

func (c *serviceClient) Redeploy(ctx context.Context, serviceID string) error {
	// 创建请求体的实例
	requestBody := RedeployRequest{
		ID: serviceID,
	}

	err := c.restClient.Put().
		Resource("service/deployment/restart").
		Body(requestBody).
		Do(ctx).
		Into(nil)

	return err
}

func (c *serviceClient) ValidateName(ctx context.Context, opts *ValidateNameOptions) (*ValidationResult, error) {
	var nameExists bool

	// 开始构建请求
	req := c.restClient.Get().Resource("service/name/check")

	// 添加查询参数
	req.Param("name", opts.Name)
	if opts.ID != "" {
		req.Param("id", opts.ID)
	}

	err := req.Do(ctx).Into(&nameExists)
	if err != nil {
		return nil, err
	}

	// 将 API 返回的 "exists" (存在) 逻辑，转换为我们更通用的 "IsValid" (有效) 逻辑
	// 如果 nameExists 为 true，说明名称已存在，即名称无效 (IsValid = false)
	result := &ValidationResult{
		IsValid: !nameExists,
	}

	if nameExists {
		result.Message = fmt.Sprintf("service name '%s' already exists", opts.Name)
	}

	return result, nil

}

func (c *serviceClient) RollBack(ctx context.Context, req *RollBackRequest) (*Transaction, error) {
	result := &Transaction{}

	// 开始构建请求
	err := c.restClient.Put().
		Resource("service").
		Subresource("rollback").
		Body(req).
		Do(ctx).
		Into(result)

	return result, err
}

func (c *serviceClient) GetStatistics(ctx context.Context) (*ServiceStatistics, error) {
	result := &ServiceStatistics{}

	err := c.restClient.Get().
		Resource("service/summary").
		Do(ctx).
		Into(result)

	return result, err
}
