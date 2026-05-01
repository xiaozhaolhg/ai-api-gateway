// API client for admin endpoints
import { message } from 'antd';
import type { APIClientInterface } from './types';
import MockAPIClient from './mockClient';
import { API_CONFIG } from './config';

const TOKEN_KEY = 'auth_token';

export interface Provider {
  id: string;
  name: string;
  type: string;
  base_url: string;
  models: string[];
  status: string;
  created_at: string;
  updated_at: string;
}

export interface User {
  id: string;
  name: string;
  email: string;
  role: string;
  status: string;
  created_at: string;
}

export interface APIKey {
  id: string;
  user_id: string;
  name: string;
  scopes: string[];
  created_at: string;
  expires_at?: string;
}

export interface UsageRecord {
  id: string;
  user_id: string;
  model: string;
  provider: string;
  prompt_tokens: number;
  completion_tokens: number;
  total_tokens: number;
  cost: number;
  timestamp: string;
}

export interface RoutingRule {
  id: string;
  model_pattern: string;
  provider: string;
  adapter_type: string;
  priority: number;
  fallback_chain: string[];
  status: string;
  created_at: string;
  updated_at: string;
}

export interface Group {
  id: string;
  name: string;
  description: string;
  member_count: number;
  created_at: string;
  updated_at: string;
}

export interface Permission {
  id: string;
  group_id: string;
  model_pattern: string;
  effect: string;
  status: string;
  created_at: string;
  updated_at: string;
}

export interface Budget {
  id: string;
  name: string;
  scope: string;
  scope_id?: string;
  limit: number;
  current_spend: number;
  period: string;
  status: string;
  created_at: string;
  updated_at: string;
}

export interface PricingRule {
  id: string;
  model: string;
  provider: string;
  prompt_price: number;
  completion_price: number;
  currency: string;
  effective_date: string;
  created_at: string;
  updated_at: string;
}

export interface AlertRule {
  id: string;
  name: string;
  metric: string;
  condition: string;
  threshold: number;
  channel: string;
  status: string;
  created_at: string;
  updated_at: string;
}

export interface Alert {
  id: string;
  rule_id: string;
  severity: string;
  status: string;
  triggered_at: string;
  description: string;
  acknowledged_at?: string;
}

export interface ProviderHealth {
  id: string;
  name: string;
  status: string;
  latency_ms: number;
  error_rate: number;
  last_check: string;
}

class RealAPIClient implements APIClientInterface {
  private baseURL: string;

  constructor(baseURL: string) {
    this.baseURL = baseURL;
  }

  private getAuthHeader(): Record<string, string> {
    const token = localStorage.getItem(TOKEN_KEY);
    if (token) {
      return { Authorization: `Bearer ${token}` };
    }
    return {};
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`;
    const authHeader = this.getAuthHeader();

    try {
      const response = await fetch(url, {
        ...options,
        headers: {
          'Content-Type': 'application/json',
          ...authHeader,
          ...options.headers,
        },
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        const errorMessage = errorData.message || `API error: ${response.status} ${response.statusText}`;
        message.error(errorMessage);
        throw new Error(errorMessage);
      }

      return response.json();
    } catch (error) {
      if (error instanceof Error) {
        message.error(error.message);
      }
      throw error;
    }
  }

  // Provider endpoints
  async getProviders(): Promise<Provider[]> {
    return this.request<Provider[]>('/admin/providers');
  }

  async createProvider(provider: Omit<Provider, 'id' | 'created_at' | 'updated_at'>): Promise<Provider> {
    return this.request<Provider>('/admin/providers', {
      method: 'POST',
      body: JSON.stringify(provider),
    });
  }

  async updateProvider(id: string, provider: Partial<Provider>): Promise<Provider> {
    return this.request<Provider>(`/admin/providers/${id}`, {
      method: 'PUT',
      body: JSON.stringify(provider),
    });
  }

  async deleteProvider(id: string): Promise<void> {
    await this.request<void>(`/admin/providers/${id}`, {
      method: 'DELETE',
    });
  }

  // User endpoints
  async getUsers(): Promise<User[]> {
    return this.request<User[]>('/admin/users');
  }

  async createUser(user: Omit<User, 'id' | 'created_at'>): Promise<User> {
    return this.request<User>('/admin/users', {
      method: 'POST',
      body: JSON.stringify(user),
    });
  }

  async updateUser(id: string, user: Partial<User>): Promise<User> {
    return this.request<User>(`/admin/users/${id}`, {
      method: 'PUT',
      body: JSON.stringify(user),
    });
  }

  async deleteUser(id: string): Promise<void> {
    await this.request<void>(`/admin/users/${id}`, {
      method: 'DELETE',
    });
  }

  // API Key endpoints
  async getAPIKeys(userId: string): Promise<APIKey[]> {
    return this.request<APIKey[]>(`/admin/api-keys/${userId}`);
  }

  async createAPIKey(userId: string, name: string): Promise<{ api_key_id: string; api_key: string }> {
    return this.request<{ api_key_id: string; api_key: string }>('/admin/api-keys', {
      method: 'POST',
      body: JSON.stringify({ user_id: userId, name }),
    });
  }

  async deleteAPIKey(id: string): Promise<void> {
    await this.request<void>(`/admin/api-keys/${id}`, {
      method: 'DELETE',
    });
  }

  // Usage endpoints
  async getUsage(userId?: string, startDate?: string, endDate?: string): Promise<UsageRecord[]> {
    const params = new URLSearchParams();
    if (userId) params.append('user_id', userId);
    if (startDate) params.append('start_date', startDate);
    if (endDate) params.append('end_date', endDate);

    const query = params.toString();
    return this.request<UsageRecord[]>(`/admin/usage${query ? `?${query}` : ''}`);
  }

  // Authentication endpoints
  async login(emailOrUsername: string, password: string): Promise<{ token: string; user: User }> {
    const url = `${this.baseURL}/admin/auth/login`;
    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ email: emailOrUsername, password }),
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      const errorMessage = errorData.message || `API error: ${response.status} ${response.statusText}`;
      throw new Error(errorMessage);
    }

    return response.json();
  }

  async register(name: string, username: string, email: string, password: string): Promise<{ token: string; user: User }> {
    const url = `${this.baseURL}/admin/auth/register`;
    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ name, username, email, password }),
    });

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      const errorMessage = errorData.error || `API error: ${response.status} ${response.statusText}`;
      throw new Error(errorMessage);
    }

    return response.json();
  }

  async logout(): Promise<void> {
    await this.request<void>('/admin/auth/logout', {
      method: 'POST',
    });
  }

  async getCurrentUser(): Promise<User> {
    return this.request<User>('/admin/auth/me');
  }

  // Routing rule endpoints
  async getRoutingRules(): Promise<RoutingRule[]> {
    return this.request<RoutingRule[]>('/admin/routing-rules');
  }

  async createRoutingRule(rule: Omit<RoutingRule, 'id' | 'created_at' | 'updated_at'>): Promise<RoutingRule> {
    return this.request<RoutingRule>('/admin/routing-rules', {
      method: 'POST',
      body: JSON.stringify(rule),
    });
  }

  async updateRoutingRule(id: string, rule: Partial<RoutingRule>): Promise<RoutingRule> {
    return this.request<RoutingRule>(`/admin/routing-rules/${id}`, {
      method: 'PUT',
      body: JSON.stringify(rule),
    });
  }

  async deleteRoutingRule(id: string): Promise<void> {
    await this.request<void>(`/admin/routing-rules/${id}`, {
      method: 'DELETE',
    });
  }

  // Group endpoints
  async getGroups(): Promise<Group[]> {
    return this.request<Group[]>('/admin/groups');
  }

  async createGroup(group: Omit<Group, 'id' | 'created_at' | 'updated_at' | 'member_count'>): Promise<Group> {
    return this.request<Group>('/admin/groups', {
      method: 'POST',
      body: JSON.stringify(group),
    });
  }

  async updateGroup(id: string, group: Partial<Group>): Promise<Group> {
    return this.request<Group>(`/admin/groups/${id}`, {
      method: 'PUT',
      body: JSON.stringify(group),
    });
  }

  async deleteGroup(id: string): Promise<void> {
    await this.request<void>(`/admin/groups/${id}`, {
      method: 'DELETE',
    });
  }

  async addGroupMember(groupId: string, userId: string): Promise<void> {
    await this.request<void>(`/admin/groups/${groupId}/members`, {
      method: 'POST',
      body: JSON.stringify({ user_id: userId }),
    });
  }

  async removeGroupMember(groupId: string, userId: string): Promise<void> {
    await this.request<void>(`/admin/groups/${groupId}/members/${userId}`, {
      method: 'DELETE',
    });
  }

  // Permission endpoints
  async getPermissions(): Promise<Permission[]> {
    return this.request<Permission[]>('/admin/permissions');
  }

  async createPermission(permission: Omit<Permission, 'id' | 'created_at' | 'updated_at'>): Promise<Permission> {
    return this.request<Permission>('/admin/permissions', {
      method: 'POST',
      body: JSON.stringify(permission),
    });
  }

  async updatePermission(id: string, permission: Partial<Permission>): Promise<Permission> {
    return this.request<Permission>(`/admin/permissions/${id}`, {
      method: 'PUT',
      body: JSON.stringify(permission),
    });
  }

  async deletePermission(id: string): Promise<void> {
    await this.request<void>(`/admin/permissions/${id}`, {
      method: 'DELETE',
    });
  }

  // Budget endpoints
  async getBudgets(): Promise<Budget[]> {
    return this.request<Budget[]>('/admin/budgets');
  }

  async createBudget(budget: Omit<Budget, 'id' | 'created_at' | 'updated_at' | 'current_spend'>): Promise<Budget> {
    return this.request<Budget>('/admin/budgets', {
      method: 'POST',
      body: JSON.stringify(budget),
    });
  }

  async updateBudget(id: string, budget: Partial<Budget>): Promise<Budget> {
    return this.request<Budget>(`/admin/budgets/${id}`, {
      method: 'PUT',
      body: JSON.stringify(budget),
    });
  }

  async deleteBudget(id: string): Promise<void> {
    await this.request<void>(`/admin/budgets/${id}`, {
      method: 'DELETE',
    });
  }

  // Pricing rule endpoints
  async getPricingRules(): Promise<PricingRule[]> {
    return this.request<PricingRule[]>('/admin/pricing-rules');
  }

  async createPricingRule(rule: Omit<PricingRule, 'id' | 'created_at' | 'updated_at'>): Promise<PricingRule> {
    return this.request<PricingRule>('/admin/pricing-rules', {
      method: 'POST',
      body: JSON.stringify(rule),
    });
  }

  async updatePricingRule(id: string, rule: Partial<PricingRule>): Promise<PricingRule> {
    return this.request<PricingRule>(`/admin/pricing-rules/${id}`, {
      method: 'PUT',
      body: JSON.stringify(rule),
    });
  }

  async deletePricingRule(id: string): Promise<void> {
    await this.request<void>(`/admin/pricing-rules/${id}`, {
      method: 'DELETE',
    });
  }

  // Alert endpoints
  async getAlertRules(): Promise<AlertRule[]> {
    return this.request<AlertRule[]>('/admin/alert-rules');
  }

  async createAlertRule(rule: Omit<AlertRule, 'id' | 'created_at' | 'updated_at'>): Promise<AlertRule> {
    return this.request<AlertRule>('/admin/alert-rules', {
      method: 'POST',
      body: JSON.stringify(rule),
    });
  }

  async updateAlertRule(id: string, rule: Partial<AlertRule>): Promise<AlertRule> {
    return this.request<AlertRule>(`/admin/alert-rules/${id}`, {
      method: 'PUT',
      body: JSON.stringify(rule),
    });
  }

  async deleteAlertRule(id: string): Promise<void> {
    await this.request<void>(`/admin/alert-rules/${id}`, {
      method: 'DELETE',
    });
  }

  async getAlerts(): Promise<Alert[]> {
    return this.request<Alert[]>('/admin/alerts');
  }

  async acknowledgeAlert(id: string): Promise<void> {
    await this.request<void>(`/admin/alerts/${id}/acknowledge`, {
      method: 'PUT',
    });
  }

  // Health endpoint
  async getProviderHealth(): Promise<ProviderHealth[]> {
    return this.request<ProviderHealth[]>('/admin/health');
  }
}

// Unified APIClient that switches between Real and Mock based on configuration
class UnifiedAPIClient implements APIClientInterface {
  private realClient: RealAPIClient;
  private mockClient: MockAPIClient;
  private useMock: boolean;

  constructor() {
    this.realClient = new RealAPIClient(API_CONFIG.baseURL);
    this.mockClient = new MockAPIClient(API_CONFIG.mockDelay);
    this.useMock = API_CONFIG.useMock;
  }

  private getActiveClient(): APIClientInterface {
    return this.useMock ? this.mockClient : this.realClient;
  }

  // Authentication
  async login(email: string, password: string) {
    return this.getActiveClient().login(email, password);
  }

  async register(name: string, username: string, email: string, password: string) {
    return this.getActiveClient().register(name, username, email, password);
  }

  async logout() {
    return this.getActiveClient().logout();
  }

  async getCurrentUser() {
    return this.getActiveClient().getCurrentUser();
  }

  // Providers
  async getProviders() {
    return this.getActiveClient().getProviders();
  }

  async createProvider(provider: Omit<Provider, 'id' | 'created_at' | 'updated_at'>) {
    return this.getActiveClient().createProvider(provider);
  }

  async updateProvider(id: string, provider: Partial<Provider>) {
    return this.getActiveClient().updateProvider(id, provider);
  }

  async deleteProvider(id: string) {
    return this.getActiveClient().deleteProvider(id);
  }

  // Users
  async getUsers() {
    return this.getActiveClient().getUsers();
  }

  async createUser(user: Omit<User, 'id' | 'created_at'>) {
    return this.getActiveClient().createUser(user);
  }

  async updateUser(id: string, user: Partial<User>) {
    return this.getActiveClient().updateUser(id, user);
  }

  async deleteUser(id: string) {
    return this.getActiveClient().deleteUser(id);
  }

  // API Keys
  async getAPIKeys(userId: string) {
    return this.getActiveClient().getAPIKeys(userId);
  }

  async createAPIKey(userId: string, name: string) {
    return this.getActiveClient().createAPIKey(userId, name);
  }

  async deleteAPIKey(id: string) {
    return this.getActiveClient().deleteAPIKey(id);
  }

  // Usage
  async getUsage(userId?: string, startDate?: string, endDate?: string) {
    return this.getActiveClient().getUsage(userId, startDate, endDate);
  }

  // Routing Rules
  async getRoutingRules() {
    return this.getActiveClient().getRoutingRules();
  }

  async createRoutingRule(rule: Omit<any, 'id' | 'created_at' | 'updated_at'>) {
    return this.getActiveClient().createRoutingRule(rule);
  }

  async updateRoutingRule(id: string, rule: Partial<any>) {
    return this.getActiveClient().updateRoutingRule(id, rule);
  }

  async deleteRoutingRule(id: string) {
    return this.getActiveClient().deleteRoutingRule(id);
  }

  // Groups
  async getGroups() {
    return this.getActiveClient().getGroups();
  }

  async createGroup(group: Omit<any, 'id' | 'created_at' | 'updated_at' | 'member_count'>) {
    return this.getActiveClient().createGroup(group);
  }

  async updateGroup(id: string, group: Partial<any>) {
    return this.getActiveClient().updateGroup(id, group);
  }

  async deleteGroup(id: string) {
    return this.getActiveClient().deleteGroup(id);
  }

  async addGroupMember(groupId: string, userId: string) {
    return this.getActiveClient().addGroupMember(groupId, userId);
  }

  async removeGroupMember(groupId: string, userId: string) {
    return this.getActiveClient().removeGroupMember(groupId, userId);
  }

  // Permissions
  async getPermissions() {
    return this.getActiveClient().getPermissions();
  }

  async createPermission(permission: Omit<any, 'id' | 'created_at' | 'updated_at'>) {
    return this.getActiveClient().createPermission(permission);
  }

  async updatePermission(id: string, permission: Partial<any>) {
    return this.getActiveClient().updatePermission(id, permission);
  }

  async deletePermission(id: string) {
    return this.getActiveClient().deletePermission(id);
  }

  // Budgets
  async getBudgets() {
    return this.getActiveClient().getBudgets();
  }

  async createBudget(budget: Omit<any, 'id' | 'created_at' | 'updated_at' | 'current_spend'>) {
    return this.getActiveClient().createBudget(budget);
  }

  async updateBudget(id: string, budget: Partial<any>) {
    return this.getActiveClient().updateBudget(id, budget);
  }

  async deleteBudget(id: string) {
    return this.getActiveClient().deleteBudget(id);
  }

  // Pricing Rules
  async getPricingRules() {
    return this.getActiveClient().getPricingRules();
  }

  async createPricingRule(rule: Omit<any, 'id' | 'created_at' | 'updated_at'>) {
    return this.getActiveClient().createPricingRule(rule);
  }

  async updatePricingRule(id: string, rule: Partial<any>) {
    return this.getActiveClient().updatePricingRule(id, rule);
  }

  async deletePricingRule(id: string) {
    return this.getActiveClient().deletePricingRule(id);
  }

  // Alert Rules
  async getAlertRules() {
    return this.getActiveClient().getAlertRules();
  }

  async createAlertRule(rule: Omit<any, 'id' | 'created_at' | 'updated_at'>) {
    return this.getActiveClient().createAlertRule(rule);
  }

  async updateAlertRule(id: string, rule: Partial<any>) {
    return this.getActiveClient().updateAlertRule(id, rule);
  }

  async deleteAlertRule(id: string) {
    return this.getActiveClient().deleteAlertRule(id);
  }

  // Alerts
  async getAlerts() {
    return this.getActiveClient().getAlerts();
  }

  async acknowledgeAlert(id: string) {
    return this.getActiveClient().acknowledgeAlert(id);
  }

  // Health
  async getProviderHealth() {
    return this.getActiveClient().getProviderHealth();
  }

  // Method to switch between mock and real mode
  setMockMode(enabled: boolean) {
    this.useMock = enabled;
  }

  getMockMode(): boolean {
    return this.useMock;
  }
}

export const apiClient = new UnifiedAPIClient();
