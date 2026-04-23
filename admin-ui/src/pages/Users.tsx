import { useState, useEffect } from 'react';
import { Table, Button, Modal, Form, Input, Select, Popconfirm, Tag, Spin, Empty, message } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { apiClient, type User } from '../api/client';

export default function Users() {
  const { t } = useTranslation(['users', 'common']);
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingUser, setEditingUser] = useState<User | null>(null);
  const [form] = Form.useForm();

  useEffect(() => {
    loadUsers();
  }, []);

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

  const handleAdd = () => {
    setEditingUser(null);
    form.resetFields();
    setModalVisible(true);
  };

  const handleEdit = (user: User) => {
    setEditingUser(user);
    form.setFieldsValue(user);
    setModalVisible(true);
  };

  const handleDelete = async (id: string) => {
    try {
      await apiClient.deleteUser(id);
      message.success('User deleted successfully');
      loadUsers();
    } catch (error) {
      message.error('Failed to delete user');
    }
  };

  const handleModalOk = async () => {
    try {
      const values = await form.validateFields();

      if (editingUser) {
        await apiClient.updateUser(editingUser.id, values);
        message.success('User updated successfully');
      } else {
        await apiClient.createUser(values);
        message.success('User created successfully');
      }

      setModalVisible(false);
      loadUsers();
    } catch (error) {
      message.error('Failed to save user');
    }
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
      title: t('users:fields.role'),
      dataIndex: 'role',
      key: 'role',
      render: (role: string) => (
        <Tag color={role === 'admin' ? 'blue' : 'default'}>{role}</Tag>
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
        <h2>{t('users:title')}</h2>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
          {t('users:addUser')}
        </Button>
      </div>

      {users.length === 0 ? (
        <Empty description="No users found" />
      ) : (
        <Table
          dataSource={users}
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
            label={t('users:fields.email')}
            name="email"
            rules={[{ required: true, message: 'Please input email' }]}
          >
            <Input type="email" />
          </Form.Item>

          <Form.Item
            label={t('users:fields.role')}
            name="role"
            rules={[{ required: true, message: 'Please select role' }]}
          >
            <Select>
              <Select.Option value="user">User</Select.Option>
              <Select.Option value="admin">Admin</Select.Option>
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
