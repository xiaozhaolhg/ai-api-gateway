import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { RealAPIClient } from '../client';

describe('RealAPIClient 401 Handling', () => {
  let client: RealAPIClient;
  let originalFetch: typeof fetch;

  beforeEach(() => {
    // Mock window object for Node.js test environment
    global.window = {
      fetch: vi.fn(),
    } as any;
    originalFetch = global.fetch;
    global.fetch = global.window.fetch;
    client = new RealAPIClient('http://localhost:8080');
    localStorage.clear();
  });

  afterEach(() => {
    global.fetch = originalFetch;
    vi.restoreAllMocks();
  });

  it('calls onUnauthorized callback when 401 response received', async () => {
    const mockCallback = vi.fn();
    client.setOnUnauthorized(mockCallback);

    const mockResponse = {
      ok: false,
      status: 401,
      statusText: 'Unauthorized',
      json: () => Promise.resolve({ error: 'Unauthorized' }),
      headers: new Headers({ 'Content-Type': 'application/json' })
    };

    global.fetch = vi.fn().mockResolvedValue(mockResponse);

    try {
      await client.getCurrentUser();
    } catch (error) {}

    expect(mockCallback).toHaveBeenCalledTimes(1);
  });

  it('does not call onUnauthorized for non-401 errors', async () => {
    const mockCallback = vi.fn();
    client.setOnUnauthorized(mockCallback);

    const mockResponse = {
      ok: false,
      status: 500,
      statusText: 'Internal Server Error',
      json: () => Promise.resolve({ error: 'Server Error' }),
      headers: new Headers({ 'Content-Type': 'application/json' })
    };

    global.fetch = vi.fn().mockResolvedValue(mockResponse);

    try {
      await client.getCurrentUser();
    } catch (error) {}

    expect(mockCallback).not.toHaveBeenCalled();
  });

  it('does not call onUnauthorized for successful responses', async () => {
    const mockCallback = vi.fn();
    client.setOnUnauthorized(mockCallback);

    const mockResponse = {
      ok: true,
      status: 200,
      statusText: 'OK',
      json: () => Promise.resolve({ id: '1', name: 'Test User', email: 'test@example.com', role: 'user' }),
      headers: new Headers({ 'Content-Type': 'application/json' })
    };

    global.fetch = vi.fn().mockResolvedValue(mockResponse);

    const result = await client.getCurrentUser();
    expect(result).toEqual({ id: '1', name: 'Test User', email: 'test@example.com', role: 'user' });
    expect(mockCallback).not.toHaveBeenCalled();
  });

  it('allows clearing onUnauthorized callback', async () => {
    const mockCallback = vi.fn();
    client.setOnUnauthorized(mockCallback);
    client.setOnUnauthorized(undefined);

    const mockResponse = {
      ok: false,
      status: 401,
      statusText: 'Unauthorized',
      json: () => Promise.resolve({ error: 'Unauthorized' }),
      headers: new Headers({ 'Content-Type': 'application/json' })
    };

    global.fetch = vi.fn().mockResolvedValue(mockResponse);

    try {
      await client.getCurrentUser();
    } catch (error) {}

    expect(mockCallback).not.toHaveBeenCalled();
  });
});
