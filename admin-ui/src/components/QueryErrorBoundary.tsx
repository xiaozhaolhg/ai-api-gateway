import React from 'react';
import { QueryErrorResetBoundary } from '@tanstack/react-query';

interface QueryErrorBoundaryProps {
  children: React.ReactNode;
}

export const QueryErrorBoundary: React.FC<QueryErrorBoundaryProps> = ({ children }) => {
  return (
    <QueryErrorResetBoundary>
      {() => (
        <div>
          {children}
        </div>
      )}
    </QueryErrorResetBoundary>
  );
};
