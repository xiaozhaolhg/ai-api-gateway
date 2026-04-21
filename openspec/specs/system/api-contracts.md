# System API Contracts

> Shared message definitions used across services

## Common Messages

### UserIdentity

Returned by auth-service after API key validation.

```protobuf
message UserIdentity {
  string user_id = 1;
  string role = 2;             // "admin" | "user"
  repeated string group_ids = 3;  // user's group memberships
  repeated string scopes = 4;     // API key scopes
}
```

### AuthorizationResult

Returned by auth-service for model authorization check.

```protobuf
message AuthorizationResult {
  bool allowed = 1;
  string reason = 2;
  repeated string authorized_models = 3;  // models this user/group can access
}
```

### RouteResult

Returned by router-service for route resolution.

```protobuf
message RouteResult {
  string provider_id = 1;
  string adapter_type = 2;   // "openai" | "anthropic" | "gemini"
  repeated string fallback_provider_ids = 3;
}
```

### TokenCounts

Token counts extracted from provider response.

```protobuf
message TokenCounts {
  int64 prompt_tokens = 1;
  int64 completion_tokens = 2;
}
```

### ProviderResponseCallback

Dispatched by provider-service to subscribers after each response.

```protobuf
message ProviderResponseCallback {
  string request_id = 1;
  string user_id = 2;
  string group_id = 3;
  string provider_id = 4;
  string model = 5;
  int64 prompt_tokens = 6;
  int64 completion_tokens = 7;
  int64 latency_ms = 8;
  string status = 9;      // "success" | "error"
  string error_code = 10;
  int64 timestamp = 11;
}
```

### BudgetStatus

Returned by billing-service for budget check.

```protobuf
message BudgetStatus {
  double current_spend = 1;
  double budget_limit = 2;
  double remaining = 3;
  bool soft_cap_exceeded = 4;
  bool hard_cap_exceeded = 5;
}
```

## Common Error Codes

| Code | Description |
|---|---|
| `NOT_FOUND` | Resource not found |
| `UNAUTHORIZED` | Authentication failed |
| `FORBIDDEN` | Authorization failed |
| `INVALID_ARGUMENT` | Invalid request parameters |
| `RESOURCE_EXHAUSTED` | Rate limit or budget exceeded |
| `UNAVAILABLE` | Service unavailable |