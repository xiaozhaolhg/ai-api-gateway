import React, { useState } from 'react';
import { Table, Button, Modal, Form, Input, Select, InputNumber, Popconfirm, Tag, Empty, message } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient, type Budget } from '../api/client';

export const Budgets: React.FC = () => {
  const { t } = useTranslation(['budgets', 'common']);
  const queryClient = useQueryClient();
  const [modalVisible, setModalVisible] = useState(false);
  const [editingBudget, setEditingBudget] = useState<Budget | null>(null);
  const [form] = Form.useForm();

  const { data: budgets = [], isLoading } = useQuery({
    queryKey: ['budgets'],
    queryFn: () => apiClient.getBudgets(),
  });

  const createMutation = useMutation({
    mutationFn: (data: Omit<Budget, 'id' | 'created_at' | 'updated_at' | 'current_spend'>) => apiClient.createBudget(data),
    onMutate: async (newBudget) => {
      await queryClient.cancelQueries({ queryKey: ['budgets'] });
      const previous = queryClient.getQueryData<Budget[]>(['budgets']);
      queryClient.setQueryData<Budget[]>(['budgets'], (old = []) => [
        ...old, { ...newBudget, id: `temp-${Date.now()}`, current_spend: 0, created_at: new Date().toISOString(), updated_at: new Date().toISOString() } as Budget,
      ]);
      return { previous };
    },
    onError: (_err, _new, context) => {
      if (context?.previous) queryClient.setQueryData(['budgets'], context.previous);
      message.error('Failed to create budget');
    },
    onSuccess: () => {
      message.success('Budget created successfully');
      queryClient.invalidateQueries({ queryKey: ['budgets'] });
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<Budget> }) => apiClient.updateBudget(id, data),
    onMutate: async ({ id, data }) => {
      await queryClient.cancelQueries({ queryKey: ['budgets'] });
      const previous = queryClient.getQueryData<Budget[]>(['budgets']);
      queryClient.setQueryData<Budget[]>(['budgets'], (old = []) => old.map(b => b.id === id ? { ...b, ...data } : b));
      return { previous };
    },
    onError: (_err, _vars, context) => {
      if (context?.previous) queryClient.setQueryData(['budgets'], context.previous);
      message.error('Failed to update budget');
    },
    onSuccess: () => {
      message.success('Budget updated successfully');
      queryClient.invalidateQueries({ queryKey: ['budgets'] });
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => apiClient.deleteBudget(id),
    onMutate: async (id) => {
      await queryClient.cancelQueries({ queryKey: ['budgets'] });
      const previous = queryClient.getQueryData<Budget[]>(['budgets']);
      queryClient.setQueryData<Budget[]>(['budgets'], (old = []) => old.filter(b => b.id !== id));
      return { previous };
    },
    onError: (_err, _id, context) => {
      if (context?.previous) queryClient.setQueryData(['budgets'], context.previous);
      message.error('Failed to delete budget');
    },
    onSuccess: () => {
      message.success('Budget deleted successfully');
    },
  });

  const handleAdd = () => {
    setEditingBudget(null);
    form.resetFields();
    setModalVisible(true);
  };

  const handleEdit = (budget: Budget) => {
    setEditingBudget(budget);
    form.setFieldsValue(budget);
    setModalVisible(true);
  };

  const handleModalOk = async () => {
    const values = await form.validateFields();
    if (editingBudget) {
      updateMutation.mutate({ id: editingBudget.id, data: values });
    } else {
      createMutation.mutate(values);
    }
    setModalVisible(false);
  };

  const handleModalCancel = () => {
    setModalVisible(false);
    form.resetFields();
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active': return 'green';
      case 'warning': return 'orange';
      case 'exceeded': return 'red';
      default: return 'default';
    }
  };

  const columns = [
    {
      title: t('budgets:fields.name'),
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: t('budgets:fields.scope'),
      dataIndex: 'scope',
      key: 'scope',
    },
    {
      title: t('budgets:fields.limit'),
      dataIndex: 'limit',
      key: 'limit',
      render: (limit: number) => `$${limit.toFixed(2)}`,
    },
    {
      title: t('budgets:fields.currentSpend'),
      dataIndex: 'current_spend',
      key: 'current_spend',
      render: (spend: number) => `$${spend.toFixed(2)}`,
    },
    {
      title: t('budgets:fields.remaining'),
      key: 'remaining',
      render: (_: any, record: Budget) => {
        const remaining = record.limit - record.current_spend;
        return `$${remaining.toFixed(2)}`;
      },
    },
    {
      title: t('budgets:fields.period'),
      dataIndex: 'period',
      key: 'period',
    },
    {
      title: t('budgets:fields.status'),
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={getStatusColor(status)}>{status}</Tag>
      ),
    },
    {
      title: t('common:common.actions'),
      key: 'actions',
      render: (_: any, record: Budget) => (
        <div>
          <Button type="link" icon={<EditOutlined />} onClick={() => handleEdit(record)}>
            {t('common:edit')}
          </Button>
          <Popconfirm
            title="Are you sure you want to delete this budget?"
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
        <h2>{t('budgets:title')}</h2>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
          {t('budgets:addBudget')}
        </Button>
      </div>

      <Table
        dataSource={budgets}
        columns={columns}
        rowKey="id"
        pagination={{ pageSize: 10 }}
        locale={{ emptyText: <Empty description="No budgets found" /> }}
      />

      <Modal
        title={editingBudget ? t('budgets:editBudget') : t('budgets:addBudget')}
        open={modalVisible}
        onOk={handleModalOk}
        onCancel={handleModalCancel}
        confirmLoading={createMutation.isPending || updateMutation.isPending}
        width={600}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            label={t('budgets:fields.name')}
            name="name"
            rules={[{ required: true, message: 'Please input budget name' }]}
          >
            <Input />
          </Form.Item>

          <Form.Item
            label={t('budgets:fields.scope')}
            name="scope"
            rules={[{ required: true, message: 'Please select scope' }]}
          >
            <Select>
              <Select.Option value="global">Global</Select.Option>
              <Select.Option value="user">User</Select.Option>
              <Select.Option value="group">Group</Select.Option>
            </Select>
          </Form.Item>

          <Form.Item
            label={t('budgets:fields.limit')}
            name="limit"
            rules={[{ required: true, message: 'Please input limit' }]}
          >
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>

          <Form.Item
            label={t('budgets:fields.period')}
            name="period"
            rules={[{ required: true, message: 'Please select period' }]}
          >
            <Select>
              <Select.Option value="monthly">Monthly</Select.Option>
              <Select.Option value="quarterly">Quarterly</Select.Option>
              <Select.Option value="yearly">Yearly</Select.Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};
