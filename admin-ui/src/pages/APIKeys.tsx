import { useState, useEffect } from 'react';
import Layout from '../components/Layout';
import { apiClient, type APIKey } from '../api/client';

export default function APIKeys() {
  const [apiKeys, setApiKeys] = useState<APIKey[]>([]);
  const [selectedUserId, setSelectedUserId] = useState('');
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [newKeyName, setNewKeyName] = useState('');
  const [createdKey, setCreatedKey] = useState<{ api_key_id: string; api_key: string } | null>(null);

  useEffect(() => {
    if (selectedUserId) {
      loadAPIKeys();
    }
  }, [selectedUserId]);

  const loadAPIKeys = async () => {
    try {
      const data = await apiClient.getAPIKeys(selectedUserId);
      setApiKeys(data);
    } catch (error) {
      console.error('Failed to load API keys:', error);
    }
  };

  const handleCreateAPIKey = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const result = await apiClient.createAPIKey(selectedUserId, newKeyName);
      setCreatedKey(result);
      setShowCreateForm(false);
      loadAPIKeys();
    } catch (error) {
      console.error('Failed to create API key:', error);
    }
  };

  const handleDeleteAPIKey = async (id: string) => {
    if (!confirm('Are you sure you want to delete this API key?')) return;
    try {
      await apiClient.deleteAPIKey(id);
      loadAPIKeys();
    } catch (error) {
      console.error('Failed to delete API key:', error);
    }
  };

  return (
    <Layout>
      <div className="space-y-6">
        <h2 className="text-2xl font-bold text-gray-900">API Keys</h2>

        <div>
          <label className="block text-sm font-medium text-gray-700">Select User</label>
          <select
            value={selectedUserId}
            onChange={(e) => setSelectedUserId(e.target.value)}
            className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2"
          >
            <option value="">Select a user...</option>
            {/* In a real app, you'd fetch users here */}
            <option value="user-1">User 1</option>
            <option value="user-2">User 2</option>
          </select>
        </div>

        {selectedUserId && (
          <>
            <div className="flex justify-between items-center">
              <h3 className="text-lg font-medium text-gray-900">API Keys for {selectedUserId}</h3>
              <button
                onClick={() => setShowCreateForm(true)}
                className="bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700"
              >
                Issue API Key
              </button>
            </div>

            {showCreateForm && (
              <div className="bg-white shadow rounded-lg p-6">
                <form onSubmit={handleCreateAPIKey} className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Key Name</label>
                    <input
                      type="text"
                      value={newKeyName}
                      onChange={(e) => setNewKeyName(e.target.value)}
                      className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2"
                      required
                    />
                  </div>
                  <div className="flex space-x-2">
                    <button
                      type="submit"
                      className="bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700"
                    >
                      Create
                    </button>
                    <button
                      type="button"
                      onClick={() => setShowCreateForm(false)}
                      className="bg-gray-300 text-gray-700 px-4 py-2 rounded-md hover:bg-gray-400"
                    >
                      Cancel
                    </button>
                  </div>
                </form>
              </div>
            )}

            {createdKey && (
              <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
                <h4 className="font-medium text-yellow-800">API Key Created</h4>
                <p className="text-sm text-yellow-700 mt-1">This key will only be shown once. Copy it now.</p>
                <div className="mt-2 flex items-center space-x-2">
                  <code className="bg-yellow-100 px-2 py-1 rounded text-sm">{createdKey.api_key}</code>
                  <button
                    onClick={() => {
                      navigator.clipboard.writeText(createdKey.api_key);
                      setCreatedKey(null);
                    }}
                    className="text-blue-600 hover:text-blue-800 text-sm"
                  >
                    Copy & Close
                  </button>
                </div>
              </div>
            )}

            <div className="bg-white shadow rounded-lg overflow-hidden">
              <table className="min-w-full divide-y divide-gray-200">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Name</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Scopes</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Created At</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Expires At</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {apiKeys.map((key) => (
                    <tr key={key.id}>
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{key.name}</td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{key.scopes.join(', ')}</td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{new Date(key.created_at).toLocaleString()}</td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {key.expires_at ? new Date(key.expires_at).toLocaleString() : 'Never'}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                        <button
                          onClick={() => handleDeleteAPIKey(key.id)}
                          className="text-red-600 hover:text-red-900"
                        >
                          Revoke
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </>
        )}
      </div>
    </Layout>
  );
}
