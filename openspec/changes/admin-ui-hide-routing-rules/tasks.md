## 1. Remove Route Registration

- [ ] 1.1 Remove `<Route path="routing">` from App.tsx

## 2. Remove Sidebar Navigation

- [ ] 2.1 Remove Routing Rules sidebar item from AppShell.tsx
- [ ] 2.2 Remove `/routing` entry from roleAccess in AppShell.tsx

## 3. Verify

- [ ] 3.1 Verify `RoutingRules` import removal in App.tsx (if no longer used elsewhere)
- [ ] 3.2 Verify TypeScript compilation passes
- [ ] 3.3 Verify `BranchesOutlined` icon import is still needed (check other usages)

## 4. OpenSpec Documentation

- [x] 4.1 Create proposal.md
- [x] 4.2 Create design.md
- [x] 4.3 Create tasks.md
