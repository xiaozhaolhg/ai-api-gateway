# Admin UI

React-based admin dashboard for AI API Gateway management.

## Features

- **Authentication**: JWT-based login with role management (admin/user/viewer)
- **Dashboard**: System overview with key metrics
- **Provider Management**: Configure LLM providers (Ollama, OpenCode Zen)
- **User Management**: Create and manage users with roles
- **API Keys**: Generate and manage API keys
- **Usage Analytics**: Track token usage and costs
- **Health Monitoring**: Service status monitoring

## Tech Stack

- React 18 + TypeScript
- Vite
- TanStack Query (data fetching)
- React Hook Form (forms)
- Tailwind CSS
- i18n (English/Chinese)

## Getting Started

### Prerequisites

- Node.js 18+
- Running gateway-service (port 8080)

### Install Dependencies

```bash
cd admin-ui
npm install
```

### Development

```bash
# Start dev server
npm run dev

# Build for production
npm run build
```

### Docker

```bash
# Build image
docker build -t ai-api-gateway/admin-ui:latest .

# Run container
docker run -p 3000:80 ai-api-gateway/admin-ui:latest
```

## Routes

| Route | Description | Required Role |
|-------|-------------|--------------|
| /admin/login | Login page | - |
| /admin/dashboard | Dashboard | any |
| /admin/providers | Providers | admin |
| /admin/users | Users | admin |
| /admin/api-keys | API Keys | user |
| /admin/usage | Usage | user |
| /admin/health | Health | user |
| /admin/settings | Settings | user |

## Roles

| Role | Permissions |
|------|-------------|
| admin | Full access |
| user | CRUD own resources |
| viewer | Read-only |

## API Integration

The admin UI communicates with gateway-service (port 8080) which proxies to backend gRPC services.

### Authentication

1. POST /admin/login with email/password
2. Returns JWT in HTTP-only cookie (path: /admin, secure)
3. Middleware validates JWT for protected routes
4. POST /admin/logout clears cookie

### Data Fetching

Uses TanStack Query with:
- 5-minute cache TTL
- Automatic retry on failure
- Optimistic updates for mutations

## Environment

| Variable | Description | Default |
|----------|-------------|---------|
| VITE_API_URL | API base URL | http://localhost:8080 |

## i18n

Supported locales:
- English (en)
- Chinese (zh)

Language switcher in header.

## Code Organization

```
src/
├── api/           # API client
├── components/    # Reusable components
│   ├── AuthContext.tsx
│   ├── ProtectedRoute.tsx
│   └── AppShell.tsx
├── contexts/      # React contexts
├── i18n/          # Translations
├── pages/         # Route pages
│   ├── Login.tsx
│   ├── Dashboard.tsx
│   ├── Providers.tsx
│   ├── Users.tsx
│   ├── APIKeys.tsx
│   ├── Usage.tsx
│   └── Health.tsx
└── main.tsx       # Entry point
```

## License

MIT