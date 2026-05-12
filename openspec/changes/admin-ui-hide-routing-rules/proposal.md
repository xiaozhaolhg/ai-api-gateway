## Why

The Routing Rules page in the admin UI is not required for the current operational workflow. Model-to-provider routing is handled automatically by the router service through tier-based permissions and the routing configuration managed via provider setup. The standalone Routing Rules page adds UI complexity without serving a functional need, and removing it simplifies navigation.

## What Changes

- Remove the Routing Rules page route from App.tsx
- Remove the Routing Rules sidebar navigation entry from AppShell.tsx
- Keep the page component source file for potential future re-enablement
- No backend code changes

## Capabilities

### Removed Capabilities
- `admin-ui-routing-rules-page`: Standalone Routing Rules management page

### Preserved Capabilities
- Backend routing rule API endpoints remain intact (accessible via API)
- Router service continues automatic model-to-provider resolution
