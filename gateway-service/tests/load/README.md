# Gateway Load Tests (k6)

This directory contains k6 load tests for the AI API Gateway.

## Prerequisites

Install k6:
- **Mac**: `brew install k6`
- **Linux**: `sudo gpg -k && sudo gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C802E315D07D0CD8A4B881D7E55F7B2A9 && echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list && sudo apt-get update && sudo apt-get install k6`
- **Windows**: `winget install k6`
- **Docker**: `docker pull grafana/k6`

See [k6 Installation Guide](https://k6.io/docs/get-started/installation/) for more options.

## Running Tests

### Basic Run

```bash
# Run with default settings (localhost:8080)
k6 run gateway_load_test.js

# Run against specific gateway URL
GATEWAY_URL=http://gateway.example.com k6 run gateway_load_test.js

# Run with API key
GATEWAY_URL=http://localhost:8080 API_KEY=your-api-key k6 run gateway_load_test.js
```

### Docker

```bash
docker run -v $(pwd):/tests -e GATEWAY_URL=http://host.docker.internal:8080 grafana/k6 run /tests/gateway_load_test.js
```

### Cloud (k6 Cloud)

```bash
# Login to k6 Cloud
k6 login cloud

# Run test on k6 Cloud
k6 cloud gateway_load_test.js
```

## Test Scenarios

### Non-Streaming Load Test

- **Target**: 100 virtual users
- **Duration**: 5 minutes total (ramp up/down included)
- **Endpoints**:
  - GET /health
  - GET /gateway/health
  - GET /gateway/models
  - POST /v1/chat/completions (non-streaming)

### Streaming Load Test

- **Target**: 50 virtual users
- **Duration**: 5 minutes total
- **Endpoint**: POST /v1/chat/completions (streaming)

## Thresholds

The test will fail if:
- **p95 response time** > 2 seconds
- **Error rate** > 1%

## Test Results

Example successful output:

```
✓ http_req_duration..............: avg=145.12ms min=12.34ms med=132.45ms max=1.89s p(95)=1.23s p(99)=1.56s
✓ http_req_failed................: 0.00% ✓ 0 ✗ 5000
✓ errors.........................: 0.00% ✓ 0 ✗ 5000
```

## Customizing Tests

Edit `gateway_load_test.js` to modify:
- Virtual users and ramp stages
- Request payloads
- Thresholds
- Test duration

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `GATEWAY_URL` | `http://localhost:8080` | Gateway base URL |
| `API_KEY` | `test-api-key` | API key for authentication |

## Troubleshooting

### Connection Refused

Ensure the gateway service is running:
```bash
curl http://localhost:8080/health
```

### High Error Rate

Check gateway logs:
```bash
docker logs gateway-service
```

### Timeout Issues

Increase timeout in test options if testing slow providers.

## CI/CD Integration

### GitHub Actions

```yaml
- name: Run k6 Load Test
  uses: grafana/k6-action@v0.3.1
  with:
    filename: tests/load/gateway_load_test.js
  env:
    GATEWAY_URL: http://localhost:8080
```

### GitLab CI

```yaml
load-test:
  image: grafana/k6:latest
  script:
    - k6 run tests/load/gateway_load_test.js
  variables:
    GATEWAY_URL: http://gateway:8080
```
