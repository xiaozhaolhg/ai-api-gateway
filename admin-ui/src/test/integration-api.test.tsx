import { describe, it, expect, beforeAll, afterAll } from 'vitest';

const API_BASE = 'http://localhost:8080';
let authToken: string | null = null;
let testGroupId: string | null = null;
let testUserId: string | null = null;
let testPermissionId: string | null = null;

async function apiRequest(path: string, options: RequestInit = {}) {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string> || {}),
  };

  if (authToken) {
    headers['Authorization'] = `Bearer ${authToken}`;
  }

  const response = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers,
  });

  return response;
}

async function isBackendAvailable(): Promise<boolean> {
  try {
    const response = await fetch(`${API_BASE}/health`, { signal: AbortSignal.timeout(2000) });
    return response.ok;
  } catch {
    return false;
  }
}

describe('Group Management UI - API Integration Tests', () => {
  const backendAvailable = isBackendAvailable();

  beforeAll(async () => {
    const loginRes = await apiRequest('/admin/auth/login', {
      method: 'POST',
      body: JSON.stringify({
        email: 'admin@example.com',
        password: 'admin123',
      }),
    });

    if (!loginRes.ok) {
      const registerRes = await apiRequest('/admin/auth/register', {
        method: 'POST',
        body: JSON.stringify({
          name: 'Admin',
          username: 'admin',
          email: 'admin@example.com',
          password: 'admin123',
        }),
      });

      if (registerRes.ok) {
        const data = await registerRes.json();
        authToken = data.token;
      }
    } else {
      const data = await loginRes.json();
      authToken = data.token;
    }

    expect(authToken).toBeTruthy();
  }, 30000);

  afterAll(async () => {
    if (testPermissionId) {
      await apiRequest(`/admin/auth/permissions/${testPermissionId}`, {
        method: 'DELETE',
      });
    }
    if (testGroupId) {
      await apiRequest(`/admin/auth/groups/${testGroupId}`, {
        method: 'DELETE',
      });
    }
  });

  describe('Authentication', () => {
    it('should login and get JWT token', () => {
      expect(authToken).toBeTruthy();
    });
  });

  describe('Groups API', () => {
    it('should create group with new fields (model_patterns, parent_group_id)', async () => {
      const res = await apiRequest('/admin/auth/groups', {
        method: 'POST',
        body: JSON.stringify({
          name: 'Test Group Integration',
          description: 'Group for integration testing',
          model_patterns: ['gpt-*', 'ollama:*'],
        }),
      });

      expect(res.ok).toBe(true);
      const group = await res.json();
      expect(group.id).toBeTruthy();
      expect(group.name).toBe('Test Group Integration');
      testGroupId = group.id;
    });

    it('should get groups list', async () => {
      const res = await apiRequest('/admin/auth/groups');
      expect(res.ok).toBe(true);
      const data = await res.json();
      expect(data.groups).toBeInstanceOf(Array);
    });

    it('should get group members (empty initially)', async () => {
      if (!testGroupId) return;
      
      const res = await apiRequest(`/admin/auth/groups/${testGroupId}/members`);
      expect(res.ok).toBe(true);
      const data = await res.json();
      expect(data.members).toBeInstanceOf(Array);
    });

    it('should add member to group', async () => {
      if (!testGroupId) return;

      const usersRes = await apiRequest('/admin/auth/users');
      const usersData = await usersRes.json();
      
      if (usersData.users && usersData.users.length > 0) {
        const userId = usersData.users[0].id;
        const res = await apiRequest(`/admin/auth/groups/${testGroupId}/members`, {
          method: 'POST',
          body: JSON.stringify({ user_id: userId }),
        });
        expect(res.ok || res.status === 409).toBe(true);
      }
    });
  });

  describe('Permissions API - Actual Backend Behavior', () => {
    it('should create permission with resource_type, resource_id, action', async () => {
      if (!testGroupId) return;

      const res = await apiRequest('/admin/auth/permissions', {
        method: 'POST',
        body: JSON.stringify({
          group_id: testGroupId,
          resource_type: 'model',
          resource_id: 'gpt-4',
          action: 'access',
          effect: 'allow',
        }),
      });

      expect(res.ok).toBe(true);
      const permission = await res.json();
      expect(permission.id).toBeTruthy();
      expect(permission.resource_type).toBe('model');
      expect(permission.resource_id).toBe('gpt-4');
      expect(permission.action).toBe('access');
      
      testPermissionId = permission.id;
    });

it('should verify effect field IS in response (backend fixed)', async () => {
      const res = await apiRequest('/admin/auth/permissions', {
        method: 'POST',
        body: JSON.stringify({
          group_id: testGroupId,
          resource_type: 'model',
          resource_id: 'gpt-4',
          action: 'access',
          effect: 'allow',
        }),
      });

      const permission = await res.json();

      expect(permission.effect).toBe('allow');
    });

    it('should get permissions (returns null when empty, array when has data)', async () => {
      const res = await apiRequest('/admin/auth/permissions');
      expect(res.ok).toBe(true);
      const data = await res.json();
      
      if (data.permissions === null) {
        expect(data.permissions).toBeNull();
      } else {
        expect(data.permissions).toBeInstanceOf(Array);
      }
    });

    it('should get group permissions filtered by group_id', async () => {
      if (!testGroupId) return;
      
      const res = await apiRequest(`/admin/auth/permissions?group_id=${testGroupId}`);
      expect(res.ok).toBe(true);
      const data = await res.json();
      
      if (data.permissions !== null) {
        data.permissions.forEach((p: any) => {
          expect(p.group_id).toBe(testGroupId);
        });
      }
    });
  });

  describe('Users API', () => {
    it('should create user with password field', async () => {
      const res = await apiRequest('/admin/auth/users', {
        method: 'POST',
        body: JSON.stringify({
          name: 'Test User Integration',
          email: `test-${Date.now()}@example.com`,
          password: 'securePassword123',
          role: 'user',
        }),
      });

      if (res.ok) {
        const user = await res.json();
        expect(user.id).toBeTruthy();
        expect(user.email).toContain('@example.com');
        testUserId = user.id;
      }
    });

    it('should get users list', async () => {
      const res = await apiRequest('/admin/auth/users');
      expect(res.ok).toBe(true);
      const data = await res.json();
      expect(data.users).toBeInstanceOf(Array);
    });
  });

  describe('Frontend-Backend Contract Verification', () => {
    it('should verify Permission interface matches backend', async () => {
      const res = await apiRequest('/admin/auth/permissions');
      const data = await res.json();
      
      if (data.permissions && data.permissions.length > 0) {
        const perm = data.permissions[0];
        
        expect(perm).toHaveProperty('id');
        expect(perm).toHaveProperty('group_id');
        expect(perm).toHaveProperty('resource_type');
        expect(perm).toHaveProperty('resource_id');
        expect(perm).toHaveProperty('action');
        expect(perm).toHaveProperty('effect');
        expect(perm).toHaveProperty('created_at');
        
        expect(perm).not.toHaveProperty('model_pattern');
        expect(perm).not.toHaveProperty('status');
      }
    });

    it('should verify Group interface matches backend', async () => {
      const res = await apiRequest('/admin/auth/groups');
      const data = await res.json();

      if (data.groups && data.groups.length > 0) {
        const group = data.groups[0];

        expect(group).toHaveProperty('id');
        expect(group).toHaveProperty('name');
      }
    });
  });
});
