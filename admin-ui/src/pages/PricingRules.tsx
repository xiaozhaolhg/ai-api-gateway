import React, { useState, useEffect } from 'react';
import { Table, Button, Modal, Form, Input, Select, InputNumber, Popconfirm, Spin, Empty, message } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { apiClient, type PricingRule } from '../api/client';

export const PricingRules: React.FC = () => {
  const { t } = useTranslation(['pricing', 'common']);
  const [rules, setRules] = useState<PricingRule[]>([]);
  const [loading, setLoading] = useState(true);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingRule, setEditingRule] = useState<PricingRule | null>(null);
  const [form] = Form.useForm();

  useEffect(() => {
    loadRules();
  }, []);

  const loadRules = async () => {
    try {
      setLoading(true);
      const data = await apiClient.getPricingRules();
      setRules(data);
    } catch (error) {
      message.error('Failed to load pricing rules');
    } finally {
      setLoading(false);
    }
  };

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

  const handleDelete = async (id: string) => {
    try {
      await apiClient.deletePricingRule(id);
      message.success('Pricing rule deleted successfully');
      loadRules();
    } catch (error) {
      message.error('Failed to delete pricing rule');
    }
  };

  const handleModalOk = async () => {
    try {
      const values = await form.validateFields();

      if (editingRule) {
        await apiClient.updatePricingRule(editingRule.id, values);
        message.success('Pricing rule updated successfully');
      } else {
        await apiClient.createPricingRule(values);
        message.success('Pricing rule created successfully');
      }

      setModalVisible(false);
      loadRules();
    } catch (error) {
      message.error('Failed to save pricing rule');
    }
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
      title: t('common:actions'),
      key: 'actions',
      render: (_: any, record: PricingRule) => (
        <div>
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
          >
            {t('common:edit')}
          </Button>
          <Popconfirm
            title="Are you sure you want to delete this pricing rule?"
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
        <h2>{t('pricing:title')}</h2>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
          {t('pricing:addRule')}
        </Button>
      </div>

      {rules.length === 0 ? (
        <Empty description="No pricing rules found" />
      ) : (
        <Table
          dataSource={rules}
          columns={columns}
          rowKey="id"
          pagination={{ pageSize: 10 }}
        />
      )}

      <Modal
        title={editingRule ? t('pricing:editRule') : t('pricing:addRule')}
        open={modalVisible}
        onOk={handleModalOk}
        onCancel={handleModalCancel}
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
