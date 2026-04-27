import React, { useState } from 'react';
import { Table, Button, Modal, Form, Input, Select, Popconfirm, Tag, Empty, message } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient, type RoutingRule } from '../api/client';

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
      fallback_chain: rule.fallback_chain.join(', '),
    });
    setModalVisible(true);
  };

  const handleModalOk = async () => {
    const values = await form.validateFields();
    const ruleData = {
      ...values,
      fallback_chain: values.fallback_chain.split(',').map((m: string) => m.trim()),
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
      dataIndex: 'provider',
      key: 'provider',
    },
    {
      title: t('routing:fields.adapterType'),
      dataIndex: 'adapter_type',
      key: 'adapter_type',
    },
    {
      title: t('routing:fields.priority'),
      dataIndex: 'priority',
      key: 'priority',
    },
    {
      title: t('routing:fields.fallbackChain'),
      dataIndex: 'fallback_chain',
      key: 'fallback_chain',
      render: (chain: string[]) => chain.join(', '),
    },
    {
      title: t('routing:fields.status'),
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={status === 'active' ? 'green' : 'red'}>{status}</Tag>
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
            <Input placeholder="gpt-*" />
          </Form.Item>

          <Form.Item
            label={t('routing:fields.provider')}
            name="provider"
            rules={[{ required: true, message: 'Please select provider' }]}
          >
            <Select>
              <Select.Option value="openai">OpenAI</Select.Option>
              <Select.Option value="anthropic">Anthropic</Select.Option>
              <Select.Option value="gemini">Gemini</Select.Option>
              <Select.Option value="ollama">Ollama</Select.Option>
            </Select>
          </Form.Item>

          <Form.Item
            label={t('routing:fields.adapterType')}
            name="adapter_type"
            rules={[{ required: true, message: 'Please select adapter type' }]}
          >
            <Select>
              <Select.Option value="openai">OpenAI</Select.Option>
              <Select.Option value="anthropic">Anthropic</Select.Option>
              <Select.Option value="gemini">Gemini</Select.Option>
              <Select.Option value="ollama">Ollama</Select.Option>
            </Select>
          </Form.Item>

          <Form.Item
            label={t('routing:fields.priority')}
            name="priority"
            rules={[{ required: true, message: 'Please input priority' }]}
          >
            <Input type="number" />
          </Form.Item>

          <Form.Item
            label={t('routing:fields.fallbackChain')}
            name="fallback_chain"
            rules={[{ required: true, message: 'Please input fallback chain' }]}
          >
            <Input placeholder="provider1, provider2" />
          </Form.Item>

          <Form.Item
            label={t('routing:fields.status')}
            name="status"
            initialValue="active"
          >
            <Select>
              <Select.Option value="active">Active</Select.Option>
              <Select.Option value="inactive">Inactive</Select.Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};
