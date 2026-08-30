# Task 7 Report: shell layout — nav, topbar, role home pages

**Date:** 2026-08-31
**Branch:** `feat/web-ui`
**Commit:** `840d198` — `feat(web): authenticated app shell with role nav`

## What was implemented

Authenticated app shell for all role-scoped pages, implementing the brief verbatim
plus controller rulings D1 (async `createServerClient()`) and D2 (`(app)` route group).

**Files created:**
- `web/src/components/nav.tsx` — topbar nav: role-aware links from `DASHBOARD_ACCESS`,
  brand link to `HOME_BY_ROLE[role]`, `TITLES` label map, `ROLE_LABEL[role]` badge,
  `ThemeToggle`, `LogoutButton`. Verbatim from the plan brief.
- `web/src/app/(app)/layout.tsx` — `(app)` route-group layout. `await createServerClient()`
  (D1), reads `data.user.app_metadata` → `roleFromAppMetadata`, `notFound()` if no role,
  renders `<Nav role>` header + `<main>` content wrapper.
- `web/src/app/(app)/ops/page.tsx` — `scheme_operator` gate, `<h1>Operations</h1>` placeholder.
- `web/src/app/(app)/issuer/page.tsx` — `issuer` gate, `<h1>Issuer</h1>` placeholder.
- `web/src/app/(app)/acquirer/page.tsx` — `acquirer` gate, `<h1>Acquirer</h1>` placeholder.
- `web/src/app/(app)/merchant/page.tsx` — `merchant` gate, `<h1>Merchant</h1>` placeholder.
- `web/src/app/(app)/overview/page.tsx` — no role check (viewers + any authed user), `<h1>Overview</h1>` placeholder.

**Files modified:**
- `web/src/app/page.tsx` — the create-next-app scaffold. Per the task's "keep it simple
  and consistent with the demo" option, this now redirects: authenticated user → their
  role home via `HOME_BY_ROLE[role]`; otherwise → `/login` (belts-and-braces; middleware
  already handles unauthenticated). This change IS noted here per the brief's requirement.

**Dynamic rendering decision (beyond verbatim code, required for a clean build):**
- `web/src/app/(app)/layout.tsx` and `web/src/app/page.tsx` both needed
  `export const dynamic = "force-dynamic"` — the brief anticipated this pattern ("login
  page may need its force-dynamic — already present from Task 5") and the `(app)` layout
  + root page read cookies/auth, so Next.js production build otherwise failed during
  static prerender with `missing env SUPABASE_URL`. This matches the login page's existing
  `force-dynamic`. This is a minimal, consistent addition — no other changes.

## What was tested

Commands run (Windows PowerShell → `cmd /c`):

1. `cmd /c "cd web && npx.cmd tsc --noEmit"` — **passed** (no output).
2. `cmd /c "cd web && npx.cmd next build"` — first run failed prerendering `/acquirer`
   (then `/`) with `missing env SUPABASE_URL`; after adding `force-dynamic` to the `(app)`
   layout and root page: **passed**. Route table from successful build:
   ```
   ┌ ƒ /
   ├ ○ /_not-found
   ├ ƒ /acquirer
   ├ ƒ /api/data/[...path]
   ├ ƒ /issuer
   ├ ƒ /login
   ├ ƒ /merchant
   ├ ƒ /ops
   └ ƒ /overview
   ```
   All shell routes render dynamically (ƒ), confirming the `(app)` group layout applies
   and no page is statically prerendered.
3. `cmd /c "cd web && npm.cmd run lint"` — **passed** (1 pre-existing warning in Task 6's
   `src/app/api/data/[...path]/route.ts`: unused `getEnv`; not part of this task, untouched).

Note: `next lint` is removed in Next 16 (treated the arg as a directory); the repo's
`npm run lint` (= `eslint`) is the correct command.

## Files changed

Created:
- `web/src/components/nav.tsx`
- `web/src/app/(app)/layout.tsx`
- `web/src/app/(app)/ops/page.tsx`
- `web/src/app/(app)/issuer/page.tsx`
- `web/src/app/(app)/acquirer/page.tsx`
- `web/src/app/(app)/merchant/page.tsx`
- `web/src/app/(app)/overview/page.tsx`

Modified:
- `web/src/app/page.tsx` (scaffold → role-aware redirect; N/A's remaining
  create-next-app boilerplate)

Commit: `840d198 feat(web): authenticated app shell with role nav` (8 files,
+109/−68). Only web files were staged; unrelated working-tree files
(`.github/**`, `AGENTS.md`, `docs/28-*.md`, `ROADMAP.md`, etc.) were not touched.

## Self-review findings

- **D1 satisfied:** `await createServerClient()` in `(app)/layout.tsx` and in every
  role-home page (all server components).
- **D2 satisfied:** all five role homes live under `src/app/(app)/`, so the group layout
  applies; URLs are `/ops`, `/issuer`, `/acquirer`, `/merchant`, `/overview` (no `(app)`).
- Role checks match the brief exactly (`scheme_operator`/`issuer`/`acquirer`/`merchant`;
  `/overview` ungated).
- `nav.tsx` uses `DASHBOARD_ACCESS[role] ?? []` per plan; `TITLES` map as-is.
- Placeholders only — no dashboard content (Tasks 8–9).
- Hrefs like `/cards`, `/tokens`, `/funding` etc. will 404 until Tasks 9 creates them —
  acceptable per brief; final review checks no dangling 404s.

## Issues / concerns

- **Known, brief-sanctioned:** nav links to not-yet-created routes (`/transactions`,
  `/clearing`, `/settlement`, `/ledger`, `/cards`, `/tokens`, `/merchants`, `/funding`,
  `/disputes`) 404 until implemented in later tasks.
- None blocking. Everything in scope typechecks, builds, and lints clean.