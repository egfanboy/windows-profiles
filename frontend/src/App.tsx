import { useState, useEffect } from 'react';
import { Layout, Typography, message as antMessage, Space, Alert, Row, Col, Card, Button } from 'antd';
import { DesktopOutlined } from '@ant-design/icons';
import './App.css';
import { MonitorsTable } from './components/monitors/MonitorsTable';
import { AudioDevicesTable } from './components/audio/AudioDevicesTable';
import { ProfileManagement } from './components/profiles/ProfileManagement';
import { 
  GetMonitors, GetProfiles, RefreshMonitors,
  GetAudioDevicesWithIgnoreStatus, RefreshAudioDevices
} from "../wailsjs/go/main/App";

const { Header, Content } = Layout;
const { Title } = Typography;

interface Monitor {
  deviceName: string;
  displayName: string;
  isPrimary: boolean;
  monitorId: string;
  isActive: boolean;
  isEnabled: boolean; // user-controlled enable/disable state
  nickname: string;
}

interface AudioDevice {
  id: string;
  name: string;
  isDefault: boolean;
  isEnabled: boolean;
  deviceType: string; // "output" or "input"
  state: string;      // "active", "disabled", "notpresent", "unplugged"
  selected: boolean;  // whether this device is selected for the profile
  nickname: string;   // optional custom nickname
}

interface Profile {
  name: string;
  monitors?: Monitor[];
  audioDevices?: AudioDevice[];
}

function App() {
  const [monitors, setMonitors] = useState<Monitor[]>([]);
  const [audioDevices, setAudioDevices] = useState<{filtered: AudioDevice[], ignored: AudioDevice[]}>({filtered: [], ignored: []});
  const [showIgnoredAudio, setShowIgnoredAudio] = useState<boolean>(false);
  const [profiles, setProfiles] = useState<Profile[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  // Error boundary catch
  if (error) {
    return (
      <Layout style={{ minHeight: '100vh', background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' }}>
        <Content style={{ padding: '50px', textAlign: 'center' }}>
          <Card style={{ maxWidth: '500px', margin: '0 auto' }}>
            <Title level={2} style={{ color: '#ff4d4f' }}>Application Error</Title>
            <p>Something went wrong while loading the application:</p>
            <p><code>{error}</code></p>
            <Button type="primary" onClick={() => window.location.reload()}>
              Reload Application
            </Button>
          </Card>
        </Content>
      </Layout>
    );
  }

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      setLoading(true);
      console.log('Starting to load data...');
      
      // Load data individually to better identify which part is failing
      let monitorsData: Monitor[] = [];
      let audioData: {filtered: AudioDevice[], ignored: AudioDevice[]} = {filtered: [], ignored: []};
      let profilesData: Profile[] = [];
      
      try {
        console.log('Loading monitors...');
        monitorsData = await GetMonitors();
        console.log('Monitors loaded:', monitorsData);
      } catch (error) {
        console.error('Error loading monitors:', error);
        antMessage.error(`Error loading monitors: ${error}`);
        monitorsData = [];
      }
      
      try {
        console.log('Loading audio devices...');
        const audioResult = await GetAudioDevicesWithIgnoreStatus();
        console.log('Audio devices loaded:', audioResult);
        audioData = audioResult as {filtered: AudioDevice[], ignored: AudioDevice[]};
        
        
        // Validate the structure
        if (!audioData || !Array.isArray(audioData.filtered) || !Array.isArray(audioData.ignored)) {
          console.warn('Invalid audio data structure, using defaults');
          audioData = {filtered: [], ignored: []};
        }
        
        console.log('Audio devices loaded:', audioData);
      } catch (error) {
        console.error('Error loading audio devices:', error);
        antMessage.error(`Error loading audio devices: ${error}`);
        audioData = {filtered: [], ignored: []};
      }
      
      try {
        console.log('Loading profiles...');
        profilesData = await GetProfiles();
        console.log('Profiles loaded:', profilesData);
      } catch (error) {
        console.error('Error loading profiles:', error);
        antMessage.error(`Error loading profiles: ${error}`);
        profilesData = [];
      }
      
      setMonitors(monitorsData);
      setAudioDevices(audioData as {filtered: AudioDevice[], ignored: AudioDevice[]});
      setProfiles(profilesData);
      console.log('All data loaded successfully');
    } catch (error) {
      console.error('Critical error in loadData:', error);
      const errorMessage = error instanceof Error ? error.message : String(error);
      setError(errorMessage);
      antMessage.error(`Critical error loading data: ${error}`);
      // Set empty fallbacks
      setMonitors([]);
      setAudioDevices({filtered: [], ignored: []});
      setProfiles([]);
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

  const handleProfilesChange = async () => {
    try {
      const profilesData = await GetProfiles();
      setProfiles(profilesData);
    } catch (error) {
      antMessage.error(`Error loading profiles: ${error}`);
    }
  };

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
              <MonitorsTable 
                monitors={monitors}
                loading={loading}
                onRefresh={handleRefreshMonitors}
              />
              
              <AudioDevicesTable 
                audioDevices={audioDevices}
                showIgnoredAudio={showIgnoredAudio}
                setShowIgnoredAudio={setShowIgnoredAudio}
                loading={loading}
                onRefresh={handleRefreshAudio}
              />
            </Space>
          </Col>

          <Col xs={24} xl={8}>
            <ProfileManagement 
              profiles={profiles}
              loading={loading}
              onProfilesChange={handleProfilesChange}
              audioDevices={audioDevices}
            />
          </Col>
        </Row>
      </Content>
    </Layout>
  );
}

export default App;
