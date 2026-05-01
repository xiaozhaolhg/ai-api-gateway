import type {
  MockDataStore,
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
} from '../../api/types';
import { defaultMockData } from '../data';

const MOCK_DATA_STORAGE_KEY = 'mockDataStore';

class MockDataHandler {
  private static instance: MockDataHandler;
  private dataStore: MockDataStore;

  private constructor() {
    this.dataStore = this.loadData();
  }

  static getInstance(): MockDataHandler {
    if (!MockDataHandler.instance) {
      MockDataHandler.instance = new MockDataHandler();
    }
    return MockDataHandler.instance;
  }

  private loadData(): MockDataStore {
    if (typeof window === 'undefined') {
      return JSON.parse(JSON.stringify(defaultMockData));
    }

    const stored = localStorage.getItem(MOCK_DATA_STORAGE_KEY);
    if (stored) {
      try {
        return JSON.parse(stored);
      } catch (error) {
        console.error('Failed to parse stored mock data:', error);
        return JSON.parse(JSON.stringify(defaultMockData));
      }
    }
    return JSON.parse(JSON.stringify(defaultMockData));
  }

  private saveData(): void {
    if (typeof window !== 'undefined') {
      localStorage.setItem(MOCK_DATA_STORAGE_KEY, JSON.stringify(this.dataStore));
    }
  }

  // Reset to default data
  resetToDefaults(): void {
    this.dataStore = JSON.parse(JSON.stringify(defaultMockData));
    this.saveData();
  }

  // Export data
  exportData(): string {
    return JSON.stringify(this.dataStore, null, 2);
  }

  // Import data
  importData(jsonData: string): void {
    try {
      const parsed = JSON.parse(jsonData);
      this.dataStore = parsed;
      this.saveData();
    } catch (error) {
      throw new Error('Invalid JSON data');
    }
  }

  // Get entire data store
  getDataStore(): MockDataStore {
    return JSON.parse(JSON.stringify(this.dataStore));
  }

  // ===== User Operations =====
  getUsers(): User[] {
    return [...this.dataStore.users];
  }

  getUserById(id: string): User | undefined {
    return this.dataStore.users.find(u => u.id === id);
  }

  addUser(user: User): void {
    this.dataStore.users.push(user);
    this.saveData();
  }

  updateUser(id: string, updates: Partial<User>): void {
    const index = this.dataStore.users.findIndex(u => u.id === id);
    if (index !== -1) {
      this.dataStore.users[index] = { ...this.dataStore.users[index], ...updates };
      this.saveData();
    }
  }

  deleteUser(id: string): void {
    this.dataStore.users = this.dataStore.users.filter(u => u.id !== id);
    this.saveData();
  }

  // ===== Provider Operations =====
  getProviders(): Provider[] {
    return [...this.dataStore.providers];
  }

  getProviderById(id: string): Provider | undefined {
    return this.dataStore.providers.find(p => p.id === id);
  }

  addProvider(provider: Provider): void {
    this.dataStore.providers.push(provider);
    this.saveData();
  }

  updateProvider(id: string, updates: Partial<Provider>): void {
    const index = this.dataStore.providers.findIndex(p => p.id === id);
    if (index !== -1) {
      this.dataStore.providers[index] = { ...this.dataStore.providers[index], ...updates };
      this.saveData();
    }
  }

  deleteProvider(id: string): void {
    this.dataStore.providers = this.dataStore.providers.filter(p => p.id !== id);
    this.saveData();
  }

  // ===== API Key Operations =====
  getAPIKeys(userId?: string): APIKey[] {
    if (userId) {
      return this.dataStore.apiKeys.filter(k => k.user_id === userId);
    }
    return [...this.dataStore.apiKeys];
  }

  getAPIKeyById(id: string): APIKey | undefined {
    return this.dataStore.apiKeys.find(k => k.id === id);
  }

  addAPIKey(apiKey: APIKey): void {
    this.dataStore.apiKeys.push(apiKey);
    this.saveData();
  }

  deleteAPIKey(id: string): void {
    this.dataStore.apiKeys = this.dataStore.apiKeys.filter(k => k.id !== id);
    this.saveData();
  }

  // ===== Usage Operations =====
  getUsage(userId?: string, startDate?: string, endDate?: string): UsageRecord[] {
    let records = [...this.dataStore.usage];

    if (userId) {
      records = records.filter(r => r.user_id === userId);
    }

    if (startDate) {
      records = records.filter(r => r.timestamp >= startDate);
    }

    if (endDate) {
      records = records.filter(r => r.timestamp <= endDate);
    }

    return records;
  }

  addUsageRecord(record: UsageRecord): void {
    this.dataStore.usage.push(record);
    this.saveData();
  }

  // ===== Routing Rule Operations =====
  getRoutingRules(): RoutingRule[] {
    return [...this.dataStore.routingRules];
  }

  getRoutingRuleById(id: string): RoutingRule | undefined {
    return this.dataStore.routingRules.find(r => r.id === id);
  }

  addRoutingRule(rule: RoutingRule): void {
    this.dataStore.routingRules.push(rule);
    this.saveData();
  }

  updateRoutingRule(id: string, updates: Partial<RoutingRule>): void {
    const index = this.dataStore.routingRules.findIndex(r => r.id === id);
    if (index !== -1) {
      this.dataStore.routingRules[index] = { ...this.dataStore.routingRules[index], ...updates };
      this.saveData();
    }
  }

  deleteRoutingRule(id: string): void {
    this.dataStore.routingRules = this.dataStore.routingRules.filter(r => r.id !== id);
    this.saveData();
  }

  // ===== Group Operations =====
  getGroups(): Group[] {
    return [...this.dataStore.groups];
  }

  getGroupById(id: string): Group | undefined {
    return this.dataStore.groups.find(g => g.id === id);
  }

  addGroup(group: Group): void {
    this.dataStore.groups.push(group);
    this.saveData();
  }

  updateGroup(id: string, updates: Partial<Group>): void {
    const index = this.dataStore.groups.findIndex(g => g.id === id);
    if (index !== -1) {
      this.dataStore.groups[index] = { ...this.dataStore.groups[index], ...updates };
      this.saveData();
    }
  }

  deleteGroup(id: string): void {
    this.dataStore.groups = this.dataStore.groups.filter(g => g.id !== id);
    this.saveData();
  }

  // ===== Permission Operations =====
  getPermissions(groupId?: string): Permission[] {
    if (groupId) {
      return this.dataStore.permissions.filter(p => p.group_id === groupId);
    }
    return [...this.dataStore.permissions];
  }

  getPermissionById(id: string): Permission | undefined {
    return this.dataStore.permissions.find(p => p.id === id);
  }

  addPermission(permission: Permission): void {
    this.dataStore.permissions.push(permission);
    this.saveData();
  }

  updatePermission(id: string, updates: Partial<Permission>): void {
    const index = this.dataStore.permissions.findIndex(p => p.id === id);
    if (index !== -1) {
      this.dataStore.permissions[index] = { ...this.dataStore.permissions[index], ...updates };
      this.saveData();
    }
  }

  deletePermission(id: string): void {
    this.dataStore.permissions = this.dataStore.permissions.filter(p => p.id !== id);
    this.saveData();
  }

  // ===== Budget Operations =====
  getBudgets(): Budget[] {
    return [...this.dataStore.budgets];
  }

  getBudgetById(id: string): Budget | undefined {
    return this.dataStore.budgets.find(b => b.id === id);
  }

  addBudget(budget: Budget): void {
    this.dataStore.budgets.push(budget);
    this.saveData();
  }

  updateBudget(id: string, updates: Partial<Budget>): void {
    const index = this.dataStore.budgets.findIndex(b => b.id === id);
    if (index !== -1) {
      this.dataStore.budgets[index] = { ...this.dataStore.budgets[index], ...updates };
      this.saveData();
    }
  }

  deleteBudget(id: string): void {
    this.dataStore.budgets = this.dataStore.budgets.filter(b => b.id !== id);
    this.saveData();
  }

  // ===== Pricing Rule Operations =====
  getPricingRules(): PricingRule[] {
    return [...this.dataStore.pricingRules];
  }

  getPricingRuleById(id: string): PricingRule | undefined {
    return this.dataStore.pricingRules.find(p => p.id === id);
  }

  addPricingRule(rule: PricingRule): void {
    this.dataStore.pricingRules.push(rule);
    this.saveData();
  }

  updatePricingRule(id: string, updates: Partial<PricingRule>): void {
    const index = this.dataStore.pricingRules.findIndex(p => p.id === id);
    if (index !== -1) {
      this.dataStore.pricingRules[index] = { ...this.dataStore.pricingRules[index], ...updates };
      this.saveData();
    }
  }

  deletePricingRule(id: string): void {
    this.dataStore.pricingRules = this.dataStore.pricingRules.filter(p => p.id !== id);
    this.saveData();
  }

  // ===== Alert Rule Operations =====
  getAlertRules(): AlertRule[] {
    return [...this.dataStore.alertRules];
  }

  getAlertRuleById(id: string): AlertRule | undefined {
    return this.dataStore.alertRules.find(a => a.id === id);
  }

  addAlertRule(rule: AlertRule): void {
    this.dataStore.alertRules.push(rule);
    this.saveData();
  }

  updateAlertRule(id: string, updates: Partial<AlertRule>): void {
    const index = this.dataStore.alertRules.findIndex(a => a.id === id);
    if (index !== -1) {
      this.dataStore.alertRules[index] = { ...this.dataStore.alertRules[index], ...updates };
      this.saveData();
    }
  }

  deleteAlertRule(id: string): void {
    this.dataStore.alertRules = this.dataStore.alertRules.filter(a => a.id !== id);
    this.saveData();
  }

  // ===== Alert Operations =====
  getAlerts(): Alert[] {
    return [...this.dataStore.alerts];
  }

  getAlertById(id: string): Alert | undefined {
    return this.dataStore.alerts.find(a => a.id === id);
  }

  addAlert(alert: Alert): void {
    this.dataStore.alerts.push(alert);
    this.saveData();
  }

  updateAlert(id: string, updates: Partial<Alert>): void {
    const index = this.dataStore.alerts.findIndex(a => a.id === id);
    if (index !== -1) {
      this.dataStore.alerts[index] = { ...this.dataStore.alerts[index], ...updates };
      this.saveData();
    }
  }

  deleteAlert(id: string): void {
    this.dataStore.alerts = this.dataStore.alerts.filter(a => a.id !== id);
    this.saveData();
  }

  // ===== Provider Health Operations =====
  getProviderHealth(): ProviderHealth[] {
    return [...this.dataStore.providerHealth];
  }

  updateProviderHealth(id: string, updates: Partial<ProviderHealth>): void {
    const index = this.dataStore.providerHealth.findIndex(p => p.id === id);
    if (index !== -1) {
      this.dataStore.providerHealth[index] = { ...this.dataStore.providerHealth[index], ...updates };
      this.saveData();
    }
  }
}

export default MockDataHandler;
