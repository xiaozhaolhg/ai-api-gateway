import React, { useState } from 'react';
import { Table, Button, Modal, Form, Input, Select, InputNumber, Popconfirm, Empty, message } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient, type PricingRule } from '../api/client';

export const PricingRules: React.FC = () => {
  const { t } = useTranslation(['pricing', 'common']);
  const queryClient = useQueryClient();
  const [modalVisible, setModalVisible] = useState(false);
  const [editingRule, setEditingRule] = useState<PricingRule | null>(null);
  const [form] = Form.useForm();

  const { data: rules = [], isLoading } = useQuery({
    queryKey: ['pricingRules'],
    queryFn: () => apiClient.getPricingRules(),
  });

  const createMutation = useMutation({
    mutationFn: (data: Omit<PricingRule, 'id' | 'created_at' | 'updated_at'>) => apiClient.createPricingRule(data),
    onMutate: async (newRule) => {
      await queryClient.cancelQueries({ queryKey: ['pricingRules'] });
      const previous = queryClient.getQueryData<PricingRule[]>(['pricingRules']);
      queryClient.setQueryData<PricingRule[]>(['pricingRules'], (old = []) => [
        ...old, { ...newRule, id: `temp-${Date.now()}`, created_at: new Date().toISOString(), updated_at: new Date().toISOString() } as PricingRule,
      ]);
      return { previous };
    },
    onError: (_err, _new, context) => {
      if (context?.previous) queryClient.setQueryData(['pricingRules'], context.previous);
      message.error('Failed to create pricing rule');
    },
    onSuccess: () => {
      message.success('Pricing rule created successfully');
      queryClient.invalidateQueries({ queryKey: ['pricingRules'] });
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<PricingRule> }) => apiClient.updatePricingRule(id, data),
    onMutate: async ({ id, data }) => {
      await queryClient.cancelQueries({ queryKey: ['pricingRules'] });
      const previous = queryClient.getQueryData<PricingRule[]>(['pricingRules']);
      queryClient.setQueryData<PricingRule[]>(['pricingRules'], (old = []) => old.map(r => r.id === id ? { ...r, ...data } : r));
      return { previous };
    },
    onError: (_err, _vars, context) => {
      if (context?.previous) queryClient.setQueryData(['pricingRules'], context.previous);
      message.error('Failed to update pricing rule');
    },
    onSuccess: () => {
      message.success('Pricing rule updated successfully');
      queryClient.invalidateQueries({ queryKey: ['pricingRules'] });
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => apiClient.deletePricingRule(id),
    onMutate: async (id) => {
      await queryClient.cancelQueries({ queryKey: ['pricingRules'] });
      const previous = queryClient.getQueryData<PricingRule[]>(['pricingRules']);
      queryClient.setQueryData<PricingRule[]>(['pricingRules'], (old = []) => old.filter(r => r.id !== id));
      return { previous };
    },
    onError: (_err, _id, context) => {
      if (context?.previous) queryClient.setQueryData(['pricingRules'], context.previous);
      message.error('Failed to delete pricing rule');
    },
    onSuccess: () => {
      message.success('Pricing rule deleted successfully');
    },
  });

  const handleAdd = () => {
    setEditingRule(null);
    form.resetFields();
    setModalVisible(true);
  };

  const handleEdit = (rule: PricingRule) => {
    setEditingRule(rule);
    form.setFieldsValue(rule);
    setModalVisible(true);
  };

  const handleModalOk = async () => {
    const values = await form.validateFields();
    if (editingRule) {
      updateMutation.mutate({ id: editingRule.id, data: values });
    } else {
      createMutation.mutate(values);
    }
    setModalVisible(false);
  };

  const handleModalCancel = () => {
    setModalVisible(false);
    form.resetFields();
  };

  const columns = [
    {
      title: t('pricing:fields.model'),
      dataIndex: 'model',
      key: 'model',
    },
    {
      title: t('pricing:fields.provider'),
      dataIndex: 'provider',
      key: 'provider',
    },
    {
      title: t('pricing:fields.promptPrice'),
      dataIndex: 'prompt_price',
      key: 'prompt_price',
      render: (price: number) => `$${price.toFixed(4)}`,
    },
    {
      title: t('pricing:fields.completionPrice'),
      dataIndex: 'completion_price',
      key: 'completion_price',
      render: (price: number) => `$${price.toFixed(4)}`,
    },
    {
      title: t('pricing:fields.currency'),
      dataIndex: 'currency',
      key: 'currency',
    },
    {
      title: t('pricing:fields.effectiveDate'),
      dataIndex: 'effective_date',
      key: 'effective_date',
    },
    {
      title: t('common:common.actions'),
      key: 'actions',
      render: (_: any, record: PricingRule) => (
        <div>
          <Button type="link" icon={<EditOutlined />} onClick={() => handleEdit(record)}>
            {t('common:edit')}
          </Button>
          <Popconfirm
            title="Are you sure you want to delete this pricing rule?"
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
        <h2>{t('pricing:title')}</h2>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
          {t('pricing:addRule')}
        </Button>
      </div>

      <Table
        dataSource={rules}
        columns={columns}
        rowKey="id"
        pagination={{ pageSize: 10 }}
        locale={{ emptyText: <Empty description="No pricing rules found" /> }}
      />

      <Modal
        title={editingRule ? t('pricing:editRule') : t('pricing:addRule')}
        open={modalVisible}
        onOk={handleModalOk}
        onCancel={handleModalCancel}
        confirmLoading={createMutation.isPending || updateMutation.isPending}
        width={600}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            label={t('pricing:fields.model')}
            name="model"
            rules={[{ required: true, message: 'Please input model' }]}
          >
            <Input placeholder="gpt-4" />
          </Form.Item>

          <Form.Item
            label={t('pricing:fields.provider')}
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
            label={t('pricing:fields.promptPrice')}
            name="prompt_price"
            rules={[{ required: true, message: 'Please input prompt price' }]}
          >
            <InputNumber min={0} step={0.0001} style={{ width: '100%' }} />
          </Form.Item>

          <Form.Item
            label={t('pricing:fields.completionPrice')}
            name="completion_price"
            rules={[{ required: true, message: 'Please input completion price' }]}
          >
            <InputNumber min={0} step={0.0001} style={{ width: '100%' }} />
          </Form.Item>

          <Form.Item
            label={t('pricing:fields.currency')}
            name="currency"
            initialValue="USD"
          >
            <Select>
              <Select.Option value="USD">USD</Select.Option>
              <Select.Option value="EUR">EUR</Select.Option>
              <Select.Option value="CNY">CNY</Select.Option>
            </Select>
          </Form.Item>

          <Form.Item
            label={t('pricing:fields.effectiveDate')}
            name="effective_date"
            rules={[{ required: true, message: 'Please select effective date' }]}
          >
            <Input type="date" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};
