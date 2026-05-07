import { useState } from 'react';
import { Table, Button, Select, Popconfirm, message } from 'antd';
import { PlusOutlined, DeleteOutlined } from '@ant-design/icons';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient, type User } from '../api/client';

interface GroupMembersTabProps {
  groupId: string;
}

export const GroupMembersTab: React.FC<GroupMembersTabProps> = ({ groupId }) => {
  const queryClient = useQueryClient();
  const [userSelectorVisible, setUserSelectorVisible] = useState(false);

  const { data: members = [], isLoading } = useQuery({
    queryKey: ['groupMembers', groupId],
    queryFn: () => apiClient.getGroupMembers(groupId),
  });

  const { data: allUsers = [] } = useQuery({
    queryKey: ['users'],
    queryFn: () => apiClient.getUsers(),
  });

  const removeMutation = useMutation({
    mutationFn: (userId: string) => apiClient.removeGroupMember(groupId, userId),
    onSuccess: () => {
      message.success('Member removed successfully');
      queryClient.invalidateQueries({ queryKey: ['groupMembers', groupId] });
    },
    onError: () => {
      message.error('Failed to remove member');
    },
  });

  const addMutation = useMutation({
    mutationFn: (userId: string) => apiClient.addGroupMember(groupId, userId),
    onSuccess: () => {
      message.success('Member added successfully');
      queryClient.invalidateQueries({ queryKey: ['groupMembers', groupId] });
      setUserSelectorVisible(false);
    },
    onError: () => {
      message.error('Failed to add member');
    },
  });

  const columns = [
    {
      title: 'Name',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: 'Email',
      dataIndex: 'email',
      key: 'email',
    },
    {
      title: 'Role',
      dataIndex: 'role',
      key: 'role',
    },
    {
      title: 'Actions',
      key: 'actions',
      render: (_: any, record: User) => (
        <Popconfirm
          title="Are you sure you want to remove this member?"
          onConfirm={() => removeMutation.mutate(record.id)}
          okText="Yes"
          cancelText="No"
        >
          <Button
            type="link"
            danger
            icon={<DeleteOutlined />}
            loading={removeMutation.isPending}
          >
            Remove
          </Button>
        </Popconfirm>
      ),
    },
  ];

  const availableUsers = allUsers.filter(
    (u) => !members.some((m) => m.id === u.id)
  );

  return (
    <div>
      <div style={{ marginBottom: 16 }}>
        <Button
          type="primary"
          icon={<PlusOutlined />}
          onClick={() => setUserSelectorVisible(true)}
        >
          Add Member
        </Button>
      </div>

      {userSelectorVisible && (
        <div style={{ marginBottom: 16 }}>
          <Select
            style={{ width: 300 }}
            placeholder="Select a user to add"
            onChange={(userId) => addMutation.mutate(userId)}
            loading={addMutation.isPending}
          >
            {availableUsers.map((user) => (
              <Select.Option key={user.id} value={user.id}>
                {user.name} ({user.email})
              </Select.Option>
            ))}
          </Select>
        </div>
      )}

      <Table
        dataSource={members}
        columns={columns}
        rowKey="id"
        loading={isLoading}
        pagination={false}
        size="small"
      />
    </div>
  );
};
