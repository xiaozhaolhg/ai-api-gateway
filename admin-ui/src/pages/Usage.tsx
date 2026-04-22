import { useState, useEffect } from 'react';
import Layout from '../components/Layout';
import { apiClient, type UsageRecord } from '../api/client';

export default function Usage() {
  const [usage, setUsage] = useState<UsageRecord[]>([]);
  const [filters, setFilters] = useState({
    userId: '',
    startDate: '',
    endDate: '',
  });

  useEffect(() => {
    loadUsage();
  }, []);

  const loadUsage = async () => {
    try {
      const data = await apiClient.getUsage(
        filters.userId || undefined,
        filters.startDate || undefined,
        filters.endDate || undefined
      );
      setUsage(data);
    } catch (error) {
      console.error('Failed to load usage:', error);
    }
  };

  const handleFilter = (e: React.FormEvent) => {
    e.preventDefault();
    loadUsage();
  };

  const totalTokens = usage.reduce((sum, record) => sum + record.total_tokens, 0);
  const totalCost = usage.reduce((sum, record) => sum + record.cost, 0);

  return (
    <Layout>
      <div className="space-y-6">
        <h2 className="text-2xl font-bold text-gray-900">Usage Dashboard</h2>

        <form onSubmit={handleFilter} className="bg-white shadow rounded-lg p-6 space-y-4">
          <div className="grid grid-cols-3 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700">User ID</label>
              <input
                type="text"
                value={filters.userId}
                onChange={(e) => setFilters({ ...filters, userId: e.target.value })}
                className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">Start Date</label>
              <input
                type="date"
                value={filters.startDate}
                onChange={(e) => setFilters({ ...filters, startDate: e.target.value })}
                className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">End Date</label>
              <input
                type="date"
                value={filters.endDate}
                onChange={(e) => setFilters({ ...filters, endDate: e.target.value })}
                className="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2"
              />
            </div>
          </div>
          <button
            type="submit"
            className="bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700"
          >
            Apply Filters
          </button>
        </form>

        <div className="grid grid-cols-3 gap-4">
          <div className="bg-white shadow rounded-lg p-6">
            <h3 className="text-lg font-medium text-gray-900">Total Requests</h3>
            <p className="text-3xl font-bold text-gray-900 mt-2">{usage.length}</p>
          </div>
          <div className="bg-white shadow rounded-lg p-6">
            <h3 className="text-lg font-medium text-gray-900">Total Tokens</h3>
            <p className="text-3xl font-bold text-gray-900 mt-2">{totalTokens.toLocaleString()}</p>
          </div>
          <div className="bg-white shadow rounded-lg p-6">
            <h3 className="text-lg font-medium text-gray-900">Total Cost</h3>
            <p className="text-3xl font-bold text-gray-900 mt-2">${totalCost.toFixed(2)}</p>
          </div>
        </div>

        <div className="bg-white shadow rounded-lg overflow-hidden">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">User</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Model</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Prompt Tokens</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Completion Tokens</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Total Tokens</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Cost</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Timestamp</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {usage.map((record) => (
                <tr key={record.id}>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{record.user_id}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{record.model}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{record.prompt_tokens}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{record.completion_tokens}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{record.total_tokens}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">${record.cost.toFixed(4)}</td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{new Date(record.timestamp).toLocaleString()}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </Layout>
  );
}
