import { useState, useEffect } from 'react';
import { Table, Button, Modal, Form, Input, Select, InputNumber, Popconfirm, Tag, Tabs, Spin, message } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, CheckOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { apiClient, type AlertRule, type Alert } from '../api/client';

export const Alerts: React.FC = () => {
  const { t } = useTranslation(['alerts', 'common']);
  const [alertRules, setAlertRules] = useState<AlertRule[]>([]);
  const [alerts, setAlerts] = useState<Alert[]>([]);
  const [loading, setLoading] = useState(true);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingRule, setEditingRule] = useState<AlertRule | null>(null);
  const [form] = Form.useForm();

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      setLoading(true);
      const [rulesData, alertsData] = await Promise.all([
        apiClient.getAlertRules(),
        apiClient.getAlerts(),
      ]);
      setAlertRules(rulesData);
      setAlerts(alertsData);
    } catch (error) {
      message.error('Failed to load alerts');
    } finally {
      setLoading(false);
    }
  };

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
    try {
      await apiClient.deleteAlertRule(id);
      message.success('Alert rule deleted successfully');
      loadData();
    } catch (error) {
      message.error('Failed to delete alert rule');
    }
  };

  const handleAcknowledge = async (id: string) => {
    try {
      await apiClient.acknowledgeAlert(id);
      message.success('Alert acknowledged successfully');
      loadData();
    } catch (error) {
      message.error('Failed to acknowledge alert');
    }
  };

  const handleModalOk = async () => {
    try {
      const values = await form.validateFields();

      if (editingRule) {
        await apiClient.updateAlertRule(editingRule.id, values);
        message.success('Alert rule updated successfully');
      } else {
        await apiClient.createAlertRule(values);
        message.success('Alert rule created successfully');
      }

      setModalVisible(false);
      loadData();
    } catch (error) {
      message.error('Failed to save alert rule');
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
        <Tag color={status === 'active' ? 'green' : 'red'}>{status}</Tag>
      ),
    },
    {
      title: t('common:actions'),
      key: 'actions',
      render: (_: any, record: AlertRule) => (
        <div>
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
          >
            {t('common:edit')}
          </Button>
          <Popconfirm
            title="Are you sure you want to delete this alert rule?"
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

  const alertColumns = [
    {
      title: t('alerts:fields.severity'),
      dataIndex: 'severity',
      key: 'severity',
      render: (severity: string) => {
        const color = severity === 'critical' ? 'red' : severity === 'warning' ? 'orange' : 'blue';
        return <Tag color={color}>{severity}</Tag>;
      },
    },
    {
      title: t('alerts:fields.status'),
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={status === 'firing' ? 'red' : 'green'}>{status}</Tag>
      ),
    },
    {
      title: t('alerts:fields.triggeredAt'),
      dataIndex: 'triggered_at',
      key: 'triggered_at',
    },
    {
      title: t('alerts:fields.description'),
      dataIndex: 'description',
      key: 'description',
    },
    {
      title: t('common:actions'),
      key: 'actions',
      render: (_: any, record: Alert) => (
        <div>
          {record.status === 'firing' && (
            <Button
              type="link"
              icon={<CheckOutlined />}
              onClick={() => handleAcknowledge(record.id)}
            >
              {t('alerts:acknowledge')}
            </Button>
          )}
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
        <h2>{t('alerts:title')}</h2>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
          {t('alerts:addRule')}
        </Button>
      </div>

      <Tabs
        defaultActiveKey="rules"
        items={[
          {
            key: 'rules',
            label: t('alerts:tabs.rules'),
            children: (
              <Table
                dataSource={alertRules}
                columns={ruleColumns}
                rowKey="id"
                pagination={{ pageSize: 10 }}
              />
            ),
          },
          {
            key: 'alerts',
            label: t('alerts:tabs.activeAlerts'),
            children: (
              <Table
                dataSource={alerts}
                columns={alertColumns}
                rowKey="id"
                pagination={{ pageSize: 10 }}
              />
            ),
          },
        ]}
      />

      <Modal
        title={editingRule ? t('alerts:editRule') : t('alerts:addRule')}
        open={modalVisible}
        onOk={handleModalOk}
        onCancel={handleModalCancel}
        width={600}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            label={t('alerts:fields.name')}
            name="name"
            rules={[{ required: true, message: 'Please input rule name' }]}
          >
            <Input />
          </Form.Item>

          <Form.Item
            label={t('alerts:fields.metric')}
            name="metric"
            rules={[{ required: true, message: 'Please select metric' }]}
          >
            <Select>
              <Select.Option value="spend">Spend</Select.Option>
              <Select.Option value="tokens">Tokens</Select.Option>
              <Select.Option value="requests">Requests</Select.Option>
              <Select.Option value="errors">Errors</Select.Option>
            </Select>
          </Form.Item>

          <Form.Item
            label={t('alerts:fields.condition')}
            name="condition"
            rules={[{ required: true, message: 'Please select condition' }]}
          >
            <Select>
              <Select.Option value="greater_than">Greater Than</Select.Option>
              <Select.Option value="less_than">Less Than</Select.Option>
              <Select.Option value="equals">Equals</Select.Option>
            </Select>
          </Form.Item>

          <Form.Item
            label={t('alerts:fields.threshold')}
            name="threshold"
            rules={[{ required: true, message: 'Please input threshold' }]}
          >
            <InputNumber style={{ width: '100%' }} />
          </Form.Item>

          <Form.Item
            label={t('alerts:fields.channel')}
            name="channel"
            rules={[{ required: true, message: 'Please select channel' }]}
          >
            <Select>
              <Select.Option value="email">Email</Select.Option>
              <Select.Option value="slack">Slack</Select.Option>
              <Select.Option value="webhook">Webhook</Select.Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};
