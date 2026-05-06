package client

import (
	"context"
	"fmt"

	billingv1 "github.com/ai-api-gateway/api/gen/billing/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// BillingClient wraps the billing-service gRPC client
type BillingClient struct {
	client billingv1.BillingServiceClient
	conn   *grpc.ClientConn
}

// NewBillingClient creates a new billing service client
func NewBillingClient(address string) (*BillingClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to billing service: %w", err)
	}

	return &BillingClient{
		client: billingv1.NewBillingServiceClient(conn),
		conn:   conn,
	}, nil
}

// Close closes the connection
func (c *BillingClient) Close() error {
	return c.conn.Close()
}

// GetUsage retrieves usage records
func (c *BillingClient) GetUsage(ctx context.Context, userID string, page, pageSize int32) (*billingv1.ListUsageResponse, error) {
	req := &billingv1.GetUsageRequest{
		UserId:   userID,
		Page:     page,
		PageSize: pageSize,
	}

	resp, err := c.client.GetUsage(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get usage: %w", err)
	}

	return resp, nil
}

func (c *BillingClient) RecordUsage(ctx context.Context, userID, groupID, providerID, model string, promptTokens, completionTokens int64) error {
	req := &billingv1.RecordUsageRequest{
		UserId:           userID,
		GroupId:          groupID,
		ProviderId:       providerID,
		Model:            model,
		PromptTokens:     promptTokens,
		CompletionTokens: completionTokens,
	}

	_, err := c.client.RecordUsage(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to record usage: %w", err)
	}

	return nil
}

// Budget CRUD operations
func (c *BillingClient) ListBudgets(ctx context.Context, page, pageSize int32) (*billingv1.ListBudgetsResponse, error) {
	req := &billingv1.ListBudgetsRequest{
		UserId:   "",
		Page:     page,
		PageSize: pageSize,
	}

	resp, err := c.client.ListBudgets(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list budgets: %w", err)
	}

	return resp, nil
}

func (c *BillingClient) CreateBudget(ctx context.Context, userID string, limit float64, period string, softCapPct, hardCapPct float64) (*billingv1.Budget, error) {
	req := &billingv1.CreateBudgetRequest{
		UserId:      userID,
		Limit:       limit,
		Period:      period,
		SoftCapPct:  softCapPct,
		HardCapPct:  hardCapPct,
	}

	resp, err := c.client.CreateBudget(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create budget: %w", err)
	}

	return resp, nil
}

func (c *BillingClient) UpdateBudget(ctx context.Context, id, userID string, limit float64, period string, softCapPct, hardCapPct float64, status string) (*billingv1.Budget, error) {
	req := &billingv1.UpdateBudgetRequest{
		Id:         id,
		UserId:     userID,
		Limit:      limit,
		Period:     period,
		SoftCapPct:  softCapPct,
		HardCapPct:  hardCapPct,
		Status:     status,
	}

	resp, err := c.client.UpdateBudget(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update budget: %w", err)
	}

	return resp, nil
}

func (c *BillingClient) DeleteBudget(ctx context.Context, id string) error {
	req := &billingv1.DeleteBudgetRequest{
		Id: id,
	}

	_, err := c.client.DeleteBudget(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete budget: %w", err)
	}

	return nil
}

// PricingRule CRUD operations
func (c *BillingClient) ListPricingRules(ctx context.Context, page, pageSize int32) (*billingv1.ListPricingRulesResponse, error) {
	req := &billingv1.ListPricingRulesRequest{
		Model:     "",
		ProviderId: "",
		Page:      page,
		PageSize:  pageSize,
	}

	resp, err := c.client.ListPricingRules(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list pricing rules: %w", err)
	}

	return resp, nil
}

func (c *BillingClient) CreatePricingRule(ctx context.Context, model, providerID string, promptPrice, completionPrice float64, currency string) (*billingv1.PricingRule, error) {
	req := &billingv1.CreatePricingRuleRequest{
		Model:                   model,
		ProviderId:              providerID,
		PricePerPromptToken:     promptPrice,
		PricePerCompletionToken: completionPrice,
		Currency:                currency,
	}

	resp, err := c.client.CreatePricingRule(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create pricing rule: %w", err)
	}

	return resp, nil
}

func (c *BillingClient) UpdatePricingRule(ctx context.Context, id, model, providerID string, promptPrice, completionPrice float64, currency string) (*billingv1.PricingRule, error) {
	req := &billingv1.UpdatePricingRuleRequest{
		Id:                      id,
		Model:                   model,
		ProviderId:              providerID,
		PricePerPromptToken:     promptPrice,
		PricePerCompletionToken: completionPrice,
		Currency:                currency,
	}

	resp, err := c.client.UpdatePricingRule(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update pricing rule: %w", err)
	}

	return resp, nil
}

func (c *BillingClient) DeletePricingRule(ctx context.Context, id string) error {
	req := &billingv1.DeletePricingRuleRequest{
		Id: id,
	}

	_, err := c.client.DeletePricingRule(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete pricing rule: %w", err)
	}

	return nil
}
