import { useState, useEffect } from 'react';
import { Table, Tag, Spin, Empty, Card, Statistic, Row, Col, message, Switch, InputNumber, Space, Button } from 'antd';
import { CheckCircleOutlined, CloseCircleOutlined, ClockCircleOutlined, ReloadOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { apiClient, type ProviderHealth } from '../api/client';

export default function Health() {
  const { t } = useTranslation(['health', 'common']);
  const [health, setHealth] = useState<ProviderHealth[]>([]);
  const [loading, setLoading] = useState(true);
  const [autoRefresh, setAutoRefresh] = useState(true);
  const [refreshInterval, setRefreshInterval] = useState(30); // seconds

  useEffect(() => {
    loadHealth(true);
  }, []);

  useEffect(() => {
    let intervalId: ReturnType<typeof setInterval> | null = null;

    if (autoRefresh) {
      intervalId = setInterval(() => {
        loadHealth(false);
      }, refreshInterval * 1000);
    }

    return () => {
      if (intervalId) {
        clearInterval(intervalId);
      }
    };
  }, [autoRefresh, refreshInterval]);

  const loadHealth = async (isInitial: boolean = false) => {
    try {
      if (isInitial) {
        setLoading(true);
      }
      const data = await apiClient.getProviderHealth();
      setHealth(data);
    } catch (error) {
      message.error('Failed to load provider health data');
    } finally {
      if (isInitial) {
        setLoading(false);
      }
    }
  };

  const handleManualRefresh = () => {
    loadHealth(true);
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

  if (loading) {
    return <Spin size="large" />;
  }

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
        <h2 style={{ margin: 0 }}>{t('health:title')}</h2>
        <Space>
          <Button icon={<ReloadOutlined />} onClick={handleManualRefresh} loading={loading}>
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

      {health.length === 0 ? (
        <Empty description="No provider health data found" />
      ) : (
        <Table
          dataSource={health}
          columns={columns}
          rowKey="id"
          pagination={false}
        />
      )}
    </div>
  );
}
