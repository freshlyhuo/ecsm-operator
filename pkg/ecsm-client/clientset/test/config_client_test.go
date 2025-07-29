// file: pkg/ecsm-client/clientset/test/config_client_test.go

package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/fx147/ecsm-operator/pkg/ecsm-client/clientset"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- newTestClientset() 辅助函数 (已在 service_client_test.go 中定义) ---

// TestConfigClient_Lifecycle 测试 Config 资源的完整生命周期 (Create, Get, List, Update, Delete)。
// 这个测试会创建、修改和删除真实的配置项，因此需要一个正在运行的 ECSM 后端。
func TestConfigClient_Lifecycle(t *testing.T) {
	// --- Arrange (Setup) ---
	cs := newTestClientset(t)
	configClient := cs.Configs()
	ctx := context.Background()

	// 为不同类型的配置创建独立的生命周期测试
	runConfigLifecycleTest(t, configClient, ctx, "string", "hello world")
	runConfigLifecycleTest(t, configClient, ctx, "number", 123.45)
	runConfigLifecycleTest(t, configClient, ctx, "json", map[string]interface{}{"user": "admin", "enabled": true})
}

// runConfigLifecycleTest 是一个辅助函数，用于执行单个配置项的完整生命周期测试。
// 这样可以方便地为不同类型的配置 (string, number, json) 复用测试逻辑。
func runConfigLifecycleTest(t *testing.T, client clientset.ConfigInterface, ctx context.Context, configType clientset.ConfigItemType, initialValue interface{}) {
	// 使用 t.Run 将每种类型的测试封装起来，方便阅读测试报告
	t.Run(fmt.Sprintf("LifecycleForType_%s", configType), func(t *testing.T) {
		// --- 1. 创建 (Create) ---
		uniqueKey := fmt.Sprintf("test-config-%s-%d", configType, time.Now().UnixNano())
		createReq := &clientset.CreateConfigRequest{
			Key:   uniqueKey,
			Type:  configType,
			Value: initialValue,
		}

		t.Logf("正在创建配置项, Key: %s, Type: %s", uniqueKey, configType)
		err := client.CreateConfig(ctx, createReq)
		require.NoError(t, err, "创建配置项失败")

		// --- 使用 List 获取刚创建的项的 ID，为后续 Update 和 Delete 做准备 ---
		listOpts := clientset.ListConfigsOptions{PageNum: 1, PageSize: 10, Key: uniqueKey}
		// 注意: 你的 config.go 中函数名为 ListAllConfig
		list, err := client.ListAllConfig(ctx, listOpts)
		require.NoError(t, err, "创建后列出配置项失败")
		require.Len(t, list, 1, "应该只找到一个刚刚创建的配置项")
		createdItem := list[0]
		require.NotEmpty(t, createdItem.ID, "创建的配置项 ID 不能为空")

		// --- 设置延迟清理（关键步骤！） ---
		// 无论后续测试是否失败，defer 都会在函数返回时执行，确保资源被删除。
		defer func() {
			t.Logf("清理中: 正在删除配置项, ID: %s, Key: %s", createdItem.ID, uniqueKey)
			deleteErr := client.DeleteConfig(ctx, createdItem.ID)
			assert.NoError(t, deleteErr, "清理配置项失败")
		}()

		// --- 2. 获取 (Get) ---
		t.Logf("正在获取配置项, Key: %s", uniqueKey)
		retrievedValue, err := client.GetConfig(ctx, uniqueKey)
		require.NoError(t, err, "根据 Key 获取配置项失败")
		// 由于 JSON 数字在 unmarshal 后可能变成 float64，所以我们用 assert.ObjectsAreEqual
		assert.ObjectsAreEqual(initialValue, retrievedValue)

		// --- 3. 更新 (Update) ---
		var updatedValue interface{}
		switch configType {
		case clientset.ConfigItemTypeString:
			updatedValue = "updated hello world"
		case clientset.ConfigItemTypeNumber:
			updatedValue = 999.0
		case clientset.ConfigItemTypeJSON:
			updatedValue = map[string]interface{}{"status": "disabled", "level": 10}
		}
		
		updateReq := &clientset.ConfigItem{
			ID:    createdItem.ID,
			Key:   uniqueKey, // 通常更新时也需要 Key
			Type:  configType,
			Value: updatedValue,
		}

		t.Logf("正在更新配置项, ID: %s", createdItem.ID)
		err = client.UpdateConfig(ctx, updateReq)
		require.NoError(t, err, "更新配置项失败")

		// 验证更新是否成功
		t.Log("验证更新结果...")
		retrievedAfterUpdate, err := client.GetConfig(ctx, uniqueKey)
		require.NoError(t, err, "获取更新后的配置项失败")
		assert.ObjectsAreEqual(updatedValue, retrievedAfterUpdate)
		
		// --- 4. 列出 (List) 和分页测试 ---
		t.Log("正在测试 ListConfig 分页功能...")
		pagedListOpts := clientset.ListConfigsOptions{PageNum: 1, PageSize: 1, Key: uniqueKey}
		pagedList, err := client.ListConfig(ctx, pagedListOpts)
		require.NoError(t, err, "分页列出配置项失败")
		require.Len(t, pagedList, 1, "分页查询应返回一个结果")
		assert.Equal(t, createdItem.ID, pagedList[0].ID)

		// 验证一个不存在的页码
		pagedListOpts.PageNum = 99
		emptyList, err := client.ListConfig(ctx, pagedListOpts)
		require.NoError(t, err, "查询一个不存在的页码不应报错")
		assert.Empty(t, emptyList, "查询一个不存在的页码应返回空列表")
	})
}
