## Context

**Current State:**
The admin-ui API key creation page displays the generated API key in a modal/dialog after successful creation. The key remains stored in component state and persists even when users navigate away and back to the page. This creates a security vulnerability where sensitive credentials remain visible in the browser.

**Constraints:**
- Must work with existing React Router navigation
- Cannot rely on server-side state (key is only shown once at creation)
- Must maintain good UX - users need time to copy the key
- Should work across browser tabs/windows (if user opens new tab)

## Goals / Non-Goals

**Goals:**
1. Ensure API keys are only displayed once immediately after creation
2. Clear the key from browser memory when user navigates away or closes modal
3. Provide clear UX messaging about the single-display behavior
4. Prevent key persistence across page reloads or navigation

**Non-Goals:**
1. Implement server-side key revocation or re-display
2. Store keys in any persistent browser storage (localStorage, sessionStorage)
3. Modify the backend API key generation process

## Decisions

### D1: Use Component State + Navigation Detection

**Decision:** Store the API key in component state only, clear it on navigation events and modal close.

**Rationale:**
- Component state is ephemeral - cleared on unmount
- Navigation detection ensures key disappears when user leaves the page
- Modal close provides immediate clearing for better UX
- No persistent storage means key truly disappears

**Implementation:**
```typescript
const [newApiKey, setNewApiKey] = useState<string | null>(null);
const [keyDisplayed, setKeyDisplayed] = useState(false);

// Clear key on navigation
useEffect(() => {
  const handleNavigation = () => {
    if (newApiKey) {
      setNewApiKey(null);
      setKeyDisplayed(true); // Mark as already shown
    }
  };
  
  window.addEventListener('beforeunload', handleNavigation);
  return () => window.removeEventListener('beforeunload', handleNavigation);
}, [newApiKey]);

// Clear key on modal close
const handleModalClose = () => {
  setNewApiKey(null);
  setKeyDisplayed(true);
};
```

### D2: Add "Already Shown" Flag

**Decision:** Track whether the key has been displayed once to prevent re-display even if user navigates back.

**Rationale:**
- Prevents accidental re-display if component remounts
- Provides clear signal that key was already shown
- Works across browser refreshes (use sessionStorage for flag only)

**Implementation:**
```typescript
// Track display status
const [keyDismissed, setKeyDismissed] = useState(() => {
  return sessionStorage.getItem('api-key-dismissed') === 'true';
});

// Mark as dismissed when modal closes
const handleModalClose = () => {
  setNewApiKey(null);
  setKeyDismissed(true);
  sessionStorage.setItem('api-key-dismissed', 'true');
};
```

### D3: Clear UX Messaging

**Decision:** Add prominent warning that the key will only be shown once.

**Rationale:**
- Prevents user confusion when key disappears
- Encourages immediate copying behavior
- Sets proper security expectations

**Implementation:**
```typescript
<Alert
  type="warning"
  message="Security Notice"
  description="This API key will only be shown once. Please copy it now and store it securely."
  showIcon
/>
```

## Risks / Trade-offs

**Risk:** User may lose the key if they don't copy it immediately
→ **Mitigation:** Clear messaging + copy-to-clipboard button + prominent warning

**Risk:** Component remount might re-display key if flag not persisted
→ **Mitigation:** Use sessionStorage for "dismissed" flag (not the key itself)

**Risk:** Browser tab navigation might not trigger beforeunload event
→ **Mitigation:** Also clear on route change detection using React Router

**Trade-off:** Using sessionStorage for flag vs pure component state
→ Acceptable because flag doesn't contain sensitive data, just boolean status

## Migration Plan

1. **Phase 1:** Add ephemeral state management to API key creation component
2. **Phase 2:** Implement navigation detection and modal close handlers
3. **Phase 3:** Add UX messaging and copy functionality
4. **Phase 4:** Test edge cases (tab navigation, browser refresh, multiple tabs)

**Rollback:** Changes are frontend-only. Removing the navigation detection and state clearing restores original behavior.

## Open Questions

1. Should we also clear the key when user switches to another browser tab?
2. Should we add a countdown timer showing how long until key auto-dismisses?
3. Should we implement a "I've saved this key" checkbox to let users explicitly dismiss?
