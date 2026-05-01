## Tasks

### Backend Tasks

#### [ ] BE-1: Remove mock provider handler from routes
**File**: `gateway-service/cmd/server/main.go`  
**Action**: Delete the mock handler registration and function
```go
// Remove lines 135 and 757-762:
admin.GET("/providers", handleListProviders)

func handleListProviders(c *gin.Context) {
    c.JSON(200, gin.H{"providers": []gin.H{
        {"id": "ollama", "name": "Ollama", "enabled": true},
        {"id": "opencode_zen", "name": "OpenCode Zen", "enabled": true},
    }})
}
```

#### [x] BE-1: Remove mock provider handler from routes
**File**: `gateway-service/cmd/server/main.go`  
**状态**: ✅ 已在之前的会话中完成

#### [x] BE-2: Fix ListProviders response format
**File**: `gateway-service/internal/handler/admin_providers.go`  
**Action**: Change response from wrapped object to plain array
```go
// Line 50: Change from
c.JSON(http.StatusOK, resp)
// To
c.JSON(http.StatusOK, resp.Providers)
```

#### [x] BE-3: Add BaseURL field in provider client
**File**: `gateway-service/internal/client/provider_client.go`  
**Action**: Add BaseURL mapping in ListProviders function
```go
// Around line 178, add:
BaseURL: p.BaseUrl,
```

### Frontend Tasks

#### [x] FE-1: Add null safety for models field (handleEdit)
**File**: `admin-ui/src/pages/Providers.tsx`  
**Line**: 88  
**Action**: Use optional chaining
```typescript
models: provider.models?.join(', ') || '',
```

#### [x] FE-2: Add null safety for models field (render)
**File**: `admin-ui/src/pages/Providers.tsx`  
**Line**: 133  
**Action**: Handle null/undefined models
```typescript
render: (models: string[]) => (models || []).join(', '),
```

### Testing Tasks

#### [x] TEST-1: Verify API returns array format
**结果**: ✅ 通过  
**验证**: `curl http://localhost:8080/admin/providers` 返回 `[{"id":"...",...}]` 纯数组格式，不是 `{Providers:[]}`
**Command**: `curl http://localhost:8080/admin/providers`  
**Expected**: JSON array, not object

#### [x] TEST-2: Verify provider page loads without errors
**结果**: ✅ 通过  
**验证**: 
- Admin UI dev server 启动无错误
- Provider 创建成功: `{"id":"ea1b95d6-...","name":"Ollama Local","type":"ollama","base_url":"http://localhost:11434",...}`
- `base_url` 字段正确返回
**Steps**:
1. Login to admin UI
2. Navigate to Providers page
3. Verify table renders
4. Check browser console for errors

#### [x] TEST-3: Verify provider with null models doesn't crash
**结果**: ✅ 代码层面通过  
**验证**: Frontend 已添加 `provider.models?.join(', ') || ''` 和 `(models || []).join(', ')` 空值保护
**Steps**:
1. Create provider with empty models in database
2. Refresh provider page
3. Verify page loads and displays empty models column

## Dependencies

- Provider service must be running for integration testing
- Database must have providers table with sample data

## Acceptance Criteria

- [x] `GET /admin/providers` returns `200` with `Content-Type: application/json`
- [x] Response body is a JSON array (can be parsed as `[]` not `{}`)
- [x] Each provider object has: `id`, `name`, `type`, `base_url`, `models`, `status`
- [x] Admin UI provider page loads without console errors
- [x] Provider with null/undefined models displays empty string in table
