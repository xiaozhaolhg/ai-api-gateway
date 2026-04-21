# provider-service API Contracts

> Provider management and request forwarding service gRPC API

## Service Definition

```protobuf
service ProviderService {
  // Request Forwarding
  rpc ForwardRequest(ForwardRequestRequest) returns (ForwardRequestResponse);
  rpc StreamRequest(StreamRequestRequest) returns (stream ProviderChunk);
  
  // Provider Management
  rpc GetProvider(GetProviderRequest) returns (Provider);
  rpc CreateProvider(CreateProviderRequest) returns (Provider);
  rpc UpdateProvider(UpdateProviderRequest) returns (Provider);
  rpc DeleteProvider(DeleteProviderRequest) returns (Empty);
  rpc ListProviders(ListProvidersRequest) returns (ListProvidersResponse);
  rpc ListModels(ListModelsRequest) returns (ListModelsResponse);
  
  // Provider Discovery
  rpc GetProviderByType(GetProviderByTypeRequest) returns (Provider);
  
  // Callback Subscription
  rpc RegisterSubscriber(RegisterSubscriberRequest) returns (Empty);
  rpc UnregisterSubscriber(UnregisterSubscriberRequest) returns (Empty);
}
```

## Request/Response Messages

### Request Forwarding

```protobuf
message ForwardRequestRequest {
  string provider_id = 1;
  bytes request_body = 2;      // JSON-encoded request
  map<string, string> headers = 3;
}

message ForwardRequestResponse {
  bytes response_body = 1;     // JSON-encoded response
  TokenCounts token_counts = 2;
  int32 status_code = 3;
}

message StreamRequestRequest {
  string provider_id = 1;
  bytes request_body = 2;
  map<string, string> headers = 3;
}

message ProviderChunk {
  bytes chunk_data = 1;              // SSE chunk
  TokenCounts accumulated_tokens = 2;
  bool done = 3;
}

message TokenCounts {
  int64 prompt_tokens = 1;
  int64 completion_tokens = 2;
}
```

### Provider

```protobuf
message Provider {
  string id = 1;
  string name = 2;
  string type = 3;               // "openai" | "anthropic" | "gemini" | "custom"
  string base_url = 4;
  string credentials = 5;        // encrypted
  repeated string models = 6;
  string status = 7;             // "active" | "inactive"
  int64 created_at = 8;
  int64 updated_at = 9;
}
```

### Callback Subscription

```protobuf
message RegisterSubscriberRequest {
  string service_name = 1;       // "billing-service" | "monitor-service"
  string callback_endpoint = 2;  // gRPC endpoint
}

message UnregisterSubscriberRequest {
  string service_name = 1;
}
```

### Response Callback (dispatched TO subscribers)

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
  string status = 9;            // "success" | "error"
  string error_code = 10;
  int64 timestamp = 11;
}
```