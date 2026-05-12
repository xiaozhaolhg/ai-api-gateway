import { useState, useEffect } from 'react';
import { Table, Button, Modal, Form, Input, Select, Popconfirm, Empty, Alert, message } from 'antd';
import { PlusOutlined, DeleteOutlined, CopyOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient, type APIKey } from '../api/client';
import { useAuth } from '../contexts/AuthContext';

export default function APIKeys() {
  const { t } = useTranslation(['apiKeys', 'common']);
  const queryClient = useQueryClient();
  const { user } = useAuth();
  const [selectedUserId, setSelectedUserId] = useState(user?.id || '');
  const [modalVisible, setModalVisible] = useState(false);
  const [newApiKey, setNewApiKey] = useState<string | null>(null);
  const [keyDismissed, setKeyDismissed] = useState(() => {
    return sessionStorage.getItem('api-key-dismissed') === 'true';
  });
  const [newKeyName, setNewKeyName] = useState('');

  // Clear key on navigation events and tab switching
  useEffect(() => {
    const handleNavigation = () => {
      if (newApiKey) {
        setNewApiKey(null);
        setKeyDismissed(true);
        sessionStorage.setItem('api-key-dismissed', 'true');
      }
    };

    const handleVisibilityChange = () => {
      // Clear key when user switches to another tab
      if (document.hidden && newApiKey) {
        setNewApiKey(null);
        setKeyDismissed(true);
        sessionStorage.setItem('api-key-dismissed', 'true');
      }
    };
    
    window.addEventListener('beforeunload', handleNavigation);
    document.addEventListener('visibilitychange', handleVisibilityChange);
    
    return () => {
      window.removeEventListener('beforeunload', handleNavigation);
      document.removeEventListener('visibilitychange', handleVisibilityChange);
    };
  }, [newApiKey]);

  useEffect(() => {
    setNewApiKey(null);
    setKeyDismissed(false);
    sessionStorage.removeItem('api-key-dismissed');
    setNewKeyName('');
  }, [selectedUserId]);

  const { data: users = [], isLoading: usersLoading } = useQuery({
    queryKey: ['users'],
    queryFn: () => apiClient.getUsers(),
  });

  const { data: apiKeys = [], isLoading: keysLoading } = useQuery({
    queryKey: ['apiKeys', selectedUserId],
    queryFn: () => apiClient.getAPIKeys(selectedUserId),
    enabled: !!selectedUserId,
  });

  const createMutation = useMutation({
    mutationFn: ({ userId, name }: { userId: string; name: string }) => apiClient.createAPIKey(userId, name),
    onSuccess: (result) => {
      setNewApiKey(result.api_key);
      setModalVisible(false);
      setNewKeyName('');
      message.success('API key created successfully');
      queryClient.invalidateQueries({ queryKey: ['apiKeys', selectedUserId] });
    },
    onError: () => {
      message.error('Failed to create API key');
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => apiClient.deleteAPIKey(id),
    onMutate: async (id) => {
      await queryClient.cancelQueries({ queryKey: ['apiKeys', selectedUserId] });
      const previous = queryClient.getQueryData<APIKey[]>(['apiKeys', selectedUserId]);
      queryClient.setQueryData<APIKey[]>(['apiKeys', selectedUserId], (old = []) => old.filter(k => k.id !== id));
      return { previous };
    },
    onError: (_err, _id, context) => {
      if (context?.previous) queryClient.setQueryData(['apiKeys', selectedUserId], context.previous);
      message.error('Failed to revoke API key');
    },
    onSuccess: () => {
      message.success('API key revoked successfully');
    },
  });

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
      title: t('common:common.actions'),
      key: 'actions',
      render: (_: any, record: APIKey) => (
        <Popconfirm
          title="Are you sure you want to revoke this API key?"
          onConfirm={() => deleteMutation.mutate(record.id)}
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

  if (usersLoading) {
    return <Table loading={true} columns={columns} dataSource={[]} rowKey="id" />;
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

          {newApiKey && !keyDismissed && (
            <Alert
              type="success"
              message={t('apiKeys:title')}
              description={
                <div>
                  <p style={{ marginBottom: 8, color: '#d46b08' }}>This API key will only be shown once. Please copy it now and store it securely.</p>
                  <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                    <code style={{ background: '#f6ffed', padding: '4px 8px', borderRadius: 4, wordBreak: 'break-all' }}>{newApiKey}</code>
                    <Button
                      icon={<CopyOutlined />}
                      onClick={() => handleCopyKey(newApiKey)}
                    >
                      {t('apiKeys:copyKey')}
                    </Button>
                  </div>
                </div>
              }
              closable
              onClose={() => {
                setNewApiKey(null);
                setKeyDismissed(true);
                sessionStorage.setItem('api-key-dismissed', 'true');
              }}
              style={{ marginBottom: 16 }}
            />
          )}

          {keyDismissed && (
            <Alert
              type="info"
              message="API Key Previously Shown"
              description="Previous API key was shown once. Generate a new key if needed."
              style={{ marginBottom: 16 }}
            />
          )}

          <Table
            dataSource={apiKeys}
            columns={columns}
            rowKey="id"
            loading={keysLoading}
            pagination={{ pageSize: 10 }}
            locale={{ emptyText: <Empty description="No API keys found" /> }}
          />
        </>
      )}

      <Modal
        title={t('apiKeys:addKey')}
        open={modalVisible}
        onOk={() => {
          if (newKeyName) createMutation.mutate({ userId: selectedUserId, name: newKeyName });
        }}
        onCancel={() => {
          setModalVisible(false);
          setNewKeyName('');
        }}
        confirmLoading={createMutation.isPending}
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
