package clientset

import (
	"context"
	"strconv"

	"github.com/fx147/ecsm-operator/pkg/ecsm-client/rest"
)

type RecordGetter interface {
	Records() RecordInterface
}

type RecordInterface interface {

	// Get 根据部署记录 ID 获取一个部署记录的详细信息。
	GetRecord(ctx context.Context, RecordId string) (*RecordGet, error)

	// Delete 根据部署记录 ID 删除一个部署记录。
	DeleteRecord(ctx context.Context, RecordId string) error

	//List列出指定ServiceId的部署记录
	ListRecord(ctx context.Context, opts ListRecordOptions) (*RecordList, error)

	ListAllRecord(ctx context.Context, opts ListRecordOptions) ([]DeployRecord, error)
}

type recordClient struct {
	restClient rest.Interface
}

func newRecords(restClient rest.Interface) *recordClient {
	return &recordClient{restClient: restClient}
}

func (c *recordClient) GetRecord(ctx context.Context, RecordId string) (*RecordGet, error) {
	result := &RecordGet{}

	// 开始构建请求
	err := c.restClient.Get().
		Resource("service/record").
		Name(RecordId).
		Do(ctx).
		Into(result)

	return result, err
}

func (c *recordClient) DeleteRecord(ctx context.Context, recordID string) error {
	// 开始构建请求
	req := c.restClient.Delete().
		Resource("service/record")

	// 将部署记录 ID 添加为查询参数
	req.Param("id", recordID)

	// 执行请求。因为成功时我们不关心返回的 "success" 字符串，
	// 所以使用 Into(nil) 来忽略响应体。
	err := req.Do(ctx).Into(nil)

	return err
}

func (c *recordClient) ListRecord(ctx context.Context, opts ListRecordOptions) (*RecordList, error) {
	result := &RecordList{}

	// 开始构建请求
	req := c.restClient.Get().Resource("service/record")

	// 将 Options 结构体翻译成 URL Query 参数
	req.Param("id", opts.ServiceID)
	req.Param("pageNum", strconv.Itoa(opts.PageNum))
	req.Param("pageSize", strconv.Itoa(opts.PageSize))

	// 执行请求并解码结果
	err := req.Do(ctx).Into(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *recordClient) ListAllRecord(ctx context.Context, opts ListRecordOptions) ([]DeployRecord, error) {
	var allItems []DeployRecord
	opts.PageNum = 1
	if opts.PageSize == 0 {
		opts.PageSize = 100
	}

	for {
		list, err := c.ListRecord(ctx, opts)
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
