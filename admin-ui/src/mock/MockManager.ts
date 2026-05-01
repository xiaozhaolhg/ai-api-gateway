import MockDataHandler from './handlers/dataHandler';
import { apiClient } from '../api/client';

class MockManager {
  private static instance: MockManager;
  private dataHandler: MockDataHandler;

  private constructor() {
    this.dataHandler = MockDataHandler.getInstance();
  }

  static getInstance(): MockManager {
    if (!MockManager.instance) {
      MockManager.instance = new MockManager();
    }
    return MockManager.instance;
  }

  // Reset mock data to defaults
  resetToDefaults(): void {
    this.dataHandler.resetToDefaults();
    if (typeof window !== 'undefined') {
      window.location.reload();
    }
  }

  // Export mock data as JSON
  exportData(): string {
    return this.dataHandler.exportData();
  }

  // Download mock data as file
  downloadData(): void {
    const data = this.exportData();
    const blob = new Blob([data], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `mock-data-${new Date().toISOString().split('T')[0]}.json`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  }

  // Import mock data from JSON string
  importData(jsonData: string): void {
    try {
      this.dataHandler.importData(jsonData);
      if (typeof window !== 'undefined') {
        window.location.reload();
      }
    } catch (error) {
      throw new Error('Invalid JSON data format');
    }
  }

  // Import mock data from file
  importFromFile(file: File): Promise<void> {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = (e) => {
        try {
          const content = e.target?.result as string;
          this.importData(content);
          resolve();
        } catch (error) {
          reject(error);
        }
      };
      reader.onerror = () => reject(new Error('Failed to read file'));
      reader.readAsText(file);
    });
  }

  // Get current mock mode
  getMockMode(): boolean {
    return (apiClient as any).getMockMode();
  }

  // Set mock mode
  setMockMode(enabled: boolean): void {
    (apiClient as any).setMockMode(enabled);
  }

  // Get mock data statistics
  getDataStats(): Record<string, number> {
    const dataStore = this.dataHandler.getDataStore();
    return {
      users: dataStore.users.length,
      providers: dataStore.providers.length,
      apiKeys: dataStore.apiKeys.length,
      usage: dataStore.usage.length,
      routingRules: dataStore.routingRules.length,
      groups: dataStore.groups.length,
      permissions: dataStore.permissions.length,
      budgets: dataStore.budgets.length,
      pricingRules: dataStore.pricingRules.length,
      alertRules: dataStore.alertRules.length,
      alerts: dataStore.alerts.length,
    };
  }

  // Check if mock mode is available
  isMockAvailable(): boolean {
    return typeof window !== 'undefined' && typeof localStorage !== 'undefined';
  }
}

export default MockManager;
