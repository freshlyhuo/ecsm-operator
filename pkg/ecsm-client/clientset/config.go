package clientset

import (
	"context"
	"fmt"
	"reflect"
	"strconv"

	"github.com/fx147/ecsm-operator/pkg/ecsm-client/rest"
)

type ConfigGetter interface {
	Configs() ConfigInterface
}

// ServiceInterface 提供了所有操作 Service 核心资源的方法。
type ConfigInterface interface {
	//新建配置
	CreateConfig(ctx context.Context, config *CreateConfigRequest) error

	UpdateConfig(ctx context.Context, config *ConfigItem) error

	DeleteConfig(ctx context.Context, configID string) error

	// 根据 Key 查询配置项。
	GetConfig(ctx context.Context, key string) (interface{}, error)

	// List 列出所有服务，支持通过 Options 进行过滤。
	ListConfig(ctx context.Context, opts ListConfigsOptions) ([]ConfigItem, error)

	ListAllConfig(ctx context.Context, opts ListConfigsOptions) ([]ConfigItem, error)
}

type configClient struct {
	restClient rest.Interface
}

func newConfigs(restClient rest.Interface) *configClient {
	return &configClient{restClient: restClient}
}

// Create 实现了创建配置项的逻辑, 包含客户端校验。
func (c *configClient) CreateConfig(ctx context.Context, config *CreateConfigRequest) error {
	// 检查 Value 的实际类型是否与 'Type' 字段匹配
	switch config.Type {
	case ConfigItemTypeString:
		if _, ok := config.Value.(string); !ok {
			return fmt.Errorf("config type is 'string', but provided value is of type %T, not a string", config.Value)
		}
	case ConfigItemTypeNumber:
		valKind := reflect.ValueOf(config.Value).Kind()
		if valKind != reflect.Float64 && valKind != reflect.Int && valKind != reflect.Int32 && valKind != reflect.Int64 {
			return fmt.Errorf("config type is 'number', but provided value is of type %T, not a number", config.Value)
		}
	case ConfigItemTypeJSON:
		valKind := reflect.ValueOf(config.Value).Kind()
		if valKind != reflect.Map && valKind != reflect.Slice {
			return fmt.Errorf("config type is 'json', but provided value is a %s, not a map or slice", valKind)
		}
	default:
		return fmt.Errorf("unsupported config type: '%s'. Must be one of: %s, %s, %s",
			config.Type, ConfigItemTypeString, ConfigItemTypeNumber, ConfigItemTypeJSON)
	}

	err := c.restClient.Post().
		Resource("configmap").
		Body(config).
		Do(ctx).
		Into(nil)

	return err
}

func (c *configClient) UpdateConfig(ctx context.Context, config *ConfigItem) error {
	// 检查 Value 的实际类型是否与 'Type' 字段匹配
	switch config.Type {
	case ConfigItemTypeString:
		if _, ok := config.Value.(string); !ok {
			return fmt.Errorf("config type is 'string', but provided value is of type %T, not a string", config.Value)
		}
	case ConfigItemTypeNumber:
		valKind := reflect.ValueOf(config.Value).Kind()
		if valKind != reflect.Float64 && valKind != reflect.Int && valKind != reflect.Int32 && valKind != reflect.Int64 {
			return fmt.Errorf("config type is 'number', but provided value is of type %T, not a number", config.Value)
		}
	case ConfigItemTypeJSON:
		valKind := reflect.ValueOf(config.Value).Kind()
		if valKind != reflect.Map && valKind != reflect.Slice {
			return fmt.Errorf("config type is 'json', but provided value is a %s, not a map or slice", valKind)
		}
	default:
		return fmt.Errorf("unsupported config type: '%s'. Must be one of: %s, %s, %s",
			config.Type, ConfigItemTypeString, ConfigItemTypeNumber, ConfigItemTypeJSON)
	}

	err := c.restClient.Put().
		Resource("configmap").
		Body(config).
		Do(ctx).
		Into(nil)

	return err
}

func (c *configClient) DeleteConfig(ctx context.Context, configID string) error {
	err := c.restClient.Delete().
		Resource("configmap").
		Name(configID).
		Do(ctx).
		Into(nil)

	return err
}

func (c *configClient) GetConfig(ctx context.Context, key string) (interface{}, error) {
	var result interface{}

	err := c.restClient.Get().
		Resource("configmap/key").
		Param("key", key).
		Do(ctx).
		Into(&result)

	return result, err
}

func (c *configClient) ListConfig(ctx context.Context, opts ListConfigsOptions) ([]ConfigItem, error) {
	var result []ConfigItem

	req := c.restClient.Get().Resource("configmap")

	req.Param("pageNum", strconv.Itoa(opts.PageNum))
	req.Param("pageSize", strconv.Itoa(opts.PageSize))
	if opts.Key != "" {
		req.Param("key", opts.Key)
	}

	err := req.Do(ctx).Into(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *configClient) ListAllConfig(ctx context.Context, opts ListConfigsOptions) ([]ConfigItem, error) {
	var allItems []ConfigItem
	opts.PageNum = 1
	if opts.PageSize == 0 {
		opts.PageSize = 100
	}

	for {
		pageItems, err := c.ListConfig(ctx, opts)
		if err != nil {
			return nil, err
		}

		if len(pageItems) == 0 {
			break
		}

		allItems = append(allItems, pageItems...)

		if len(pageItems) < opts.PageSize {
			break
		}

		opts.PageNum++
	}
	return allItems, nil
}
