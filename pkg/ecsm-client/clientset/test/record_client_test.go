// file: pkg/ecsm-client/clientset/test/record_client_test.go

package test

import (
	"context"
	"testing"

	"github.com/fx147/ecsm-operator/pkg/ecsm-client/clientset"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- newTestClientset() 辅助函数 (已在 service_client_test.go 中定义) ---

// TestRecordClient_ReadOperations 测试部署记录的只读操作 (List 和 Get)。
// 这个测试是安全的，因为它不会修改任何外部系统状态。
// 它依赖于你的 ECSM 环境中至少有一个服务，并且该服务至少有一条部署记录。
func TestRecordClient_ReadOperations(t *testing.T) {
	// --- Setup ---
	cs := newTestClientset(t)
	recordClient := cs.Records()
	serviceClient := cs.Services() // 我们需要服务客户端来找到一个有效的 Service ID
	ctx := context.Background()

	// 1. 找到一个有部署记录的服务
	// 首先，列出所有服务
	serviceList, err := serviceClient.List(ctx, clientset.ListServicesOptions{PageNum: 1, PageSize: 100})
	require.NoError(t, err)
	require.NotEmpty(t, serviceList.Items, "测试失败：ECSM 环境中必须至少存在一个服务")

	var targetServiceID string
	// 遍历所有服务，找到第一个后就跳出
	for _, service := range serviceList.Items {
		// 检查该服务是否有部署记录
		records, listErr := recordClient.ListRecord(ctx, clientset.ListRecordOptions{ServiceID: service.ID, PageNum: 1, PageSize: 1})
		if listErr == nil && len(records.Items) > 0 {
			targetServiceID = service.ID
			t.Logf("找到一个有部署记录的服务: Name=%s, ID=%s", service.Name, service.ID)
			break
		}
	}
	require.NotEmpty(t, targetServiceID, "测试失败：在所有服务中都找不到任何部署记录")

	// --- Test: ListRecord 和 ListAllRecord ---
	var firstRecord clientset.DeployRecord
	t.Run("ListAndListAll", func(t *testing.T) {
		// 2. 使用找到的服务 ID 列出部署记录
		opts := clientset.ListRecordOptions{
			ServiceID: targetServiceID,
			PageNum:   1,
			PageSize:  5,
		}
		list, err := recordClient.ListRecord(ctx, opts)
		require.NoError(t, err)
		require.NotNil(t, list)
		require.NotEmpty(t, list.Items, "服务 (ID: %s) 应该有部署记录", targetServiceID)
		firstRecord = list.Items[0] // 保存第一个记录用于后续 Get 测试

		// 验证分页信息
		assert.Equal(t, opts.PageNum, list.PageNum)
		assert.Equal(t, opts.PageSize, list.PageSize)
		assert.GreaterOrEqual(t, list.Total, len(list.Items))

		// 测试 ListAllRecord
		allRecords, err := recordClient.ListAllRecord(ctx, clientset.ListRecordOptions{ServiceID: targetServiceID, PageSize: 5})
		require.NoError(t, err)
		assert.Len(t, allRecords, list.Total, "ListAllRecord 返回的记录总数应与分页接口中的 Total 字段一致")
	})

	// --- Test: GetRecord ---
	t.Run("GetRecord", func(t *testing.T) {
		require.NotEmpty(t, firstRecord.ID, "无法从 List 中获取有效的部署记录 ID")

		t.Logf("正在获取部署记录详情, ID: %s", firstRecord.ID)
		details, err := recordClient.GetRecord(ctx, firstRecord.ID)

		require.NoError(t, err)
		require.NotNil(t, details)

		// 验证获取到的详情
		assert.Equal(t, firstRecord.Name, details.Name, "Get返回的记录名称应与List中的名称一致")
		assert.NotEmpty(t, details.Action)
		assert.NotNil(t, details.Image)
		assert.NotEmpty(t, details.Image.Ref)
		assert.NotEmpty(t, details.CreatedTime)
	})
}

// TestRecordClient_Delete 测试删除部署记录。
// 这是一个写操作，会修改外部系统。
func TestRecordClient_Delete(t *testing.T) {
	// --- Setup ---
	cs := newTestClientset(t)
	ctx := context.Background()
	recordClient := cs.Records()
	serviceClient := cs.Services()

	// 1. 找到一个可以被删除的部署记录。逻辑与上面类似。
	serviceList, err := serviceClient.List(ctx, clientset.ListServicesOptions{PageNum: 1, PageSize: 100})
	require.NoError(t, err)
	require.NotEmpty(t, serviceList.Items)

	var recordToDeleteID string
	for _, service := range serviceList.Items {
		records, listErr := recordClient.ListRecord(ctx, clientset.ListRecordOptions{ServiceID: service.ID, PageNum: 1, PageSize: 10})
		if listErr == nil && len(records.Items) > 0 {
			// 为了安全，我们选择删除最后一个记录，而不是第一个
			recordToDeleteID = records.Items[len(records.Items)-1].ID
			t.Logf("找到一个可删除的部署记录: ID=%s, 所属服务: %s", recordToDeleteID, service.Name)
			break
		}
	}
	require.NotEmpty(t, recordToDeleteID, "测试失败：找不到任何可以删除的部署记录")

	// --- Test: DeleteRecord ---
	t.Run("DeleteAndVerify", func(t *testing.T) {
		// 2. 执行删除操作
		err := recordClient.DeleteRecord(ctx, recordToDeleteID)
		require.NoError(t, err, "删除部署记录 (ID: %s) 失败", recordToDeleteID)
		t.Logf("成功提交删除部署记录 (ID: %s) 的请求", recordToDeleteID)

		// 3. 验证删除
		// 尝试再次获取已删除的记录，应该会失败或返回空
		_, err = recordClient.GetRecord(ctx, recordToDeleteID)
		// 我们期望这里返回一个错误（通常是 Not Found 错误）
		assert.Error(t, err, "获取已删除的部署记录应该返回错误")
		t.Logf("验证成功：获取已删除的记录返回错误: %v", err)
	})
}
