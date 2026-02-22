import { useState } from 'react';
import { 
  Table, Button, Input, Space, Tag, Tooltip, Switch, Radio, Card
} from 'antd';
import { 
  DesktopOutlined, EditOutlined, ReloadOutlined, CheckCircleOutlined, CloseCircleOutlined
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import { 
  SetMonitorNickname, SetMonitorPrimary, SetMonitorEnabled
} from "../../../wailsjs/go/main/App";

interface Monitor {
  deviceName: string;
  displayName: string;
  isPrimary: boolean;
  isActive: boolean;
  isEnabled: boolean;
  nickname: string;
}

interface MonitorsTableProps {
  monitors: Monitor[];
  loading: boolean;
  onRefresh: () => void;
}

export function MonitorsTable({ monitors, loading, onRefresh }: MonitorsTableProps) {
  const [editingMonitor, setEditingMonitor] = useState<string | null>(null);
  const [tempNickname, setTempNickname] = useState<string>('');

  const startEditingMonitorNickname = (deviceName: string, currentNickname: string) => {
    setEditingMonitor(deviceName);
    setTempNickname(currentNickname);
  };

  const saveMonitorNickname = async (deviceName: string) => {
    try {
      await SetMonitorNickname(deviceName, tempNickname);
      window.location.reload(); // Simple refresh for now
    } catch (error) {
      console.error('Error saving monitor nickname:', error);
    }
  };

  const cancelEditing = () => {
    setEditingMonitor(null);
    setTempNickname('');
  };

  const handleSetMonitorPrimary = async (deviceName: string) => {
    try {
      await SetMonitorPrimary(deviceName);
      window.location.reload();
    } catch (error) {
      console.error('Error setting primary monitor:', error);
    }
  };

  const handleSetMonitorEnabled = async (deviceName: string, enabled: boolean) => {
    try {
      await SetMonitorEnabled(deviceName, enabled);
      window.location.reload();
    } catch (error) {
      console.error('Error setting monitor enabled:', error);
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
      key: 'deviceName',
      width: 350,
      render: (_, record: Monitor) => {
        const isEditing = editingMonitor === record.deviceName;
        const displayName = record.nickname || record.deviceName;
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
                onClick={() => saveMonitorNickname(record.deviceName)}
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
            <Tooltip title={hasNickname ? `Original: ${record.deviceName}` : ''}>
              <span>{displayName}</span>
            </Tooltip>
            <Button 
              size="small" 
              type="text" 
              icon={<EditOutlined />}
              onClick={() => startEditingMonitorNickname(record.deviceName, record.nickname || '')}
            />
          </Space>
        );
      }
    },
    {
      title: 'Primary',
      key: 'isPrimary',
      width: 100,
      render: (_, record: Monitor) => (
        <Radio
          checked={record.isPrimary}
          onChange={() => handleSetMonitorPrimary(record.deviceName)}
        />
      )
    },
    {
      title: 'Enabled',
      key: 'isEnabled',
      width: 100,
      render: (_, record: Monitor) => (
        <Switch
          checked={record.isEnabled}
          onChange={(enabled) => handleSetMonitorEnabled(record.deviceName, enabled)}
          disabled={record.isPrimary && !record.isEnabled}
        />
      )
    },
  ];

  return (
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
          onClick={onRefresh}
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
  );
}
