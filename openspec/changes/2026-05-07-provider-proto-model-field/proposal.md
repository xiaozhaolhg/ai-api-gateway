# Proposal: Provider Proto Model Field Enhancement

## Problem Statement

The current `provider.proto` (v1) lacks a dedicated `model` field in request/response messages, causing:

1. **Type-unsafe model passing**: Gateway-service extracts model from JSON body via string parsing instead of typed proto fields
2. **Inconsistent model references**: `Provider` message has `repeated string models` (available models), but request/response messages lack the selected `model` field
3. **Billing inaccuracy risk**: Billing-service relies on gateway to parse model from JSON, introducing coupling and potential errors
4. **Proto versioning gap**: No forward-compatible way to pass model information in gRPC calls

The four affected messages are:
- `ForwardRequestRequest` (line 34-38) — missing selected model
- `StreamRequestRequest` (line 46-50) — missing selected model
- `ForwardRequestResponse` (line 40-44) — missing model in response
- `ProviderChunk` (line 52-56) — missing model in streaming chunks

## Proposed Solution

Add a `string model` field to all four messages in `provider.proto`, regenerate Go code with `buf generate`, then update gateway-service, provider-service, and billing-service to use the typed field instead of JSON parsing.

## Scope

- **In Scope**:
  - Add `model` field to `ForwardRequestRequest`, `StreamRequestRequest`, `ForwardRequestResponse`, `ProviderChunk`
  - Regenerate proto with `buf generate`
  - Update gateway-service to extract model from proto request (remove JSON parsing dependency)
  - Update provider-service to return model in `ForwardRequestResponse` and `ProviderChunk`
  - Update billing-service call to use model from proto response

- **Out of Scope**:
  - Changes to other proto files (auth, router, billing, monitor)
  - Database schema changes
  - New provider adapter implementations
  - Frontend changes (admin-ui uses REST API, unaffected)

## Success Criteria

- `buf generate` completes without errors after proto changes
- Gateway-service `ForwardRequest` handler extracts model from `req.Model` (not JSON parsing)
- Provider-service includes model in `ForwardRequestResponse` and `ProviderChunk`
- Billing-service `RecordUsage` calls use model from proto response
- All existing tests pass after changes
- No JSON parsing dependency for model extraction in gateway-service

## Dependencies

- **buf CLI**: Must be installed and configured (`buf.gen.yaml` exists)
- **gateway-service**: Depends on regenerated `providerv1` package
- **provider-service**: Depends on regenerated `providerv1` package
- **billing-service**: Depends on updated gateway-service response handling

## Owner

- **Primary**: Developer B (cynkiller)
- **Collaborators**: Developer A (gateway-service wiring), Developer C (testing)

## References

- Original task: `docs/phase1_work_division.md` → Developer B → Week 4 → Proto Schema Improvements
- Current proto: `api/proto/provider/v1/provider.proto`
- Related: `api/proto/common/v1/common.proto` (TokenCounts message)
