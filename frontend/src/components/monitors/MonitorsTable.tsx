import { useState } from 'react';
import { 
  Table, Button, Input, Space, Tag, Tooltip, Switch, Radio, Card, Dropdown, Menu
} from 'antd';
import { 
  DesktopOutlined, EditOutlined, ReloadOutlined, CheckCircleOutlined, CloseCircleOutlined, EllipsisOutlined, StarOutlined
} from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import { 
  SetMonitorNickname, SetMonitorPrimary, SetMonitorEnabledState
} from "../../../wailsjs/go/main/App";

interface Monitor {
  deviceName: string;
  displayName: string;
  isPrimary: boolean;
  isActive: boolean;
  isEnabled: boolean;
  nickname: string;
  monitorId: string;
}

interface MonitorsTableProps {
  monitors: Monitor[];
  loading: boolean;
  onRefresh: () => void;
}

export function MonitorsTable({ monitors, loading, onRefresh }: MonitorsTableProps) {
  const [editingMonitor, setEditingMonitor] = useState<string | null>(null);
  const [tempNickname, setTempNickname] = useState<string>('');

  // Define actions for the dropdown menu
  const getMonitorActions = (record: Monitor) => [
    {
      key: 'toggle-enabled',
      label: record.isActive ? 'Disable Monitor' : 'Enable Monitor',
      onClick: () => handleSetMonitorEnabled(record.monitorId, !record.isActive),
      disabled: record.isPrimary && !record.isEnabled
    },
    {
      key: 'set-primary',
      label: 'Set as Primary',
      onClick: () => handleSetMonitorPrimary(record.monitorId),
      disabled: record.isPrimary || !record.isEnabled || !record.isActive
    }
  ];

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

  const handleSetMonitorEnabled = async (monitorId: string, enabled: boolean) => {
    try {
      await SetMonitorEnabledState(monitorId, enabled);
      window.location.reload();
    } catch (error) {
      console.error('Error setting monitor enabled:', error);
    }
  };

  const monitorColumns: ColumnsType<Monitor> = [
    {
      title: 'State',
      dataIndex: 'isActive',
      key: 'state',
      width: 50,
      render: (_, record: Monitor) => (
        <div style={{ display: 'flex', flexDirection: 'column', gap: 4 }}>
          {record.isPrimary && (
            <Tag color="purple" icon={<StarOutlined />}>
              Primary
            </Tag>
          )}
          <Tag color={record.isActive ? 'success' : 'error'} icon={record.isActive ? <CheckCircleOutlined /> : <CloseCircleOutlined />}>
            {record.isActive ? 'Active' : 'Inactive'}
          </Tag>
        </div>
      )
    },
    {
      title: 'Display Name',
      dataIndex: 'displayName',
      key: 'displayName',
      width: 150
    },
    {
      title: 'Device Name',
      key: 'deviceName',
      width: 200,
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
      title: 'Actions',
      key: 'actions',
      width: 50,
      render: (_, record: Monitor) => {
        const actions = getMonitorActions(record);
        const menuItems = actions.map(action => ({
          key: action.key,
          label: action.label,
          onClick: action.onClick,
          disabled: action.disabled
        }));

        return (
          <Dropdown
            menu={{ items: menuItems }}
            trigger={['click']}
            placement="bottomRight"
          >
            <Button
              type="text"
              icon={<EllipsisOutlined />}
              onClick={(e) => e.preventDefault()}
            />
          </Dropdown>
        );
      }
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
        size="middle"
      />
    </Card>
  );
}
