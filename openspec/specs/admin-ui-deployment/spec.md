## Purpose

Deployment specifications for the admin-ui React SPA with nginx serving and Docker Compose integration.
## Requirements
### Requirement: Docker build with nginx
The admin-ui SHALL be built as a Docker image that serves static assets via nginx.

#### Scenario: Docker build produces nginx image
- **WHEN** docker build is executed on the admin-ui Dockerfile
- **THEN** the resulting image SHALL be a multi-stage build: node for building, nginx for serving

#### Scenario: nginx serves SPA with fallback routing
- **WHEN** a request is made to nginx for a client-side route (e.g., /providers)
- **THEN** nginx SHALL serve index.html so client-side routing works

### Requirement: CORS and proxy configuration
The admin-ui SHALL be configured to communicate with gateway-service.

#### Scenario: nginx proxies API calls in Docker Compose
- **WHEN** the admin-ui makes an API call in Docker Compose
- **THEN** nginx SHALL proxy /api/* requests to gateway-service:8080
- **AND** CORS headers SHALL NOT be needed (same-origin via proxy)

#### Scenario: Vite dev proxy for local development
- **WHEN** the admin-ui dev server is running locally
- **THEN** Vite SHALL proxy API calls to the gateway-service at localhost:8080

### Requirement: Docker Compose integration
The admin-ui SHALL be deployable via Docker Compose.

#### Scenario: Service starts in Docker Compose
- **WHEN** `docker-compose up` is run
- **THEN** the admin-ui SHALL start on HTTP port 3000
- **AND** it SHALL be accessible at http://localhost:3000

