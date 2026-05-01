import { useState } from 'react';
import { Table, Button, Modal, Form, Input, Select, Popconfirm, Tag, Empty, message } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient, type Provider } from '../api/client';

export default function Providers() {
  const { t } = useTranslation(['providers', 'common']);
  const queryClient = useQueryClient();
  const [modalVisible, setModalVisible] = useState(false);
  const [editingProvider, setEditingProvider] = useState<Provider | null>(null);
  const [form] = Form.useForm();

  const { data: providers = [], isLoading } = useQuery({
    queryKey: ['providers'],
    queryFn: () => apiClient.getProviders(),
  });

  const createMutation = useMutation({
    mutationFn: (data: Omit<Provider, 'id' | 'created_at' | 'updated_at'>) => apiClient.createProvider(data),
    onMutate: async (newProvider) => {
      await queryClient.cancelQueries({ queryKey: ['providers'] });
      const previous = queryClient.getQueryData<Provider[]>(['providers']);
      queryClient.setQueryData<Provider[]>(['providers'], (old = []) => [
        ...old,
        { ...newProvider, id: `temp-${Date.now()}`, created_at: new Date().toISOString(), updated_at: new Date().toISOString() } as Provider,
      ]);
      return { previous };
    },
    onError: (_err, _new, context) => {
      if (context?.previous) queryClient.setQueryData(['providers'], context.previous);
      message.error('Failed to create provider');
    },
    onSuccess: () => {
      message.success('Provider created successfully');
      queryClient.invalidateQueries({ queryKey: ['providers'] });
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<Provider> }) => apiClient.updateProvider(id, data),
    onMutate: async ({ id, data }) => {
      await queryClient.cancelQueries({ queryKey: ['providers'] });
      const previous = queryClient.getQueryData<Provider[]>(['providers']);
      queryClient.setQueryData<Provider[]>(['providers'], (old = []) =>
        old.map(p => p.id === id ? { ...p, ...data } : p)
      );
      return { previous };
    },
    onError: (_err, _vars, context) => {
      if (context?.previous) queryClient.setQueryData(['providers'], context.previous);
      message.error('Failed to update provider');
    },
    onSuccess: () => {
      message.success('Provider updated successfully');
      queryClient.invalidateQueries({ queryKey: ['providers'] });
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => apiClient.deleteProvider(id),
    onMutate: async (id) => {
      await queryClient.cancelQueries({ queryKey: ['providers'] });
      const previous = queryClient.getQueryData<Provider[]>(['providers']);
      queryClient.setQueryData<Provider[]>(['providers'], (old = []) => old.filter(p => p.id !== id));
      return { previous };
    },
    onError: (_err, _id, context) => {
      if (context?.previous) queryClient.setQueryData(['providers'], context.previous);
      message.error('Failed to delete provider');
    },
    onSuccess: () => {
      message.success('Provider deleted successfully');
    },
  });

  const handleAdd = () => {
    setEditingProvider(null);
    form.resetFields();
    setModalVisible(true);
  };

  const handleEdit = (provider: Provider) => {
    setEditingProvider(provider);
    form.setFieldsValue({
      ...provider,
      models: provider.models?.join(', ') || '',
    });
    setModalVisible(true);
  };

  const handleModalOk = async () => {
    const values = await form.validateFields();
    const providerData = {
      ...values,
      models: values.models.split(',').map((m: string) => m.trim()),
    };

    if (editingProvider) {
      updateMutation.mutate({ id: editingProvider.id, data: providerData });
    } else {
      createMutation.mutate(providerData);
    }
    setModalVisible(false);
  };

  const handleModalCancel = () => {
    setModalVisible(false);
    form.resetFields();
  };

  const columns = [
    {
      title: t('providers:fields.name'),
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: t('providers:fields.type'),
      dataIndex: 'type',
      key: 'type',
    },
    {
      title: 'Base URL',
      dataIndex: 'base_url',
      key: 'base_url',
    },
    {
      title: 'Models',
      dataIndex: 'models',
      key: 'models',
      render: (models: string[]) => (models || []).join(', '),
    },
    {
      title: t('providers:fields.status'),
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={status === 'active' ? 'green' : 'red'}>{status}</Tag>
      ),
    },
    {
      title: t('common:actions'),
      key: 'actions',
      render: (_: any, record: Provider) => (
        <div>
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
          >
            {t('common:edit')}
          </Button>
          <Popconfirm
            title="Are you sure you want to delete this provider?"
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
        <h2>{t('providers:title')}</h2>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
          {t('providers:addProvider')}
        </Button>
      </div>

      {providers.length === 0 ? (
        <Empty description="No providers found" />
      ) : (
        <Table
          dataSource={providers}
          columns={columns}
          rowKey="id"
          pagination={{ pageSize: 10 }}
        />
      )}

      <Modal
        title={editingProvider ? t('providers:editProvider') : t('providers:addProvider')}
        open={modalVisible}
        onOk={handleModalOk}
        onCancel={handleModalCancel}
        confirmLoading={createMutation.isPending || updateMutation.isPending}
        width={600}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            label={t('providers:fields.name')}
            name="name"
            rules={[{ required: true, message: 'Please input provider name' }]}
          >
            <Input />
          </Form.Item>

          <Form.Item
            label={t('providers:fields.type')}
            name="type"
            rules={[{ required: true, message: 'Please select provider type' }]}
          >
            <Select>
              <Select.Option value="openai">OpenAI</Select.Option>
              <Select.Option value="anthropic">Anthropic</Select.Option>
              <Select.Option value="gemini">Gemini</Select.Option>
              <Select.Option value="ollama">Ollama</Select.Option>
              <Select.Option value="custom">Custom</Select.Option>
            </Select>
          </Form.Item>

          <Form.Item
            label="Base URL"
            name="base_url"
            rules={[{ required: true, message: 'Please input base URL' }]}
          >
            <Input />
          </Form.Item>

          <Form.Item
            label="Models (comma-separated)"
            name="models"
            rules={[{ required: true, message: 'Please input models' }]}
          >
            <Input placeholder="gpt-4, gpt-3.5-turbo" />
          </Form.Item>

          <Form.Item
            label={t('providers:fields.enabled')}
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
}
