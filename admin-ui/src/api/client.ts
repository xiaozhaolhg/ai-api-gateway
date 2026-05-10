// API client for admin endpoints
import { message } from 'antd';
import type {
  APIClientInterface,
  UnauthorizedCallback,
  Provider,
  User,
  APIKey,
  UsageRecord,
  RoutingRule,
  Group,
  Permission,
  Budget,
  PricingRule,
  AlertRule,
  Alert,
  ProviderHealth,
  Tier,
} from './types';

export type {
  Provider,
  User,
  APIKey,
  UsageRecord,
  RoutingRule,
  Group,
  Permission,
  Budget,
  PricingRule,
  AlertRule,
  Alert,
  ProviderHealth,
  Tier,
} from './types';

import MockAPIClient from './mockClient';
import { API_CONFIG } from './config';

const TOKEN_KEY = 'auth_token';

export class RealAPIClient implements APIClientInterface {
  private baseURL: string;
  onUnauthorized?: UnauthorizedCallback;

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

        if (response.status === 401) {
          message.error('Session expired. Please login again.');
          this.onUnauthorized?.();
        }

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
    try {
      return await this.request<Provider[]>('/admin/providers');
    } catch (error) {
      console.error('Failed to fetch providers:', error);
      return [];
    }
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
    try {
      const response = await this.request<{users: User[]; total: number}>('/admin/auth/users');
      return response.users || [];
    } catch (error) {
      console.error('Failed to fetch users:', error);
      return [];
    }
  }

  async createUser(user: Omit<User, 'id' | 'created_at'>): Promise<User> {
    return this.request<User>('/admin/auth/users', {
      method: 'POST',
      body: JSON.stringify(user),
    });
  }

  async checkUsernameAvailability(username: string): Promise<{ available: boolean }> {
    return this.request<{ available: boolean }>('/admin/auth/check-username', {
      method: 'POST',
      body: JSON.stringify({ username }),
    });
  }

  async updateUser(id: string, user: Partial<User>): Promise<User> {
    return this.request<User>(`/admin/auth/users/${id}`, {
      method: 'PUT',
      body: JSON.stringify(user),
    });
  }

  async deleteUser(id: string): Promise<void> {
    await this.request<void>(`/admin/auth/users/${id}`, {
      method: 'DELETE',
    });
  }

  // API Key endpoints
  async getAPIKeys(userId: string): Promise<APIKey[]> {
    try {
      const response = await this.request<{api_keys: APIKey[]; total: number}>(`/admin/auth/api-keys/${userId}`);
      return response.api_keys || [];
    } catch (error) {
      console.error('Failed to fetch API keys:', error);
      return [];
    }
  }

  async createAPIKey(userId: string, name: string): Promise<{ api_key_id: string; api_key: string }> {
    return this.request<{ api_key_id: string; api_key: string }>('/admin/auth/api-keys', {
      method: 'POST',
      body: JSON.stringify({ user_id: userId, name }),
    });
  }

  async deleteAPIKey(id: string): Promise<void> {
    await this.request<void>(`/admin/auth/api-keys/${id}`, {
      method: 'DELETE',
    });
  }

  // Usage endpoints
  async getUsage(userId?: string, startDate?: string, endDate?: string): Promise<UsageRecord[]> {
    try {
      const params = new URLSearchParams();
      if (userId) params.append('user_id', userId);
      if (startDate) params.append('start_date', startDate);
      if (endDate) params.append('end_date', endDate);

      const query = params.toString();
      const response = await this.request<{usage: UsageRecord[]; total: number}>(`/admin/auth/usage${query ? `?${query}` : ''}`);
      return response.usage || [];
    } catch (error) {
      console.error('Failed to fetch usage:', error);
      return [];
    }
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
    try {
      return await this.request<RoutingRule[]>('/admin/routing-rules');
    } catch (error) {
      console.error('Failed to fetch routing rules:', error);
      return [];
    }
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
    try {
      const response = await this.request<{groups: Group[]; total: number}>('/admin/auth/groups');
      return response.groups || [];
    } catch (error) {
      console.error('Failed to fetch groups:', error);
      return [];
    }
  }

  async createGroup(group: Omit<Group, 'id' | 'created_at' | 'updated_at' | 'member_count'>): Promise<Group> {
    return this.request<Group>('/admin/auth/groups', {
      method: 'POST',
      body: JSON.stringify(group),
    });
  }

  async updateGroup(id: string, group: Partial<Group>): Promise<Group> {
    return this.request<Group>(`/admin/auth/groups/${id}`, {
      method: 'PUT',
      body: JSON.stringify(group),
    });
  }

  async deleteGroup(id: string): Promise<void> {
    await this.request<void>(`/admin/auth/groups/${id}`, {
      method: 'DELETE',
    });
  }

  async getGroupMembers(groupId: string): Promise<User[]> {
    try {
      const response = await this.request<{members: User[]; total: number}>(`/admin/auth/groups/${groupId}/members`);
      return response.members || [];
    } catch (error) {
      console.error('Failed to fetch group members:', error);
      return [];
    }
  }

  async addGroupMember(groupId: string, userId: string): Promise<void> {
    await this.request<void>(`/admin/auth/groups/${groupId}/members`, {
      method: 'POST',
      body: JSON.stringify({ user_id: userId }),
    });
  }

  async removeGroupMember(groupId: string, userId: string): Promise<void> {
    await this.request<void>(`/admin/auth/groups/${groupId}/members/${userId}`, {
      method: 'DELETE',
    });
  }

  async getGroupPermissions(groupId: string): Promise<Permission[]> {
    try {
      const response = await this.request<{permissions: Permission[]; total: number}>(`/admin/auth/permissions?group_id=${groupId}`);
      return response.permissions || [];
    } catch (error) {
      console.error('Failed to fetch group permissions:', error);
      return [];
    }
  }

  async getUserGroups(_userId: string): Promise<Group[]> {
    try {
      const allGroups = await this.getGroups();
      return allGroups;
    } catch (error) {
      console.error('Failed to fetch user groups:', error);
      return [];
    }
  }

  // Permission endpoints
  async getPermissions(): Promise<Permission[]> {
    try {
      const response = await this.request<{permissions: Permission[]; total: number}>('/admin/auth/permissions');
      return response.permissions || [];
    } catch (error) {
      console.error('Failed to fetch permissions:', error);
      return [];
    }
  }

  async createPermission(permission: Omit<Permission, 'id' | 'created_at' | 'updated_at'>): Promise<Permission> {
    return this.request<Permission>('/admin/auth/permissions', {
      method: 'POST',
      body: JSON.stringify(permission),
    });
  }

  async updatePermission(id: string, permission: Partial<Permission>): Promise<Permission> {
    return this.request<Permission>(`/admin/auth/permissions/${id}`, {
      method: 'PUT',
      body: JSON.stringify(permission),
    });
  }

  async deletePermission(id: string): Promise<void> {
    await this.request<void>(`/admin/auth/permissions/${id}`, {
      method: 'DELETE',
    });
  }

  // Budget endpoints
  async getBudgets(): Promise<Budget[]> {
    try {
      return await this.request<Budget[]>('/admin/budgets');
    } catch (error) {
      console.error('Failed to fetch budgets:', error);
      return [];
    }
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
    try {
      return await this.request<PricingRule[]>('/admin/pricing-rules');
    } catch (error) {
      console.error('Failed to fetch pricing rules:', error);
      return [];
    }
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
    try {
      return await this.request<AlertRule[]>('/admin/alert-rules');
    } catch (error) {
      console.error('Failed to fetch alert rules:', error);
      return [];
    }
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
    try {
      return await this.request<Alert[]>('/admin/alerts');
    } catch (error) {
      console.error('Failed to fetch alerts:', error);
      return [];
    }
  }

  async acknowledgeAlert(id: string): Promise<void> {
    await this.request<void>(`/admin/alerts/${id}/acknowledge`, {
      method: 'PUT',
    });
  }

  async getProviderHealth(): Promise<ProviderHealth[]> {
    try {
      const providers = await this.request<Provider[]>('/admin/providers');
      const healthPromises = providers.map(async (p) => {
        try {
          const result = await this.request<{ status: string; latency_ms: number; error_rate: number }>(`/admin/providers/${p.id}/health`);
          return {
            id: p.id,
            name: p.name,
            status: result.status,
            latency_ms: result.latency_ms,
            error_rate: result.error_rate,
            last_check: new Date().toISOString(),
          } as ProviderHealth;
        } catch {
          return {
            id: p.id,
            name: p.name,
            status: 'unhealthy',
            latency_ms: 0,
            error_rate: 1,
            last_check: new Date().toISOString(),
          } as ProviderHealth;
        }
      });
      return Promise.all(healthPromises);
    } catch (error) {
      console.error('Failed to fetch provider health:', error);
      return [];
    }
  }

  async getTiers(): Promise<Tier[]> {
    try {
      const response = await this.request<{tiers: Tier[]; total: number}>('/admin/auth/tiers');
      return response.tiers || [];
    } catch (error) {
      console.error('Failed to fetch tiers:', error);
      return [];
    }
  }

  async createTier(tier: Omit<Tier, 'id' | 'created_at' | 'updated_at'>): Promise<Tier> {
    return this.request<Tier>('/admin/auth/tiers', {
      method: 'POST',
      body: JSON.stringify(tier),
    });
  }

  async updateTier(id: string, tier: Partial<Tier>): Promise<Tier> {
    return this.request<Tier>(`/admin/auth/tiers/${id}`, {
      method: 'PUT',
      body: JSON.stringify(tier),
    });
  }

  async deleteTier(id: string): Promise<void> {
    await this.request<void>(`/admin/auth/tiers/${id}`, {
      method: 'DELETE',
    });
  }

  async assignTierToGroup(groupId: string, tierId: string): Promise<void> {
    await this.request<void>(`/admin/auth/groups/${groupId}/tier`, {
      method: 'POST',
      body: JSON.stringify({ tier_id: tierId }),
    });
  }

  async removeTierFromGroup(groupId: string): Promise<void> {
    await this.request<void>(`/admin/auth/groups/${groupId}/tier`, {
      method: 'DELETE',
    });
  }

  setOnUnauthorized(callback?: UnauthorizedCallback) {
    this.onUnauthorized = callback;
  }

  getOnUnauthorized(): UnauthorizedCallback | undefined {
    return this.onUnauthorized;
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

  async checkUsernameAvailability(username: string) {
    return this.getActiveClient().checkUsernameAvailability(username);
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

  async createGroup(group: Omit<Group, 'id' | 'created_at' | 'updated_at' | 'member_count'>) {
    return this.getActiveClient().createGroup(group);
  }

  async updateGroup(id: string, group: Partial<Group>) {
    return this.getActiveClient().updateGroup(id, group);
  }

  async deleteGroup(id: string) {
    return this.getActiveClient().deleteGroup(id);
  }

  async getGroupMembers(groupId: string): Promise<User[]> {
    return this.getActiveClient().getGroupMembers(groupId);
  }

  async addGroupMember(groupId: string, userId: string) {
    return this.getActiveClient().addGroupMember(groupId, userId);
  }

  async removeGroupMember(groupId: string, userId: string) {
    return this.getActiveClient().removeGroupMember(groupId, userId);
  }

  async getGroupPermissions(groupId: string): Promise<Permission[]> {
    return this.getActiveClient().getGroupPermissions(groupId);
  }

  async getUserGroups(userId: string): Promise<Group[]> {
    return this.getActiveClient().getUserGroups(userId);
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

  async getTiers() {
    return this.getActiveClient().getTiers();
  }

  async createTier(tier: Omit<Tier, 'id' | 'created_at' | 'updated_at'>) {
    return this.getActiveClient().createTier(tier);
  }

  async updateTier(id: string, tier: Partial<Tier>) {
    return this.getActiveClient().updateTier(id, tier);
  }

  async deleteTier(id: string) {
    return this.getActiveClient().deleteTier(id);
  }

  async assignTierToGroup(groupId: string, tierId: string) {
    return this.getActiveClient().assignTierToGroup(groupId, tierId);
  }

  async removeTierFromGroup(groupId: string) {
    return this.getActiveClient().removeTierFromGroup(groupId);
  }

  setMockMode(enabled: boolean) {
    this.useMock = enabled;
  }

  getMockMode(): boolean {
    return this.useMock;
  }

  setOnUnauthorized(callback?: UnauthorizedCallback) {
    this.realClient.setOnUnauthorized(callback);
  }

  getOnUnauthorized(): UnauthorizedCallback | undefined {
    return this.realClient.getOnUnauthorized();
  }
}

export const apiClient = new UnifiedAPIClient();
