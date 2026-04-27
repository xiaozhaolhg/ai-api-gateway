import { useState } from 'react';
import { Table, Button, Modal, Form, Input, Select, InputNumber, Popconfirm, Tag, Tabs, Spin, message } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, CheckOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient, type AlertRule, type Alert } from '../api/client';

export const Alerts: React.FC = () => {
  const { t } = useTranslation(['alerts', 'common']);
  const queryClient = useQueryClient();
  const [modalVisible, setModalVisible] = useState(false);
  const [editingRule, setEditingRule] = useState<AlertRule | null>(null);
  const [form] = Form.useForm();

  const { data: alertRules = [], isLoading: rulesLoading } = useQuery({
    queryKey: ['alertRules'],
    queryFn: () => apiClient.getAlertRules(),
  });

  const { data: alerts = [], isLoading: alertsLoading } = useQuery({
    queryKey: ['alerts'],
    queryFn: () => apiClient.getAlerts(),
  });

  const createMutation = useMutation({
    mutationFn: (data: Omit<AlertRule, 'id' | 'created_at' | 'updated_at'>) => apiClient.createAlertRule(data),
    onMutate: async (newRule) => {
      await queryClient.cancelQueries({ queryKey: ['alertRules'] });
      const previous = queryClient.getQueryData<AlertRule[]>(['alertRules']);
      queryClient.setQueryData<AlertRule[]>(['alertRules'], (old = []) => [
        ...old,
        { ...newRule, id: `temp-${Date.now()}`, created_at: new Date().toISOString(), updated_at: new Date().toISOString() } as AlertRule,
      ]);
      return { previous };
    },
    onError: (_err, _new, context) => {
      if (context?.previous) queryClient.setQueryData(['alertRules'], context.previous);
      message.error('Failed to create alert rule');
    },
    onSuccess: () => {
      message.success('Alert rule created successfully');
      queryClient.invalidateQueries({ queryKey: ['alertRules'] });
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<AlertRule> }) => apiClient.updateAlertRule(id, data),
    onMutate: async ({ id, data }) => {
      await queryClient.cancelQueries({ queryKey: ['alertRules'] });
      const previous = queryClient.getQueryData<AlertRule[]>(['alertRules']);
      queryClient.setQueryData<AlertRule[]>(['alertRules'], (old = []) =>
        old.map(r => r.id === id ? { ...r, ...data } : r)
      );
      return { previous };
    },
    onError: (_err, _vars, context) => {
      if (context?.previous) queryClient.setQueryData(['alertRules'], context.previous);
      message.error('Failed to update alert rule');
    },
    onSuccess: () => {
      message.success('Alert rule updated successfully');
      queryClient.invalidateQueries({ queryKey: ['alertRules'] });
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => apiClient.deleteAlertRule(id),
    onMutate: async (id) => {
      await queryClient.cancelQueries({ queryKey: ['alertRules'] });
      const previous = queryClient.getQueryData<AlertRule[]>(['alertRules']);
      queryClient.setQueryData<AlertRule[]>(['alertRules'], (old = []) => old.filter(r => r.id !== id));
      return { previous };
    },
    onError: (_err, _id, context) => {
      if (context?.previous) queryClient.setQueryData(['alertRules'], context.previous);
      message.error('Failed to delete alert rule');
    },
    onSuccess: () => {
      message.success('Alert rule deleted successfully');
      queryClient.invalidateQueries({ queryKey: ['alertRules'] });
    },
  });

  const acknowledgeMutation = useMutation({
    mutationFn: (id: string) => apiClient.acknowledgeAlert(id),
    onMutate: async (id) => {
      await queryClient.cancelQueries({ queryKey: ['alerts'] });
      const previous = queryClient.getQueryData<Alert[]>(['alerts']);
      queryClient.setQueryData<Alert[]>(['alerts'], (old = []) =>
        old.map(a => a.id === id ? { ...a, status: 'acknowledged' } : a)
      );
      return { previous };
    },
    onError: (_err, _id, context) => {
      if (context?.previous) queryClient.setQueryData(['alerts'], context.previous);
      message.error('Failed to acknowledge alert');
    },
    onSuccess: () => {
      message.success('Alert acknowledged successfully');
      queryClient.invalidateQueries({ queryKey: ['alerts'] });
    },
  });

  const handleAdd = () => {
    setEditingRule(null);
    form.resetFields();
    setModalVisible(true);
  };

  const handleEdit = (rule: AlertRule) => {
    setEditingRule(rule);
    form.setFieldsValue(rule);
    setModalVisible(true);
  };

  const handleDelete = async (id: string) => {
    deleteMutation.mutate(id);
  };

  const handleAcknowledge = async (id: string) => {
    acknowledgeMutation.mutate(id);
  };

  const handleModalOk = async () => {
    try {
      const values = await form.validateFields();

      if (editingRule) {
        updateMutation.mutate({ id: editingRule.id, data: values });
      } else {
        createMutation.mutate(values);
      }

      setModalVisible(false);
      form.resetFields();
    } catch (error) {
      // Validation error - do nothing
    }
  };

  const handleModalCancel = () => {
    setModalVisible(false);
    form.resetFields();
  };

  const ruleColumns = [
    {
      title: t('alerts:fields.name'),
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: t('alerts:fields.metric'),
      dataIndex: 'metric',
      key: 'metric',
    },
    {
      title: t('alerts:fields.condition'),
      dataIndex: 'condition',
      key: 'condition',
    },
    {
      title: t('alerts:fields.threshold'),
      dataIndex: 'threshold',
      key: 'threshold',
    },
    {
      title: t('alerts:fields.channel'),
      dataIndex: 'channel',
      key: 'channel',
    },
    {
      title: t('alerts:fields.status'),
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={status === 'active' ? 'green' : 'default'}>
          {status}
        </Tag>
      ),
    },
    {
      title: t('common:actions'),
      key: 'actions',
      render: (_: unknown, rule: AlertRule) => (
        <div style={{ display: 'flex', gap: 8 }}>
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => handleEdit(rule)}
          />
          <Popconfirm
            title={t('common:confirmDelete')}
            onConfirm={() => handleDelete(rule.id)}
            okText={t('common:ok')}
            cancelText={t('common:cancel')}
          >
            <Button type="link" danger icon={<DeleteOutlined />} />
          </Popconfirm>
        </div>
      ),
    },
  ];

  const alertColumns = [
    {
      title: t('alerts:fields.severity'),
      dataIndex: 'severity',
      key: 'severity',
      render: (severity: string) => (
        <Tag color={severity === 'critical' ? 'red' : severity === 'warning' ? 'orange' : 'blue'}>
          {severity}
        </Tag>
      ),
    },
    {
      title: t('alerts:fields.description'),
      dataIndex: 'description',
      key: 'description',
    },
    {
      title: t('alerts:fields.triggered_at'),
      dataIndex: 'triggered_at',
      key: 'triggered_at',
      render: (date: string) => new Date(date).toLocaleString(),
    },
    {
      title: t('alerts:fields.status'),
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={status === 'acknowledged' ? 'green' : 'red'}>
          {status}
        </Tag>
      ),
    },
    {
      title: t('common:actions'),
      key: 'actions',
      render: (_: unknown, alert: Alert) => (
        <div style={{ display: 'flex', gap: 8 }}>
          {alert.status === 'triggered' && (
            <Button
              type="link"
              icon={<CheckOutlined />}
              onClick={() => handleAcknowledge(alert.id)}
            >
              {t('alerts:actions.acknowledge')}
            </Button>
          )}
        </div>
      ),
    },
  ];

  const tabItems = [
    {
      key: 'rules',
      label: t('alerts:tabs.rules'),
      children: (
        <Table
          columns={ruleColumns}
          dataSource={alertRules}
          rowKey="id"
          loading={rulesLoading}
          pagination={{ pageSize: 10 }}
        />
      ),
    },
    {
      key: 'alerts',
      label: t('alerts:tabs.alerts'),
      children: (
        <Table
          columns={alertColumns}
          dataSource={alerts}
          rowKey="id"
          loading={alertsLoading}
          pagination={{ pageSize: 10 }}
        />
      ),
    },
  ];

  if (rulesLoading && alertsLoading) {
    return <Spin size="large" style={{ display: 'block', margin: '48px auto' }} />;
  }

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h1>{t('alerts:title')}</h1>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
          {t('alerts:actions.addRule')}
        </Button>
      </div>

      <Tabs items={tabItems} />

      <Modal
        title={editingRule ? t('alerts:actions.editRule') : t('alerts:actions.addRule')}
        open={modalVisible}
        onOk={handleModalOk}
        onCancel={handleModalCancel}
        okText={t('common:ok')}
        cancelText={t('common:cancel')}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="name"
            label={t('alerts:fields.name')}
            rules={[{ required: true, message: t('alerts:validation.nameRequired') }]}
          >
            <Input placeholder={t('alerts:placeholders.name')} />
          </Form.Item>
          <Form.Item
            name="metric"
            label={t('alerts:fields.metric')}
            rules={[{ required: true, message: t('alerts:validation.metricRequired') }]}
          >
            <Select placeholder={t('alerts:placeholders.metric')}>
              <Select.Option value="error_rate">{t('alerts:metrics.errorRate')}</Select.Option>
              <Select.Option value="latency">{t('alerts:metrics.latency')}</Select.Option>
              <Select.Option value="cost">{t('alerts:metrics.cost')}</Select.Option>
              <Select.Option value="usage">{t('alerts:metrics.usage')}</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item
            name="condition"
            label={t('alerts:fields.condition')}
            rules={[{ required: true, message: t('alerts:validation.conditionRequired') }]}
          >
            <Select placeholder={t('alerts:placeholders.condition')}>
              <Select.Option value="gt">{t('alerts:conditions.greaterThan')}</Select.Option>
              <Select.Option value="lt">{t('alerts:conditions.lessThan')}</Select.Option>
              <Select.Option value="eq">{t('alerts:conditions.equals')}</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item
            name="threshold"
            label={t('alerts:fields.threshold')}
            rules={[{ required: true, message: t('alerts:validation.thresholdRequired') }]}
          >
            <InputNumber style={{ width: '100%' }} placeholder={t('alerts:placeholders.threshold')} />
          </Form.Item>
          <Form.Item
            name="channel"
            label={t('alerts:fields.channel')}
            rules={[{ required: true, message: t('alerts:validation.channelRequired') }]}
          >
            <Select placeholder={t('alerts:placeholders.channel')}>
              <Select.Option value="email">{t('alerts:channels.email')}</Select.Option>
              <Select.Option value="slack">{t('alerts:channels.slack')}</Select.Option>
              <Select.Option value="webhook">{t('alerts:channels.webhook')}</Select.Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};
