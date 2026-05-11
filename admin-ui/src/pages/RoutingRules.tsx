import React, { useState } from 'react';
import { Table, Button, Modal, Form, Input, Select, Popconfirm, Tag, Empty, message } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient, type RoutingRule, type Provider } from '../api/client';

export const RoutingRules: React.FC = () => {
  const { t } = useTranslation(['routing', 'common']);
  const queryClient = useQueryClient();
  const [modalVisible, setModalVisible] = useState(false);
  const [editingRule, setEditingRule] = useState<RoutingRule | null>(null);
  const [form] = Form.useForm();

  const { data: rules = [], isLoading } = useQuery({
    queryKey: ['routingRules'],
    queryFn: () => apiClient.getRoutingRules(),
  });

  const { data: providers = [] } = useQuery({
    queryKey: ['providers'],
    queryFn: () => apiClient.getProviders(),
  });

  const createMutation = useMutation({
    mutationFn: (data: Omit<RoutingRule, 'id' | 'created_at' | 'updated_at'>) => apiClient.createRoutingRule(data),
    onMutate: async (newRule) => {
      await queryClient.cancelQueries({ queryKey: ['routingRules'] });
      const previous = queryClient.getQueryData<RoutingRule[]>(['routingRules']);
      queryClient.setQueryData<RoutingRule[]>(['routingRules'], (old = []) => [
        ...old, { ...newRule, id: `temp-${Date.now()}`, created_at: new Date().toISOString(), updated_at: new Date().toISOString() } as RoutingRule,
      ]);
      return { previous };
    },
    onError: (_err, _new, context) => {
      if (context?.previous) queryClient.setQueryData(['routingRules'], context.previous);
      message.error('Failed to create routing rule');
    },
    onSuccess: () => {
      message.success('Routing rule created successfully');
      queryClient.invalidateQueries({ queryKey: ['routingRules'] });
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<RoutingRule> }) => apiClient.updateRoutingRule(id, data),
    onMutate: async ({ id, data }) => {
      await queryClient.cancelQueries({ queryKey: ['routingRules'] });
      const previous = queryClient.getQueryData<RoutingRule[]>(['routingRules']);
      queryClient.setQueryData<RoutingRule[]>(['routingRules'], (old = []) => old.map(r => r.id === id ? { ...r, ...data } : r));
      return { previous };
    },
    onError: (_err, _vars, context) => {
      if (context?.previous) queryClient.setQueryData(['routingRules'], context.previous);
      message.error('Failed to update routing rule');
    },
    onSuccess: () => {
      message.success('Routing rule updated successfully');
      queryClient.invalidateQueries({ queryKey: ['routingRules'] });
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => apiClient.deleteRoutingRule(id),
    onMutate: async (id) => {
      await queryClient.cancelQueries({ queryKey: ['routingRules'] });
      const previous = queryClient.getQueryData<RoutingRule[]>(['routingRules']);
      queryClient.setQueryData<RoutingRule[]>(['routingRules'], (old = []) => old.filter(r => r.id !== id));
      return { previous };
    },
    onError: (_err, _id, context) => {
      if (context?.previous) queryClient.setQueryData(['routingRules'], context.previous);
      message.error('Failed to delete routing rule');
    },
    onSuccess: () => {
      message.success('Routing rule deleted successfully');
    },
  });

  const handleAdd = () => {
    setEditingRule(null);
    form.resetFields();
    setModalVisible(true);
  };

  const handleEdit = (rule: RoutingRule) => {
    setEditingRule(rule);
    form.setFieldsValue({
      ...rule,
      fallback_provider_ids: rule.fallback_provider_ids?.join(', ') || '',
    });
    setModalVisible(true);
  };

  const handleModalOk = async () => {
    const values = await form.validateFields();
    const ruleData = {
      model_pattern: values.model_pattern,
      provider_id: values.provider_id,
      priority: Number(values.priority) || 1,
      fallback_provider_ids: values.fallback_provider_ids ? values.fallback_provider_ids.split(',').map((m: string) => m.trim()).filter(Boolean) : [],
      fallback_models: [],
      is_system_default: values.is_system_default ?? true,
    };

    if (editingRule) {
      updateMutation.mutate({ id: editingRule.id, data: ruleData });
    } else {
      createMutation.mutate(ruleData);
    }
    setModalVisible(false);
  };

  const handleModalCancel = () => {
    setModalVisible(false);
    form.resetFields();
  };

  const columns = [
    {
      title: t('routing:fields.modelPattern'),
      dataIndex: 'model_pattern',
      key: 'model_pattern',
    },
    {
      title: t('routing:fields.provider'),
      dataIndex: 'provider_id',
      key: 'provider_id',
    },
    {
      title: t('routing:fields.priority'),
      dataIndex: 'priority',
      key: 'priority',
    },
    {
      title: t('routing:fields.fallbackChain'),
      dataIndex: 'fallback_provider_ids',
      key: 'fallback_provider_ids',
      render: (chain: string[]) => chain?.join(', ') || '-',
    },
    {
      title: t('routing:fields.status'),
      dataIndex: 'is_system_default',
      key: 'is_system_default',
      render: (isDefault: boolean) => (
        <Tag color={isDefault ? 'green' : 'red'}>{isDefault ? 'active' : 'inactive'}</Tag>
      ),
    },
    {
      title: t('common:actions'),
      key: 'actions',
      render: (_: any, record: RoutingRule) => (
        <div>
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
          >
            {t('common:edit')}
          </Button>
          <Popconfirm
            title="Are you sure you want to delete this routing rule?"
            onConfirm={() => deleteMutation.mutate(record.id)}
            okText="Yes"
            cancelText="No"
          >
            <Button type="link" danger icon={<DeleteOutlined />}>
              {t('common:delete')}
            </Button>
          </Popconfirm>
        </div>
      ),
    },
  ];

  if (isLoading) {
    return <Table loading={true} columns={columns} dataSource={[]} rowKey="id" />;
  }

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <h2>{t('routing:title')}</h2>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
          {t('routing:addRule')}
        </Button>
      </div>

      <Table
        dataSource={rules}
        columns={columns}
        rowKey="id"
        pagination={{ pageSize: 10 }}
        locale={{ emptyText: <Empty description="No routing rules found" /> }}
      />

      <Modal
        title={editingRule ? t('routing:editRule') : t('routing:addRule')}
        open={modalVisible}
        onOk={handleModalOk}
        onCancel={handleModalCancel}
        confirmLoading={createMutation.isPending || updateMutation.isPending}
        width={600}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            label={t('routing:fields.modelPattern')}
            name="model_pattern"
            rules={[{ required: true, message: 'Please input model pattern' }]}
          >
            <Input placeholder="llama2, gpt-4, *" />
          </Form.Item>

          <Form.Item
            label={t('routing:fields.provider')}
            name="provider_id"
            rules={[{ required: true, message: 'Please select provider' }]}
          >
            <Select>
              {providers
                .filter((p: Provider) => p.status === 'active')
                .map((p: Provider) => (
                  <Select.Option key={p.id} value={p.id}>
                    {p.name} ({p.type})
                  </Select.Option>
                ))}
            </Select>
          </Form.Item>

          <Form.Item
            label={t('routing:fields.priority')}
            name="priority"
            initialValue={1}
            rules={[{ required: true, message: 'Please input priority' }]}
          >
            <Input type="number" placeholder="1" />
          </Form.Item>

          <Form.Item
            label={t('routing:fields.fallbackChain')}
            name="fallback_provider_ids"
            rules={[{ required: false, message: 'Please input fallback chain' }]}
          >
            <Input placeholder="ollama_1_211, another_provider" />
          </Form.Item>

          <Form.Item
            label={t('routing:fields.status')}
            name="is_system_default"
            initialValue={true}
          >
            <Select>
              <Select.Option value={true}>Active</Select.Option>
              <Select.Option value={false}>Inactive</Select.Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};
