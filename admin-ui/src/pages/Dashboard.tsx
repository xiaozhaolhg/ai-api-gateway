import React from 'react';
import { Card, Statistic, Row, Col, Button, Alert, Space, Spin } from 'antd';
import {
  UserOutlined,
  CloudServerOutlined,
  DollarOutlined,
  BellOutlined,
  PlusOutlined,
} from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { apiClient } from '../api/client';
import { useAuth } from '../contexts/AuthContext';

export const Dashboard: React.FC = () => {
  const { t } = useTranslation(['dashboard', 'common']);
  const navigate = useNavigate();
  const { user } = useAuth();

  const { data: users = [], isLoading: usersLoading, error: usersError } = useQuery({
    queryKey: ['users'],
    queryFn: () => apiClient.getUsers(),
  });

  const { data: providers = [], isLoading: providersLoading } = useQuery({
    queryKey: ['providers'],
    queryFn: () => apiClient.getProviders(),
  });

  const { data: alerts = [], isLoading: alertsLoading } = useQuery({
    queryKey: ['alerts'],
    queryFn: () => apiClient.getAlerts(),
  });

  const isLoading = usersLoading || providersLoading || alertsLoading;

  const safeAlerts = Array.isArray(alerts) ? alerts : [];

  const stats = {
    totalUsers: users.length,
    totalProviders: providers.length,
    monthlySpend: 0,
    activeAlerts: safeAlerts.filter((a) => a.status === 'firing').length,
  };

  const userRole = user?.role || 'viewer';

  const quickActions = [
    ...(userRole === 'admin' ? [
      {
        key: 'add-provider',
        icon: <CloudServerOutlined />,
        label: t('quickActions.addProvider', { ns: 'dashboard' }),
        onClick: () => navigate('/providers'),
      },
      {
        key: 'create-user',
        icon: <UserOutlined />,
        label: t('quickActions.createUser', { ns: 'dashboard' }),
        onClick: () => navigate('/users'),
      },
    ] : []),
    {
      key: 'issue-api-key',
      icon: <PlusOutlined />,
      label: t('quickActions.issueApiKey', { ns: 'dashboard' }),
      onClick: () => navigate('/api-keys'),
    },
    ...(userRole === 'admin' ? [
      {
        key: 'view-alerts',
        icon: <BellOutlined />,
        label: t('quickActions.viewAlerts', { ns: 'dashboard' }),
        onClick: () => navigate('/alerts'),
      },
    ] : []),
  ];

  if (isLoading) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', padding: '50px' }}>
        <Spin size="large" />
      </div>
    );
  }

  if (usersError) {
    return <Alert message="Failed to load dashboard data" type="error" showIcon />;
  }

  return (
    <div>
      <h2 style={{ marginBottom: 24 }}>{t('title', { ns: 'dashboard' })}</h2>

      <Row gutter={16} style={{ marginBottom: 24 }}>
        <Col span={6}>
          <Card>
            <Statistic
              title={t('summary.totalUsers', { ns: 'dashboard' })}
              value={stats.totalUsers}
              prefix={<UserOutlined />}
              valueStyle={{ color: '#3f8600' }}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title={t('summary.totalProviders', { ns: 'dashboard' })}
              value={stats.totalProviders}
              prefix={<CloudServerOutlined />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title={t('summary.monthlySpend', { ns: 'dashboard' })}
              value={stats.monthlySpend}
              prefix={<DollarOutlined />}
              precision={2}
              valueStyle={{ color: '#cf1322' }}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title={t('summary.activeAlerts', { ns: 'dashboard' })}
              value={stats.activeAlerts}
              prefix={<BellOutlined />}
              valueStyle={{ color: stats.activeAlerts > 0 ? '#cf1322' : '#3f8600' }}
            />
          </Card>
        </Col>
      </Row>

      <Card title={t('quickActions.title', { ns: 'dashboard' })} style={{ marginBottom: 24 }}>
        <Space size="middle">
          {quickActions.map((action) => (
            <Button
              key={action.key}
              icon={action.icon}
              onClick={action.onClick}
              size="large"
            >
              {action.label}
            </Button>
          ))}
        </Space>
      </Card>
    </div>
  );
};
