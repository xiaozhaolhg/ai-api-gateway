import { useState } from 'react';
import { Table, Button, Modal, Form, Input, Select, Popconfirm, Tag, Empty, message } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient, type User, type Group } from '../api/client';

export default function Users() {
  const { t } = useTranslation(['users', 'common']);
  const queryClient = useQueryClient();
  const [modalVisible, setModalVisible] = useState(false);
  const [editingUser, setEditingUser] = useState<User | null>(null);
  const [form] = Form.useForm();
  const [selectedGroups, setSelectedGroups] = useState<string[]>([]);

  const { data: users = [], isLoading } = useQuery({
    queryKey: ['users'],
    queryFn: () => apiClient.getUsers(),
  });

  const [searchText, setSearchText] = useState('');
  const [roleFilter, setRoleFilter] = useState<string | undefined>(undefined);

  const filteredUsers = users.filter((u: User) => {
    const matchesSearch = !searchText ||
      u.name.toLowerCase().includes(searchText.toLowerCase()) ||
      u.email.toLowerCase().includes(searchText.toLowerCase()) ||
      (u.username && u.username.toLowerCase().includes(searchText.toLowerCase()));
    const matchesRole = !roleFilter || u.role === roleFilter;
    return matchesSearch && matchesRole;
  });

  const { data: groups = [] } = useQuery({
    queryKey: ['groups'],
    queryFn: () => apiClient.getGroups(),
  });

  const createMutation = useMutation({
    mutationFn: (data: Omit<User, 'id' | 'created_at'>) => apiClient.createUser(data),
    onSuccess: async (newUser) => {
      message.success('User created successfully');
      for (const groupId of selectedGroups) {
        try {
          await apiClient.addGroupMember(groupId, newUser.id);
        } catch (error) {
          console.error(`Failed to add user to group ${groupId}:`, error);
        }
      }
      queryClient.invalidateQueries({ queryKey: ['users'] });
      queryClient.invalidateQueries({ queryKey: ['groups'] });
    },
    onError: () => {
      message.error('Failed to create user');
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<User> }) => apiClient.updateUser(id, data),
    onSuccess: async (updatedUser) => {
      message.success('User updated successfully');
      const oldGroups = editingUser?.groups || [];
      const newGroups = selectedGroups;
      const added = newGroups.filter(g => !oldGroups.includes(g));
      const removed = oldGroups.filter(g => !newGroups.includes(g));

      for (const groupId of added) {
        try {
          await apiClient.addGroupMember(groupId, updatedUser.id);
        } catch (error) {
          console.error(`Failed to add user to group ${groupId}:`, error);
        }
      }
      for (const groupId of removed) {
        try {
          await apiClient.removeGroupMember(groupId, updatedUser.id);
        } catch (error) {
          console.error(`Failed to remove user from group ${groupId}:`, error);
        }
      }

      queryClient.invalidateQueries({ queryKey: ['users'] });
      queryClient.invalidateQueries({ queryKey: ['groups'] });
    },
    onError: () => {
      message.error('Failed to update user');
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => apiClient.deleteUser(id),
    onMutate: async (id) => {
      await queryClient.cancelQueries({ queryKey: ['users'] });
      const previous = queryClient.getQueryData<User[]>(['users']);
      queryClient.setQueryData<User[]>(['users'], (old = []) => old.filter(u => u.id !== id));
      return { previous };
    },
    onError: (_err, _id, context) => {
      if (context?.previous) queryClient.setQueryData(['users'], context.previous);
      message.error('Failed to delete user');
    },
    onSuccess: () => {
      message.success('User deleted successfully');
    },
  });

  const handleAdd = () => {
    setEditingUser(null);
    setSelectedGroups([]);
    form.resetFields();
    setModalVisible(true);
  };

  const handleEdit = (user: User) => {
    setEditingUser(user);
    setSelectedGroups(user.groups || []);
    form.setFieldsValue(user);
    setModalVisible(true);
  };

  const handleModalOk = async () => {
    const values = await form.validateFields();
    const { password, ...userData } = values;

    // Validate username uniqueness for new users
    if (!editingUser && userData.username) {
      try {
        const response = await apiClient.checkUsernameAvailability(userData.username);
        if (!response.available) {
          message.error('Username already taken');
          return;
        }
      } catch (error) {
        message.error('Failed to check username availability');
        return;
      }
    }

    if (editingUser) {
      updateMutation.mutate({ id: editingUser.id, data: userData });
    } else {
      const payload = { ...userData, password };
      createMutation.mutate(payload);
    }
    setModalVisible(false);
  };

  const handleModalCancel = () => {
    setModalVisible(false);
    form.resetFields();
  };

  const columns = [
    {
      title: t('users:fields.name'),
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: t('users:fields.email'),
      dataIndex: 'email',
      key: 'email',
    },
    {
      title: 'Groups',
      key: 'groups',
      render: (_: any, record: User) => (
        <>
          {record.groups?.map((groupId: string) => {
            const group = groups.find((g: Group) => g.id === groupId);
            return group ? (
              <Tag key={groupId} style={{ marginBottom: 4 }}>{group.name}</Tag>
            ) : null;
          })}
        </>
      ),
    },
    {
      title: t('users:fields.role'),
      dataIndex: 'role',
      key: 'role',
      render: (role: string) => (
        <Tag color={role === 'admin' ? 'blue' : role === 'viewer' ? 'orange' : 'default'}>{role}</Tag>
      ),
    },
    {
      title: t('users:fields.status'),
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={status === 'active' ? 'green' : 'red'}>{status}</Tag>
      ),
    },
    {
      title: t('common:createdAt'),
      dataIndex: 'created_at',
      key: 'created_at',
    },
    {
      title: t('common:actions'),
      key: 'actions',
      render: (_: any, record: User) => (
        <div>
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
          >
            {t('common:edit')}
          </Button>
          <Popconfirm
            title="Are you sure you want to delete this user?"
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
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h2>{t('users:title')}</h2>
        <div style={{ display: 'flex', gap: 8 }}>
          <Input
            placeholder="Search by name or email..."
            value={searchText}
            onChange={(e) => setSearchText(e.target.value)}
            style={{ width: 250 }}
            allowClear
          />
          <Select
            style={{ width: 150 }}
            placeholder="Filter by role"
            value={roleFilter}
            onChange={(value) => setRoleFilter(value)}
            allowClear
          >
            <Select.Option value="admin">Admin</Select.Option>
            <Select.Option value="user">User</Select.Option>
            <Select.Option value="viewer">Viewer</Select.Option>
          </Select>
          <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
            {t('users:addUser')}
          </Button>
        </div>
      </div>

      {filteredUsers.length === 0 ? (
        <Empty description="No users found" />
      ) : (
        <Table
          dataSource={filteredUsers}
          columns={columns}
          rowKey="id"
          pagination={{ pageSize: 10 }}
        />
      )}

      <Modal
        title={editingUser ? t('users:editUser') : t('users:addUser')}
        open={modalVisible}
        onOk={handleModalOk}
        onCancel={handleModalCancel}
        confirmLoading={createMutation.isPending || updateMutation.isPending}
        width={600}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            label={t('users:fields.name')}
            name="name"
            rules={[{ required: true, message: 'Please input user name' }]}
          >
            <Input />
          </Form.Item>

          <Form.Item
            label={t('users:fields.username')}
            name="username"
            rules={[
              { required: true, message: 'Please input username' },
              { pattern: /^[a-zA-Z0-9_]+$/, message: 'Username can only contain letters, numbers, and underscores' },
            ]}
          >
            <Input />
          </Form.Item>

          <Form.Item
            label={t('users:fields.email')}
            name="email"
            rules={[{ required: true, message: 'Please input email' }]}
          >
            <Input type="email" />
          </Form.Item>

          {!editingUser && (
            <Form.Item
              label="Password"
              name="password"
              rules={[
                { required: true, message: 'Please input password' },
                { min: 8, message: 'Password must be at least 8 characters' },
              ]}
            >
              <Input.Password placeholder="Minimum 8 characters" />
            </Form.Item>
          )}

          <Form.Item
            label={t('users:fields.role')}
            name="role"
            rules={[{ required: true, message: 'Please select role' }]}
          >
            <Select>
              <Select.Option value="user">User</Select.Option>
              <Select.Option value="admin">Admin</Select.Option>
              <Select.Option value="viewer">Viewer</Select.Option>
            </Select>
          </Form.Item>

          <Form.Item
            label="Groups"
            name="groups"
            valuePropName="value"
          >
            <Select
              mode="multiple"
              placeholder="Select groups"
              value={selectedGroups}
              onChange={(values) => setSelectedGroups(values)}
            >
              {groups.map((group: Group) => (
                <Select.Option key={group.id} value={group.id}>
                  {group.name}
                </Select.Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item
            label={t('users:fields.status')}
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
