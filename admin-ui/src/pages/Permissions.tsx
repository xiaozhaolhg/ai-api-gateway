import React, { useState } from 'react';
import { Table, Button, Modal, Form, Input, Select, Popconfirm, Tag, Empty, message } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient, type Permission, type Group } from '../api/client';

export const Permissions: React.FC = () => {
  const { t } = useTranslation(['permissions', 'common']);
  const queryClient = useQueryClient();
  const [modalVisible, setModalVisible] = useState(false);
  const [editingPermission, setEditingPermission] = useState<Permission | null>(null);
  const [form] = Form.useForm();

  const { data: permissions = [], isLoading: permsLoading } = useQuery({
    queryKey: ['permissions'],
    queryFn: () => apiClient.getPermissions(),
  });

  const { data: groups = [] } = useQuery({
    queryKey: ['groups'],
    queryFn: () => apiClient.getGroups(),
  });

  const isLoading = permsLoading;

  const createMutation = useMutation({
    mutationFn: (data: Omit<Permission, 'id' | 'created_at' | 'updated_at'>) => apiClient.createPermission(data),
    onMutate: async (newPerm) => {
      await queryClient.cancelQueries({ queryKey: ['permissions'] });
      const previous = queryClient.getQueryData<Permission[]>(['permissions']);
      queryClient.setQueryData<Permission[]>(['permissions'], (old = []) => [
        ...old, { ...newPerm, id: `temp-${Date.now()}`, created_at: new Date().toISOString(), updated_at: new Date().toISOString() } as Permission,
      ]);
      return { previous };
    },
    onError: (_err, _new, context) => {
      if (context?.previous) queryClient.setQueryData(['permissions'], context.previous);
      message.error('Failed to create permission');
    },
    onSuccess: () => {
      message.success('Permission created successfully');
      queryClient.invalidateQueries({ queryKey: ['permissions'] });
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<Permission> }) => apiClient.updatePermission(id, data),
    onMutate: async ({ id, data }) => {
      await queryClient.cancelQueries({ queryKey: ['permissions'] });
      const previous = queryClient.getQueryData<Permission[]>(['permissions']);
      queryClient.setQueryData<Permission[]>(['permissions'], (old = []) => old.map(p => p.id === id ? { ...p, ...data } : p));
      return { previous };
    },
    onError: (_err, _vars, context) => {
      if (context?.previous) queryClient.setQueryData(['permissions'], context.previous);
      message.error('Failed to update permission');
    },
    onSuccess: () => {
      message.success('Permission updated successfully');
      queryClient.invalidateQueries({ queryKey: ['permissions'] });
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => apiClient.deletePermission(id),
    onMutate: async (id) => {
      await queryClient.cancelQueries({ queryKey: ['permissions'] });
      const previous = queryClient.getQueryData<Permission[]>(['permissions']);
      queryClient.setQueryData<Permission[]>(['permissions'], (old = []) => old.filter(p => p.id !== id));
      return { previous };
    },
    onError: (_err, _id, context) => {
      if (context?.previous) queryClient.setQueryData(['permissions'], context.previous);
      message.error('Failed to delete permission');
    },
    onSuccess: () => {
      message.success('Permission deleted successfully');
    },
  });

  const handleAdd = () => {
    setEditingPermission(null);
    form.resetFields();
    setModalVisible(true);
  };

  const handleEdit = (permission: Permission) => {
    setEditingPermission(permission);
    form.setFieldsValue(permission);
    setModalVisible(true);
  };

  const handleModalOk = async () => {
    const values = await form.validateFields();
    if (editingPermission) {
      updateMutation.mutate({ id: editingPermission.id, data: values });
    } else {
      createMutation.mutate(values);
    }
    setModalVisible(false);
  };

  const handleModalCancel = () => {
    setModalVisible(false);
    form.resetFields();
  };

  const columns = [
    {
      title: t('permissions:fields.group'),
      dataIndex: 'group_id',
      key: 'group_id',
      render: (groupId: string) => {
        const group = groups.find((g: Group) => g.id === groupId);
        return group ? group.name : groupId;
      },
    },
    {
      title: t('permissions:fields.modelPattern'),
      dataIndex: 'model_pattern',
      key: 'model_pattern',
    },
    {
      title: t('permissions:fields.effect'),
      dataIndex: 'effect',
      key: 'effect',
      render: (effect: string) => (
        <Tag color={effect === 'allow' ? 'green' : 'red'}>{effect}</Tag>
      ),
    },
    {
      title: t('permissions:fields.status'),
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
      render: (_: any, record: Permission) => (
        <div>
          <Button type="link" icon={<EditOutlined />} onClick={() => handleEdit(record)}>
            {t('common:edit')}
          </Button>
          <Popconfirm
            title="Are you sure you want to delete this permission?"
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
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <h2>{t('permissions:title')}</h2>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
          {t('permissions:addPermission')}
        </Button>
      </div>

      <Table
        dataSource={permissions}
        columns={columns}
        rowKey="id"
        pagination={{ pageSize: 10 }}
        locale={{ emptyText: <Empty description="No permissions found" /> }}
      />

      <Modal
        title={editingPermission ? t('permissions:editPermission') : t('permissions:addPermission')}
        open={modalVisible}
        onOk={handleModalOk}
        onCancel={handleModalCancel}
        confirmLoading={createMutation.isPending || updateMutation.isPending}
        width={600}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            label={t('permissions:fields.group')}
            name="group_id"
            rules={[{ required: true, message: 'Please select group' }]}
          >
            <Select>
              {groups.map((group: Group) => (
                <Select.Option key={group.id} value={group.id}>
                  {group.name}
                </Select.Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item
            label={t('permissions:fields.modelPattern')}
            name="model_pattern"
            rules={[{ required: true, message: 'Please input model pattern' }]}
          >
            <Input placeholder="gpt-*" />
          </Form.Item>

          <Form.Item
            label={t('permissions:fields.effect')}
            name="effect"
            rules={[{ required: true, message: 'Please select effect' }]}
          >
            <Select>
              <Select.Option value="allow">Allow</Select.Option>
              <Select.Option value="deny">Deny</Select.Option>
            </Select>
          </Form.Item>

          <Form.Item
            label={t('permissions:fields.status')}
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
};
