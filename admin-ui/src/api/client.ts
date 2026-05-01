// API client for admin endpoints
import { message } from 'antd';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '';
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

class APIClient {
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
    return this.request<User[]>('/admin/auth/users');
  }

  async createUser(user: Omit<User, 'id' | 'created_at'>): Promise<User> {
    return this.request<User>('/admin/auth/users', {
      method: 'POST',
      body: JSON.stringify(user),
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
    return this.request<APIKey[]>(`/admin/auth/api-keys/${userId}`);
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
    const params = new URLSearchParams();
    if (userId) params.append('user_id', userId);
    if (startDate) params.append('start_date', startDate);
    if (endDate) params.append('end_date', endDate);

    const query = params.toString();
    return this.request<UsageRecord[]>(`/admin/auth/usage${query ? `?${query}` : ''}`);
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
    return this.request<Group[]>('/admin/auth/groups');
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

  // Permission endpoints
  async getPermissions(): Promise<Permission[]> {
    return this.request<Permission[]>('/admin/auth/permissions');
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

  async getProviderHealth(): Promise<ProviderHealth[]> {
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
  }
}

export const apiClient = new APIClient(API_BASE_URL);
