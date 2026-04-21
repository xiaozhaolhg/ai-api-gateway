# billing-service API Contracts

> Usage and billing service gRPC API

## Service Definition

```protobuf
service BillingService {
  // Callback (subscriber)
  rpc OnProviderResponse(ProviderResponseCallback) returns (Empty);
  
  // Direct Usage Recording
  rpc RecordUsage(RecordUsageRequest) returns (Empty);
  
  // Usage Queries
  rpc GetUsage(GetUsageRequest) returns (ListUsageResponse);
  rpc GetUsageAggregation(GetUsageAggregationRequest) returns (ListUsageAggregationResponse);
  
  // Cost Estimation
  rpc EstimateCost(EstimateCostRequest) returns (CostEstimate);
  
  // Budget
  rpc CheckBudget(CheckBudgetRequest) returns (BudgetStatus);
  rpc CreateBudget(CreateBudgetRequest) returns (Budget);
  rpc UpdateBudget(UpdateBudgetRequest) returns (Budget);
  rpc DeleteBudget(DeleteBudgetRequest) returns (Empty);
  rpc ListBudgets(ListBudgetsRequest) returns (ListBudgetsResponse);
  
  // Pricing
  rpc CreatePricingRule(CreatePricingRuleRequest) returns (PricingRule);
  rpc UpdatePricingRule(UpdatePricingRuleRequest) returns (PricingRule);
  rpc DeletePricingRule(DeletePricingRuleRequest) returns (Empty);
  rpc ListPricingRules(ListPricingRulesRequest) returns (ListPricingRulesResponse);
  
  // Billing Account
  rpc GetBillingAccount(GetBillingAccountRequest) returns (BillingAccount);
  rpc CreateBillingAccount(CreateBillingAccountRequest) returns (BillingAccount);
  rpc UpdateBillingAccount(UpdateBillingAccountRequest) returns (BillingAccount);
  
  // Invoice (Phase 3+)
  rpc GenerateInvoice(GenerateInvoiceRequest) returns (Invoice);
  rpc GetInvoices(GetInvoicesRequest) returns (ListInvoicesResponse);
}
```

## Request/Response Messages

### Usage Recording

```protobuf
message RecordUsageRequest {
  string user_id = 1;
  string group_id = 2;
  string provider_id = 3;
  string model = 4;
  int64 prompt_tokens = 5;
  int64 completion_tokens = 6;
}

message UsageRecord {
  string id = 1;
  string user_id = 2;
  string group_id = 3;
  string provider_id = 4;
  string model = 5;
  int64 prompt_tokens = 6;
  int64 completion_tokens = 7;
  double cost = 8;
  int64 timestamp = 9;
}
```

### Usage Aggregation

```protobuf
message UsageAggregation {
  string group_key = 1;              // user_id, model, or provider_id
  int64 total_prompt_tokens = 2;
  int64 total_completion_tokens = 3;
  double total_cost = 4;
  int64 request_count = 5;
}
```

### Cost Estimation

```protobuf
message EstimateCostRequest {
  string model = 1;
  int64 prompt_tokens = 2;
  int64 completion_tokens = 3;
}

message CostEstimate {
  double estimated_cost = 1;
  string currency = 2;
  double price_per_prompt_token = 3;
  double price_per_completion_token = 4;
}
```

### Budget

```protobuf
message BudgetStatus {
  double current_spend = 1;
  double budget_limit = 2;
  double remaining = 3;
  double soft_cap_pct = 4;
  double hard_cap_pct = 5;
  bool soft_cap_exceeded = 6;
  bool hard_cap_exceeded = 7;
}

message Budget {
  string id = 1;
  string account_id = 2;
  double limit = 3;
  string period = 4;           // "daily" | "weekly" | "monthly"
  double soft_cap_pct = 5;
  double hard_cap_pct = 6;
  string status = 7;             // "active" | "paused"
}
```

### Pricing Rule

```protobuf
message PricingRule {
  string id = 1;
  string model = 2;
  string provider_id = 3;
  double price_per_prompt_token = 4;
  double price_per_completion_token = 5;
  string currency = 6;
}
```

### Invoice (Phase 3+)

```protobuf
message Invoice {
  string id = 1;
  string account_id = 2;
  int64 period_start = 3;
  int64 period_end = 4;
  double total_cost = 5;
  repeated InvoiceLineItem line_items = 6;
  string status = 7;           // "draft" | "finalized" | "paid"
}
```