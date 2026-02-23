import { Modal, Button, Typography, Space } from 'antd';
import { DeleteOutlined } from '@ant-design/icons';

const { Title, Text } = Typography;

interface ConfirmProfileDeleteProps {
  visible: boolean;
  profileName: string;
  onConfirm: () => void;
  onCancel: () => void;
  loading?: boolean;
}

export function ConfirmProfileDelete({ 
  visible, 
  profileName, 
  onConfirm, 
  onCancel, 
  loading = false 
}: ConfirmProfileDeleteProps) {
  return (
    <Modal
      title="Delete Profile"
      open={visible}
      onCancel={onCancel}
      footer={[
        <Button key="cancel" onClick={onCancel} disabled={loading}>
          Cancel
        </Button>,
        <Button
          key="delete"
          type="primary"
          danger
          onClick={onConfirm}
          loading={loading}
        >
          Delete Profile
        </Button>
      ]}
      centered
    >
      <div style={{ textAlign: 'center', padding: '20px 0' }}>
        <Title level={4}>Are you sure you want to delete this profile?</Title>
        <Text strong style={{ fontSize: '16px', color: '#ff4d4f' }}>
          "{profileName}"
        </Text>
        <div style={{ marginTop: '16px' }}>
          <Text type="secondary">
            This action cannot be undone. The profile will be permanently removed.
          </Text>
        </div>
      </div>
    </Modal>
  );
}
