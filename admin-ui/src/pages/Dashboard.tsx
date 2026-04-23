import React, { useState, useEffect } from 'react';
import { Card, Statistic, Row, Col, Button, Spin, Alert, Space } from 'antd';
import {
  UserOutlined,
  CloudServerOutlined,
  DollarOutlined,
  BellOutlined,
  PlusOutlined,
} from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import { apiClient } from '../api/client';

export const Dashboard: React.FC = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [stats, setStats] = useState({
    totalUsers: 0,
    totalProviders: 0,
    monthlySpend: 0,
    activeAlerts: 0,
  });

  useEffect(() => {
    const fetchDashboardData = async () => {
      try {
        setLoading(true);
        setError(null);

        const [users, providers, alerts] = await Promise.all([
          apiClient.getUsers(),
          apiClient.getProviders(),
          apiClient.getAlerts(),
        ]);

        setStats({
          totalUsers: users.length,
          totalProviders: providers.length,
          monthlySpend: 0, // This would come from a dedicated endpoint
          activeAlerts: alerts.filter((a) => a.status === 'firing').length,
        });
      } catch (err) {
        setError('Failed to load dashboard data');
      } finally {
        setLoading(false);
      }
    };

    fetchDashboardData();
  }, []);

  const quickActions = [
    {
      key: 'add-provider',
      icon: <CloudServerOutlined />,
      label: t('dashboard.quickActions.addProvider'),
      onClick: () => navigate('/providers'),
    },
    {
      key: 'create-user',
      icon: <UserOutlined />,
      label: t('dashboard.quickActions.createUser'),
      onClick: () => navigate('/users'),
    },
    {
      key: 'issue-api-key',
      icon: <PlusOutlined />,
      label: t('dashboard.quickActions.issueApiKey'),
      onClick: () => navigate('/api-keys'),
    },
    {
      key: 'view-alerts',
      icon: <BellOutlined />,
      label: t('dashboard.quickActions.viewAlerts'),
      onClick: () => navigate('/alerts'),
    },
  ];

  if (loading) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', padding: '50px' }}>
        <Spin size="large" />
      </div>
    );
  }

  if (error) {
    return <Alert message={error} type="error" showIcon />;
  }

  return (
    <div>
      <h2 style={{ marginBottom: 24 }}>{t('dashboard.title')}</h2>

      <Row gutter={16} style={{ marginBottom: 24 }}>
        <Col span={6}>
          <Card>
            <Statistic
              title={t('dashboard.summary.totalUsers')}
              value={stats.totalUsers}
              prefix={<UserOutlined />}
              valueStyle={{ color: '#3f8600' }}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title={t('dashboard.summary.totalProviders')}
              value={stats.totalProviders}
              prefix={<CloudServerOutlined />}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title={t('dashboard.summary.monthlySpend')}
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
              title={t('dashboard.summary.activeAlerts')}
              value={stats.activeAlerts}
              prefix={<BellOutlined />}
              valueStyle={{ color: stats.activeAlerts > 0 ? '#cf1322' : '#3f8600' }}
            />
          </Card>
        </Col>
      </Row>

      <Card title={t('dashboard.quickActions.title')} style={{ marginBottom: 24 }}>
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
