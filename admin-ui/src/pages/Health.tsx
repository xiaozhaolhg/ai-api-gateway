import { useState } from 'react';
import { Table, Tag, Empty, Card, Statistic, Row, Col, Switch, InputNumber, Space, Button } from 'antd';
import { CheckCircleOutlined, CloseCircleOutlined, ClockCircleOutlined, ReloadOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useQuery } from '@tanstack/react-query';
import { apiClient } from '../api/client';

export default function Health() {
  const { t } = useTranslation(['health', 'common']);
  const [autoRefresh, setAutoRefresh] = useState(true);
  const [refreshInterval, setRefreshInterval] = useState(30);

  const { data: health = [], isLoading, refetch } = useQuery({
    queryKey: ['health'],
    queryFn: () => apiClient.getProviderHealth(),
    refetchInterval: autoRefresh ? refreshInterval * 1000 : false,
  });

  const handleManualRefresh = () => {
    refetch();
  };

  const getStatusTag = (status: string) => {
    switch (status) {
      case 'healthy':
        return <Tag icon={<CheckCircleOutlined />} color="success">{t('health:status.healthy')}</Tag>;
      case 'unhealthy':
        return <Tag icon={<CloseCircleOutlined />} color="error">{t('health:status.unhealthy')}</Tag>;
      default:
        return <Tag icon={<ClockCircleOutlined />} color="default">{t('health:status.unknown')}</Tag>;
    }
  };

  const columns = [
    {
      title: t('health:fields.provider'),
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: t('health:fields.status'),
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => getStatusTag(status),
    },
    {
      title: t('health:fields.latency'),
      dataIndex: 'latency_ms',
      key: 'latency_ms',
      render: (latency: number) => `${latency}ms`,
    },
    {
      title: t('health:fields.errorRate'),
      dataIndex: 'error_rate',
      key: 'error_rate',
      render: (errorRate: number) => `${(errorRate * 100).toFixed(1)}%`,
    },
    {
      title: t('health:fields.lastCheck'),
      dataIndex: 'last_check',
      key: 'last_check',
      render: (lastCheck: string) => new Date(lastCheck).toLocaleString(),
    },
  ];

  const healthyCount = health.filter(h => h.status === 'healthy').length;
  const unhealthyCount = health.filter(h => h.status === 'unhealthy').length;
  const avgLatency = health.length > 0
    ? health.reduce((sum, h) => sum + h.latency_ms, 0) / health.length
    : 0;

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
        <h2 style={{ margin: 0 }}>{t('health:title')}</h2>
        <Space>
          <Button icon={<ReloadOutlined />} onClick={handleManualRefresh} loading={isLoading}>
            {t('common:refresh')}
          </Button>
          <Space size="middle">
            <span>Auto-refresh:</span>
            <Switch checked={autoRefresh} onChange={setAutoRefresh} />
            <InputNumber
              min={5}
              max={300}
              value={refreshInterval}
              onChange={(value) => setRefreshInterval(value || 30)}
              disabled={!autoRefresh}
              addonAfter="s"
              style={{ width: 100 }}
            />
          </Space>
        </Space>
      </div>

      <Row gutter={16} style={{ marginBottom: 16 }}>
        <Col span={8}>
          <Card>
            <Statistic
              title={t('health:stats.healthyProviders')}
              value={healthyCount}
              valueStyle={{ color: '#52c41a' }}
              prefix={<CheckCircleOutlined />}
            />
          </Card>
        </Col>
        <Col span={8}>
          <Card>
            <Statistic
              title={t('health:stats.unhealthyProviders')}
              value={unhealthyCount}
              valueStyle={{ color: '#ff4d4f' }}
              prefix={<CloseCircleOutlined />}
            />
          </Card>
        </Col>
        <Col span={8}>
          <Card>
            <Statistic
              title={t('health:stats.avgLatency')}
              value={avgLatency}
              precision={0}
              suffix="ms"
            />
          </Card>
        </Col>
      </Row>

      <Table
        dataSource={health}
        columns={columns}
        rowKey="id"
        loading={isLoading}
        pagination={false}
        locale={{ emptyText: <Empty description="No provider health data found" /> }}
      />
    </div>
  );
}
