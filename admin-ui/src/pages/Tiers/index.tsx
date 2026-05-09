import { useState } from 'react';
import { Table, Button, Modal, Form, Input, Select, Popconfirm, Tag, Space, message, Card, Descriptions, Collapse } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, EyeOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient, type Tier } from '../../api/client';

export const Tiers: React.FC = () => {
  const { t } = useTranslation(['tiers', 'common']);
  const queryClient = useQueryClient();
  const [modalVisible, setModalVisible] = useState(false);
  const [detailVisible, setDetailVisible] = useState(false);
  const [editingTier, setEditingTier] = useState<Tier | null>(null);
  const [selectedTier, setSelectedTier] = useState<Tier | null>(null);
  const [form] = Form.useForm();

  const { data: tiers = [], isLoading } = useQuery<Tier[]>({
    queryKey: ['tiers'],
    queryFn: () => apiClient.getTiers(),
  });

  const createMutation = useMutation({
    mutationFn: async (data: Partial<Tier>) => {
      return apiClient.createTier({
        name: data.name || '',
        description: data.description || '',
        is_default: false,
        allowed_models: data.allowed_models || [],
        allowed_providers: data.allowed_providers || [],
      });
    },
    onSuccess: () => {
      message.success('Tier created successfully');
      queryClient.invalidateQueries({ queryKey: ['tiers'] });
      setModalVisible(false);
    },
    onError: () => {
      message.error('Failed to create tier');
    },
  });

  const updateMutation = useMutation({
    mutationFn: async ({ id, data }: { id: string; data: Partial<Tier> }) => {
      return apiClient.updateTier(id, data);
    },
    onSuccess: () => {
      message.success('Tier updated successfully');
      queryClient.invalidateQueries({ queryKey: ['tiers'] });
      setModalVisible(false);
    },
    onError: () => {
      message.error('Failed to update tier');
    },
  });

  const deleteMutation = useMutation({
    mutationFn: async (id: string) => {
      await apiClient.deleteTier(id);
    },
    onSuccess: () => {
      message.success('Tier deleted successfully');
      queryClient.invalidateQueries({ queryKey: ['tiers'] });
    },
    onError: () => {
    },
  });

  const handleAdd = () => {
    setEditingTier(null);
    form.resetFields();
    setModalVisible(true);
  };

  const handleEdit = (tier: Tier) => {
    setEditingTier(tier);
    form.setFieldsValue({
      name: tier.name,
      description: tier.description,
      allowed_models: tier.allowed_models,
      allowed_providers: tier.allowed_providers,
    });
    setModalVisible(true);
  };

  const handleView = (tier: Tier) => {
    setSelectedTier(tier);
    setDetailVisible(true);
  };

  const handleModalOk = async () => {
    const values = await form.validateFields();

    if (editingTier) {
      updateMutation.mutate({
        id: editingTier.id,
        data: {
          name: values.name,
          description: values.description,
          allowed_models: values.allowed_models || [],
          allowed_providers: values.allowed_providers || [],
        },
      });
    } else {
      createMutation.mutate({
        name: values.name,
        description: values.description,
        allowed_models: values.allowed_models || [],
        allowed_providers: values.allowed_providers || [],
      });
    }
  };

  const handleDelete = (id: string) => {
    deleteMutation.mutate(id);
  };

  const columns = [
    {
      title: t('name'),
      dataIndex: 'name',
      key: 'name',
      render: (name: string, record: Tier) => (
        <Space>
          {name}
          {record.is_default && <Tag color="blue">Default</Tag>}
        </Space>
      ),
    },
    {
      title: t('description'),
      dataIndex: 'description',
      key: 'description',
    },
    {
      title: t('models'),
      key: 'models',
      render: (_: unknown, record: Tier) => (
        <span>{record.allowed_models?.length || 0} models</span>
      ),
    },
    {
      title: t('providers'),
      key: 'providers',
      render: (_: unknown, record: Tier) => (
        <span>{record.allowed_providers?.length || 0} providers</span>
      ),
    },
    {
      title: t('actions'),
      key: 'actions',
      render: (_: unknown, record: Tier) => (
        <Space>
          <Button
            icon={<EyeOutlined />}
            size="small"
            onClick={() => handleView(record)}
          >
            {t('view')}
          </Button>
          {!record.is_default && (
            <>
              <Button
                icon={<EditOutlined />}
                size="small"
                onClick={() => handleEdit(record)}
              />
              <Popconfirm
                title={t('confirmDelete')}
                onConfirm={() => handleDelete(record.id)}
              >
                <Button icon={<DeleteOutlined />} size="small" danger />
              </Popconfirm>
            </>
          )}
        </Space>
      ),
    },
  ];

  const mockProviders = [
    { id: 'ollama', name: 'Ollama', models: ['llama2', 'mistral', 'codellama', 'orca-mini'] },
    { id: 'openai', name: 'OpenAI', models: ['gpt-4', 'gpt-4-turbo', 'gpt-3.5-turbo'] },
    { id: 'anthropic', name: 'Anthropic', models: ['claude-3', 'claude-3-sonnet', 'claude-3-opus'] },
    { id: 'gemini', name: 'Google Gemini', models: ['gemini-pro', 'gemini-pro-vision'] },
  ];

  const groupedModels = mockProviders.reduce((acc, provider) => {
    acc[provider.id] = {
      name: provider.name,
      models: provider.models.map(m => `${provider.id}:${m}`),
    };
    return acc;
  }, {} as Record<string, { name: string; models: string[] }>);

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <h2>{t('title')}</h2>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
          {t('createTier')}
        </Button>
      </div>

      <Table
        columns={columns}
        dataSource={tiers}
        rowKey="id"
        loading={isLoading}
      />

      <Modal
        title={editingTier ? t('editTier') : t('createTier')}
        open={modalVisible}
        onOk={handleModalOk}
        onCancel={() => setModalVisible(false)}
        width={700}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="name"
            label={t('name')}
            rules={[{ required: true, message: t('nameRequired') }]}
          >
            <Input />
          </Form.Item>

          <Form.Item name="description" label={t('description')}>
            <Input.TextArea rows={2} />
          </Form.Item>

          <Form.Item
            name="allowed_models"
            label={t('allowedModels')}
          >
            <Select
              mode="multiple"
              placeholder={t('selectModels')}
              style={{ width: '100%' }}
            >
              {Object.entries(groupedModels).map(([providerId, provider]) => (
                <Select.OptGroup key={providerId} label={provider.name}>
                  {provider.models.map(model => (
                    <Select.Option key={model} value={model}>
                      {model}
                    </Select.Option>
                  ))}
                </Select.OptGroup>
              ))}
            </Select>
          </Form.Item>

          <Form.Item
            name="allowed_providers"
            label={t('allowedProviders')}
          >
            <Select
              mode="multiple"
              placeholder={t('selectProviders')}
              style={{ width: '100%' }}
            >
              {mockProviders.map(provider => (
                <Select.Option key={provider.id} value={provider.id}>
                  {provider.name}
                </Select.Option>
              ))}
            </Select>
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title={selectedTier?.name}
        open={detailVisible}
        onCancel={() => setDetailVisible(false)}
        footer={[
          <Button key="close" onClick={() => setDetailVisible(false)}>
            {t('close')}
          </Button>,
        ]}
        width={600}
      >
        {selectedTier && (
          <div>
            <Descriptions column={1} bordered>
              <Descriptions.Item label={t('name')}>{selectedTier.name}</Descriptions.Item>
              <Descriptions.Item label={t('description')}>{selectedTier.description}</Descriptions.Item>
              <Descriptions.Item label={t('isDefault')}>
                {selectedTier.is_default ? t('yes') : t('no')}
              </Descriptions.Item>
            </Descriptions>

            <Card title={t('allowedModels')} style={{ marginTop: 16 }}>
              <Collapse
                items={Object.entries(groupedModels).map(([providerId, provider]) => ({
                  key: providerId,
                  label: provider.name,
                  children: (
                    <div>
                      {provider.models.map(model => {
                        const fullModel = model;
                        const isAllowed = selectedTier.allowed_models?.includes(fullModel) ||
                          selectedTier.allowed_models?.includes(`${providerId}:*`);
                        return (
                          <Tag key={fullModel} color={isAllowed ? 'green' : 'default'}>
                            {fullModel} {isAllowed ? '✓' : ''}
                          </Tag>
                        );
                      })}
                    </div>
                  ),
                }))}
              />
            </Card>

            <Card title={t('allowedProviders')} style={{ marginTop: 16 }}>
              {mockProviders.map(provider => {
                const isAllowed = selectedTier.allowed_providers?.includes(provider.id) ||
                  selectedTier.allowed_providers?.includes('*');
                return (
                  <Tag key={provider.id} color={isAllowed ? 'green' : 'default'}>
                    {provider.name} {isAllowed ? '✓' : ''}
                  </Tag>
                );
              })}
            </Card>
          </div>
        )}
      </Modal>
    </div>
  );
};