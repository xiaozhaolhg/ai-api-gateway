# monitor-service API Contracts

> Observability and monitoring service gRPC API

## Service Definition

```protobuf
service MonitorService {
  // Callback (subscriber)
  rpc OnProviderResponse(ProviderResponseCallback) returns (Empty);
  
  // Metrics
  rpc RecordMetric(RecordMetricRequest) returns (Empty);
  rpc GetMetrics(GetMetricsRequest) returns (ListMetricsResponse);
  rpc GetMetricAggregation(GetMetricAggregationRequest) returns (MetricAggregationResponse);
  
  // Provider Health
  rpc GetProviderHealth(GetProviderHealthRequest) returns (ProviderHealthStatus);
  rpc ListProviderHealth(ListProviderHealthRequest) returns (ListProviderHealthResponse);
  rpc ReportProviderHealth(ReportProviderHealthRequest) returns (Empty);
  
  // Alert Rules
  rpc CreateAlertRule(CreateAlertRuleRequest) returns (AlertRule);
  rpc UpdateAlertRule(UpdateAlertRuleRequest) returns (AlertRule);
  rpc DeleteAlertRule(DeleteAlertRuleRequest) returns (Empty);
  rpc ListAlertRules(ListAlertRulesRequest) returns (ListAlertRulesResponse);
  
  // Alerts
  rpc GetAlerts(GetAlertsRequest) returns (ListAlertsResponse);
  rpc AcknowledgeAlert(AcknowledgeAlertRequest) returns (Alert);
}
```

## Request/Response Messages

### Metrics

```protobuf
message RecordMetricRequest {
  string metric_type = 1;           // "request_latency" | "error_rate" | "token_throughput"
  map<string, string> labels = 2;     // {"provider": "openai", "model": "gpt-4o"}
  double value = 3;
  int64 timestamp = 4;
}

message Metric {
  string id = 1;
  string type = 2;
  map<string, string> labels = 3;
  double value = 4;
  int64 timestamp = 5;
}

message MetricAggregation {
  string metric_type = 1;
  map<string, string> labels = 2;
  double avg = 3;
  double min = 4;
  double max = 5;
  double p50 = 6;
  double p95 = 7;
  double p99 = 8;
  int64 count = 9;
}
```

### Provider Health

```protobuf
message ProviderHealthStatus {
  string provider_id = 1;
  double latency_p50 = 2;
  double latency_p95 = 3;
  double latency_p99 = 4;
  double error_rate = 5;
  double uptime_pct = 6;
  int64 last_check = 7;
  string status = 8;                // "healthy" | "degraded" | "down"
}

message ReportProviderHealthRequest {
  string provider_id = 1;
  double latency = 2;
  double error_rate = 3;
  double uptime_pct = 4;
  int64 timestamp = 5;
}
```

### Alert Rules

```protobuf
message AlertRule {
  string id = 1;
  string metric_type = 2;
  string condition = 3;              // "gt" | "lt" | "eq"
  double threshold = 4;
  string channel = 5;               // "email" | "webhook" | "slack"
  string channel_config = 6;       // JSON: email, webhook URL, etc.
  string status = 7;               // "active" | "paused"
}
```

### Alerts

```protobuf
message Alert {
  string id = 1;
  string rule_id = 2;
  int64 triggered_at = 3;
  string status = 4;               // "firing" | "acknowledged" | "resolved"
  int64 acknowledged_at = 5;
  int64 resolved_at = 6;
}
```