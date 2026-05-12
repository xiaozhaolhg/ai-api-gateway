import { useState } from 'react';
import { Table, Button, Modal, Form, Input, Select, Popconfirm, Empty, message, Tabs } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient, type Group } from '../api/client';
import { GroupMembersTab } from '../components/GroupMembersTab';
import { GroupTierTab } from '../components/GroupTierTab';

export const Groups: React.FC = () => {
  const { t } = useTranslation(['groups', 'common']);
  const queryClient = useQueryClient();
  const [modalVisible, setModalVisible] = useState(false);
  const [editingGroup, setEditingGroup] = useState<Group | null>(null);
  const [form] = Form.useForm();

  const { data: groups = [], isLoading } = useQuery({
    queryKey: ['groups'],
    queryFn: () => apiClient.getGroups(),
  });

  const { data: tiers = [] } = useQuery({
    queryKey: ['tiers'],
    queryFn: () => apiClient.getTiers(),
  });

  const [searchText, setSearchText] = useState('');

  const filteredGroups = groups.filter((g: Group) =>
    g.name.toLowerCase().includes(searchText.toLowerCase())
  );

  const createMutation = useMutation({
    mutationFn: (data: Omit<Group, 'id' | 'created_at' | 'updated_at' | 'member_count'>) => apiClient.createGroup(data),
    onMutate: async (newGroup) => {
      await queryClient.cancelQueries({ queryKey: ['groups'] });
      const previous = queryClient.getQueryData<Group[]>(['groups']);
      queryClient.setQueryData<Group[]>(['groups'], (old = []) => [
        ...old, { ...newGroup, id: `temp-${Date.now()}`, member_count: 0, created_at: new Date().toISOString(), updated_at: new Date().toISOString() } as Group,
      ]);
      return { previous };
    },
    onError: (_err, _new, context) => {
      if (context?.previous) queryClient.setQueryData(['groups'], context.previous);
      message.error('Failed to create group');
    },
    onSuccess: () => {
      message.success('Group created successfully');
      queryClient.invalidateQueries({ queryKey: ['groups'] });
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<Group> }) => apiClient.updateGroup(id, data),
    onMutate: async ({ id, data }) => {
      await queryClient.cancelQueries({ queryKey: ['groups'] });
      const previous = queryClient.getQueryData<Group[]>(['groups']);
      queryClient.setQueryData<Group[]>(['groups'], (old = []) => old.map(g => g.id === id ? { ...g, ...data } : g));
      return { previous };
    },
    onError: (_err, _vars, context) => {
      if (context?.previous) queryClient.setQueryData(['groups'], context.previous);
      message.error('Failed to update group');
    },
    onSuccess: () => {
      message.success('Group updated successfully');
      queryClient.invalidateQueries({ queryKey: ['groups'] });
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => apiClient.deleteGroup(id),
    onMutate: async (id) => {
      await queryClient.cancelQueries({ queryKey: ['groups'] });
      const previous = queryClient.getQueryData<Group[]>(['groups']);
      queryClient.setQueryData<Group[]>(['groups'], (old = []) => old.filter(g => g.id !== id));
      return { previous };
    },
    onError: (_err, _id, context) => {
      if (context?.previous) queryClient.setQueryData(['groups'], context.previous);
      message.error('Failed to delete group');
    },
    onSuccess: () => {
      message.success('Group deleted successfully');
    },
  });

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

  const handleModalOk = async () => {
    const values = await form.validateFields();
    const payload = {
      name: values.name,
      description: values.description,
      parent_group_id: values.parent_group_id || undefined,
      tier_id: values.tier_id || undefined,
    };
    if (editingGroup) {
      updateMutation.mutate({ id: editingGroup.id, data: payload });
    } else {
      createMutation.mutate(payload as any);
    }
    setModalVisible(false);
  };

  const handleModalCancel = () => {
    setModalVisible(false);
    form.resetFields();
  };

  const expandedRowRender = (record: Group) => {
    return (
      <Tabs
        defaultActiveKey="members"
        items={[
          {
            key: 'members',
            label: 'Members',
            children: <GroupMembersTab groupId={record.id} />,
          },
          {
            key: 'tier',
            label: 'Tier',
            children: <GroupTierTab groupId={record.id} currentTierId={record.tier_id} />,
          },
        ]}
      />
    );
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
      title: t('common:common.actions'),
      key: 'actions',
      render: (_: any, record: Group) => (
        <div>
          <Button type="link" icon={<EditOutlined />} onClick={() => handleEdit(record)}>
            {t('common:edit')}
          </Button>
          <Popconfirm
            title="Are you sure you want to delete this group?"
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
        <h2>{t('groups:title')}</h2>
        <div style={{ display: 'flex', gap: 8 }}>
          <Input
            placeholder="Search groups..."
            value={searchText}
            onChange={(e) => setSearchText(e.target.value)}
            style={{ width: 250 }}
            allowClear
          />
          <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
            {t('groups:addGroup')}
          </Button>
        </div>
      </div>

      <Table
        dataSource={filteredGroups}
        columns={columns}
        rowKey="id"
        pagination={{ pageSize: 10 }}
        locale={{ emptyText: <Empty description="No groups found" /> }}
        expandable={{
          expandedRowRender,
          expandRowByClick: false,
        }}
      />

      <Modal
        title={editingGroup ? t('groups:editGroup') : t('groups:addGroup')}
        open={modalVisible}
        onOk={handleModalOk}
        onCancel={handleModalCancel}
        confirmLoading={createMutation.isPending || updateMutation.isPending}
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

          <Form.Item
            label="Parent Group"
            name="parent_group_id"
          >
            <Select
              placeholder="Select parent group (optional)"
              allowClear
            >
              {groups
                .filter((g: Group) => !editingGroup || g.id !== editingGroup.id)
                .map((group: Group) => (
                  <Select.Option key={group.id} value={group.id}>
                    {group.name}
                  </Select.Option>
                ))}
            </Select>
          </Form.Item>

          <Form.Item
            label="Tier"
            name="tier_id"
          >
            <Select
              placeholder="Select tier (optional)"
              allowClear
            >
              {tiers.map((tier) => (
                <Select.Option key={tier.id} value={tier.id}>
                  {tier.name}
                </Select.Option>
              ))}
            </Select>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};
