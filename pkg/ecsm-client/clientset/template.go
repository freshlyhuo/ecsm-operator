package clientset

import (
	"context"
	"fmt"
	"strconv"

	"github.com/fx147/ecsm-operator/pkg/ecsm-client/rest"
)

type TemplateGetter interface {
	Templates() TemplateInterface
}

// TemplateInterface 提供了所有操作 Template 核心资源的方法。
type TemplateInterface interface {

	// CreateTemplate 创建一个新的模板。
	CreateTemplate(ctx context.Context, template *CreateTemplateRequest) (*CreateTemplateResponse, error)

	CreateDictory(ctx context.Context, dictory *CreateDictoryRequest) (*CreateDictoryResponse, error)

	// 移动资源模板或模板目录
	MoveTempOrDict(ctx context.Context, move *MoveRequest) (*MoveResponse, error)

	UpdateTemplate(ctx context.Context, TemplateID string, req *UpdateTemplatesRequest) (*UpdateTemplateResult, error)

	// 根据ID批量更新模板
	UpdateTemplatesByID(ctx context.Context, req *UpdateTemplatesBatchRequest) ([]UpdateTemplateBatchResult, error)

	//根据Name批量更新模板
	UpdateTemplatesByName(ctx context.Context, req *UpdateTemplatesBatchRequest) ([]UpdateTemplateBatchResult, error)

	// GetTemplateTree 获取指定路径下的资源模板和模板目录树。
	GetTemplateTree(ctx context.Context, opts GetTemplateTreeOptions) (*ProvisionTmplTree, error)

	GetTemplateByID(ctx context.Context, templateID string) (*TemplateGet, error)

	GetTemplateByPath(ctx context.Context, templatePath string) (*TemplateGet, error)

	//搜索指定路径下的资源模板和模板目录，支持指定搜索关键字、目录和类型。
	SearchTempOrDict(ctx context.Context, opts SearchTemplateOptions) (*SearchTemplateResult, error)

	DeleteTempOrDict(ctx context.Context, path string) (*DeleteTempalteResult, error)

	DeleteTempOrDictByID(ctx context.Context, templateID string) (*DeleteTempalteResult, error)

	//根据ID列表删除模板或目录
	DeleteTempOrDictByIDs(ctx context.Context, templateIDs []string) (*DeleteTempaltesResult, error)
}

type templateClient struct {
	restClient rest.Interface
}

func newTemplates(restClient rest.Interface) *templateClient {
	return &templateClient{restClient: restClient}
}

func (c *templateClient) CreateTemplate(ctx context.Context, template *CreateTemplateRequest) (*CreateTemplateResponse, error) {
	result := &CreateTemplateResponse{}

	// 开始构建请求
	err := c.restClient.Post().
		Resource("provision-template/path-label/service/batch").
		Body(template).
		Do(ctx).
		Into(result)

	return result, err
}

func (c *templateClient) CreateDictory(ctx context.Context, dictory *CreateDictoryRequest) (*CreateDictoryResponse, error) {
	result := &CreateDictoryResponse{}

	// 开始构建请求
	err := c.restClient.Post().
		Resource("provision-template/path-label/folder").
		Body(dictory).
		Do(ctx).
		Into(result)

	return result, err
}

func (c *templateClient) MoveTempOrDict(ctx context.Context, move *MoveRequest) (*MoveResponse, error) {
	result := &MoveResponse{}

	// 开始构建请求
	err := c.restClient.Put().
		Resource("provision-template/path-label/move").
		Body(move).
		Do(ctx).
		Into(result)

	return result, err
}

func (c *templateClient) UpdateTemplate(ctx context.Context, TemplateID string, req *UpdateTemplatesRequest) (*UpdateTemplateResult, error) {
	result := &UpdateTemplateResult{}

	// 构建 PUT 请求
	err := c.restClient.Put().
		Resource("provision-templates").
		Name(TemplateID).
		Body(req).
		Do(ctx).
		Into(&result)

	return result, err
}

func (c *templateClient) UpdateTemplatesByID(ctx context.Context, req *UpdateTemplatesBatchRequest) ([]UpdateTemplateBatchResult, error) {
	// 前置校验：确保每个模板都提供了ID
	for i, t := range req.Templates {
		if t.ID == "" {
			return nil, fmt.Errorf("template ID is required for UpdateTemplatesByID at index %d", i)
		}
	}

	var result []UpdateTemplateBatchResult

	// 构建 PUT 请求
	err := c.restClient.Put().
		Resource("provision-templates").
		Body(req).
		Do(ctx).
		Into(&result)

	return result, err
}

func (c *templateClient) UpdateTemplatesByName(ctx context.Context, req *UpdateTemplatesBatchRequest) ([]UpdateTemplateBatchResult, error) {
	// 前置校验：确保每个模板都提供了Name
	for i, t := range req.Templates {
		if t.Name == "" {
			return nil, fmt.Errorf("template Name is required for UpdateTemplatesByName at index %d", i)
		}
	}

	var result []UpdateTemplateBatchResult

	// 构建 PUT 请求
	err := c.restClient.Put().
		Resource("provision-templates/images").
		Body(req).
		Do(ctx).
		Into(&result)

	return result, err
}

// GetTemplateTree 实现了获取模板树的接口
// API: GET /api/v1/provision-template/path-label/tree
func (c *templateClient) GetTemplateTree(ctx context.Context, opts GetTemplateTreeOptions) (*ProvisionTmplTree, error) {
	result := &ProvisionTmplTree{}

	req := c.restClient.Get().
		Resource("provision-template/path-label/tree").
		Param("path", opts.Path)

	// 添加可选参数 level
	if opts.Level > 0 {
		req.Param("level", strconv.Itoa(opts.Level))
	}

	// 添加可选参数 model
	if opts.Model != "" {
		req.Param("model", opts.Model)
	}

	err := req.Do(ctx).Into(result)

	return result, err
}

func (c *templateClient) GetTemplateByID(ctx context.Context, templateID string) (*TemplateGet, error) {
	result := &TemplateGet{}

	err := c.restClient.Get().
		Resource("provision-template").
		Name(templateID).
		Do(ctx).
		Into(result)

	return result, err
}

func (c *templateClient) GetTemplateByPath(ctx context.Context, templatePath string) (*TemplateGet, error) {
	result := &TemplateGet{}

	err := c.restClient.Get().
		Resource("provision-template/path-label").
		Param("path", templatePath).
		Do(ctx).
		Into(result)

	return result, err
}

func (c *templateClient) SearchTempOrDict(ctx context.Context, opts SearchTemplateOptions) (*SearchTemplateResult, error) {
	result := &SearchTemplateResult{}

	req := c.restClient.Get().Resource("provision-template/path-label/search")

	if opts.Key != "" {
		req.Param("key", opts.Key)
	}
	if opts.Path != "" {
		req.Param("path", opts.Path)
	}
	if opts.Kind != "" {
		req.Param("kind", opts.Kind)
	}

	err := req.Do(ctx).Into(result)

	return result, err
}

func (c *templateClient) DeleteTempOrDict(ctx context.Context, path string) (*DeleteTempalteResult, error) {
	// 准备请求体
	requestBody := struct {
		Path string `json:"path"`
	}{
		Path: path,
	}

	result := &DeleteTempalteResult{}

	err := c.restClient.Delete().
		Resource("provision-template/path-label").
		Body(requestBody).
		Do(ctx).
		Into(&result)

	return result, err
}

func (c *templateClient) DeleteTempOrDictByID(ctx context.Context, templateID string) (*DeleteTempalteResult, error) {
	result := &DeleteTempalteResult{}

	err := c.restClient.Delete().
		Resource("provision-template/path-label").
		Name(templateID).
		Do(ctx).
		Into(&result)

	return result, err
}

func (c *templateClient) DeleteTempOrDictByIDs(ctx context.Context, templateIDs []string) (*DeleteTempaltesResult, error) {
	// 准备请求体
	requestBody := &DeleteTempaltesResult{
		IDs: templateIDs,
	}

	result := &DeleteTempaltesResult{}

	err := c.restClient.Delete().
		Resource("provision-template/path-label").
		Body(requestBody).
		Do(ctx).
		Into(&result)

	return result, err
}
