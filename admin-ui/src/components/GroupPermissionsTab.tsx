import { useState } from 'react';
import { Table, Button, Select, Popconfirm, message } from 'antd';
import { PlusOutlined, DeleteOutlined } from '@ant-design/icons';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient, type Permission } from '../api/client';

interface GroupPermissionsTabProps {
  groupId: string;
}

export const GroupPermissionsTab: React.FC<GroupPermissionsTabProps> = ({ groupId }) => {
  const queryClient = useQueryClient();
  const [formVisible, setFormVisible] = useState(false);
  const [formData, setFormData] = useState({
    resource_type: '',
    resource_id: '',
    action: '',
    effect: 'allow',
  });

  const { data: permissions = [], isLoading } = useQuery({
    queryKey: ['groupPermissions', groupId],
    queryFn: () => apiClient.getGroupPermissions(groupId),
  });

  const createMutation = useMutation({
    mutationFn: (data: Omit<Permission, 'id' | 'created_at' | 'updated_at'>) =>
      apiClient.createPermission({ ...data, group_id: groupId }),
    onSuccess: () => {
      message.success('Permission added successfully');
      queryClient.invalidateQueries({ queryKey: ['groupPermissions', groupId] });
      setFormVisible(false);
      setFormData({ resource_type: '', resource_id: '', action: '', effect: 'allow' });
    },
    onError: () => {
      message.error('Failed to add permission');
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (permId: string) => apiClient.deletePermission(permId),
    onSuccess: () => {
      message.success('Permission removed successfully');
      queryClient.invalidateQueries({ queryKey: ['groupPermissions', groupId] });
    },
    onError: () => {
      message.error('Failed to remove permission');
    },
  });

  const columns = [
    {
      title: 'Resource Type',
      dataIndex: 'resource_type',
      key: 'resource_type',
    },
    {
      title: 'Resource ID',
      dataIndex: 'resource_id',
      key: 'resource_id',
    },
    {
      title: 'Action',
      dataIndex: 'action',
      key: 'action',
    },
    {
      title: 'Effect',
      dataIndex: 'effect',
      key: 'effect',
      render: (effect: string) => (
        <span style={{ color: effect === 'allow' ? 'green' : 'red' }}>{effect}</span>
      ),
    },
    {
      title: 'Actions',
      key: 'actions',
      render: (_: any, record: Permission) => (
        <Popconfirm
          title="Are you sure you want to remove this permission?"
          onConfirm={() => deleteMutation.mutate(record.id)}
          okText="Yes"
          cancelText="No"
        >
          <Button
            type="link"
            danger
            icon={<DeleteOutlined />}
            loading={deleteMutation.isPending}
          >
            Remove
          </Button>
        </Popconfirm>
      ),
    },
  ];

  const handleAddPermission = () => {
    createMutation.mutate(formData as any);
  };

  return (
    <div>
      <div style={{ marginBottom: 16 }}>
        <Button
          type="primary"
          icon={<PlusOutlined />}
          onClick={() => setFormVisible(true)}
        >
          Add Permission
        </Button>
      </div>

      {formVisible && (
        <div style={{ marginBottom: 16, padding: 16, border: '1px solid #d9d9d9', borderRadius: 8 }}>
          <div style={{ display: 'flex', gap: 8, marginBottom: 8 }}>
            <Select
              style={{ width: 150 }}
              placeholder="Resource Type"
              value={formData.resource_type || undefined}
              onChange={(value) => setFormData({ ...formData, resource_type: value })}
            >
              <Select.Option value="model">Model</Select.Option>
              <Select.Option value="provider">Provider</Select.Option>
              <Select.Option value="system">System</Select.Option>
            </Select>

            <input
              style={{ flex: 1, padding: '4px 11px', border: '1px solid #d9d9d9', borderRadius: 6 }}
              placeholder="Resource ID (e.g., gpt-4, ollama:*)"
              value={formData.resource_id}
              onChange={(e) => setFormData({ ...formData, resource_id: e.target.value })}
            />

            <Select
              style={{ width: 120 }}
              placeholder="Action"
              value={formData.action || undefined}
              onChange={(value) => setFormData({ ...formData, action: value })}
            >
              <Select.Option value="access">Access</Select.Option>
              <Select.Option value="manage">Manage</Select.Option>
            </Select>

            <Select
              style={{ width: 120 }}
              placeholder="Effect"
              value={formData.effect}
              onChange={(value) => setFormData({ ...formData, effect: value })}
            >
              <Select.Option value="allow">Allow</Select.Option>
              <Select.Option value="deny">Deny</Select.Option>
            </Select>

            <Button
              type="primary"
              onClick={handleAddPermission}
              loading={createMutation.isPending}
            >
              Add
            </Button>
            <Button onClick={() => setFormVisible(false)}>Cancel</Button>
          </div>
        </div>
      )}

      <Table
        dataSource={permissions}
        columns={columns}
        rowKey="id"
        loading={isLoading}
        pagination={false}
        size="small"
      />
    </div>
  );
};
