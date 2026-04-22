// API client for admin endpoints
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

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
  prompt_tokens: number;
  completion_tokens: number;
  total_tokens: number;
  cost: number;
  timestamp: string;
}

class APIClient {
  private baseURL: string;

  constructor(baseURL: string) {
    this.baseURL = baseURL;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`;
    const response = await fetch(url, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
    });

    if (!response.ok) {
      throw new Error(`API error: ${response.status} ${response.statusText}`);
    }

    return response.json();
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
    return this.request<APIKey[]>(`/admin/users/${userId}/api-keys`);
  }

  async createAPIKey(userId: string, name: string): Promise<{ api_key_id: string; api_key: string }> {
    return this.request<{ api_key_id: string; api_key: string }>(`/admin/users/${userId}/api-keys`, {
      method: 'POST',
      body: JSON.stringify({ name }),
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
}

export const apiClient = new APIClient(API_BASE_URL);
