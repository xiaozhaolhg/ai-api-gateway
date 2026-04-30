import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { Rate, Trend, Counter } from 'k6/metrics';

// Custom metrics
const errorRate = new Rate('errors');
const responseTime = new Trend('response_time');
const streamingTime = new Trend('streaming_time');

// Test configuration
export const options = {
  scenarios: {
    // Non-streaming load test
    non_streaming: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '1m', target: 50 },    // Ramp up to 50 VUs
        { duration: '3m', target: 100 },   // Ramp to 100 VUs
        { duration: '1m', target: 0 },     // Ramp down
      ],
      gracefulRampDown: '30s',
      exec: 'nonStreamingTest',
    },
    // Streaming load test
    streaming: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '1m', target: 20 },    // Ramp up to 20 VUs
        { duration: '3m', target: 50 },      // Ramp to 50 VUs
        { duration: '1m', target: 0 },       // Ramp down
      ],
      gracefulRampDown: '30s',
      exec: 'streamingTest',
    },
  },
  thresholds: {
    http_req_duration: ['p(95)<2000'],     // 95% of requests under 2s
    http_req_failed: ['rate<0.01'],         // Error rate under 1%
    errors: ['rate<0.01'],
  },
};

const BASE_URL = __ENV.GATEWAY_URL || 'http://localhost:8080';

// Sample chat completion request
const chatCompletionPayload = JSON.stringify({
  model: 'ollama:llama2',
  messages: [
    { role: 'user', content: 'Hello, how are you?' }
  ],
  max_tokens: 150,
  temperature: 0.7,
});

// Sample streaming chat completion request
const streamingPayload = JSON.stringify({
  model: 'ollama:llama2',
  messages: [
    { role: 'user', content: 'Tell me a short story' }
  ],
  max_tokens: 500,
  temperature: 0.7,
  stream: true,
});

// Health check test
export function healthCheck() {
  group('Health Endpoints', () => {
    // Simple health check
    const healthRes = http.get(`${BASE_URL}/health`);
    check(healthRes, {
      'health status is 200': (r) => r.status === 200,
      'health response has status': (r) => r.json('status') === 'ok',
    });
    errorRate.add(healthRes.status !== 200);
    responseTime.add(healthRes.timings.duration);

    // Deep health check
    const gatewayHealthRes = http.get(`${BASE_URL}/gateway/health`);
    check(gatewayHealthRes, {
      'gateway health status is 200 or 503': (r) => r.status === 200 || r.status === 503,
      'gateway health has services': (r) => r.json('services') !== undefined,
    });
    errorRate.add(gatewayHealthRes.status !== 200 && gatewayHealthRes.status !== 503);
    responseTime.add(gatewayHealthRes.timings.duration);
  });
}

// Models endpoint test
export function modelsTest() {
  group('Models Endpoint', () => {
    const res = http.get(`${BASE_URL}/gateway/models`);
    check(res, {
      'models status is 200': (r) => r.status === 200,
      'models has data array': (r) => Array.isArray(r.json('data')),
    });
    errorRate.add(res.status !== 200);
    responseTime.add(res.timings.duration);
  });
}

// Non-streaming chat completion test
export function nonStreamingTest() {
  // Run health check first
  healthCheck();
  
  // Run models test
  modelsTest();

  group('Chat Completion (Non-Streaming)', () => {
    const headers = {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${__ENV.API_KEY || 'test-api-key'}`,
    };

    const startTime = Date.now();
    const res = http.post(
      `${BASE_URL}/v1/chat/completions`,
      chatCompletionPayload,
      { headers }
    );
    const duration = Date.now() - startTime;

    check(res, {
      'non-streaming status is 200': (r) => r.status === 200,
      'non-streaming has choices': (r) => r.json('choices') !== undefined || r.json('message') !== undefined,
      'non-streaming response time < 2s': () => duration < 2000,
    });

    errorRate.add(res.status !== 200);
    responseTime.add(res.timings.duration);
  });

  sleep(1);
}

// Streaming chat completion test
export function streamingTest() {
  group('Chat Completion (Streaming)', () => {
    const headers = {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${__ENV.API_KEY || 'test-api-key'}`,
      'Accept': 'text/event-stream',
    };

    const startTime = Date.now();
    const res = http.post(
      `${BASE_URL}/v1/chat/completions`,
      streamingPayload,
      { headers }
    );
    const duration = Date.now() - startTime;

    check(res, {
      'streaming status is 200': (r) => r.status === 200,
      'streaming content-type is SSE': (r) => r.headers['Content-Type'] && r.headers['Content-Type'].includes('text/event-stream'),
      'streaming has data': (r) => r.body && r.body.includes('data:'),
    });

    errorRate.add(res.status !== 200);
    streamingTime.add(duration);
    responseTime.add(res.timings.duration);
  });

  sleep(2);
}

// Default function for simple tests
export default function() {
  healthCheck();
  sleep(1);
}
