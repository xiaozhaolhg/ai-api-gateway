import { useState } from 'react';
import { Select, message } from 'antd';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { apiClient, type Tier } from '../api/client';

interface GroupTierTabProps {
  groupId: string;
  currentTierId?: string;
}

export const GroupTierTab: React.FC<GroupTierTabProps> = ({ groupId, currentTierId }) => {
  const queryClient = useQueryClient();
  const [selectedTierId, setSelectedTierId] = useState<string | undefined>(currentTierId);

  const { data: tiers = [], isLoading } = useQuery<Tier[]>({
    queryKey: ['tiers'],
    queryFn: () => apiClient.getTiers(),
  });

  const assignMutation = useMutation({
    mutationFn: (tierId: string) => apiClient.assignTierToGroup(groupId, tierId),
    onSuccess: () => {
      message.success('Tier assigned successfully');
      queryClient.invalidateQueries({ queryKey: ['groups'] });
      queryClient.invalidateQueries({ queryKey: ['groupTier', groupId] });
    },
    onError: () => {
      message.error('Failed to assign tier');
    },
  });

  const removeMutation = useMutation({
    mutationFn: () => apiClient.removeTierFromGroup(groupId),
    onSuccess: () => {
      message.success('Tier removed from group');
      setSelectedTierId(undefined);
      queryClient.invalidateQueries({ queryKey: ['groups'] });
      queryClient.invalidateQueries({ queryKey: ['groupTier', groupId] });
    },
    onError: () => {
      message.error('Failed to remove tier');
    },
  });

  const handleAssign = () => {
    if (selectedTierId) {
      assignMutation.mutate(selectedTierId);
    }
  };

  const handleRemove = () => {
    removeMutation.mutate();
  };

  const currentTier = tiers.find(t => t.id === currentTierId);

  return (
    <div>
      <div style={{ marginBottom: 16 }}>
        <p style={{ marginBottom: 16, color: '#666' }}>
          Assign a tier to this group to control which models and providers members can access.
          {currentTier && (
            <span style={{ display: 'block', marginTop: 8 }}>
              Current tier: <strong>{currentTier.name}</strong>
            </span>
          )}
        </p>

        <div style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
          <Select
            placeholder="Select a tier"
            value={selectedTierId}
            onChange={setSelectedTierId}
            style={{ width: 250 }}
            loading={isLoading}
            allowClear
            options={tiers.map(tier => ({
              value: tier.id,
              label: (
                <span>
                  {tier.name}
                  {tier.is_default && ' (Default)'}
                </span>
              ),
            }))}
          />
          <button
            type="button"
            onClick={handleAssign}
            disabled={!selectedTierId || assignMutation.isPending}
            style={{
              padding: '4px 16px',
              background: selectedTierId ? '#1890ff' : '#ccc',
              color: '#fff',
              border: 'none',
              borderRadius: 4,
              cursor: selectedTierId ? 'pointer' : 'not-allowed',
            }}
          >
            {assignMutation.isPending ? 'Assigning...' : 'Assign'}
          </button>
          {currentTierId && (
            <button
              type="button"
              onClick={handleRemove}
              disabled={removeMutation.isPending}
              style={{
                padding: '4px 16px',
                background: '#ff4d4f',
                color: '#fff',
                border: 'none',
                borderRadius: 4,
                cursor: 'pointer',
              }}
            >
              {removeMutation.isPending ? 'Removing...' : 'Remove Tier'}
            </button>
          )}
        </div>
      </div>
    </div>
  );
};
