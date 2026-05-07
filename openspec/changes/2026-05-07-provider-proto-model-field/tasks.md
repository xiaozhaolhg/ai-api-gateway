# Tasks: Provider Proto Model Field Enhancement

**Owner**: Developer B (cynkiller)  
**Collaborators**: Developer A (gateway-service wiring), Developer C (testing)  
**Status**: Planning

## Phase 1: Proto Modification & Code Generation (High Priority)

- [ ] **Task 1.1**: Update `provider.proto` with model fields
  - [ ] Add `string model = 4` to `ForwardRequestRequest` message
  - [ ] Add `string model = 4` to `StreamRequestRequest` message
  - [ ] Add `string model = 4` to `ForwardRequestResponse` message
  - [ ] Add `string model = 4` to `ProviderChunk` message
  - **Acceptance**: `buf lint api/proto/provider/v1/provider.proto` passes without errors

- [ ] **Task 1.2**: Regenerate Go code with buf
  - [ ] Run `cd api/proto && buf generate`
  - [ ] Verify `api/gen/provider/v1/provider.pb.go` contains Model fields
  - [ ] Verify field numbers are correct (field 4 for all new fields)
  - **Acceptance**: Generated code compiles; `grep -c "Model " api/gen/provider/v1/provider.pb.go` returns 4+

## Phase 2: Gateway-Service Updates (High Priority)

- [ ] **Task 2.1**: Update proxy middleware to use typed model field
  - [ ] Modify `gateway-service/internal/middleware/proxy.go` (or equivalent)
  - [ ] Extract model from request JSON → set `forwardReq.Model = model`
  - [ ] Remove JSON parsing dependency for model extraction
  - [ ] Pass `req.Model` to billing call instead of parsing response JSON
  - **Acceptance**: Gateway compiles; model passed via proto field (not JSON)

- [ ] **Task 2.2**: Update streaming handler
  - [ ] Modify streaming request to include `StreamRequestRequest.Model`
  - [ ] Update billing intermediate recording to use `req.Model`
  - **Acceptance**: Streaming requests include model in proto; billing records correct model

## Phase 3: Provider-Service Updates (High Priority)

- [ ] **Task 3.1**: Return model in ForwardRequestResponse
  - [ ] Modify `provider-service/internal/handler/grpc_handler.go`
  - [ ] Set `resp.Model = req.Model` in `ForwardRequest` handler
  - **Acceptance**: Response contains model field; verify with grpcurl

- [ ] **Task 3.2**: Return model in ProviderChunk (streaming)
  - [ ] Set `chunk.Model = req.Model` in `StreamRequest` handler
  - [ ] Ensure model is set in every chunk (including final chunk)
  - **Acceptance**: Streaming chunks contain model field; verify with grpcurl

## Phase 4: Billing-Service Integration (Medium Priority)

- [ ] **Task 4.1**: Update billing call for non-streaming requests
  - [ ] Modify gateway-service billing call location
  - [ ] Use `forwardResp.Model` from proto response (not JSON parsing)
  - [ ] Handle case where model field is empty (fallback to "unknown")
  - **Acceptance**: Billing records show correct model from proto field

- [ ] **Task 4.2**: Update billing call for streaming requests
  - [ ] Store model in stream context when initiating stream
  - [ ] Use stored model (from `StreamRequestRequest`) for intermediate billing records
  - [ ] Use `chunk.Model` from final chunk (verify consistency)
  - [ ] Implement warning if `chunk.Model != stored_model`
  - **Acceptance**: All streaming billing records use consistent model; interval billing works correctly

## Phase 5: Testing & Verification (High Priority)

- [ ] **Task 5.1**: Unit tests for proto changes
  - [ ] Write test: create `ForwardRequestRequest` with model field
  - [ ] Write test: create `StreamRequestRequest` with model field
  - [ ] Write test: verify `ForwardRequestResponse.Model` is set
  - [ ] Write test: verify `ProviderChunk.Model` is set
  - **Acceptance**: All new tests pass; >80% coverage for proto struct usage

- [ ] **Task 5.2**: Integration test - end-to-end model passing (non-streaming)
  - [ ] Test: Send request with `"model": "ollama:llama2"`
  - [ ] Verify model appears in provider-service logs (via proto)
  - [ ] Verify billing-service records show `model: "ollama:llama2"`
  - **Acceptance**: Full flow: request → gateway → provider → billing all have correct model

- [ ] **Task 5.2b**: Integration test - streaming model consistency
  - [ ] Test: Send streaming request with `"model": "gpt-4"`
  - [ ] Verify intermediate billing records (every N tokens) use `gpt-4`
  - [ ] Verify final billing record uses `gpt-4`
  - [ ] Verify `ProviderChunk.Model` consistency across all chunks
  - **Acceptance**: Streaming flow: all billing records use consistent model

- [ ] **Task 5.3**: Backward compatibility verification
  - [ ] Test: Old client (no model field) → new server → provider uses JSON fallback
  - [ ] Test: New client (with model field) → new server → model passed via proto
  - [ ] Test: Streaming with old client → provider extracts model from first request body
  - [ ] Test: New client → old server (if test environment available) → graceful degradation
  - [ ] Verify: Empty model string handled gracefully (not panics, uses fallback)
  - **Acceptance**: No breaking changes; all scenarios handled gracefully

## Phase 6: Documentation & Cleanup (Low Priority)

- [ ] **Task 6.1**: Update API documentation
  - [ ] Update `openspec/specs/provider-service/spec.md` with new proto fields
  - [ ] Update architecture diagrams if needed
  - **Acceptance**: Documentation reflects new model field in all messages

- [ ] **Task 6.2**: Code cleanup
  - [ ] Remove deprecated JSON parsing code for model extraction
  - [ ] Search for TODO comments related to model parsing
  - [ ] Update comments to reflect type-safe model passing
  - **Acceptance**: No leftover JSON parsing for model; code comments are accurate

## Summary

| Phase | Tasks | Priority | Dependencies |
|-------|-------|----------|--------------|
| Phase 1: Proto Modification | 2 | **High** | None |
| Phase 2: Gateway-Service | 2 | **High** | Phase 1 |
| Phase 3: Provider-Service | 2 | **High** | Phase 1 |
| Phase 4: Billing Integration | 2 | **Medium** | Phase 2, 3 |
| Phase 5: Testing | 5 | **High** | Phase 1, 2, 3, 4 |
| Phase 6: Documentation | 2 | **Low** | Phase 1-5 |
| **Total** | **15** | | |

## Timeline Estimate

| Phase | Estimated Time | Owner |
|-------|----------------|--------|
| Phase 1 | 1 hour | Developer B |
| Phase 2 | 2 hours | Developer B + Dev A |
| Phase 3 | 2 hours | Developer B |
| Phase 4 | 1 hour | Developer B + Dev A |
| Phase 5 | 5 hours | Developer B + Dev C | 新增流式测试和向后兼容验证 |
| Phase 6 | 1 hour | Developer B |
| **Total** | **10 hours** | |

## Critical Path

```
Phase 1 (Proto) → Phase 2 (Gateway) → Phase 5.2 (Integration Test)
                            ↓
                    Phase 3 (Provider) → Phase 4 (Billing)
```

**Blocker**: Phase 5.2 (integration test) cannot start until Phases 1-4 are complete.
