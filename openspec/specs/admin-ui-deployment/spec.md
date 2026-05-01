# admin-ui-deployment

## Purpose

Deployment configuration and documentation for admin-ui.

## Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `VITE_API_BASE_URL` | API base URL | `http://localhost:8080` | Yes |
| `VITE_USE_MOCK` | Enable Mock API mode | `false` | No |
| `VITE_MOCK_DELAY` | Mock API network delay (ms) | `500` | No |

## Deployment Scenarios

### Development
```bash
npm run dev
# Uses .env.development with VITE_USE_MOCK=true by default
```

### Production (Docker)
```dockerfile
FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
```

### Production (Standalone)
```bash
npm run build
npm run preview
# Or serve dist/ with any static file server
```

## Build Configuration

### Vite Config
- Proxy `/admin/*` and `/v1/*` to gateway-service (dev only)
- Build output to `dist/` directory
- TypeScript strict mode enabled

### Nginx Config (Docker)
```nginx
server {
    listen 80;
    server_name _;
    
    root /usr/share/nginx/html;
    index index.html;
    
    location / {
        try_files $uri $uri/ /index.html;
    }
}
```

## Build Optimization

- Mock API code SHALL be tree-shaken in production (DevTools component checks `import.meta.env.PROD`)
- Bundle size SHALL be monitored and kept under 500KB gzipped

## CI/CD

### Build Verification
```bash
npm run build
npm run test
npm run lint
```

### Environment-Specific Builds
- Development: `npm run dev` with hot reload
- Staging: `npm run build` with `.env.staging`
- Production: `npm run build` with `.env.production`
