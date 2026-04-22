import { useState, useEffect } from 'react';
import Layout from '../components/Layout';

interface ProviderHealth {
  id: string;
  name: string;
  status: 'healthy' | 'unhealthy' | 'unknown';
  latency_ms: number;
  error_rate: number;
  last_check: string;
}

export default function Health() {
  const [health, setHealth] = useState<ProviderHealth[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Mock data - in real app, fetch from monitor service
    setHealth([
      {
        id: 'provider-1',
        name: 'OpenAI',
        status: 'healthy',
        latency_ms: 150,
        error_rate: 0.01,
        last_check: new Date().toISOString(),
      },
      {
        id: 'provider-2',
        name: 'Ollama',
        status: 'healthy',
        latency_ms: 50,
        error_rate: 0.0,
        last_check: new Date().toISOString(),
      },
      {
        id: 'provider-3',
        name: 'Anthropic',
        status: 'unhealthy',
        latency_ms: 0,
        error_rate: 1.0,
        last_check: new Date().toISOString(),
      },
    ]);
    setLoading(false);
  }, []);

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'healthy':
        return 'bg-green-100 text-green-800';
      case 'unhealthy':
        return 'bg-red-100 text-red-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  return (
    <Layout>
      <div className="space-y-6">
        <h2 className="text-2xl font-bold text-gray-900">Provider Health</h2>

        {loading ? (
          <div>Loading...</div>
        ) : (
          <div className="bg-white shadow rounded-lg overflow-hidden">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Provider</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Latency</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Error Rate</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Last Check</th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {health.map((provider) => (
                  <tr key={provider.id}>
                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{provider.name}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm">
                      <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${getStatusColor(provider.status)}`}>
                        {provider.status}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{provider.latency_ms}ms</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{(provider.error_rate * 100).toFixed(1)}%</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{new Date(provider.last_check).toLocaleString()}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </Layout>
  );
}
