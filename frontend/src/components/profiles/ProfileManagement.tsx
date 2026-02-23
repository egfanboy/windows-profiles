import { useState } from 'react';
import { 
  Card, Typography, Button, Input, Select, Space, Divider, Tag
} from 'antd';
import { 
  SaveOutlined, PlayCircleOutlined, EditOutlined,
  DeleteOutlined
} from '@ant-design/icons';
import { ConfirmProfileDelete } from './ConfirmProfileDelete';
import { 
  SaveProfile, ApplyProfile, GetProfiles, DeleteProfile
} from "../../../wailsjs/go/main/App";

const { Title } = Typography;

interface Monitor {
  deviceName: string;
  displayName: string;
  isPrimary: boolean;
  isActive: boolean;
  isEnabled: boolean;
  nickname: string;
}

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

interface Profile {
  name: string;
  monitors?: Monitor[];
  audioDevices?: AudioDevice[];
}

interface ProfileManagementProps {
  profiles: Profile[];
  loading: boolean;
  onProfilesChange: () => void;
}

export function ProfileManagement({ profiles, loading, onProfilesChange }: ProfileManagementProps) {
  const [selectedProfile, setSelectedProfile] = useState<string>('');
  const [profileName, setProfileName] = useState<string>('');
  const [editingProfile, setEditingProfile] = useState<string | null>(null);
  const [deleteModalVisible, setDeleteModalVisible] = useState<boolean>(false);
  const [profileToDelete, setProfileToDelete] = useState<string>('');

  const handleSaveProfile = async () => {
    if (!profileName.trim()) {
      return;
    }

    try {
      await SaveProfile(profileName);
      setProfileName('');
      setEditingProfile(null);
      onProfilesChange();
    } catch (error) {
      console.error('Error saving profile:', error);
    }
  };

  const startEditingProfile = (profileName: string) => {
    setEditingProfile(profileName);
    setProfileName(profileName);
  };

  const cancelEditingProfile = () => {
    setEditingProfile(null);
    setProfileName('');
  };

  const handleApplyProfile = async () => {
    if (!selectedProfile) {
      return;
    }

    try {
      await ApplyProfile(selectedProfile);
      window.location.reload();
    } catch (error) {
      console.error('Error applying profile:', error);
    }
  };

  const handleApplyProfileByName = async (profileName: string) => {
    try {
      await ApplyProfile(profileName);
      window.location.reload();
    } catch (error) {
      console.error('Error applying profile:', error);
    }
  };

  const showDeleteConfirm = (profileName: string) => {
    setProfileToDelete(profileName);
    setDeleteModalVisible(true);
  };

  const handleDeleteProfile = async () => {
    if (!profileToDelete) {
      return;
    }

    try {
      await DeleteProfile(profileToDelete);
      setDeleteModalVisible(false);
      setProfileToDelete('');
      onProfilesChange();
    } catch (error) {
      console.error('Error deleting profile:', error);
    }
  };

  const cancelDeleteProfile = () => {
    setDeleteModalVisible(false);
    setProfileToDelete('');
  };

  return (
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
              style={{ width: 'calc(100% - 170px)' }}
              placeholder={editingProfile ? "Edit profile name" : "Enter profile name"}
              value={profileName}
              onChange={(e) => setProfileName(e.target.value)}
              onPressEnter={handleSaveProfile}
            />
            <Button 
              type="primary" 
              icon={<SaveOutlined />}
              onClick={handleSaveProfile}
              loading={loading}
              style={{ width: '70px' }}
            >
              {editingProfile ? 'Update' : 'Save'}
            </Button>
            {editingProfile && (
              <Button 
                onClick={cancelEditingProfile}
                style={{ width: '70px' }}
              >
                Cancel
              </Button>
            )}
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
                  <Card 
                    key={profile.name} 
                    size="small" 
                    style={{ backgroundColor: '#fafafa' }}
                    actions={[
                      <Button 
                        key="apply"
                        type="primary" 
                        size="small"
                        icon={<PlayCircleOutlined />}
                        onClick={() => handleApplyProfileByName(profile.name)}
                        loading={loading}
                      >
                        Apply
                      </Button>
                    ]}
                    extra={
                      <>
                        <Button 
                        size="small" 
                        type="text" 
                        icon={<EditOutlined />}
                        onClick={() => startEditingProfile(profile.name)}
                      />
                      <Button 
                        size="small" 
                        type="text" 
                        danger
                        icon={<DeleteOutlined />}
                        onClick={() => showDeleteConfirm(profile.name)}
                      />
                         </>
                    }
                  >
                    <Space direction="vertical" size="small" style={{ width: '100%' }}>
                      <strong>{profile.name}</strong>
                      <div>
                        <Tag color="blue">{(profile.monitors || []).length} monitor(s)</Tag>
                        <Tag color="green">{(profile.audioDevices || []).length} audio device(s)</Tag>
                        <Tag color="orange">{(profile.monitors || []).filter(m => m.isActive).length} active</Tag>
                        <Tag color="purple">{(profile.audioDevices || []).filter(d => d.isDefault).length} default</Tag>
                      </div>
                    </Space>
                  </Card>
                ))}
              </Space>
            </div>
          </>
        )}
      </Space>
      <ConfirmProfileDelete
        visible={deleteModalVisible}
        profileName={profileToDelete}
        onConfirm={handleDeleteProfile}
        onCancel={cancelDeleteProfile}
        loading={loading}
      />
    </Card>
  );
}
