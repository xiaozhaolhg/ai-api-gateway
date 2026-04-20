# router-service API Contracts

> Routing service gRPC API

## Service Definition

```protobuf
service RouterService {
  // Route Resolution
  rpc ResolveRoute(ResolveRouteRequest) returns (RouteResult);
  
  // Routing Rule Management
  rpc GetRoutingRules(GetRoutingRulesRequest) returns (ListRoutingRulesResponse);
  rpc CreateRoutingRule(CreateRoutingRuleRequest) returns (RoutingRule);
  rpc UpdateRoutingRule(UpdateRoutingRuleRequest) returns (RoutingRule);
  rpc DeleteRoutingRule(DeleteRoutingRuleRequest) returns (Empty);
  
  // Cache
  rpc RefreshRoutingTable(Empty) returns (Empty);
  
  // Fallback (Phase 2+)
  rpc ResolveFallback(ResolveFallbackRequest) returns (RouteResult);
}
```

## Request/Response Messages

### Route Resolution

```protobuf
message ResolveRouteRequest {
  string model = 1;
  repeated string authorized_models = 2;    // from auth-service
}

message RouteResult {
  string provider_id = 1;
  string adapter_type = 2;               // "openai" | "anthropic" | "gemini"
  repeated string fallback_provider_ids = 3;
}
```

### Routing Rule

```protobuf
message RoutingRule {
  string id = 1;
  string model_pattern = 2;     // e.g., "gpt-4*" or "claude-*"
  string provider_id = 3;
  int32 priority = 4;
  string fallback_provider_id = 5;
}
```

### Fallback (Phase 2+)

```protobuf
message ResolveFallbackRequest {
  string model = 1;
  string failed_provider_id = 2;
}
```