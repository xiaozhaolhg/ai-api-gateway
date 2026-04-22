package client

import (
	"context"
	"fmt"

	"github.com/ai-api-gateway/api/gen/billing/v1"
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

// OnProviderResponse handles provider response callback
func (c *BillingClient) OnProviderResponse(ctx context.Context, callback *billingv1.ProviderResponseCallback) error {
	req := &billingv1.OnProviderResponseRequest{
		Callback: callback,
	}

	_, err := c.client.OnProviderResponse(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to handle provider response callback: %w", err)
	}

	return nil
}

// RecordUsage records usage for a request
func (c *BillingClient) RecordUsage(ctx context.Context, userID, providerID, model string, promptTokens, completionTokens int64) error {
	req := &billingv1.RecordUsageRequest{
		UserId:           userID,
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

// GetUsageAggregation retrieves aggregated usage statistics
func (c *BillingClient) GetUsageAggregation(ctx context.Context, userID string, startDate, endDate string) (*billingv1.UsageAggregationResponse, error) {
	req := &billingv1.GetUsageAggregationRequest{
		UserId:    userID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	resp, err := c.client.GetUsageAggregation(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get usage aggregation: %w", err)
	}

	return resp, nil
}

// EstimateCost estimates the cost of a request
func (c *BillingClient) EstimateCost(ctx context.Context, providerID, model string, promptTokens, completionTokens int64) (*billingv1.CostEstimate, error) {
	req := &billingv1.EstimateCostRequest{
		ProviderId:       providerID,
		Model:            model,
		PromptTokens:     promptTokens,
		CompletionTokens: completionTokens,
	}

	resp, err := c.client.EstimateCost(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to estimate cost: %w", err)
	}

	return resp, nil
}

// CheckBudget checks if a user has budget available
func (c *BillingClient) CheckBudget(ctx context.Context, userID string) (*billingv1.BudgetStatus, error) {
	req := &billingv1.CheckBudgetRequest{
		UserId: userID,
	}

	resp, err := c.client.CheckBudget(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to check budget: %w", err)
	}

	return resp, nil
}

// CreateBudget creates a new budget
func (c *BillingClient) CreateBudget(ctx context.Context, budget *billingv1.Budget) (*billingv1.Budget, error) {
	req := &billingv1.CreateBudgetRequest{
		Budget: budget,
	}

	resp, err := c.client.CreateBudget(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create budget: %w", err)
	}

	return resp, nil
}

// UpdateBudget updates an existing budget
func (c *BillingClient) UpdateBudget(ctx context.Context, budget *billingv1.Budget) (*billingv1.Budget, error) {
	req := &billingv1.UpdateBudgetRequest{
		Budget: budget,
	}

	resp, err := c.client.UpdateBudget(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update budget: %w", err)
	}

	return resp, nil
}

// DeleteBudget deletes a budget
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

// CreatePricingRule creates a new pricing rule
func (c *BillingClient) CreatePricingRule(ctx context.Context, rule *billingv1.PricingRule) (*billingv1.PricingRule, error) {
	req := &billingv1.CreatePricingRuleRequest{
		Rule: rule,
	}

	resp, err := c.client.CreatePricingRule(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create pricing rule: %w", err)
	}

	return resp, nil
}

// UpdatePricingRule updates an existing pricing rule
func (c *BillingClient) UpdatePricingRule(ctx context.Context, rule *billingv1.PricingRule) (*billingv1.PricingRule, error) {
	req := &billingv1.UpdatePricingRuleRequest{
		Rule: rule,
	}

	resp, err := c.client.UpdatePricingRule(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update pricing rule: %w", err)
	}

	return resp, nil
}

// DeletePricingRule deletes a pricing rule
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

// GetBillingAccount retrieves a billing account
func (c *BillingClient) GetBillingAccount(ctx context.Context, userID string) (*billingv1.BillingAccount, error) {
	req := &billingv1.GetBillingAccountRequest{
		UserId: userID,
	}

	resp, err := c.client.GetBillingAccount(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get billing account: %w", err)
	}

	return resp, nil
}

// CreateBillingAccount creates a new billing account
func (c *BillingClient) CreateBillingAccount(ctx context.Context, account *billingv1.BillingAccount) (*billingv1.BillingAccount, error) {
	req := &billingv1.CreateBillingAccountRequest{
		Account: account,
	}

	resp, err := c.client.CreateBillingAccount(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create billing account: %w", err)
	}

	return resp, nil
}

// UpdateBillingAccount updates an existing billing account
func (c *BillingClient) UpdateBillingAccount(ctx context.Context, account *billingv1.BillingAccount) (*billingv1.BillingAccount, error) {
	req := &billingv1.UpdateBillingAccountRequest{
		Account: account,
	}

	resp, err := c.client.UpdateBillingAccount(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update billing account: %w", err)
	}

	return resp, nil
}

// GenerateInvoice generates an invoice
func (c *BillingClient) GenerateInvoice(ctx context.Context, userID string, startDate, endDate string) (*billingv1.Invoice, error) {
	req := &billingv1.GenerateInvoiceRequest{
		UserId:    userID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	resp, err := c.client.GenerateInvoice(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate invoice: %w", err)
	}

	return resp, nil
}

// GetInvoices retrieves invoices for a user
func (c *BillingClient) GetInvoices(ctx context.Context, userID string, page, pageSize int32) (*billingv1.ListInvoicesResponse, error) {
	req := &billingv1.GetInvoicesRequest{
		UserId:   userID,
		Page:     page,
		PageSize: pageSize,
	}

	resp, err := c.client.GetInvoices(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get invoices: %w", err)
	}

	return resp, nil
}
