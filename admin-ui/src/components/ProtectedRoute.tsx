import React from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import { Result, Button } from 'antd';

type Role = 'admin' | 'user' | 'viewer';

interface ProtectedRouteProps {
  children: React.ReactNode;
  requiredRole?: Role;
}

export const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ children, requiredRole }) => {
  const { isAuthenticated, user } = useAuth();
  const location = useLocation();

  if (!isAuthenticated) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  if (requiredRole && user?.role !== requiredRole && user?.role !== 'admin') {
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
