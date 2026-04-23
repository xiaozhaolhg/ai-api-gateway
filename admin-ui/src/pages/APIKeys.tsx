import { useState, useEffect } from 'react';
import { Table, Button, Modal, Form, Input, Select, Popconfirm, Spin, Empty, Alert, message } from 'antd';
import { PlusOutlined, DeleteOutlined, CopyOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { apiClient, type APIKey, type User } from '../api/client';

export default function APIKeys() {
  const { t } = useTranslation(['apiKeys', 'common']);
  const [apiKeys, setApiKeys] = useState<APIKey[]>([]);
  const [users, setUsers] = useState<User[]>([]);
  const [selectedUserId, setSelectedUserId] = useState('');
  const [loading, setLoading] = useState(true);
  const [modalVisible, setModalVisible] = useState(false);
  const [createdKey, setCreatedKey] = useState<{ api_key_id: string; api_key: string } | null>(null);
  const [newKeyName, setNewKeyName] = useState('');

  useEffect(() => {
    loadUsers();
  }, []);

  useEffect(() => {
    if (selectedUserId) {
      loadAPIKeys();
    }
  }, [selectedUserId]);

  const loadUsers = async () => {
    try {
      setLoading(true);
      const data = await apiClient.getUsers();
      setUsers(data);
    } catch (error) {
      message.error('Failed to load users');
    } finally {
      setLoading(false);
    }
  };

  const loadAPIKeys = async () => {
    try {
      const data = await apiClient.getAPIKeys(selectedUserId);
      setApiKeys(data);
    } catch (error) {
      message.error('Failed to load API keys');
    }
  };

  const handleCreateAPIKey = async () => {
    try {
      const result = await apiClient.createAPIKey(selectedUserId, newKeyName);
      setCreatedKey(result);
      setModalVisible(false);
      setNewKeyName('');
      loadAPIKeys();
    } catch (error) {
      message.error('Failed to create API key');
    }
  };

  const handleDeleteAPIKey = async (id: string) => {
    try {
      await apiClient.deleteAPIKey(id);
      message.success('API key revoked successfully');
      loadAPIKeys();
    } catch (error) {
      message.error('Failed to revoke API key');
    }
  };

  const handleCopyKey = (key: string) => {
    navigator.clipboard.writeText(key);
    message.success(t('apiKeys:keyCopied'));
  };

  const columns = [
    {
      title: t('apiKeys:fields.name'),
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: t('apiKeys:fields.scopes'),
      dataIndex: 'scopes',
      key: 'scopes',
      render: (scopes: string[]) => scopes.join(', '),
    },
    {
      title: t('common:createdAt'),
      dataIndex: 'created_at',
      key: 'created_at',
    },
    {
      title: t('apiKeys:fields.expiresAt'),
      dataIndex: 'expires_at',
      key: 'expires_at',
      render: (expiresAt: string) => expiresAt ? new Date(expiresAt).toLocaleString() : 'Never',
    },
    {
      title: t('common:actions'),
      key: 'actions',
      render: (_: any, record: APIKey) => (
        <Popconfirm
          title="Are you sure you want to revoke this API key?"
          onConfirm={() => handleDeleteAPIKey(record.id)}
          okText="Yes"
          cancelText="No"
        >
          <Button type="link" danger icon={<DeleteOutlined />}>
            {t('common:delete')}
          </Button>
        </Popconfirm>
      ),
    },
  ];

  if (loading) {
    return <Spin size="large" />;
  }

  return (
    <div>
      <h2 style={{ marginBottom: 16 }}>{t('apiKeys:title')}</h2>

      <div style={{ marginBottom: 16 }}>
        <Select
          style={{ width: 300 }}
          placeholder={t('apiKeys:fields.user')}
          value={selectedUserId}
          onChange={setSelectedUserId}
          options={users.map(user => ({ label: user.name, value: user.id }))}
        />
      </div>

      {selectedUserId && (
        <>
          <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
            <h3>API Keys for {users.find(u => u.id === selectedUserId)?.name}</h3>
            <Button type="primary" icon={<PlusOutlined />} onClick={() => setModalVisible(true)}>
              {t('apiKeys:addKey')}
            </Button>
          </div>

          {createdKey && (
            <Alert
              type="success"
              message={t('apiKeys:title')}
              description={
                <div>
                  <p style={{ marginBottom: 8 }}>This key will only be shown once. Copy it now.</p>
                  <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                    <code style={{ background: '#f6ffed', padding: '4px 8px', borderRadius: 4 }}>{createdKey.api_key}</code>
                    <Button
                      icon={<CopyOutlined />}
                      onClick={() => handleCopyKey(createdKey.api_key)}
                    >
                      {t('apiKeys:copyKey')}
                    </Button>
                  </div>
                </div>
              }
              closable
              onClose={() => setCreatedKey(null)}
              style={{ marginBottom: 16 }}
            />
          )}

          {apiKeys.length === 0 ? (
            <Empty description="No API keys found" />
          ) : (
            <Table
              dataSource={apiKeys}
              columns={columns}
              rowKey="id"
              pagination={{ pageSize: 10 }}
            />
          )}
        </>
      )}

      <Modal
        title={t('apiKeys:addKey')}
        open={modalVisible}
        onOk={handleCreateAPIKey}
        onCancel={() => {
          setModalVisible(false);
          setNewKeyName('');
        }}
      >
        <Form layout="vertical">
          <Form.Item
            label={t('apiKeys:fields.name')}
            rules={[{ required: true, message: 'Please input key name' }]}
          >
            <Input
              value={newKeyName}
              onChange={(e) => setNewKeyName(e.target.value)}
              placeholder="Production Key"
            />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}
