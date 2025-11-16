import { useState, useEffect } from 'react';
import { 
  Layout, Typography, Button, Table, Input, Select, Card, Row, Col, 
  message as antMessage, Space, Alert, Tag, Divider, Switch, Tooltip, Checkbox
} from 'antd';
import { 
  ReloadOutlined, SaveOutlined, PlayCircleOutlined, 
  DesktopOutlined, CheckCircleOutlined, CloseCircleOutlined,
  SoundOutlined, EyeInvisibleOutlined, EyeOutlined, StopOutlined
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import './App.css';
import { 
  GetMonitors, GetProfiles, SaveProfile, ApplyProfile, RefreshMonitors,
  GetAudioDevicesWithIgnoreStatus, RefreshAudioDevices, IgnoreAudioDevice, UnignoreAudioDevice,
  SetAudioDeviceSelection, GetSelectedAudioDevices
} from "../wailsjs/go/main/App";

const { Header, Content } = Layout;
const { Title } = Typography;

interface Monitor {
  deviceName: string;
  displayName: string;
  isPrimary: boolean;
  isActive: boolean;
  bounds: {
    x: number;
    y: number;
    width: number;
    height: number;
  };
}

interface AudioDevice {
  id: string;
  name: string;
  isDefault: boolean;
  isEnabled: boolean;
  deviceType: string; // "output" or "input"
  state: string;      // "active", "disabled", "notpresent", "unplugged"
  selected: boolean;  // whether this device is selected for the profile
}

interface Profile {
  name: string;
  monitors: Monitor[];
  audioDevices: AudioDevice[];
}

function App() {
  const [monitors, setMonitors] = useState<Monitor[]>([]);
  const [audioDevices, setAudioDevices] = useState<{filtered: AudioDevice[], ignored: AudioDevice[]}>({filtered: [], ignored: []});
  const [showIgnoredAudio, setShowIgnoredAudio] = useState<boolean>(false);
  const [profiles, setProfiles] = useState<Profile[]>([]);
  const [selectedProfile, setSelectedProfile] = useState<string>('');
  const [profileName, setProfileName] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(false);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      setLoading(true);
      const [monitorsData, audioData, profilesData] = await Promise.all([
        GetMonitors(),
        GetAudioDevicesWithIgnoreStatus(),
        GetProfiles()
      ]);
      setMonitors(monitorsData);
      setAudioDevices(audioData as {filtered: AudioDevice[], ignored: AudioDevice[]});
      setProfiles(profilesData);
    } catch (error) {
      antMessage.error(`Error loading data: ${error}`);
    } finally {
      setLoading(false);
    }
  };

  const handleRefreshMonitors = async () => {
    try {
      setLoading(true);
      const monitorsData = await RefreshMonitors();
      setMonitors(monitorsData);
      antMessage.success('Monitors refreshed successfully');
    } catch (error) {
      antMessage.error(`Error refreshing monitors: ${error}`);
    } finally {
      setLoading(false);
    }
  };

  const handleRefreshAudio = async () => {
    try {
      setLoading(true);
      const audioData = await RefreshAudioDevices();
      setAudioDevices(audioData as {filtered: AudioDevice[], ignored: AudioDevice[]});
      antMessage.success('Audio devices refreshed successfully');
    } catch (error) {
      antMessage.error(`Error refreshing audio devices: ${error}`);
    } finally {
      setLoading(false);
    }
  };

  const handleSaveProfile = async () => {
    if (!profileName.trim()) {
      antMessage.warning('Please enter a profile name');
      return;
    }

    try {
      setLoading(true);
      await SaveProfile(profileName);
      const profilesData = await GetProfiles();
      setProfiles(profilesData);
      setProfileName('');
      antMessage.success('Profile saved successfully');
    } catch (error) {
      antMessage.error(`Error saving profile: ${error}`);
    } finally {
      setLoading(false);
    }
  };

  const handleApplyProfile = async () => {
    if (!selectedProfile) {
      antMessage.warning('Please select a profile to apply');
      return;
    }

    try {
      setLoading(true);
      await ApplyProfile(selectedProfile);
      antMessage.success(`Profile '${selectedProfile}' applied successfully`);
      await loadData();
    } catch (error) {
      antMessage.error(`Error applying profile: ${error}`);
    } finally {
      setLoading(false);
    }
  };

  const handleIgnoreAudioDevice = async (deviceId: string) => {
    try {
      await IgnoreAudioDevice(deviceId);
      const audioData = await GetAudioDevicesWithIgnoreStatus();
      setAudioDevices(audioData as {filtered: AudioDevice[], ignored: AudioDevice[]});
      antMessage.success('Audio device added to ignore list');
    } catch (error) {
      antMessage.error(`Error ignoring audio device: ${error}`);
    }
  };

  const handleUnignoreAudioDevice = async (deviceId: string) => {
    try {
      await UnignoreAudioDevice(deviceId);
      const audioData = await GetAudioDevicesWithIgnoreStatus();
      setAudioDevices(audioData as {filtered: AudioDevice[], ignored: AudioDevice[]});
      antMessage.success('Audio device removed from ignore list');
    } catch (error) {
      antMessage.error(`Error unignoring audio device: ${error}`);
    }
  };

  const monitorColumns: ColumnsType<Monitor> = [
    {
      title: 'Active',
      dataIndex: 'isActive',
      key: 'isActive',
      width: 100,
      render: (isActive: boolean) => (
        <Tag color={isActive ? 'success' : 'error'} icon={isActive ? <CheckCircleOutlined /> : <CloseCircleOutlined />}>
          {isActive ? 'Active' : 'Inactive'}
        </Tag>
      )
    },
    {
      title: 'Display Name',
      dataIndex: 'displayName',
      key: 'displayName',
      width: 200
    },
    {
      title: 'Device Name',
      dataIndex: 'deviceName',
      key: 'deviceName',
      width: 300
    },
    {
      title: 'Primary',
      dataIndex: 'isPrimary',
      key: 'isPrimary',
      width: 100,
      render: (isPrimary: boolean) => (
        isPrimary ? (
          <Tag color="blue" icon={<DesktopOutlined />}>Primary</Tag>
        ) : (
          <Tag>Secondary</Tag>
        )
      )
    },
    {
      title: 'Resolution',
      dataIndex: 'bounds',
      key: 'resolution',
      width: 120,
      render: (bounds: Monitor['bounds']) => (
        <Tag color="default">{bounds.width}x{bounds.height}</Tag>
      )
    },
    {
      title: 'Position',
      dataIndex: 'bounds',
      key: 'position',
      width: 120,
      render: (bounds: Monitor['bounds']) => (
        <Tag color="default">({bounds.x}, {bounds.y})</Tag>
      )
    }
  ];

  const audioColumns: ColumnsType<AudioDevice> = [
    {
      title: 'Select',
      dataIndex: 'selected',
      key: 'selected',
      width: 80,
      render: (selected: boolean, record: AudioDevice) => (
        <Checkbox
          checked={selected}
          onChange={async (e) => {
            try {
              await SetAudioDeviceSelection(record.id, e.target.checked);
              // Update local state to reflect the change
              if (!showIgnoredAudio) {
                const updatedFiltered = audioDevices.filtered.map(device => 
                  device.id === record.id ? { ...device, selected: e.target.checked } : device
                );
                setAudioDevices({ ...audioDevices, filtered: updatedFiltered });
              }
            } catch (error) {
              antMessage.error(`Error updating device selection: ${error}`);
            }
          }}
        />
      )
    },
    {
      title: 'Name',
      dataIndex: 'name',
      key: 'name',
      width: 250
    },
    {
      title: 'Type',
      dataIndex: 'deviceType',
      key: 'deviceType',
      width: 100,
      render: (deviceType: string) => (
        <Tag color={deviceType === 'output' ? 'blue' : 'green'}>
          {deviceType === 'output' ? 'Output' : 'Input'}
        </Tag>
      )
    },
    {
      title: 'State',
      dataIndex: 'state',
      key: 'state',
      width: 120,
      render: (state: string, record: AudioDevice) => {
        const colorMap: {[key: string]: string} = {
          'active': 'success',
          'disabled': 'error',
          'notpresent': 'default',
          'unplugged': 'warning'
        };
        return (
          <Tag color={colorMap[state] || 'default'} icon={!record.isEnabled ? <StopOutlined /> : undefined}>
            {state}
          </Tag>
        );
      }
    },
    {
      title: 'Default',
      dataIndex: 'isDefault',
      key: 'isDefault',
      width: 100,
      render: (isDefault: boolean) => (
        isDefault ? (
          <Tag color="gold" icon={<SoundOutlined />}>Default</Tag>
        ) : null
      )
    },
    {
      title: 'Actions',
      key: 'actions',
      width: 120,
      render: (_, record: AudioDevice) => (
        <Space>
          <Tooltip title={showIgnoredAudio ? "Remove from ignore list" : "Add to ignore list"}>
            <Button
              size="small"
              type={showIgnoredAudio ? "primary" : "default"}
              danger={showIgnoredAudio}
              icon={showIgnoredAudio ? <EyeOutlined /> : <EyeInvisibleOutlined />}
              onClick={() => showIgnoredAudio ? handleUnignoreAudioDevice(record.id) : handleIgnoreAudioDevice(record.id)}
            />
          </Tooltip>
        </Space>
      )
    }
  ];

  const allAudioDevices = showIgnoredAudio ? audioDevices.ignored : audioDevices.filtered;

  return (
    <Layout style={{ minHeight: '100vh', background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' }}>
      <Header style={{ 
        background: 'rgba(255, 255, 255, 0.95)', 
        padding: '0 24px',
        boxShadow: '0 2px 8px rgba(0, 0, 0, 0.1)',
        backdropFilter: 'blur(10px)'
      }}>
        <Title level={2} style={{ margin: '16px 0', color: '#1a1a1a' }}>
          <DesktopOutlined /> Monitor Profile Manager
        </Title>
      </Header>

      <Content style={{ padding: '24px', maxWidth: '1400px', margin: '0 auto' }}>
        <Row gutter={[24, 24]}>
          <Col xs={24} xl={16}>
            <Space direction="vertical" style={{ width: '100%' }} size="large">
              {/* Monitors Section */}
              <Card 
                title={
                  <Space>
                    <DesktopOutlined />
                    <span>Detected Monitors</span>
                  </Space>
                }
                extra={
                  <Button 
                    type="primary" 
                    icon={<ReloadOutlined />}
                    onClick={handleRefreshMonitors}
                    loading={loading}
                  >
                    Refresh Monitors
                  </Button>
                }
              >
                <Table
                  columns={monitorColumns}
                  dataSource={monitors}
                  rowKey="deviceName"
                  loading={loading}
                  pagination={false}
                  scroll={{ x: 1000 }}
                  size="middle"
                />
              </Card>

              {/* Audio Devices Section */}
              <Card 
                title={
                  <Space>
                    <SoundOutlined />
                    <span>Audio Devices</span>
                    {!showIgnoredAudio && (
                      <span style={{ fontSize: '12px', color: '#1890ff' }}>
                        ({audioDevices.filtered.filter(d => d.selected).length} selected)
                      </span>
                    )}
                    <Switch
                      checkedChildren={<EyeOutlined />}
                      unCheckedChildren={<EyeInvisibleOutlined />}
                      checked={showIgnoredAudio}
                      onChange={setShowIgnoredAudio}
                      size="small"
                    />
                    <span style={{ fontSize: '12px', color: '#666' }}>
                      {showIgnoredAudio ? 'Ignored' : 'Active'}
                    </span>
                  </Space>
                }
                extra={
                  <Button 
                    type="primary" 
                    icon={<ReloadOutlined />}
                    onClick={handleRefreshAudio}
                    loading={loading}
                  >
                    Refresh Audio
                  </Button>
                }
              >
                <Table
                  columns={audioColumns}
                  dataSource={allAudioDevices}
                  rowKey="id"
                  loading={loading}
                  pagination={false}
                  scroll={{ x: 800 }}
                  size="middle"
                />
              </Card>
            </Space>
          </Col>

          <Col xs={24} xl={8}>
            <Card 
              title={
                <Space>
                  <SaveOutlined />
                  <span>Profiles</span>
                </Space>
              }
              style={{ height: '100%' }}
            >
              <Space direction="vertical" style={{ width: '100%' }} size="large">
                <div>
                  <Title level={5}>Save Current Profile</Title>
                  <Input.Group compact>
                    <Input
                      style={{ width: 'calc(100% - 100px)' }}
                      placeholder="Enter profile name"
                      value={profileName}
                      onChange={(e) => setProfileName(e.target.value)}
                      onPressEnter={handleSaveProfile}
                    />
                    <Button 
                      type="primary" 
                      icon={<SaveOutlined />}
                      onClick={handleSaveProfile}
                      loading={loading}
                      style={{ width: '100px' }}
                    >
                      Save
                    </Button>
                  </Input.Group>
                </div>

                <Divider />

                <div>
                  <Title level={5}>Apply Profile</Title>
                  <Space direction="vertical" style={{ width: '100%' }}>
                    <Select
                      style={{ width: '100%' }}
                      placeholder="Select a profile"
                      value={selectedProfile || undefined}
                      onChange={setSelectedProfile}
                      loading={loading}
                    >
                      {profiles.map((profile) => (
                        <Select.Option key={profile.name} value={profile.name}>
                          {profile.name}
                        </Select.Option>
                      ))}
                    </Select>
                    <Button 
                      type="primary" 
                      icon={<PlayCircleOutlined />}
                      onClick={handleApplyProfile}
                      loading={loading}
                      disabled={!selectedProfile}
                      style={{ width: '100%' }}
                    >
                      Apply Selected Profile
                    </Button>
                  </Space>
                </div>

                {profiles.length > 0 && (
                  <>
                    <Divider />
                    <div>
                      <Title level={5}>Saved Profiles ({profiles.length})</Title>
                      <Space direction="vertical" style={{ width: '100%' }} size="small">
                        {profiles.map((profile) => (
                          <Card key={profile.name} size="small" style={{ backgroundColor: '#fafafa' }}>
                            <Space direction="vertical" size="small" style={{ width: '100%' }}>
                              <strong>{profile.name}</strong>
                              <div>
                                <Tag color="blue">{profile.monitors.length} monitor(s)</Tag>
                                <Tag color="green">{profile.audioDevices.length} audio device(s)</Tag>
                                <Tag color="orange">{profile.monitors.filter(m => m.isActive).length} active</Tag>
                                <Tag color="purple">{profile.audioDevices.filter(d => d.isDefault).length} default</Tag>
                              </div>
                            </Space>
                          </Card>
                        ))}
                      </Space>
                    </div>
                  </>
                )}
              </Space>
            </Card>
          </Col>
        </Row>
      </Content>
    </Layout>
  );
}

export default App;
