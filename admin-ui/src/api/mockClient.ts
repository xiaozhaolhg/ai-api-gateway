import type {
  APIClientInterface,
  LoginResponse,
  RegisterResponse,
  User,
  Provider,
  APIKey,
  UsageRecord,
  RoutingRule,
  Group,
  Permission,
  Budget,
  PricingRule,
  AlertRule,
  Alert,
  ProviderHealth
} from './types';
import MockDataHandler from '../mock/handlers/dataHandler';
import { API_CONFIG } from './config';

class MockAPIClient implements APIClientInterface {
  private dataHandler: MockDataHandler;
  private delay: number;

  constructor(delay: number = API_CONFIG.mockDelay) {
    this.dataHandler = MockDataHandler.getInstance();
    this.delay = delay;
  }

  private async simulateNetworkDelay<T>(data: T): Promise<T> {
    if (this.delay === 0) return data;
    return new Promise(resolve => setTimeout(() => resolve(data), this.delay));
  }

  private generateId(): string {
    return Date.now().toString(36) + Math.random().toString(36).substr(2);
  }

  private getCurrentTimestamp(): string {
    return new Date().toISOString();
  }

  // ===== Authentication =====
  async login(emailOrUsername: string, password: string): Promise<LoginResponse> {
    // Support login with either email or username
    const user = this.dataHandler.getUsers().find(u => 
      u.email === emailOrUsername || u.email === `${emailOrUsername}@local.dev`
    );
    
    if (!user) {
      throw new Error('User not found');
    }

    // Mock password check - accept any password for demo
    if (!password || password.length < 1) {
      throw new Error('Invalid password');
    }

    const token = `mock-jwt-token-${this.generateId()}`;
    
    return this.simulateNetworkDelay({
      token,
      user: { ...user }
    });
  }

  async register(name: string, username: string, email: string, _password: string): Promise<RegisterResponse> {
    // Check if user already exists
    const existingUser = this.dataHandler.getUsers().find(u => u.email === email);
    if (existingUser) {
      throw new Error('User already exists');
    }

    const newUser: User = {
      id: this.generateId(),
      name,
      email: email || `${username}@local.dev`,
      username: username || 'defaultuser',
      role: 'user',
      status: 'active',
      created_at: this.getCurrentTimestamp()
    };

    this.dataHandler.addUser(newUser);

    const token = `mock-jwt-token-${this.generateId()}`;
    
    return this.simulateNetworkDelay({
      token,
      user: newUser
    });
  }

  async logout(): Promise<void> {
    return this.simulateNetworkDelay(undefined);
  }

  async getCurrentUser(): Promise<User> {
    // Return first user as current user for mock
    const users = this.dataHandler.getUsers();
    if (users.length === 0) {
      throw new Error('No users found');
    }
    return this.simulateNetworkDelay(users[0]);
  }

  // ===== Providers =====
  async getProviders(): Promise<Provider[]> {
    return this.simulateNetworkDelay(this.dataHandler.getProviders());
  }

  async createProvider(provider: Omit<Provider, 'id' | 'created_at' | 'updated_at'>): Promise<Provider> {
    const newProvider: Provider = {
      ...provider,
      id: this.generateId(),
      created_at: this.getCurrentTimestamp(),
      updated_at: this.getCurrentTimestamp()
    };

    this.dataHandler.addProvider(newProvider);
    return this.simulateNetworkDelay(newProvider);
  }

  async updateProvider(id: string, provider: Partial<Provider>): Promise<Provider> {
    const existing = this.dataHandler.getProviderById(id);
    if (!existing) {
      throw new Error('Provider not found');
    }

    const updated = { ...existing, ...provider, updated_at: this.getCurrentTimestamp() };
    this.dataHandler.updateProvider(id, updated);
    return this.simulateNetworkDelay(updated);
  }

  async deleteProvider(id: string): Promise<void> {
    const existing = this.dataHandler.getProviderById(id);
    if (!existing) {
      throw new Error('Provider not found');
    }

    this.dataHandler.deleteProvider(id);
    return this.simulateNetworkDelay(undefined);
  }

  // ===== Users =====
  async getUsers(): Promise<User[]> {
    return this.simulateNetworkDelay(this.dataHandler.getUsers());
  }

  async createUser(user: Omit<User, 'id' | 'created_at'>): Promise<User> {
    const newUser: User = {
      ...user,
      id: this.generateId(),
      created_at: this.getCurrentTimestamp()
    };

    this.dataHandler.addUser(newUser);
    return this.simulateNetworkDelay(newUser);
  }

  async updateUser(id: string, user: Partial<User>): Promise<User> {
    const existing = this.dataHandler.getUserById(id);
    if (!existing) {
      throw new Error('User not found');
    }

    const updated = { ...existing, ...user };
    this.dataHandler.updateUser(id, updated);
    return this.simulateNetworkDelay(updated);
  }

  async deleteUser(id: string): Promise<void> {
    const existing = this.dataHandler.getUserById(id);
    if (!existing) {
      throw new Error('User not found');
    }

    this.dataHandler.deleteUser(id);
    return this.simulateNetworkDelay(undefined);
  }

  async checkUsernameAvailability(username: string): Promise<{ available: boolean }> {
    // Mock implementation - always return available for testing
    console.log(`Checking username availability for: ${username}`);
    return this.simulateNetworkDelay({
      available: true
    });
  }

  // ===== API Keys =====
  async getAPIKeys(userId: string): Promise<APIKey[]> {
    return this.simulateNetworkDelay(this.dataHandler.getAPIKeys(userId));
  }

  async createAPIKey(userId: string, name: string): Promise<{ api_key_id: string; api_key: string }> {
    const newAPIKey: APIKey = {
      id: this.generateId(),
      user_id: userId,
      name,
      scopes: ['read', 'write'],
      created_at: this.getCurrentTimestamp()
    };

    this.dataHandler.addAPIKey(newAPIKey);
    
    return this.simulateNetworkDelay({
      api_key_id: newAPIKey.id,
      api_key: `mock-api-key-${this.generateId()}`
    });
  }

  async deleteAPIKey(id: string): Promise<void> {
    const existing = this.dataHandler.getAPIKeyById(id);
    if (!existing) {
      throw new Error('API key not found');
    }

    this.dataHandler.deleteAPIKey(id);
    return this.simulateNetworkDelay(undefined);
  }

  // ===== Usage =====
  async getUsage(userId?: string, startDate?: string, endDate?: string): Promise<UsageRecord[]> {
    return this.simulateNetworkDelay(this.dataHandler.getUsage(userId, startDate, endDate));
  }

  // ===== Routing Rules =====
  async getRoutingRules(): Promise<RoutingRule[]> {
    return this.simulateNetworkDelay(this.dataHandler.getRoutingRules());
  }

  async createRoutingRule(rule: Omit<RoutingRule, 'id' | 'created_at' | 'updated_at'>): Promise<RoutingRule> {
    const newRule: RoutingRule = {
      ...rule,
      id: this.generateId(),
      created_at: this.getCurrentTimestamp(),
      updated_at: this.getCurrentTimestamp()
    };

    this.dataHandler.addRoutingRule(newRule);
    return this.simulateNetworkDelay(newRule);
  }

  async updateRoutingRule(id: string, rule: Partial<RoutingRule>): Promise<RoutingRule> {
    const existing = this.dataHandler.getRoutingRuleById(id);
    if (!existing) {
      throw new Error('Routing rule not found');
    }

    const updated = { ...existing, ...rule, updated_at: this.getCurrentTimestamp() };
    this.dataHandler.updateRoutingRule(id, updated);
    return this.simulateNetworkDelay(updated);
  }

  async deleteRoutingRule(id: string): Promise<void> {
    const existing = this.dataHandler.getRoutingRuleById(id);
    if (!existing) {
      throw new Error('Routing rule not found');
    }

    this.dataHandler.deleteRoutingRule(id);
    return this.simulateNetworkDelay(undefined);
  }

  // ===== Groups =====
  async getGroups(): Promise<Group[]> {
    return this.simulateNetworkDelay(this.dataHandler.getGroups());
  }

  async createGroup(group: Omit<Group, 'id' | 'created_at' | 'updated_at' | 'member_count'>): Promise<Group> {
    const newGroup: Group = {
      ...group,
      id: this.generateId(),
      member_count: 0,
      created_at: this.getCurrentTimestamp(),
      updated_at: this.getCurrentTimestamp()
    };

    this.dataHandler.addGroup(newGroup);
    return this.simulateNetworkDelay(newGroup);
  }

  async updateGroup(id: string, group: Partial<Group>): Promise<Group> {
    const existing = this.dataHandler.getGroupById(id);
    if (!existing) {
      throw new Error('Group not found');
    }

    const updated = { ...existing, ...group, updated_at: this.getCurrentTimestamp() };
    this.dataHandler.updateGroup(id, updated);
    return this.simulateNetworkDelay(updated);
  }

  async deleteGroup(id: string): Promise<void> {
    const existing = this.dataHandler.getGroupById(id);
    if (!existing) {
      throw new Error('Group not found');
    }

    this.dataHandler.deleteGroup(id);
    return this.simulateNetworkDelay(undefined);
  }

  async addGroupMember(groupId: string, _userId: string): Promise<void> {
    const group = this.dataHandler.getGroupById(groupId);
    if (!group) {
      throw new Error('Group not found');
    }

    // Increment member count
    this.dataHandler.updateGroup(groupId, { 
      member_count: group.member_count + 1 
    });
    
    return this.simulateNetworkDelay(undefined);
  }

  async removeGroupMember(groupId: string, _userId: string): Promise<void> {
    const group = this.dataHandler.getGroupById(groupId);
    if (!group) {
      throw new Error('Group not found');
    }

    // Decrement member count (but not below 0)
    const newCount = Math.max(0, group.member_count - 1);
    this.dataHandler.updateGroup(groupId, { 
      member_count: newCount 
    });
    
    return this.simulateNetworkDelay(undefined);
  }

  // ===== Permissions =====
  async getPermissions(): Promise<Permission[]> {
    return this.simulateNetworkDelay(this.dataHandler.getPermissions());
  }

  async createPermission(permission: Omit<Permission, 'id' | 'created_at' | 'updated_at'>): Promise<Permission> {
    const newPermission: Permission = {
      ...permission,
      id: this.generateId(),
      created_at: this.getCurrentTimestamp(),
      updated_at: this.getCurrentTimestamp()
    };

    this.dataHandler.addPermission(newPermission);
    return this.simulateNetworkDelay(newPermission);
  }

  async updatePermission(id: string, permission: Partial<Permission>): Promise<Permission> {
    const existing = this.dataHandler.getPermissionById(id);
    if (!existing) {
      throw new Error('Permission not found');
    }

    const updated = { ...existing, ...permission, updated_at: this.getCurrentTimestamp() };
    this.dataHandler.updatePermission(id, updated);
    return this.simulateNetworkDelay(updated);
  }

  async deletePermission(id: string): Promise<void> {
    const existing = this.dataHandler.getPermissionById(id);
    if (!existing) {
      throw new Error('Permission not found');
    }

    this.dataHandler.deletePermission(id);
    return this.simulateNetworkDelay(undefined);
  }

  // ===== Budgets =====
  async getBudgets(): Promise<Budget[]> {
    return this.simulateNetworkDelay(this.dataHandler.getBudgets());
  }

  async createBudget(budget: Omit<Budget, 'id' | 'created_at' | 'updated_at' | 'current_spend'>): Promise<Budget> {
    const newBudget: Budget = {
      ...budget,
      id: this.generateId(),
      current_spend: 0,
      created_at: this.getCurrentTimestamp(),
      updated_at: this.getCurrentTimestamp()
    };

    this.dataHandler.addBudget(newBudget);
    return this.simulateNetworkDelay(newBudget);
  }

  async updateBudget(id: string, budget: Partial<Budget>): Promise<Budget> {
    const existing = this.dataHandler.getBudgetById(id);
    if (!existing) {
      throw new Error('Budget not found');
    }

    const updated = { ...existing, ...budget, updated_at: this.getCurrentTimestamp() };
    this.dataHandler.updateBudget(id, updated);
    return this.simulateNetworkDelay(updated);
  }

  async deleteBudget(id: string): Promise<void> {
    const existing = this.dataHandler.getBudgetById(id);
    if (!existing) {
      throw new Error('Budget not found');
    }

    this.dataHandler.deleteBudget(id);
    return this.simulateNetworkDelay(undefined);
  }

  // ===== Pricing Rules =====
  async getPricingRules(): Promise<PricingRule[]> {
    return this.simulateNetworkDelay(this.dataHandler.getPricingRules());
  }

  async createPricingRule(rule: Omit<PricingRule, 'id' | 'created_at' | 'updated_at'>): Promise<PricingRule> {
    const newRule: PricingRule = {
      ...rule,
      id: this.generateId(),
      created_at: this.getCurrentTimestamp(),
      updated_at: this.getCurrentTimestamp()
    };

    this.dataHandler.addPricingRule(newRule);
    return this.simulateNetworkDelay(newRule);
  }

  async updatePricingRule(id: string, rule: Partial<PricingRule>): Promise<PricingRule> {
    const existing = this.dataHandler.getPricingRuleById(id);
    if (!existing) {
      throw new Error('Pricing rule not found');
    }

    const updated = { ...existing, ...rule, updated_at: this.getCurrentTimestamp() };
    this.dataHandler.updatePricingRule(id, updated);
    return this.simulateNetworkDelay(updated);
  }

  async deletePricingRule(id: string): Promise<void> {
    const existing = this.dataHandler.getPricingRuleById(id);
    if (!existing) {
      throw new Error('Pricing rule not found');
    }

    this.dataHandler.deletePricingRule(id);
    return this.simulateNetworkDelay(undefined);
  }

  // ===== Alert Rules =====
  async getAlertRules(): Promise<AlertRule[]> {
    return this.simulateNetworkDelay(this.dataHandler.getAlertRules());
  }

  async createAlertRule(rule: Omit<AlertRule, 'id' | 'created_at' | 'updated_at'>): Promise<AlertRule> {
    const newRule: AlertRule = {
      ...rule,
      id: this.generateId(),
      created_at: this.getCurrentTimestamp(),
      updated_at: this.getCurrentTimestamp()
    };

    this.dataHandler.addAlertRule(newRule);
    return this.simulateNetworkDelay(newRule);
  }

  async updateAlertRule(id: string, rule: Partial<AlertRule>): Promise<AlertRule> {
    const existing = this.dataHandler.getAlertRuleById(id);
    if (!existing) {
      throw new Error('Alert rule not found');
    }

    const updated = { ...existing, ...rule, updated_at: this.getCurrentTimestamp() };
    this.dataHandler.updateAlertRule(id, updated);
    return this.simulateNetworkDelay(updated);
  }

  async deleteAlertRule(id: string): Promise<void> {
    const existing = this.dataHandler.getAlertRuleById(id);
    if (!existing) {
      throw new Error('Alert rule not found');
    }

    this.dataHandler.deleteAlertRule(id);
    return this.simulateNetworkDelay(undefined);
  }

  // ===== Alerts =====
  async getAlerts(): Promise<Alert[]> {
    return this.simulateNetworkDelay(this.dataHandler.getAlerts());
  }

  async acknowledgeAlert(id: string): Promise<void> {
    const existing = this.dataHandler.getAlertById(id);
    if (!existing) {
      throw new Error('Alert not found');
    }

    this.dataHandler.updateAlert(id, {
      status: 'acknowledged',
      acknowledged_at: this.getCurrentTimestamp()
    });
    
    return this.simulateNetworkDelay(undefined);
  }

  async getGroupMembers(groupId: string): Promise<User[]> {
    const group = this.dataHandler.getGroupById(groupId);
    if (!group) return [];
    
    return this.simulateNetworkDelay(this.dataHandler.getUsers().filter((_, idx) => idx < 2));
  }

  async getGroupPermissions(groupId: string): Promise<Permission[]> {
    const permissions = this.dataHandler.getPermissions().filter(p => p.group_id === groupId);
    return this.simulateNetworkDelay(permissions);
  }

  async getUserGroups(_userId: string): Promise<Group[]> {
    return this.simulateNetworkDelay(this.dataHandler.getGroups().filter((_, idx) => idx < 1));
  }

  // ===== Health =====
  async getProviderHealth(): Promise<ProviderHealth[]> {
    return this.simulateNetworkDelay(this.dataHandler.getProviderHealth());
  }
}

export default MockAPIClient;
