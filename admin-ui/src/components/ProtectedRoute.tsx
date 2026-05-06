import React from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import { Result, Button } from 'antd';

type Role = 'admin' | 'user' | 'viewer';

interface ProtectedRouteProps {
  children: React.ReactNode;
  requiredRole?: Role;
}

// Role hierarchy: admin > user > viewer
const roleHierarchy: Record<Role, number> = {
  admin: 3,
  user: 2,
  viewer: 1,
};

const hasPermission = (userRole: string, requiredRole: Role): boolean => {
  if (!userRole || !requiredRole) return false;
  
  // Convert user role to Role type if valid
  const userRoleTyped = userRole as Role;
  if (!roleHierarchy[userRoleTyped]) return false;
  
  return roleHierarchy[userRoleTyped] >= roleHierarchy[requiredRole];
};

export const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ children, requiredRole }) => {
  const { isAuthenticated, user } = useAuth();
  const location = useLocation();

  if (!isAuthenticated) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  if (requiredRole && user?.role && !hasPermission(user.role, requiredRole)) {
    return (
      <Result
        status="403"
        title="403"
        subTitle="Sorry, you do not have permission to access this page."
        extra={<Button type="primary" href="/">Back to Dashboard</Button>}
      />
    );
  }

  return <>{children}</>;
};
