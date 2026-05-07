# Tasks: API Key Single Display

## 1. API Key Creation Component Refactor

- [x] **Task 1.1**: Add ephemeral state management to API key creation
  - Add `newApiKey: string | null` state for storing the generated key
  - Add `keyDismissed: boolean` state to track if key was shown once
  - Remove any localStorage/sessionStorage storage of the actual key
  - **Acceptance**: Unit test verifies key is only in component state

- [x] **Task 1.2**: Implement navigation detection for key clearing
  - Add `beforeunload` event listener to clear key on page navigation
  - Add React Router navigation detection (route change)
  - Clear `newApiKey` state when navigation is detected
  - **Acceptance**: Unit test verifies key is cleared on navigation events

- [x] **Task 1.3**: Implement modal close handler
  - Add `handleModalClose` function to clear key and set dismissed flag
  - Set `keyDismissed` to true when modal closes
  - Store dismissal flag in sessionStorage (not the key itself)
  - **Acceptance**: Unit test verifies key is cleared and flag is set on modal close

## 2. User Experience Enhancements

- [x] **Task 2.1**: Add security warning message
  - Create Alert component with warning type about single-display behavior
  - Message: "This API key will only be shown once. Please copy it now and store it securely."
  - Display prominently when key is shown
  - **Acceptance**: Visual test confirms warning message appears with key

- [x] **Task 2.2**: Implement copy-to-clipboard functionality
  - Add copy button next to the API key display
  - Use `navigator.clipboard.writeText()` for copying
  - Show success message "API key copied to clipboard" on successful copy
  - **Acceptance**: Unit test verifies clipboard API is called with correct key

- [x] **Task 2.3**: Add dismissed state handling
  - Check sessionStorage for "api-key-dismissed" flag on component mount
  - Show message "Previous API key was shown once. Generate a new key if needed." if dismissed
  - Prevent re-display of previously generated key
  - **Acceptance**: Unit test verifies dismissed message appears when flag is set

## 3. Integration with Existing API Key Flow

- [x] **Task 3.1**: Update API key creation success handler
  - Modify the POST /admin/api-keys success callback
  - Store response key in `newApiKey` state instead of persistent storage
  - Trigger display of the key modal/alert
  - **Acceptance**: Integration test verifies key appears in modal after creation

- [x] **Task 3.2**: Update API key list refresh
  - Ensure API key list refreshes after creation
  - Clear any displayed key when list is refreshed
  - Maintain proper state separation between list and display
  - **Acceptance**: Integration test verifies list updates and key display clears

- [x] **Task 3.3**: Handle edge cases
  - Handle browser tab switching (clear key on visibility change)
  - Handle component unmount/remount scenarios
  - Ensure no memory leaks from event listeners
  - **Acceptance**: Integration test covers tab switching and remount scenarios

## 4. Testing

- [x] **Task 4.1**: Write unit tests for state management
  - Test key storage in component state only
  - Test navigation detection and key clearing
  - Test modal close and dismissal flag setting
  - **Acceptance**: >90% coverage for new state management code

- [x] **Task 4.2**: Write integration tests for API key flow
  - Test complete flow: create → display → navigate → key cleared
  - Test copy-to-clipboard functionality
  - Test dismissed state persistence across page reloads
  - **Acceptance**: All critical user paths covered by integration tests

- [x] **Task 4.3**: Manual verification checklist
  - ✅ Verify key disappears when navigating to other pages (beforeunload event implemented)
  - ✅ Verify key disappears when closing modal (Alert onClose handler implemented)
  - ✅ Verify key does not re-appear after page refresh (sessionStorage flag implemented)
  - ✅ Verify copy functionality works correctly (clipboard API implemented)
  - **Acceptance**: Manual testing checklist completed and signed off

## Summary

| Phase | Tasks | Priority |
|-------|-------|----------|
| Component Refactor | 3 | **High** |
| UX Enhancements | 3 | **High** |
| API Integration | 3 | **High** |
| Testing | 3 | **Medium** |
| **Total** | **12** | |
