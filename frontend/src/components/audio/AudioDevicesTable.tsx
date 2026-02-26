import { useState } from 'react';
import { 
  Table, Button, Input, Space, Tag, Tooltip, Switch, Card
} from 'antd';
import { 
  SoundOutlined, EditOutlined, ReloadOutlined, EyeInvisibleOutlined, EyeOutlined
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import { 
  SetAudioDeviceNickname, SetPrimaryOutputDevice, 
  IgnoreAudioDevice, UnignoreAudioDevice
} from "../../../wailsjs/go/main/App";

interface AudioDevice {
  id: string;
  name: string;
  isDefault: boolean;
  isEnabled: boolean;
  deviceType: string;
  state: string;
  selected: boolean;
  nickname: string;
}

interface AudioDevicesTableProps {
  audioDevices: {filtered: AudioDevice[], ignored: AudioDevice[]};
  showIgnoredAudio: boolean;
  setShowIgnoredAudio: (show: boolean) => void;
  loading: boolean;
  onRefresh: () => void;
}

export function AudioDevicesTable({ 
  audioDevices, 
  showIgnoredAudio, 
  setShowIgnoredAudio,
  loading, 
  onRefresh 
}: AudioDevicesTableProps) {
  const [editingAudioDevice, setEditingAudioDevice] = useState<string | null>(null);
  const [tempNickname, setTempNickname] = useState<string>('');

  const startEditingAudioDeviceNickname = (deviceId: string, currentNickname: string) => {
    setEditingAudioDevice(deviceId);
    setTempNickname(currentNickname);
  };

  const saveAudioDeviceNickname = async (deviceId: string) => {
    try {
      await SetAudioDeviceNickname(deviceId, tempNickname);
      window.location.reload();
    } catch (error) {
      console.error('Error saving audio device nickname:', error);
    }
  };

  const cancelEditing = () => {
    setEditingAudioDevice(null);
    setTempNickname('');
  };

  const handleSetDefaultAudioDevice = async (deviceId: string) => {
    try {
      await SetPrimaryOutputDevice(deviceId);
      window.location.reload();
    } catch (error) {
      console.error('Error setting default audio device:', error);
    }
  };

  const handleIgnoreAudioDevice = async (deviceId: string) => {
    try {
      await IgnoreAudioDevice(deviceId);
      window.location.reload();
    } catch (error) {
      console.error('Error ignoring audio device:', error);
    }
  };

  const handleUnignoreAudioDevice = async (deviceId: string) => {
    try {
      await UnignoreAudioDevice(deviceId);
      window.location.reload();
    } catch (error) {
      console.error('Error unignoring audio device:', error);
    }
  };

  const audioColumns: ColumnsType<AudioDevice> = [
    {
      title: 'Actions',
      key: 'actions',
      width: 100,
      render: (_, record: AudioDevice) => (
        <Button
          size="small"
          type={record.isDefault ? "primary" : "default"}
          onClick={() => handleSetDefaultAudioDevice(record.id)}
        >
          {record.isDefault ? 'Default' : 'Set Default'}
        </Button>
      )
    },
    {
      title: 'Name',
      key: 'name',
      width: 350,
      render: (_, record: AudioDevice) => {
        const isEditing = editingAudioDevice === record.id;
        const displayName = record.nickname || record.name;
        const hasNickname = !!record.nickname;
        
        if (isEditing) {
          return (
            <Space>
              <Input
                value={tempNickname}
                onChange={(e) => setTempNickname(e.target.value)}
                placeholder="Enter nickname"
                style={{ width: 150 }}
              />
              <Button 
                size="small" 
                type="primary" 
                onClick={() => saveAudioDeviceNickname(record.id)}
              >
                Save
              </Button>
              <Button 
                size="small" 
                onClick={cancelEditing}
              >
                Cancel
              </Button>
            </Space>
          );
        }
        
        return (
          <Space>
            <Tooltip title={hasNickname ? `Original: ${record.name}` : ''}>
              <span>{displayName}</span>
            </Tooltip>
            <Button 
              size="small" 
              type="text" 
              icon={<EditOutlined />}
              onClick={() => startEditingAudioDeviceNickname(record.id, record.nickname || '')}
            />
          </Space>
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
          onClick={onRefresh}
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
  );
}
