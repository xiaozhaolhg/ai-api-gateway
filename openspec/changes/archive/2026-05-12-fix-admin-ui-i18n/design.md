## Context

The admin UI is built with React 19.2 + Vite and currently has an i18n system that is not properly configured. Translation keys like "dashboard.title" are being displayed instead of resolved text, indicating missing translation files or improper i18n initialization. The UI is otherwise functional and running on port 3000.

## Goals / Non-Goals

**Goals:**
- Implement a working i18n system that resolves translation keys to meaningful text
- Create comprehensive translation files for all UI elements
- Ensure proper fallback handling for missing translations
- Maintain existing UI functionality while fixing translations

**Non-Goals:**
- Adding new UI features or components
- Changing the overall UI design or layout
- Implementing multi-language support beyond English (initially)

## Decisions

**i18n Library Choice**: Use `react-i18next` as it's the standard for React i18n and integrates well with Vite
- Alternative considered: Custom i18n solution - rejected due to complexity and maintenance overhead
- Alternative considered: FormatJS - rejected due to steeper learning curve

**Translation File Structure**: Organize translations by feature/page rather than by component
- Rationale: More intuitive for content management and easier for non-developers to understand
- Alternative: Component-based structure - rejected as it fragments related content

**Translation Key Format**: Use nested keys with dot notation (e.g., `dashboard.title`, `navigation.providers`)
- Rationale: Matches current stub keys, minimizing code changes
- Alternative: Flat keys - rejected as it's less organized and harder to scale

## Risks / Trade-offs

**Risk**: Missing translation keys may cause runtime errors
- Mitigation: Implement proper fallback mechanisms and comprehensive testing

**Risk**: Performance impact from loading translation files
- Mitigation: Use code splitting and lazy loading for translation files

**Trade-off**: Initial implementation will be English-only
- Justification: Focus on fixing core functionality first, multi-language support can be added later

## Migration Plan

1. Install and configure `react-i18next` and dependencies
2. Create translation files with English text for all identified keys
3. Update i18n initialization code to properly load translations
4. Test all UI pages to ensure translations are resolved correctly
5. Add error handling for missing translation keys

## Open Questions

- Should we implement language detection based on browser settings?
- Do we need to support right-to-left languages in the future?
