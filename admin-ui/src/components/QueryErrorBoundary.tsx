import React from 'react';
import { QueryErrorResetBoundary, useQueryErrorResetBoundary } from '@tanstack/react-query';
import { Result, Button } from 'antd';
import { WarningOutlined } from '@ant-design/icons';

interface QueryErrorBoundaryProps {
  children: React.ReactNode;
}

const ErrorFallback: React.FC<{ error: Error; reset: () => void }> = ({ error, reset }) => {
  return (
    <Result
      status="error"
      title="Data Loading Failed"
      subTitle={error?.message || 'An unexpected error occurred while loading data.'}
      extra={[
        <Button type="primary" key="retry" onClick={() => reset()}>
          Retry
        </Button>,
      ]}
    />
  );
};

export const QueryErrorBoundary: React.FC<QueryErrorBoundaryProps> = ({ children }) => {
  return (
    <QueryErrorResetBoundary>
      {({ reset }) => (
        <ErrorBoundary
          fallbackRender={({ error }) => <ErrorFallback error={error as Error} reset={reset} />}
        >
          {children}
        </ErrorBoundary>
      )}
    </QueryErrorResetBoundary>
  );
};
