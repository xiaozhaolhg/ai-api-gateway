import { useState } from 'react';
import { Table, DatePicker, Input, Button, Statistic, Card, Row, Col, Empty } from 'antd';
import { SearchOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useQuery } from '@tanstack/react-query';
import { apiClient } from '../api/client';
import dayjs from 'dayjs';

const { RangePicker } = DatePicker;

export default function Usage() {
  const { t } = useTranslation(['usage', 'common']);
  const [filters, setFilters] = useState({
    userId: '',
    model: '',
    provider: '',
    startDate: null as dayjs.Dayjs | null,
    endDate: null as dayjs.Dayjs | null,
  });

  const { data: usage = [], isLoading, refetch } = useQuery({
    queryKey: ['usage', filters.userId, filters.startDate?.toISOString(), filters.endDate?.toISOString()],
    queryFn: () => apiClient.getUsage(
      filters.userId || undefined,
      filters.startDate?.toISOString() || undefined,
      filters.endDate?.toISOString() || undefined
    ),
  });

  const handleFilter = () => {
    refetch();
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

      <Table
        dataSource={usage}
        columns={columns}
        rowKey="id"
        loading={isLoading}
        pagination={{ pageSize: 20 }}
        locale={{ emptyText: <Empty description="No usage records found" /> }}
      />
    </div>
  );
}
