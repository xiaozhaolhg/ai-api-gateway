import React, { useState } from 'react';
import { Button, Card, Drawer, Space, Tag, Statistic, Row, Col, Upload, message, Modal, InputNumber, Switch } from 'antd';
import { 
  ReloadOutlined, 
  DownloadOutlined, 
  UploadOutlined, 
  SettingOutlined,
  DatabaseOutlined,
  ApiOutlined
} from '@ant-design/icons';
import MockManager from '../mock/MockManager';
import { API_CONFIG, getAPIMode } from '../api/config';

const DevTools: React.FC = () => {
  const [visible, setVisible] = useState(false);
  const [mockDelay, setMockDelay] = useState(API_CONFIG.mockDelay);
  const mockManager = MockManager.getInstance();

  if (import.meta.env.PROD) {
    return null;
  }

  const handleReset = () => {
    Modal.confirm({
      title: 'Reset Mock Data',
      content: 'Are you sure you want to reset all mock data to defaults? This cannot be undone.',
      onOk: () => {
        mockManager.resetToDefaults();
        message.success('Mock data reset to defaults');
      },
    });
  };

  const handleExport = () => {
    mockManager.downloadData();
    message.success('Mock data exported');
  };

  const handleImport = (file: File) => {
    mockManager.importFromFile(file)
      .then(() => message.success('Mock data imported'))
      .catch(() => message.error('Failed to import mock data'));
    return false; // Prevent automatic upload
  };

  const handleMockModeChange = (enabled: boolean) => {
    mockManager.setMockMode(enabled);
    message.success(`Switched to ${enabled ? 'Mock' : 'Real'} API mode`);
    // Reload to apply changes
    setTimeout(() => window.location.reload(), 500);
  };

  const handleDelayChange = (value: number | null) => {
    if (value !== null) {
      setMockDelay(value);
      // Update config (this would need to be stored in localStorage for persistence)
      localStorage.setItem('mock_delay', value.toString());
      message.success(`Mock delay set to ${value}ms`);
    }
  };

  const dataStats = mockManager.getDataStats();
  const currentMode = getAPIMode();

  return (
    <>
      <Button
        type="text"
        icon={<SettingOutlined />}
        onClick={() => setVisible(true)}
        style={{
          position: 'fixed',
          bottom: 20,
          right: 20,
          zIndex: 1000,
          background: 'rgba(0, 0, 0, 0.7)',
          color: 'white',
          borderRadius: '50%',
          width: 50,
          height: 50,
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
        }}
      />
      
      <Drawer
        title="Mock API Development Tools"
        placement="right"
        onClose={() => setVisible(false)}
        open={visible}
        width={400}
      >
        <Space direction="vertical" style={{ width: '100%' }} size="large">
          {/* API Mode Status */}
          <Card title="API Mode" size="small">
            <Space direction="vertical" style={{ width: '100%' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <span>Current Mode:</span>
                <Tag color={currentMode === 'Mock' ? 'green' : 'blue'}>
                  {currentMode}
                </Tag>
              </div>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <span>Switch Mode:</span>
                <Switch
                  checked={currentMode === 'Mock'}
                  onChange={handleMockModeChange}
                  checkedChildren="Mock"
                  unCheckedChildren="Real"
                />
              </div>
            </Space>
          </Card>

          {/* Mock Settings */}
          <Card title="Mock Settings" size="small">
            <Space direction="vertical" style={{ width: '100%' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <span>Network Delay:</span>
                <InputNumber
                  value={mockDelay}
                  onChange={handleDelayChange}
                  min={0}
                  max={5000}
                  step={100}
                  suffix="ms"
                  style={{ width: 120 }}
                />
              </div>
            </Space>
          </Card>

          {/* Data Management */}
          <Card title="Data Management" size="small">
            <Space direction="vertical" style={{ width: '100%' }}>
              <Button 
                icon={<ReloadOutlined />} 
                onClick={handleReset}
                block
                danger
              >
                Reset to Defaults
              </Button>
              <Button 
                icon={<DownloadOutlined />} 
                onClick={handleExport}
                block
              >
                Export Mock Data
              </Button>
              <Upload
                accept=".json"
                showUploadList={false}
                beforeUpload={handleImport}
              >
                <Button icon={<UploadOutlined />} block>
                  Import Mock Data
                </Button>
              </Upload>
            </Space>
          </Card>

          {/* Data Statistics */}
          <Card 
            title={<span><DatabaseOutlined /> Data Statistics</span>} 
            size="small"
          >
            <Row gutter={[16, 16]}>
              <Col span={12}>
                <Statistic title="Users" value={dataStats.users} />
              </Col>
              <Col span={12}>
                <Statistic title="Providers" value={dataStats.providers} />
              </Col>
              <Col span={12}>
                <Statistic title="API Keys" value={dataStats.apiKeys} />
              </Col>
              <Col span={12}>
                <Statistic title="Usage Records" value={dataStats.usage} />
              </Col>
              <Col span={12}>
                <Statistic title="Routing Rules" value={dataStats.routingRules} />
              </Col>
              <Col span={12}>
                <Statistic title="Groups" value={dataStats.groups} />
              </Col>
              <Col span={12}>
                <Statistic title="Permissions" value={dataStats.permissions} />
              </Col>
              <Col span={12}>
                <Statistic title="Budgets" value={dataStats.budgets} />
              </Col>
              <Col span={12}>
                <Statistic title="Pricing Rules" value={dataStats.pricingRules} />
              </Col>
              <Col span={12}>
                <Statistic title="Alert Rules" value={dataStats.alertRules} />
              </Col>
              <Col span={12}>
                <Statistic title="Active Alerts" value={dataStats.alerts} />
              </Col>
            </Row>
          </Card>

          {/* API Configuration */}
          <Card 
            title={<span><ApiOutlined /> API Configuration</span>} 
            size="small"
          >
            <Space direction="vertical" style={{ width: '100%' }}>
              <div>
                <div style={{ marginBottom: 4, color: '#666' }}>Base URL:</div>
                <div style={{ fontFamily: 'monospace', fontSize: 12 }}>
                  {API_CONFIG.baseURL}
                </div>
              </div>
              <div>
                <div style={{ marginBottom: 4, color: '#666' }}>Mock Enabled:</div>
                <Tag color={API_CONFIG.useMock ? 'green' : 'red'}>
                  {API_CONFIG.useMock ? 'Yes' : 'No'}
                </Tag>
              </div>
            </Space>
          </Card>
        </Space>
      </Drawer>
    </>
  );
};

export default DevTools;