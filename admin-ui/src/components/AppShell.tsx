import React, { useState } from 'react';
import { Layout, Menu, Breadcrumb, Button, Dropdown, Avatar } from 'antd';
import {
  DashboardOutlined,
  CloudServerOutlined,
  BranchesOutlined,
  UserOutlined,
  TeamOutlined,
  KeyOutlined,
  BarChartOutlined,
  ThunderboltOutlined,
  HeartOutlined,
  BellOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  LogoutOutlined,
  CrownOutlined,
} from '@ant-design/icons';
import { Outlet, useLocation, useNavigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import { LanguageSwitcher } from './LanguageSwitcher';
import { useTranslation } from 'react-i18next';

const { Header, Sider, Content } = Layout;

type Role = 'admin' | 'user' | 'viewer';

const roleAccess: Record<string, Role[]> = {
  '/': ['admin', 'user', 'viewer'],
  '/providers': ['admin'],
  '/routing': ['admin'],
  '/users': ['admin'],
  '/groups': ['admin'],
  '/tiers': ['admin'],
  '/api-keys': ['admin', 'user'],
  '/usage': ['admin', 'user', 'viewer'],
  '/pricing': ['admin'],
  '/health': ['admin', 'user', 'viewer'],
  '/alerts': ['admin'],
};

function filterMenuByRole(
  items: { type: 'group'; label: string; children: { key: string; icon: React.ReactNode; label: string }[] }[],
  role: Role
): { type: 'group'; label: string; children: { key: string; icon: React.ReactNode; label: string }[] }[] {
  return items
    .map(group => ({
      ...group,
      children: group.children.filter(item => {
        const allowed = roleAccess[item.key];
        return allowed ? allowed.includes(role) : true;
      }),
    }))
    .filter(group => group.children.length > 0);
}

interface MenuItem {
  key: string;
  icon: React.ReactNode;
  label: string;
}

interface MenuGroup {
  type: 'group';
  label: string;
  children: MenuItem[];
}

export const AppShell: React.FC = () => {
  const [collapsed, setCollapsed] = useState(false);
  const { user, logout } = useAuth();
  const location = useLocation();
  const navigate = useNavigate();
  const { t } = useTranslation();

  const userRole = (user?.role || 'viewer') as Role;

  const allMenuItems: MenuGroup[] = [
    {
      type: 'group' as const,
      label: 'Infrastructure',
      children: [
        {
          key: '/',
          icon: <DashboardOutlined />,
          label: t('navigation.dashboard'),
        },
        {
          key: '/providers',
          icon: <CloudServerOutlined />,
          label: t('navigation.providers'),
        },
        {
          key: '/routing',
          icon: <BranchesOutlined />,
          label: t('navigation.routingRules'),
        },
      ],
    },
    {
      type: 'group' as const,
      label: 'Access Control',
      children: [
        {
          key: '/users',
          icon: <UserOutlined />,
          label: t('navigation.users'),
        },
        {
          key: '/groups',
          icon: <TeamOutlined />,
          label: t('navigation.groups'),
        },
        {
          key: '/tiers',
          icon: <CrownOutlined />,
          label: t('navigation.tiers'),
        },
        {
          key: '/api-keys',
          icon: <KeyOutlined />,
          label: t('navigation.apiKeys'),
        },
      ],
    },
    {
      type: 'group' as const,
      label: 'Billing',
      children: [
        {
          key: '/usage',
          icon: <BarChartOutlined />,
          label: t('navigation.usage'),
        },
        {
          key: '/pricing',
          icon: <ThunderboltOutlined />,
          label: t('navigation.pricingRules'),
        },
      ],
    },
    {
      type: 'group' as const,
      label: 'Observability',
      children: [
        {
          key: '/health',
          icon: <HeartOutlined />,
          label: t('navigation.health'),
        },
        {
          key: '/alerts',
          icon: <BellOutlined />,
          label: t('navigation.alerts'),
        },
      ],
    },
  ];

  const menuItems = filterMenuByRole(allMenuItems, userRole);

  const handleMenuClick = ({ key }: { key: string }) => {
    navigate(key);
  };

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const userMenuItems = [
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: t('actions.logout'),
      onClick: handleLogout,
    },
  ];

  const getBreadcrumbItems = () => {
    const pathSegments = location.pathname.split('/').filter(Boolean);
    return [
      { title: <a href="/">Home</a> },
      ...pathSegments.map((segment) => ({
        title: segment.charAt(0).toUpperCase() + segment.slice(1),
      })),
    ];
  };

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider
        trigger={null}
        collapsible
        collapsed={collapsed}
        style={{
          overflow: 'auto',
          height: '100vh',
          position: 'fixed',
          left: 0,
          top: 0,
          bottom: 0,
        }}
      >
        <div
          style={{
            height: 32,
            margin: 16,
            background: 'rgba(255, 255, 255, 0.2)',
            borderRadius: 6,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            color: 'white',
            fontWeight: 'bold',
          }}
        >
          {collapsed ? 'AI' : 'AI Gateway'}
        </div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[location.pathname]}
          items={menuItems}
          onClick={handleMenuClick}
        />
      </Sider>
      <Layout style={{ marginLeft: collapsed ? 80 : 200, transition: 'all 0.2s' }}>
        <Header
          style={{
            padding: '0 24px',
            background: '#fff',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            boxShadow: '0 1px 4px rgba(0,21,41,.08)',
          }}
        >
          <Button
            type="text"
            icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
            onClick={() => setCollapsed(!collapsed)}
            style={{ fontSize: '16px', width: 64, height: 64 }}
          />
          <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
            <LanguageSwitcher />
            <Dropdown menu={{ items: userMenuItems }} placement="bottomRight">
              <div style={{ display: 'flex', alignItems: 'center', gap: 8, cursor: 'pointer' }}>
                <Avatar icon={<UserOutlined />} />
                <span>{user?.name || 'Admin'}</span>
              </div>
            </Dropdown>
          </div>
        </Header>
        <Content style={{ margin: '24px 16px 0', overflow: 'initial' }}>
          <Breadcrumb style={{ margin: '16px 0' }} items={getBreadcrumbItems()} />
          <div
            style={{
              padding: 24,
              minHeight: 360,
              background: '#fff',
              borderRadius: 8,
            }}
          >
            <Outlet />
          </div>
        </Content>
      </Layout>
    </Layout>
  );
};
