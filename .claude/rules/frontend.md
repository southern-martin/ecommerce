---
paths:
  - "apps/web/**/*.ts"
  - "apps/web/**/*.tsx"
  - "apps/web/**/*.css"
---
# Frontend Rules (React Web)

## Stack
React 18.3, Vite, TypeScript, Tailwind CSS + Shadcn/ui, TanStack Query v5, Zustand, React Router v6

## Patterns
- Server state: TanStack Query (queries + mutations)
- Client state: Zustand stores
- API client: typed fetch wrapper with auth headers
- Styling: Tailwind utility classes + Shadcn/ui components

## Build
```bash
cd apps/web && npm install && npm run dev   # dev (port 28080)
cd apps/web && npm run build                # production
```
