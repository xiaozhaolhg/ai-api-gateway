# Admin UI Deployment Guide

This guide covers deploying the admin UI to production environments.

## Prerequisites

- Node.js 18+
- Docker & Docker Compose (for containerized deployment)
- Running gateway-service instance (port 8080)

## Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `VITE_API_BASE_URL` | Base URL for gateway-service API | `http://localhost:8080` | No |

### Production Configuration

Create a `.env.production` file:

```bash
VITE_API_BASE_URL=https://api.yourdomain.com
```

## Local Development

```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Access at http://localhost:5173
```

## Docker Deployment

### Build Image

```bash
# From project root
docker build -t ai-api-gateway/admin-ui:latest ./admin-ui
```

### Run Container

```bash
docker run -d \
  -p 3000:80 \
  -e VITE_API_BASE_URL=https://api.yourdomain.com \
  --name admin-ui \
  ai-api-gateway/admin-ui:latest
```

### Docker Compose

The admin UI is included in the main `docker-compose.yaml`:

```yaml
services:
  admin-ui:
    build:
      context: ./admin-ui
      dockerfile: Dockerfile
    ports:
      - "3000:80"
    depends_on:
      - gateway-service
    environment:
      - VITE_API_BASE_URL=http://gateway-service:8080
```

Start with all services:

```bash
docker compose up -d
```

## Build for Production

```bash
# Build optimized production bundle
npm run build

# Preview production build locally
npm run preview
```

The build output is in `dist/` directory.

## Nginx Configuration

Example nginx configuration for serving the admin UI:

```nginx
server {
    listen 80;
    server_name admin.yourdomain.com;
    root /usr/share/nginx/html;
    index index.html;

    # Gzip compression
    gzip on;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/xml+rss text/javascript;

    # SPA routing - all routes serve index.html
    location / {
        try_files $uri $uri/ /index.html;
    }

    # Cache static assets
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header Content-Security-Policy "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline';" always;
}
```

## Kubernetes Deployment

Example Kubernetes manifest:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: admin-ui
spec:
  replicas: 2
  selector:
    matchLabels:
      app: admin-ui
  template:
    metadata:
      labels:
        app: admin-ui
    spec:
      containers:
      - name: admin-ui
        image: ai-api-gateway/admin-ui:latest
        ports:
        - containerPort: 80
        env:
        - name: VITE_API_BASE_URL
          value: "https://api.yourdomain.com"
        resources:
          requests:
            memory: "64Mi"
            cpu: "100m"
          limits:
            memory: "128Mi"
            cpu: "200m"
        livenessProbe:
          httpGet:
            path: /
            port: 80
          initialDelaySeconds: 10
          periodSeconds: 10
---
apiVersion: v1
kind: Service
metadata:
  name: admin-ui
spec:
  selector:
    app: admin-ui
  ports:
  - port: 80
    targetPort: 80
  type: ClusterIP
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: admin-ui
  annotations:
    kubernetes.io/ingress.class: nginx
    cert-manager.io/cluster-issuer: letsencrypt-prod
spec:
  tls:
  - hosts:
    - admin.yourdomain.com
    secretName: admin-ui-tls
  rules:
  - host: admin.yourdomain.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: admin-ui
            port:
              number: 80
```

## Health Check

The admin UI serves a health endpoint at `/health`:

```bash
curl http://localhost:3000/health
```

Response:

```json
{
  "status": "ok",
  "timestamp": "2026-04-27T12:00:00.000Z"
}
```

## Monitoring

### Logs

```bash
# Docker logs
docker logs -f admin-ui

# Kubernetes logs
kubectl logs -f deployment/admin-ui
```

### Metrics

The admin UI does not expose Prometheus metrics. Monitor at the infrastructure level:
- Request rate (nginx/ingress)
- Response time (nginx/ingress)
- Error rate (4xx, 5xx responses)
- Container resource usage (CPU, memory)

## Troubleshooting

### Issue: Blank page after deployment

**Check:**
1. Browser console for JavaScript errors
2. Network tab for failed API requests
3. `VITE_API_BASE_URL` is correctly set

### Issue: API requests failing with CORS errors

**Solution:**
Ensure gateway-service has CORS configured to allow requests from admin UI origin:

```yaml
# gateway-service config
cors:
  allowed_origins:
    - "https://admin.yourdomain.com"
```

### Issue: Login redirects to blank page

**Check:**
1. Gateway-service `/admin/auth/login` endpoint is accessible
2. JWT cookie is being set correctly
3. Browser is not blocking third-party cookies

### Issue: Build fails with TypeScript errors

**Solution:**
```bash
# Clear node_modules and rebuild
rm -rf node_modules package-lock.json
npm install
npm run build
```

## Security Best Practices

1. **Use HTTPS in production**: Always serve admin UI over HTTPS
2. **Set secure cookies**: JWT cookies should have `Secure` and `HttpOnly` flags
3. **Content Security Policy**: Implement strict CSP headers
4. **Regular updates**: Keep dependencies updated for security patches
5. **Access control**: Restrict admin UI access to trusted IPs or VPN
6. **Session timeout**: Implement automatic logout after inactivity

## Scaling Considerations

- **Horizontal scaling**: Run multiple replicas behind a load balancer
- **CDN**: Serve static assets from CDN for better performance
- **Caching**: Enable browser caching for static assets (1 year)
- **Compression**: Enable gzip/brotli compression

## Rollback

If deployment issues occur:

```bash
# Docker: rollback to previous image
docker stop admin-ui
docker rm admin-ui
docker run -d -p 3000:80 ai-api-gateway/admin-ui:previous-tag

# Kubernetes: rollback deployment
kubectl rollout undo deployment/admin-ui
```
