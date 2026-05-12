package client

import (
	"context"
	"fmt"
	"log"

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

// GetUsage retrieves usage records with optional date range filtering.
// If startTime or endTime is 0, that bound is not applied.
func (c *BillingClient) GetUsage(ctx context.Context, userID string, page, pageSize int32, startTime, endTime int64) (*billingv1.ListUsageResponse, error) {
	req := &billingv1.GetUsageRequest{
		UserId:    userID,
		Page:      page,
		PageSize:  pageSize,
		StartTime: startTime,
		EndTime:   endTime,
	}

	resp, err := c.client.GetUsage(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get usage: %w", err)
	}

	return resp, nil
}

func (c *BillingClient) RecordUsage(ctx context.Context, userID, groupID, providerID, model string, promptTokens, completionTokens int64) error {
	log.Printf("[DEBUG] BillingClient.RecordUsage: calling billing service for userID=%s", userID)
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
		log.Printf("[DEBUG] BillingClient.RecordUsage: error=%v", err)
		return fmt.Errorf("failed to record usage: %w", err)
	}

	log.Printf("[DEBUG] BillingClient.RecordUsage: success")
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

// Billing Account operations
func (c *BillingClient) GetBillingAccountByUser(ctx context.Context, userID string) (*billingv1.BillingAccount, error) {
	req := &billingv1.GetBillingAccountByUserRequest{
		UserId: userID,
	}

	resp, err := c.client.GetBillingAccountByUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get billing account: %w", err)
	}

	return resp, nil
}

func (c *BillingClient) CreateBillingAccount(ctx context.Context, userID string, initialCredit float64) (*billingv1.BillingAccount, error) {
	req := &billingv1.CreateBillingAccountRequest{
		UserId:        userID,
		InitialCredit: initialCredit,
	}

	resp, err := c.client.CreateBillingAccount(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create billing account: %w", err)
	}

	return resp, nil
}

func (c *BillingClient) UpdateBillingAccountBalance(ctx context.Context, accountID string, balance float64) (*billingv1.BillingAccount, error) {
	req := &billingv1.UpdateBillingAccountRequest{
		Id:             accountID,
		Balance:        balance,
		BalanceUpdated: true,
	}

	resp, err := c.client.UpdateBillingAccount(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update billing account: %w", err)
	}

	return resp, nil
}
