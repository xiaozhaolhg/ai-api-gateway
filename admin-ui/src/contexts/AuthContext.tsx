import React, { createContext, useContext, useState, useEffect } from 'react';
import type { ReactNode } from 'react';
import { message } from 'antd';
import { apiClient } from '../api/client';

interface User {
  id: string;
  name: string;
  email: string;
  role: string;
}

interface AuthContextType {
  user: User | null;
  token: string | null;
  login: (username: string, password: string) => Promise<void>;
  register: (name: string, username: string, email: string, password: string) => Promise<void>;
  logout: () => void;
  isAuthenticated: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

const TOKEN_KEY = 'auth_token';
const USER_KEY = 'auth_user';

export const AuthProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  const checkTokenExpiry = (tokenToCheck?: string) => {
    try {
      const tokenToUse = tokenToCheck || token;
      if (!tokenToUse) return;
      
      const payload = JSON.parse(atob(tokenToUse.split('.')[1]));
      const exp = payload.exp;
      const now = Date.now();
      const bufferTime = 30 * 1000; // 30 seconds buffer
      const isExpired = exp * 1000 <= now - bufferTime;
      
      if (isExpired) {
        logout();
        message.warning('Your session has expired. Please login again.');
      }
    } catch (error) {
      // Token may not be a JWT (e.g. opaque token) — skip client-side expiry check.
      // Server will reject truly invalid tokens via 401 → onUnauthorized → logout.
      console.debug('Token expiry check skipped:', error);
    }
  };

  // Check if server-side session is still valid by making a test request
  const checkServerSession = async () => {
    try {
      // Make a simple request to check if session is valid
      await apiClient.getCurrentUser();
    } catch (error) {
      // If we get a 401, server session is invalid
      if (error instanceof Error && error.message.includes('401')) {
        logout();
      }
    }
  };

  useEffect(() => {
    const storedToken = localStorage.getItem(TOKEN_KEY);
    const storedUser = localStorage.getItem(USER_KEY);

    if (storedToken && storedUser) {
      try {
        setToken(storedToken);
        setUser(JSON.parse(storedUser));
      } catch (error) {
        console.error('Failed to parse stored user:', error);
        localStorage.removeItem(TOKEN_KEY);
        localStorage.removeItem(USER_KEY);
      }
      // Initial expiry check - wrap in try/catch to not affect token restoration
      try {
        checkTokenExpiry(storedToken);
      } catch (error) {
        console.error('checkTokenExpiry failed:', error);
      }
    }
    // Always set loading to false, even if no stored credentials
    setLoading(false);
  }, []);

  useEffect(() => {
    if (!token) return;
    const interval = setInterval(async () => {
      // Check token expiry first
      checkTokenExpiry(token);
      // Then check if server session is still valid
      await checkServerSession();
    }, 60000);
    return () => clearInterval(interval);
  }, [token]);

  useEffect(() => {
    apiClient.setOnUnauthorized(() => {
      logout();
    });
    return () => {
      apiClient.setOnUnauthorized(undefined);
    };
  }, []);

  const login = async (email: string, password: string) => {
    const data = await apiClient.login(email, password);
    const { token: newToken, user: newUser } = data;

    const userObj = {
      id: newUser.id,
      name: newUser.name,
      email: newUser.email,
      role: newUser.role,
    };
    setToken(newToken);
    setUser(userObj);
    localStorage.setItem(TOKEN_KEY, newToken);
    localStorage.setItem(USER_KEY, JSON.stringify(userObj));
  };

  const register = async (name: string, username: string, email: string, password: string) => {
    const data = await apiClient.register(name, username, email, password);
    const { token: newToken, user: newUser } = data;

    const userObj = {
      id: newUser.id,
      name: newUser.name,
      email: newUser.email,
      role: newUser.role,
    };
    setToken(newToken);
    setUser(userObj);
    localStorage.setItem(TOKEN_KEY, newToken);
    localStorage.setItem(USER_KEY, JSON.stringify(userObj));
  };

  const logout = () => {
    setToken(null);
    setUser(null);
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(USER_KEY);
  };

  const value: AuthContextType = {
    user,
    token,
    login,
    register,
    logout,
    isAuthenticated: !!token,
  };

  if (loading) {
    return null;
  }

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
