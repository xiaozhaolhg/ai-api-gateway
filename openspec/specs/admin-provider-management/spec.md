## Requirements

### Functional Requirements

#### FR1: Remove Duplicate Provider Route
- **Given** the gateway service has both mock and real provider handlers
- **When** the service starts
- **Then** only the real provider handler should be registered
- **And** the mock handler should be removed

#### FR2: Return Provider List as Plain Array
- **Given** a GET request to `/admin/providers`
- **When** the provider service returns data
- **Then** the response should be a JSON array: `[{...}, {...}]`
- **And** NOT a wrapped object: `{Providers: [{...}, {...}]}`

#### FR3: Include BaseURL in Provider Response
- **Given** a provider with base_url configured
- **When** listing providers via the gateway
- **Then** the response should include the `base_url` field

#### FR4: Frontend Null Safety for Models
- **Given** the providers page is loading
- **When** a provider has null or undefined models
- **Then** the page should not crash
- **And** should display empty string for models

### API Contract

#### GET /admin/providers

**Request**: None

**Response** (200 OK):
```json
[
  {
    "id": "provider-1",
    "name": "Ollama Local",
    "type": "ollama",
    "base_url": "http://localhost:11434",
    "models": ["llama2", "mistral"],
    "status": "active"
  }
]
```

**Error Response** (503):
```json
{
  "error": "provider service unavailable",
  "code": "SERVICE_UNAVAILABLE"
}
```

### Non-Functional Requirements

- **NFR1**: Response time < 500ms for provider list with 100 providers
- **NFR2**: Credentials must be masked ("***") in response

### UI Requirements

- **UIR1**: Provider table displays all fields: name, type, base_url, models, status
- **UIR2**: Models column handles empty/null gracefully
- **UIR3**: Status column shows color-coded tags (green=active, red=inactive)
