## Context

The AI API Gateway currently supports system-wide routing rules configured by admins via `/admin/routing-rules`. Regular users cannot configure per-user routing preferences or fallback chains. The router-service uses a simple RoutingRule model without user association, and the gateway-service only exposes admin APIs for routing configuration.

**Current Architecture:**
```
gateway-service (:8080) → router-service (:50052) → provider-service (:50053)
                                ↑
                    System-wide rules only
                    No user_id in RoutingRule
```

**Goal:** Enable users like Harry to configure per-user routing rules with fallback chains, where user rules OVERRIDE system rules entirely.

## Goals / Non-Goals

**Goals:**
- Allow users to create/read/update/delete their own routing rules via `/v1/routing-rules` (JWT auth)
- Allow admins to configure per-user routing rules via `/admin/users/{userId}/routing-rules`
- User rules OVERRIDE system rules entirely (not merge)
- Support ordered fallback provider chains (automatic failover on 5xx, timeouts, error codes)
- Extend RoutingRule model with `user_id` and `fallback_provider_ids`

**Non-Goals:**
- Group-level routing rules (Phase 2+)
- UI implementation in admin-ui (API-only for now)
- Complex fallback logic (e.g., different fallback per error type)
- Rate limiting per user/provider (separate concern)

## Decisions

### Decision 1: Rule Resolution Strategy

**Choice:** User rules OVERRIDE system rules entirely

**Rationale:**
- User expects their configuration to be honored completely
- Simpler mental model: "My rules = my preferences"
- Avoids complex merging logic

**Alternative considered:** Merge system + user rules
- Rejected because it creates ambiguity about which rule applies
- More complex to implement and debug

**Implementation:**
```go
// In router-service ResolveRoute:
func (s *Service) ResolveRoute(ctx context.Context, req *pb.ResolveRouteRequest) (*pb.RouteResult, error) {
    // 1. Try user-specific rules first (if user_id provided)
    if req.UserId != "" {
        userRule := s.repo.FindRuleByUserAndModel(req.UserId, req.Model)
        if userRule != nil {
            return s.buildRouteResult(userRule), nil  // User rule OVERRIDES
        }
    }
    
    // 2. Fall back to system rules
    systemRule := s.repo.FindSystemRuleByModel(req.Model)
    if systemRule != nil {
        return s.buildRouteResult(systemRule), nil
    }
    
    return nil, status.Error(codes.NotFound, "no routing rule found")
}
```

### Decision 2: Fallback Chain Storage

**Choice:** `repeated string fallback_provider_ids` in protobuf

**Rationale:**
- Ordered list preserves fallback priority
- Simple to implement and understand
- Sufficient for Harry's use case (ordered failover)

**Alternative considered:** JSON blob with complex fallback logic
- Rejected: over-engineering for current requirements

**Implementation:**
```protobuf
message RoutingRule {
  string id = 1;
  string user_id = 2;              // NEW: null/empty = system rule
  string model_pattern = 3;
  string provider_id = 4;
  int32 priority = 5;
  repeated string fallback_provider_ids = 6;  // CHANGED: was single fallback_provider_id
  bool is_system_default = 7;       // NEW: distinguishes system vs user rules
}
```

### Decision 3: Fallback Trigger Conditions

**Choice:** Trigger on ALL error conditions (5xx, timeouts, specific error codes)

**Rationale:**
- User wants automatic failover for any failure
- Consistent behavior across providers

**Implementation:**
```go
// In provider-service or gateway-service proxy logic:
func (p *ProxyHandler) tryWithFallback(ctx context.Context, primaryProvider string, fallbackProviders []string, req *Request) (*Response, error) {
    providers := append([]string{primaryProvider}, fallbackProviders...)
    
    for i, provider := range providers {
        resp, err := p.forwardToProvider(ctx, provider, req)
        if err == nil && resp.StatusCode < 500 {
            return resp, nil
        }
        
        // Last provider failed, return error
        if i == len(providers)-1 {
            return resp, err
        }
        
        // Log fallback attempt
        log.Printf("Provider %s failed, trying fallback %s", provider, fallbackProviders[i])
    }
    
    return nil, errors.New("all providers failed")
}
```

### Decision 4: API Endpoint Design

**Choice:** Both user self-service and admin override APIs

**Rationale:**
- User self-service (`/v1/routing-rules`) for Harry's use case
- Admin override (`/admin/users/{userId}/routing-rules`) for support scenarios

**Endpoints:**
```
User Self-Service (JWT auth required):
  GET    /v1/routing-rules          # List own rules
  POST   /v1/routing-rules          # Create rule
  GET    /v1/routing-rules/:id      # Get specific rule
  PUT    /v1/routing-rules/:id      # Update rule
  DELETE /v1/routing-rules/:id      # Delete rule

Admin Override (admin role required):
  GET    /admin/users/:userId/routing-rules
  POST   /admin/users/:userId/routing-rules
  GET    /admin/users/:userId/routing-rules/:id
  PUT    /admin/users/:userId/routing-rules/:id
  DELETE /admin/users/:userId/routing-rules/:id
```

### Decision 5: Database Schema Migration

**Choice:** Add columns to existing `routing_rules` table

**Rationale:**
- Single table for both system and user rules
- Simple queries with `WHERE user_id = ? OR user_id IS NULL`
- Avoids complex joins

**Migration:**
```sql
ALTER TABLE routing_rules ADD COLUMN user_id VARCHAR(255) DEFAULT NULL;
ALTER TABLE routing_rules ADD COLUMN is_system_default BOOLEAN DEFAULT FALSE;
ALTER TABLE routing_rules ADD COLUMN fallback_provider_ids TEXT;  -- JSON array

-- Index for efficient user rule lookup
CREATE INDEX idx_routing_rules_user_id ON routing_rules(user_id);
CREATE INDEX idx_routing_rules_model_pattern ON routing_rules(model_pattern);
```

## Risks / Trade-offs

**[Risk] User misconfiguration could block their own access**
→ **Mitigation:** Provide GET endpoint to list current rules; admins can delete user rules via override API

**[Risk] Fallback chain could cause cascading failures (all providers down)**
→ **Mitigation:** Set timeout per provider; return error after all fallbacks exhausted

**[Risk] Database migration on existing routing_rules table**
→ **Mitigation:** Test migration on copy of production data; provide rollback script

**[Risk] Increased latency from fallback retries**
→ **Mitigation:** Set reasonable timeouts; consider async health checks for future

## Migration Plan

1. **Database Migration:**
   - Add `user_id`, `is_system_default`, `fallback_provider_ids` columns
   - Backfill `is_system_default = TRUE` for existing rules
   - No downtime required (ADD COLUMN is lightweight)

2. **Deploy router-service:**
   - Update protobuf definitions
   - Implement new ResolveRoute logic (user override)
   - Add fallback execution logic

3. **Deploy gateway-service:**
   - Add `/v1/routing-rules` endpoints (JWT auth)
   - Add `/admin/users/{userId}/routing-rules` endpoints
   - Pass `user_id` to router-service

4. **Rollback Plan:**
   - Keep old endpoints working during transition
   - Database columns nullable with defaults = backward compatible
   - Revert gateway/router services if issues detected

## Open Questions

1. ~~Should user rules override or merge with system rules?~~ **ANSWERED: Override**
2. ~~What triggers fallback?~~ **ANSWERED: All error conditions**
3. Should we limit the number of fallback providers per rule? (Recommend: max 5)
4. Should there be a default fallback chain for users who don't configure rules? (Future enhancement)
