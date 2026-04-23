import { useState, useEffect } from 'react';
import { Table, DatePicker, Input, Button, Statistic, Card, Row, Col, Spin, Empty, message } from 'antd';
import { SearchOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { apiClient, type UsageRecord } from '../api/client';
import dayjs from 'dayjs';

const { RangePicker } = DatePicker;

export default function Usage() {
  const { t } = useTranslation(['usage', 'common']);
  const [usage, setUsage] = useState<UsageRecord[]>([]);
  const [loading, setLoading] = useState(true);
  const [filters, setFilters] = useState({
    userId: '',
    model: '',
    provider: '',
    startDate: null as dayjs.Dayjs | null,
    endDate: null as dayjs.Dayjs | null,
  });

  useEffect(() => {
    loadUsage();
  }, []);

  const loadUsage = async () => {
    try {
      setLoading(true);
      const data = await apiClient.getUsage(
        filters.userId || undefined,
        filters.startDate?.toISOString() || undefined,
        filters.endDate?.toISOString() || undefined
      );
      setUsage(data);
    } catch (error) {
      message.error('Failed to load usage data');
    } finally {
      setLoading(false);
    }
  };

  const handleFilter = () => {
    loadUsage();
  };

  const totalTokens = usage.reduce((sum, record) => sum + record.total_tokens, 0);
  const totalCost = usage.reduce((sum, record) => sum + record.cost, 0);

  const columns = [
    {
      title: t('usage:fields.user'),
      dataIndex: 'user_id',
      key: 'user_id',
    },
    {
      title: t('usage:fields.model'),
      dataIndex: 'model',
      key: 'model',
    },
    {
      title: t('usage:fields.provider'),
      dataIndex: 'provider',
      key: 'provider',
    },
    {
      title: t('usage:fields.promptTokens'),
      dataIndex: 'prompt_tokens',
      key: 'prompt_tokens',
    },
    {
      title: t('usage:fields.completionTokens'),
      dataIndex: 'completion_tokens',
      key: 'completion_tokens',
    },
    {
      title: t('usage:fields.totalTokens'),
      dataIndex: 'total_tokens',
      key: 'total_tokens',
    },
    {
      title: t('usage:fields.cost'),
      dataIndex: 'cost',
      key: 'cost',
      render: (cost: number) => `$${cost.toFixed(4)}`,
    },
    {
      title: t('common:timestamp'),
      dataIndex: 'timestamp',
      key: 'timestamp',
      render: (timestamp: string) => new Date(timestamp).toLocaleString(),
    },
  ];

  if (loading) {
    return <Spin size="large" />;
  }

  return (
    <div>
      <h2 style={{ marginBottom: 16 }}>{t('usage:title')}</h2>

      <Card style={{ marginBottom: 16 }}>
        <Row gutter={16} style={{ marginBottom: 12 }}>
          <Col span={6}>
            <Input
              placeholder={t('usage:fields.userId')}
              value={filters.userId}
              onChange={(e) => setFilters({ ...filters, userId: e.target.value })}
            />
          </Col>
          <Col span={6}>
            <Input
              placeholder={t('usage:fields.model')}
              value={filters.model}
              onChange={(e) => setFilters({ ...filters, model: e.target.value })}
            />
          </Col>
          <Col span={6}>
            <Input
              placeholder={t('usage:fields.provider')}
              value={filters.provider}
              onChange={(e) => setFilters({ ...filters, provider: e.target.value })}
            />
          </Col>
        </Row>
        <Row gutter={16}>
          <Col span={16}>
            <RangePicker
              value={filters.startDate && filters.endDate ? [filters.startDate, filters.endDate] : null}
              onChange={(dates) => setFilters({
                ...filters,
                startDate: dates?.[0] || null,
                endDate: dates?.[1] || null,
              })}
              style={{ width: '100%' }}
            />
          </Col>
          <Col span={8}>
            <Button type="primary" icon={<SearchOutlined />} onClick={handleFilter} block>
              {t('common:search')}
            </Button>
          </Col>
        </Row>
      </Card>

      <Row gutter={16} style={{ marginBottom: 16 }}>
        <Col span={8}>
          <Card>
            <Statistic title={t('usage:stats.totalRequests')} value={usage.length} />
          </Card>
        </Col>
        <Col span={8}>
          <Card>
            <Statistic title={t('usage:stats.totalTokens')} value={totalTokens} />
          </Card>
        </Col>
        <Col span={8}>
          <Card>
            <Statistic title={t('usage:stats.totalCost')} value={totalCost} precision={2} prefix="$" />
          </Card>
        </Col>
      </Row>

      {usage.length === 0 ? (
        <Empty description="No usage records found" />
      ) : (
        <Table
          dataSource={usage}
          columns={columns}
          rowKey="id"
          pagination={{ pageSize: 20 }}
        />
      )}
    </div>
  );
}
