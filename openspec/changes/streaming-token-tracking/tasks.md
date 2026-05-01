## 1. Gateway-Service Streaming Middleware Enhancement

### 1.1 Add streaming token tracking state variables
- [ ] 1.1.1 Add `lastRecordedPromptTokens` and `lastRecordedCompletionTokens` variables to `handleStreamingRequest` function
- [ ] 1.1.2 Add `streamingTokenInterval` config parameter read from configuration (default: 1000)
- [ ] 1.1.3 Add check interval logic in streaming loop: `totalCompletionTokens - lastRecordedCompletionTokens >= streamingTokenInterval`

### 1.2 Implement intermediate usage recording
- [ ] 1.2.1 Calculate delta tokens: `deltaPrompt = totalPromptTokens - lastRecordedPromptTokens`, `deltaCompletion = totalCompletionTokens - lastRecordedCompletionTokens`
- [ ] 1.2.2 Call `billingClient.RecordUsage()` with delta tokens when interval threshold is reached
- [ ] 1.2.3 Update `lastRecordedPromptTokens` and `lastRecordedCompletionTokens` after successful recording
- [ ] 1.2.4 Ensure async `recordUsage` call (use goroutine) to avoid blocking streaming

### 1.3 Implement final usage recording
- [ ] 1.3.1 Calculate final delta tokens after stream completion (remaining tokens not yet recorded)
- [ ] 1.3.2 Call `billingClient.RecordUsage()` with final delta if greater than 0
- [ ] 1.3.3 Handle error case: if stream errors mid-way, still record accumulated tokens

### 1.4 Add configuration parameter
- [ ] 1.4.1 Add `StreamingTokenInterval` field to gateway-service config struct
- [ ] 1.4.2 Read from config YAML: `streaming_token_interval`
- [ ] 1.4.3 Support environment variable override: `STREAMING_TOKEN_INTERVAL`
- [ ] 1.4.4 Validate config: if 0, disable intermediate recording (legacy behavior)

## 2. Billing-Service Aggregation Verification

### 2.1 Verify aggregation handles multiple records
- [ ] 2.1.1 Check `GetUsageAggregation` repository query uses `SUM()` for token aggregation
- [ ] 2.1.2 Verify aggregation by `user_id`, `group_id`, `provider_id`, `model` works correctly
- [ ] 2.1.3 Add test case: multiple usage records from same request are correctly summed

### 2.2 Add unit tests for aggregation
- [ ] 2.2.1 Test scenario: 3 intermediate records + 1 final record from same streaming request
- [ ] 2.2.2 Assert total equals sum of all partial records
- [ ] 2.2.3 Test aggregation by different `group_by` fields (user_id, group_id, provider_id, model)

## 3. Testing

### 3.1 Unit tests for gateway streaming logic
- [ ] 3.1.1 Test interval calculation logic with various token counts
- [ ] 3.1.2 Test delta calculation: verify no double-counting
- [ ] 3.1.3 Mock billing client and verify RecordUsage call count for long streams

### 3.2 Integration tests
- [ ] 3.2.1 Integration test: send streaming request → verify multiple RecordUsage calls
- [ ] 3.2.2 Integration test: query usage aggregation → verify total matches actual tokens
- [ ] 3.2.3 Integration test with configured interval = 100 tokens (smaller for faster testing)

### 3.3 Configuration tests
- [ ] 3.3.1 Test default configuration (interval = 1000)
- [ ] 3.3.2 Test custom configuration via YAML
- [ ] 3.3.3 Test environment variable override
- [ ] 3.3.4 Test interval = 0 disables intermediate recording

## 4. Documentation

### 4.1 Update configuration documentation
- [ ] 4.1.1 Document `streaming_token_interval` parameter in gateway-service README
- [ ] 4.1.2 Add example configuration snippet
- [ ] 4.1.3 Document environment variable override

### 4.2 Update architecture documentation
- [ ] 4.2.1 Document real-time streaming usage tracking in gateway-service architecture
- [ ] 4.2.2 Add sequence diagram for intermediate usage recording
- [ ] 4.2.3 Note on backward compatibility (interval=0 for legacy behavior)

## 5. FVT (Functional Verification Tests)

### 5.1 End-to-end streaming usage test
- [ ] 5.1.1 FVT: Send 5000-token streaming request with interval=1000
- [ ] 5.1.2 Verify 5 intermediate RecordUsage calls via billing query
- [ ] 5.1.3 Verify final aggregation equals 5000 tokens

### 5.2 Configuration FVT
- [ ] 5.2.1 FVT: Change interval via config reload → verify new interval takes effect
- [ ] 5.2.2 FVT: Set interval=0 → verify single RecordUsage call at completion
