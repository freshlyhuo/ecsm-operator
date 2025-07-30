package clientset

import (
	"context"
	"fmt"
	"strconv"

	"github.com/fx147/ecsm-operator/pkg/ecsm-client/rest"
)

type MicroServiceGetter interface {
	MicroServices() MicroServiceInterface
}

// ServiceInterface 提供了所有操作 Service 核心资源的方法。
type MicroServiceInterface interface {

	// List 列出所有微服务，支持通过 Options 进行过滤。
	ListMicroService(ctx context.Context, opts ListMicroServicesOptions) (*MicroServiceList, error)

	ListAllMicroService(ctx context.Context, opts ListMicroServicesOptions) ([]MicroServiceListRow, error)

	GetMicroService(ctx context.Context, MicroServiceID string) (*MicroServiceGet, error)

	UpdateMicroService(ctx context.Context, MicroService *UpdateMicroServiceRequest) error

	// ListMicroServiceInstances 查询指定微服务的容器实例列表
	//ListMicroServiceInstances(ctx context.Context, opts ListMicroServiceInstancesOptions) (*MicroServiceInstanceList, error)
	// ListAllMicroServiceInstances 查询指定微服务的所有容器实例（自动处理分页）
	//ListAllMicroServiceInstances(ctx context.Context, opts ListMicroServiceInstancesOptions) ([]ContainerInfo, error)
}

type MicroserviceClient struct {
	restClient rest.Interface
}

func newMicroServices(restClient rest.Interface) *MicroserviceClient {
	return &MicroserviceClient{restClient: restClient}
}

func (c *MicroserviceClient) ListMicroService(ctx context.Context, opts ListMicroServicesOptions) (*MicroServiceList, error) {
	result := &MicroServiceList{}

	// 开始构建请求
	req := c.restClient.Get().Resource("micro-service")

	// 将 Options 结构体翻译成 URL Query 参数
	req.Param("pageNum", strconv.Itoa(opts.PageNum))
	req.Param("pageSize", strconv.Itoa(opts.PageSize))
	if opts.KeyWord != "" {
		req.Param("name", opts.KeyWord)
	}
	if opts.ImageID != 0 {
		// 注意：我们在 Go 结构体中叫 ImageID，但 API 参数是 projectId
		req.Param("projectId", strconv.Itoa(opts.ImageID))
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

func (c *MicroserviceClient) ListAllMicroService(ctx context.Context, opts ListMicroServicesOptions) ([]MicroServiceListRow, error) {
	var allItems []MicroServiceListRow
	opts.PageNum = 1
	if opts.PageSize == 0 {
		opts.PageSize = 100
	}

	for {
		list, err := c.ListMicroService(ctx, opts)
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

func (c *MicroserviceClient) GetMicroService(ctx context.Context, MicroServiceID string) (*MicroServiceGet, error) {
	result := &MicroServiceGet{}

	// 开始构建请求
	err := c.restClient.Get().
		Resource("micro-service").
		Name(MicroServiceID).
		Do(ctx).
		Into(result)

	return result, err
}

func (c *MicroserviceClient) UpdateMicroService(ctx context.Context, MicroService *UpdateMicroServiceRequest) error {
	// 如果负载均衡策略不是 "masterSlave"，但 loadBalanceDetail 字段却被设置了，
	// 这就是一个无效的请求。我们应该立即返回错误，而不是发送它。
	// 我们使用 len > 0 而不是 != nil，因为一个空的切片 `[]LoadBalanceDetail{}` 也不应该被发送。
	if MicroService.LoadBalance != "masterSlave" && len(MicroService.LoadBalanceDetail) > 0 {
		return fmt.Errorf(
			"validation failed: loadBalanceDetail can only be set when loadBalance strategy is 'masterSlave', but the strategy was '%s'",
			MicroService.LoadBalance,
		)
	}

	// 如果策略是 "masterSlave"，我们也应该确保 detail 字段不为空
	if MicroService.LoadBalance == "masterSlave" && len(MicroService.LoadBalanceDetail) == 0 {
		return fmt.Errorf("validation failed: loadBalanceDetail must be provided when loadBalance strategy is 'masterSlave'")
	}

	// 3. 执行请求（只有在所有校验通过后）
	err := c.restClient.Put().
		Resource("micro-service").
		Body(MicroService).
		Do(ctx).
		Into(nil)

	return err
}
