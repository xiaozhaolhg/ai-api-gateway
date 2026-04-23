import { useState, useEffect } from 'react';
import { Table, Button, Modal, Form, Input, Popconfirm, Spin, Empty, message } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { apiClient, type Group } from '../api/client';

export const Groups: React.FC = () => {
  const { t } = useTranslation(['groups', 'common']);
  const [groups, setGroups] = useState<Group[]>([]);
  const [loading, setLoading] = useState(true);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingGroup, setEditingGroup] = useState<Group | null>(null);
  const [form] = Form.useForm();

  useEffect(() => {
    loadGroups();
  }, []);

  const loadGroups = async () => {
    try {
      setLoading(true);
      const data = await apiClient.getGroups();
      setGroups(data);
    } catch (error) {
      message.error('Failed to load groups');
    } finally {
      setLoading(false);
    }
  };

  const handleAdd = () => {
    setEditingGroup(null);
    form.resetFields();
    setModalVisible(true);
  };

  const handleEdit = (group: Group) => {
    setEditingGroup(group);
    form.setFieldsValue(group);
    setModalVisible(true);
  };

  const handleDelete = async (id: string) => {
    try {
      await apiClient.deleteGroup(id);
      message.success('Group deleted successfully');
      loadGroups();
    } catch (error) {
      message.error('Failed to delete group');
    }
  };

  const handleModalOk = async () => {
    try {
      const values = await form.validateFields();

      if (editingGroup) {
        await apiClient.updateGroup(editingGroup.id, values);
        message.success('Group updated successfully');
      } else {
        await apiClient.createGroup(values);
        message.success('Group created successfully');
      }

      setModalVisible(false);
      loadGroups();
    } catch (error) {
      message.error('Failed to save group');
    }
  };

  const handleModalCancel = () => {
    setModalVisible(false);
    form.resetFields();
  };

  const columns = [
    {
      title: t('groups:fields.name'),
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: t('groups:fields.description'),
      dataIndex: 'description',
      key: 'description',
    },
    {
      title: t('groups:fields.memberCount'),
      dataIndex: 'member_count',
      key: 'member_count',
    },
    {
      title: t('common:createdAt'),
      dataIndex: 'created_at',
      key: 'created_at',
    },
    {
      title: t('common:actions'),
      key: 'actions',
      render: (_: any, record: Group) => (
        <div>
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
          >
            {t('common:edit')}
          </Button>
          <Popconfirm
            title="Are you sure you want to delete this group?"
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
        <h2>{t('groups:title')}</h2>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
          {t('groups:addGroup')}
        </Button>
      </div>

      {groups.length === 0 ? (
        <Empty description="No groups found" />
      ) : (
        <Table
          dataSource={groups}
          columns={columns}
          rowKey="id"
          pagination={{ pageSize: 10 }}
        />
      )}

      <Modal
        title={editingGroup ? t('groups:editGroup') : t('groups:addGroup')}
        open={modalVisible}
        onOk={handleModalOk}
        onCancel={handleModalCancel}
        width={600}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            label={t('groups:fields.name')}
            name="name"
            rules={[{ required: true, message: 'Please input group name' }]}
          >
            <Input />
          </Form.Item>

          <Form.Item
            label={t('groups:fields.description')}
            name="description"
          >
            <Input.TextArea rows={4} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};
