# Design: Enable go:embed for Single-Binary UI Serving

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                     gateway-service                              │
├─────────────────────────────────────────────────────────────────┤
│  HTTP Entry                                                      │
│  ├── /health (liveness)                                        │
│  ├── /gateway/health (readiness)                                │
│  ├── /v1/chat/completions (with full middleware)                │
│  ├── /admin/* (JWT auth → admin UI pages)                         │
│  └── /* (SPA fallback → embedded admin-ui)  ← ENABLE THIS   │
├─────────────────────────────────────────────────────────────────┤
│  Static File Server                                              │
│  └── go:embed admin-ui/static/ → serve at / with SPA fallback │
└─────────────────────────────────────────────────────────────────┘
```

**Current State**:
- `static/` directory exists with built admin-ui files
- `setupStaticFiles(r *gin.Engine)` function is implemented and called
- `//go:embed static` is **commented out** → static files not embedded

**Target State**:
- `//go:embed static` is active → static files embedded in binary
- Single binary can serve admin-ui without external nginx

## 1. Enable go:embed Directive

### 1.1 Location

**File**: `gateway-service/cmd/server/main.go`  
**Line**: 29

### 1.2 Current Code (Disabled)

```go
// var staticFiles embed.FS // Temporarily disabled for testing
```

### 1.3 Target Code (Enabled)

```go
//go:embed static
var staticFiles embed.FS
```

**Note**: The `//go:embed` directive must be immediately above the variable declaration.

## 2. Static File Serving (Already Implemented)

### 2.1 setupStaticFiles Function

**File**: `gateway-service/cmd/server/main.go`  
**Line**: 306-330 (approximately)

```go
func setupStaticFiles(r *gin.Engine) {
    // Serve static assets
    r.StaticFS("/assets", http.FS(staticFiles))

    // SPA fallback: serve index.html for all non-API, non-static routes
    r.NoRoute(func(c *gin.Context) {
        // Skip API routes
        if strings.HasPrefix(c.Request.URL.Path, "/admin/") ||
            strings.HasPrefix(c.Request.URL.Path, "/v1/") ||
            strings.HasPrefix(c.Request.URL.Path, "/gateway/") ||
            strings.HasPrefix(c.Request.URL.Path, "/health") {
            c.Next()
            return
        }
        data, _ := staticFiles.ReadFile("static/index.html")
        c.Data(200, "text/html; charset=utf-8", data)
    })
}
```

**Status**: ✅ Already implemented and called at line 212. No changes needed.

## 3. Build Pipeline (Already Implemented)

### 3.1 Makefile Targets

**File**: `Makefile` (root)

```makefile
build-single: build-ui embed-ui build-gateway

build-ui:
	cd admin-ui && npm ci && npm run build

embed-ui:
	rm -rf gateway-service/static
	cp -r admin-ui/dist gateway-service/static

build-gateway:
	cd gateway-service && go build -o bin/gateway ./cmd/server
```

**Status**: ✅ Already implemented. No changes needed.

### 3.2 Build Flow

```
1. npm run build (admin-ui) → dist/
2. cp -r admin-ui/dist gateway-service/static/
3. go build (with //go:embed static) → single binary
4. Binary contains embedded static files
```

## 4. Verification Strategy

### 4.1 Unit Test: Static Files Embedded

```go
// Test: Verify static files are embedded
func TestStaticFilesEmbedded(t *testing.T) {
    // After enabling go:embed, staticFiles should contain files
    files, err := fs.Glob(staticFiles, "static/*")
    if err != nil {
        t.Fatalf("Failed to glob static files: %v", err)
    }
    if len(files) == 0 {
        t.Error("No static files embedded")
    }
}
```

### 4.2 Integration Test: UI Accessible

```bash
# Build single binary
make build-single

# Start binary
./gateway-service/bin/gateway &

# Test UI is served
curl http://localhost:8080/
# Should return HTML (index.html)

curl http://localhost:8080/assets/index-*.js
# Should return JavaScript file

# Test API routes still work
curl http://localhost:8080/health
# Should return JSON health status
```

## 5. File Structure

```
gateway-service/
├── cmd/server/
│   └── main.go          # UPDATE: uncomment go:embed (line 29)
├── static/                 # EXISTS: embedded admin-ui build output
│   ├── index.html
│   ├── assets/
│   │   ├── index-*.js
│   │   └── index-*.css
│   └── favicon.svg
└── (other files unchanged)
```

## 6. Rationale

| Approach | Pros | Cons |
|----------|------|------|
| **Enable go:embed (chosen)** | Single binary, no nginx needed, simpler deployment | Requires rebuilding after UI changes |
| Keep disabled + nginx | Separate concerns, hot-reload UI | More complex deployment, extra container |
| **Decision**: Enable go:embed to fulfill Phase 1 Work Division requirement and simplify deployment.

## 7. Backward Compatibility

- **Development mode**: If `static/` directory is empty, `setupStaticFiles()` should handle gracefully (already implemented)
- **Fallback**: No fallback needed — either embed works or build fails (fail-fast)
- **Hot-reload**: For development, run `npm run dev` in `admin-ui/` separately (not affected by this change)
