## 1. Directory Creation

- [x] 1.1 Create openspec/specs/system/ directory
- [x] 1.2 Create openspec/specs/gateway-service/ directory
- [x] 1.3 Create openspec/specs/auth-service/ directory
- [x] 1.4 Create openspec/specs/router-service/ directory
- [x] 1.5 Create openspec/specs/provider-service/ directory
- [x] 1.6 Create openspec/specs/billing-service/ directory
- [x] 1.7 Create openspec/specs/monitor-service/ directory

## 2. System Specs

- [x] 2.1 Create system/spec.md (service relationships diagram, calling matrix)
- [x] 2.2 Create system/api-contracts.md (shared message definitions)

## 3. Service Specs

- [x] 3.1 Create gateway-service/spec.md (calls to auth, router, provider, billing, monitor)
- [x] 3.2 Create auth-service/spec.md (called by gateway)
- [x] 3.3 Create router-service/spec.md (reference to existing specs)
- [x] 3.4 Create provider-service/spec.md (calls to external providers, callbacks)
- [x] 3.5 Create billing-service/spec.md (receives callbacks)
- [x] 3.6 Create monitor-service/spec.md (receives callbacks)

## 4. Service API Contracts

- [x] 4.1 Create gateway-service/api-contracts.md (HTTP endpoints, middleware flow)
- [x] 4.2 Create auth-service/api-contracts.md (ValidateAPIKey, CheckModelAuthorization, User/Group/APIKey CRUD)
- [x] 4.3 Create router-service/api-contracts.md (ResolveRoute, RoutingRule CRUD)
- [x] 4.4 Create provider-service/api-contracts.md (ForwardRequest, StreamRequest, Provider CRUD, Callback registration)
- [x] 4.5 Create billing-service/api-contracts.md (OnProviderResponse, CheckBudget, RecordUsage, Usage/Budget CRUD)
- [x] 4.6 Create monitor-service/api-contracts.md (OnProviderResponse, RecordMetric, Health, Alert CRUD)

## 5. Rename Files

- [x] 5.1 Rename architecture.md → spec.md in all service folders

## 6. Validation

- [x] 6.1 Verify all directories exist
- [x] 6.2 Verify system/spec.md has calling diagram
- [x] 6.3 Verify each service has Calls To / Called By sections
- [x] 6.4 Verify each service has api-contracts.md with RPC definitions