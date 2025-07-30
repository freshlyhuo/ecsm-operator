// file: pkg/ecsm-client/clientset/test/template_test.go

package test

import (
	"context"
	"testing"
	"time"

	"github.com/fx147/ecsm-operator/pkg/ecsm-client/clientset"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- newTestClientset() 辅助函数 (已在其他测试文件中定义) ---

// TestTemplateClient_Lifecycle 测试模板和目录的完整生命周期：
// Create -> Get -> Update -> Move -> Search -> Delete
// 这是一个集成的写操作测试，会真实地在ECSM中创建和删除资源。
func TestTemplateClient_Lifecycle(t *testing.T) {
	// --- Arrange (Setup) ---
	cs := newTestClientset(t)
	templateClient := cs.Templates()
	ctx := context.Background()

	basePath := "/autotest/" + time.Now().Format("20060102150405")
	dirPath := basePath + "/test-dir"
	tmplPath := dirPath + "/test-tmpl"
	moveDestPath := basePath + "/moved-dir"

	var dirID, tmplID, moveDirID string

	// --- 1. 创建目录 ---
	t.Run("CreateDirectory", func(t *testing.T) {
		createDirReq := &clientset.CreateDictoryRequest{
			DictoryName: "test-dir",
			DictoryPath: basePath,
		}
		t.Logf("正在创建目录: %s", dirPath)
		resp, err := templateClient.CreateDictory(ctx, createDirReq)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotEmpty(t, resp.DictoryID)
		dirID = resp.DictoryID

		// **关键**: 使用 defer 确保目录最终被删除
		defer func() {
			t.Logf("清理中: 正在删除目录 (ID: %s)", dirID)
			// 注意：此时目录可能已被移动，所以删除时最好用ID
			_, err := templateClient.DeleteTempOrDictByID(ctx, dirID)
			assert.NoError(t, err, "清理目录失败")
		}()
	})

	// --- 2. 创建模板 ---
	t.Run("CreateTemplate", func(t *testing.T) {
		createTmplReq := &clientset.CreateTemplateRequest{
			ImageRefs: []string{"test-image@1.0.0#sylixos"}, // 使用一个有效的镜像引用
			Path:      dirPath,
		}
		t.Logf("正在创建模板: %s", tmplPath)
		resp, err := templateClient.CreateTemplate(ctx, createTmplReq)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, resp.ProvsionTmplList, 1)
		tmplID = resp.ProvsionTmplList[0].ID
		require.NotEmpty(t, tmplID)

		// **关键**: 使用 defer 确保模板最终被删除
		defer func() {
			t.Logf("清理中: 正在删除模板 (ID: %s)", tmplID)
			_, err := templateClient.DeleteTempOrDictByID(ctx, tmplID)
			assert.NoError(t, err, "清理模板失败")
		}()
	})

	// 等待一小段时间，确保资源在系统中完全可见
	time.Sleep(1 * time.Second)

	// --- 3. 获取模板 (By ID 和 By Path) ---
	t.Run("GetTemplateByIDAndPath", func(t *testing.T) {
		// By ID
		t.Logf("正在通过 ID 获取模板: %s", tmplID)
		tmplByID, err := templateClient.GetTemplateByID(ctx, tmplID)
		require.NoError(t, err)
		require.NotNil(t, tmplByID)
		assert.Equal(t, "test-tmpl", tmplByID.Name)

		// By Path
		t.Logf("正在通过 Path 获取模板: %s", tmplPath)
		tmplByPath, err := templateClient.GetTemplateByPath(ctx, tmplPath)
		require.NoError(t, err)
		require.NotNil(t, tmplByPath)
		assert.Equal(t, tmplID, tmplByPath.ID)
	})

	// --- 4. 批量更新 (By ID 和 By Name) ---
	t.Run("UpdateTemplatesBatch", func(t *testing.T) {
		// By ID
		updateByIDReq := &clientset.UpdateTemplatesBatchRequest{
			Templates: []clientset.TemplateUpdateBatchSpec{
				{ID: tmplID, ImageRef: "test-image@2.0.0#sylixos"},
			},
			Action: "load",
		}
		t.Logf("正在通过 ID 批量更新模板: %s", tmplID)
		results, err := templateClient.UpdateTemplatesByID(ctx, updateByIDReq)
		require.NoError(t, err)
		require.Len(t, results, 1)
		assert.True(t, results[0].Result)
		assert.Equal(t, tmplID, results[0].ID)

		// By Name (Name 是模板的全路径)
		updateByNameReq := &clientset.UpdateTemplatesBatchRequest{
			Templates: []clientset.TemplateUpdateBatchSpec{
				{Name: tmplPath, ImageRef: "test-image@3.0.0#sylixos"},
			},
			Action: "load",
		}
		t.Logf("正在通过 Name 批量更新模板: %s", tmplPath)
		results, err = templateClient.UpdateTemplatesByName(ctx, updateByNameReq)
		require.NoError(t, err)
		require.Len(t, results, 1)
		assert.True(t, results[0].Result)
		assert.Equal(t, tmplPath, results[0].Name)
	})

	// --- 5. 移动模板 ---
	t.Run("MoveTemplate", func(t *testing.T) {
		// 创建目标目录
		createMoveDirReq := &clientset.CreateDictoryRequest{
			DictoryName: "moved-dir",
			DictoryPath: basePath,
		}
		t.Logf("正在创建用于移动的目标目录: %s", moveDestPath)
		resp, err := templateClient.CreateDictory(ctx, createMoveDirReq)
		require.NoError(t, err)
		moveDirID = resp.DictoryID
		defer func() {
			t.Logf("清理中: 正在删除移动目标目录 (ID: %s)", moveDirID)
			_, err := templateClient.DeleteTempOrDictByID(ctx, moveDirID)
			assert.NoError(t, err, "清理移动目标目录失败")
		}()

		// 执行移动
		moveReq := &clientset.MoveRequest{
			Src: tmplPath,
			Dst: moveDestPath,
		}
		t.Logf("正在移动模板从 '%s' 到 '%s'", tmplPath, moveDestPath)
		_, err = templateClient.MoveTempOrDict(ctx, moveReq)
		require.NoError(t, err)

		// 验证移动是否成功
		movedTmplPath := moveDestPath + "/test-tmpl"
		t.Logf("验证移动结果，获取新路径: %s", movedTmplPath)
		tmplAfterMove, err := templateClient.GetTemplateByPath(ctx, movedTmplPath)
		require.NoError(t, err)
		assert.Equal(t, tmplID, tmplAfterMove.ID)
	})

	// --- 6. 搜索和获取树形结构 ---
	t.Run("SearchAndGetTree", func(t *testing.T) {
		// 搜索
		searchOpts := clientset.SearchTemplateOptions{
			Key:  "test-tmpl",
			Path: moveDestPath, // 在移动后的新路径里搜索
			Kind: "service",
		}
		t.Logf("正在搜索模板, Key: %s, Path: %s", searchOpts.Key, searchOpts.Path)
		searchResults, err := templateClient.SearchTempOrDict(ctx, searchOpts)
		// 注意: SearchTempOrDict API 设计似乎返回单个结果，而非列表。
		// 我们将基于此进行断言。
		require.NoError(t, err)
		require.NotNil(t, searchResults)
		assert.Equal(t, tmplID, searchResults.ID)

		// 获取树
		treeOpts := clientset.GetTemplateTreeOptions{
			Path:  basePath,
			Level: 2,
			Model: "full",
		}
		t.Logf("正在获取模板树, Path: %s", basePath)
		tree, err := templateClient.GetTemplateTree(ctx, treeOpts)
		require.NoError(t, err)
		require.NotNil(t, tree)

		// 验证树的结构
		require.Contains(t, tree.Children, "moved-dir")
		movedDirNode := tree.Children["moved-dir"]
		require.NotNil(t, movedDirNode)
		require.Contains(t, movedDirNode.Children, "test-tmpl")
		tmplNode := movedDirNode.Children["test-tmpl"]
		require.NotNil(t, tmplNode)
		require.NotNil(t, tmplNode.Data)
		assert.Equal(t, tmplID, tmplNode.Data.ID)
	})

	// --- 7. 批量删除 (By IDs) ---
	t.Run("DeleteByIDs", func(t *testing.T) {
		// 创建两个临时目录用于批量删除测试
		dir1Req := &clientset.CreateDictoryRequest{DictoryName: "del-1", DictoryPath: basePath}
		dir2Req := &clientset.CreateDictoryRequest{DictoryName: "del-2", DictoryPath: basePath}

		resp1, err := templateClient.CreateDictory(ctx, dir1Req)
		require.NoError(t, err)
		resp2, err := templateClient.CreateDictory(ctx, dir2Req)
		require.NoError(t, err)

		idsToDelete := []string{resp1.DictoryID, resp2.DictoryID}
		t.Logf("正在通过 IDs 批量删除目录: %v", idsToDelete)

		// 注意: 您的 DeleteTempOrDictByIDs 实现中，请求体和响应体共用了一个结构体，
		// 并且 json tag 可能是 `id` 而不是 `ids`。此测试基于您的代码实现。
		// 如果 API 实际需要 `{"ids": [...]}`，客户端代码可能需要修复。
		_, err = templateClient.DeleteTempOrDictByIDs(ctx, idsToDelete)
		require.NoError(t, err)

		// 验证删除
		_, err1 := templateClient.GetTemplateByID(ctx, resp1.DictoryID)
		_, err2 := templateClient.GetTemplateByID(ctx, resp2.DictoryID)
		assert.Error(t, err1, "第一个批量删除的目录不应再存在")
		assert.Error(t, err2, "第二个批量删除的目录不应再存在")
	})

	// --- 8. 单个删除 (By Path) ---
	// Defer 已经测试了 DeleteTempOrDictByID，这里我们测试 DeleteTempOrDict (by path)
	t.Run("DeleteByPath", func(t *testing.T) {
		// 创建一个临时目录
		dirReq := &clientset.CreateDictoryRequest{DictoryName: "deleteme-by-path", DictoryPath: basePath}
		resp, err := templateClient.CreateDictory(ctx, dirReq)
		require.NoError(t, err)

		pathToDel := basePath + "/deleteme-by-path"
		t.Logf("正在通过 path 删除目录: %s", pathToDel)

		_, err = templateClient.DeleteTempOrDict(ctx, pathToDel)
		require.NoError(t, err)

		// 验证删除
		_, err = templateClient.GetTemplateByID(ctx, resp.DictoryID)
		assert.Error(t, err, "通过路径删除的目录不应再存在")
	})
}
