import { useState, useEffect } from 'react';
import { Table, Button, Modal, Form, Input, Select, Popconfirm, Tag, Spin, Empty, message } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { apiClient, type Provider } from '../api/client';

export default function Providers() {
  const { t } = useTranslation(['providers', 'common']);
  const [providers, setProviders] = useState<Provider[]>([]);
  const [loading, setLoading] = useState(true);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingProvider, setEditingProvider] = useState<Provider | null>(null);
  const [form] = Form.useForm();

  useEffect(() => {
    loadProviders();
  }, []);

  const loadProviders = async () => {
    try {
      setLoading(true);
      const data = await apiClient.getProviders();
      setProviders(data);
    } catch (error) {
      message.error('Failed to load providers');
    } finally {
      setLoading(false);
    }
  };

  const handleAdd = () => {
    setEditingProvider(null);
    form.resetFields();
    setModalVisible(true);
  };

  const handleEdit = (provider: Provider) => {
    setEditingProvider(provider);
    form.setFieldsValue({
      ...provider,
      models: provider.models.join(', '),
    });
    setModalVisible(true);
  };

  const handleDelete = async (id: string) => {
    try {
      await apiClient.deleteProvider(id);
      message.success('Provider deleted successfully');
      loadProviders();
    } catch (error) {
      message.error('Failed to delete provider');
    }
  };

  const handleModalOk = async () => {
    try {
      const values = await form.validateFields();
      const providerData = {
        ...values,
        models: values.models.split(',').map((m: string) => m.trim()),
      };

      if (editingProvider) {
        await apiClient.updateProvider(editingProvider.id, providerData);
        message.success('Provider updated successfully');
      } else {
        await apiClient.createProvider(providerData);
        message.success('Provider created successfully');
      }

      setModalVisible(false);
      loadProviders();
    } catch (error) {
      message.error('Failed to save provider');
    }
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
      render: (models: string[]) => models.join(', '),
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
            onConfirm={() => handleDelete(record.id)}
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

  if (loading) {
    return <Spin size="large" />;
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
