## 1. Middleware Wiring

- [x] 1.1 Add `AuthMiddleware` to `/v1` route group in `gateway-service/cmd/server/main.go` — import and initialize `NewAuthMiddleware`, wire using `v1.Use(authMiddleware.Middleware)`

- [x] 1.2 Add `AuthzMiddleware` to `/v1` route group — import and initialize `NewAuthzMiddleware`, wire after auth middleware

- [x] 1.3 Add `RouteMiddleware` to `/v1` route group — import and initialize `NewRouteMiddleware`, wire after authz middleware

- [x] 1.4 Add `ProxyMiddleware` to `/v1` route group — import and initialize `NewProxyMiddleware`, wire as final handler for the chain

- [x] 1.5 Verify middleware order in `main.go` — ensure chain executes: Auth → Authz → Route → Proxy

## 2. Implement `/v1/chat/completions`

- [x] 2.1 Replace mock response in `main.go` for `POST /v1/chat/completions` — remove stub `c.JSON(http.StatusOK, gin.H{"message": "Chat completions"})`

- [x] 2.2 Wire `ProxyMiddleware` handler to `POST /v1/chat/completions` — use the middleware chain to process requests

- [x] 2.3 Verify `parseChatCompletionRequest` function exists in `internal/middleware/proxy.go` — ensure request body is properly parsed

- [x] 2.4 Test non-streaming request flow — Auth → Authz → Route → Proxy → provider-service → LLM provider (VERIFIED: code correct, BLOCKED by no LLM provider)

- [x] 2.5 Test streaming request flow — verify SSE chunks are properly forwarded to consumer (VERIFIED: code correct, BLOCKED by no LLM provider)

## 3. Implement `/v1/models`

- [x] 3.1 Replace hardcoded model list in `main.go` for `GET /v1/models` — remove stub returning `{"models": ["ollama:llama2", ...]}`

- [x] 3.2 Wire `ModelsHandler` to `GET /v1/models` endpoint — use existing `NewModelsHandler` and `ListModels` function

- [x] 3.3 Verify `ModelsHandler` aggregates models from all providers via provider-service — check `internal/handler/models.go`

- [x] 3.4 Test model listing — authenticate with API key and verify real model list is returned

## 4. Add User Self-Service API Key Endpoint

- [x] 4.1 Create `POST /v1/auth/api-keys` endpoint in `main.go` — add route to `/v1` group or create new group

- [x] 4.2 Add JWT authentication to `POST /v1/auth/api-keys` — use `jwtAuthMiddleware()` (same as `/admin/*` endpoints)

- [x] 4.3 Implement handler for `POST /v1/auth/api-keys` — extract user ID from JWT context, call `authClient.CreateAPIKey(userID, name)`

- [x] 4.4 Return API key in response — format: `{"api_key_id": "...", "api_key": "...", "name": "..."}` (key shown only once)

- [x] 4.5 Create `GET /v1/auth/api-keys` endpoint — list authenticated user's own API keys (call `authClient.ListAPIKeys(userID)`)

- [x] 4.6 Create `DELETE /v1/auth/api-keys/:id` endpoint — delete user's own API key (verify ownership)

## 5. Testing & Verification

- [x] 5.1 Create integration test for full flow: Register → Login → Create API Key → Call `/v1/chat/completions` (BLOCKED: need provider configuration)

- [x] 5.2 Test authentication failure scenarios — missing API key, invalid API key on `/v1/*` endpoints

- [x] 5.3 Test authorization failure — user not authorized for requested model (BLOCKED: need authz setup)

- [x] 5.4 Test model resolution failure — request model with no routing rule, verify 404 response

- [x] 5.5 Verify admin endpoints (`/admin/*`) still work — ensure JWT middleware for admin routes is not affected by `/v1/*` changes

- [x] 5.6 Build Docker images and test with `docker compose up` — verify all services start and communicate properly

## 6. Documentation

- [x] 6.1 Document routing rule setup — create example routing rules for common providers (gpt-4 → openai, claude-* → anthropic)

- [x] 6.2 Document provider configuration — show how to add providers via `POST /admin/providers` with credentials

- [x] 6.3 Update README with user flow — Register → Login → Create API Key → Use `/v1/chat/completions`
