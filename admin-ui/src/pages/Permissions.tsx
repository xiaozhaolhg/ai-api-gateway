import React, { useState, useEffect } from 'react';
import { Table, Button, Modal, Form, Input, Select, Popconfirm, Tag, Spin, Empty, message } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { apiClient, type Permission } from '../api/client';

export const Permissions: React.FC = () => {
  const { t } = useTranslation(['permissions', 'common']);
  const [permissions, setPermissions] = useState<Permission[]>([]);
  const [groups, setGroups] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingPermission, setEditingPermission] = useState<Permission | null>(null);
  const [form] = Form.useForm();

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      setLoading(true);
      const [permsData, groupsData] = await Promise.all([
        apiClient.getPermissions(),
        apiClient.getGroups(),
      ]);
      setPermissions(permsData);
      setGroups(groupsData);
    } catch (error) {
      message.error('Failed to load permissions');
    } finally {
      setLoading(false);
    }
  };

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

  const handleDelete = async (id: string) => {
    try {
      await apiClient.deletePermission(id);
      message.success('Permission deleted successfully');
      loadData();
    } catch (error) {
      message.error('Failed to delete permission');
    }
  };

  const handleModalOk = async () => {
    try {
      const values = await form.validateFields();

      if (editingPermission) {
        await apiClient.updatePermission(editingPermission.id, values);
        message.success('Permission updated successfully');
      } else {
        await apiClient.createPermission(values);
        message.success('Permission created successfully');
      }

      setModalVisible(false);
      loadData();
    } catch (error) {
      message.error('Failed to save permission');
    }
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
        const group = groups.find(g => g.id === groupId);
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
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
          >
            {t('common:edit')}
          </Button>
          <Popconfirm
            title="Are you sure you want to delete this permission?"
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
        <h2>{t('permissions:title')}</h2>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
          {t('permissions:addPermission')}
        </Button>
      </div>

      {permissions.length === 0 ? (
        <Empty description="No permissions found" />
      ) : (
        <Table
          dataSource={permissions}
          columns={columns}
          rowKey="id"
          pagination={{ pageSize: 10 }}
        />
      )}

      <Modal
        title={editingPermission ? t('permissions:editPermission') : t('permissions:addPermission')}
        open={modalVisible}
        onOk={handleModalOk}
        onCancel={handleModalCancel}
        width={600}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            label={t('permissions:fields.group')}
            name="group_id"
            rules={[{ required: true, message: 'Please select group' }]}
          >
            <Select>
              {groups.map(group => (
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
