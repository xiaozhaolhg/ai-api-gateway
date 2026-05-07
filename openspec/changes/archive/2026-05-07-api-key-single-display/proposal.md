## Why

When users create an API key on the API creation page, the key persists in the UI even after navigating away and back to the page. This creates a security risk where the API key remains visible in the browser, potentially exposing sensitive credentials to unauthorized viewers or screen sharing.

## What Changes

1. **admin-ui (API Key Creation)**: Implement ephemeral display logic where the API key is only shown once immediately after creation, then cleared from component state when the user navigates away or closes the modal
2. **admin-ui (State Management)**: Add flag to track if the key has been dismissed or user has navigated away, preventing re-display of the same key
3. **admin-ui (User Experience)**: Add clear messaging that the key will only be shown once and should be copied immediately

## Capabilities

### New Capabilities
- `admin-ui-ephemeral-secrets`: Secure handling of sensitive data in admin UI with single-display behavior for API keys and other secrets

### Modified Capabilities
- `admin-ui-architecture`: Enhance "Secure credential handling" requirement with ephemeral display pattern for sensitive data

## Impact

**Affected Services:**
- `admin-ui`: API key creation component and state management (`src/pages/APIKeys.tsx`, `src/components/APIKeyModal.tsx`)

**API Changes:** None

**Dependencies:**
- React Router (for navigation detection)
- React state management (useState, useEffect)

**Risks:**
- User may lose the key if they don't copy it immediately
- Need clear UX messaging to prevent user confusion
