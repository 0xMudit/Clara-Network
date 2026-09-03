# SDD ledger — plan: docs/superpowers/plans/2026-08-30-clara-network-web-ui.md

BASE (pre-dispatch Task 1): `48d18a4`

## Pre-flight scan

Read plan + spec (repo README, adminapi server.go, deploy/schema.sql).
Global constraints applicable: Node v24 + npm via `npm.cmd`; no Go on PATH;
Docker available; Stripe theme dark-first; role vocab; admin API read-only,
never proxied straight from browser; BFF returns 401/403/JSON; no PAN/log
secret leakage; amounts as minor units via `fmtMinor`; demo currency EUR.

Task couples (same file / interface):
- T1 (scaffold web/) → T2 (shadcn + stripe) shares `web/`; T1 produces the
  app, T2 adds components/theme. No conflict found.
- T2 (globals.css, layout) → T7 (layout.tsx root) — T7 says "Modify:
  web/src/app/layout.tsx (not needed — keep root)". Plan text self-consistent.
- T3 (env.ts, money/date, types) → T4 (supabase server + roles) consumes
  getEnv(); T6 consumes getEnv() + fetchAdmin; T8 consumes fmtMinor/fmtTs.
  Consistent names: getEnv(), fmtMinor(), fmtTs(), Page<T>, DashboardSummary.
- T4 (roles.ts) → T5 (login form + logout) consumes roleFromAppMetadata,
  HOME_BY_ROLE; T5 login has `await import("@/lib/env").then(m=>m.getEnv())`
  in a client component — `process.env` is not available in browser bundles
  for arbitrary keys. RULING below.
- T6 (api/data/[...path]/route.ts + adminapi.ts + adminapi.test.ts) →
  T8/T9 server pages consume fetchAdmin<T>. Consistent.
- T6 test file uses `@/lib/adminapi` import with `tsx` and sets a
  `globalThis.__calls` bag; tsx must resolve `@/` alias (needs tsconfig
  paths support) — see Ruling.
- T10 (seed + CLI login) → T11 (deploy) — sequential, no conflict.

Per-task self-consistency:
- T6 route uses `PUBLIC_PATHS` + `q` unused variable — minor, tsx/eslint
  will flag; leave for implementer, note below.
- T5/T6 client components reading `@/lib/env`: browser code cannot access
  arbitrary `process.env`. The plan's login form and logout button both do
  `await import("@/lib/env").then(m => m.getEnv())` inside client
  components. In Next.js, only statically-referenced `NEXT_PUBLIC_*`
  vars are inlined to the client. This is a plan defect. Ruling:
  make `getEnv()` server-only (throw if called in browser is awkward); and
  for client components use a dedicated `NEXT_PUBLIC_SUPABASE_URL` /
  `NEXT_PUBLIC_SUPABASE_ANON_KEY` pair. Vercel deploy must set both. I will
  adjust T3/T5/T6/T11 env wiring accordingly (supabase-js anon key is public
  by design, safe for client).
- T6 adminapi.test.ts: top-level `import` after setting env is fine (ESM
  hoisting: the import runs before any assertions; `fetch` is overridden
  before fetchAdmin is called in the test body, but the module import of
  fetchAdmin happens at load — harmless since only reads env at call time).
  The `globalThis.__calls = calls` line is dead — swallow it. RULING noted.

Rulings so far:
- Ruling: client components must use NEXT_PUBLIC_SUPABASE_URL /
  NEXT_PUBLIC_SUPABASE_ANON_KEY, not getEnv(). Cambia env.ts to export
  `getEnv()` (server) and add `getPublicEnv()` or NEXT_PUBLIC constants for
  client. Cost if wrong: low — an extra env var pair; correct-by-default.

Host notes:
- Model selection: harness Task tool exposes only `general`/`explore`
  agent types — cannot pick per-task model tiers. All subagents run on the
  session default model. Ruling: accept session model for everything;
  escalate via capability if a task blocks.

Worktree: no isolated worktree (user asked to build in main workspace on
branch feat/web-ui). Ruling: proceed on current checkout on feat/web-ui.

## Task state
- Task 1: not started
- Task 2: not started
- Task 3: not started
- Task 4: not started
- Task 5: not started
- Task 6: not started
- Task 7: not started
- Task 8: not started
- Task 9: not started
- Task 10: not started
- Task 11: not started
Task 1: complete (commits 48d18a4..872db5a, review clean)

Ruling: create-next-app@latest produced Next 16.3.3 (+ React 19, Tailwind v4, ESLint 9) instead of the plan's assumed Next 15. App Router + src-dir layout identical; plan's TS/@/*/Tailwind constraints all satisfied. Accept 16.3.3 as toolchain. Cost if wrong: minimal � no plan step depended on a Next 15-specific API.

Task 2: complete (commits 872db5a..5bd97ac, review clean)
Task 2: minor (deferred, final review): dead .dark block in globals.css:101-133; "Default" theme lacks default.css (breaks default-dark fallback color); 42 non-Stripe themes bundled (~220KB, trim later); per-theme fontSans inert (Geist everywhere); defaultTheme export unused; no layout smoke test for -dark suffix.

Task 3: complete (commits 5bd97ac..4c2c60a, review clean)
Task 3: minor (deferred, final review): web/.gitignore !.env.local.example is a no-op (pattern .env*.local already excludes it) and the narrowed pattern no longer protects .env, .env.production etc. Task 11 uses \ercel env pull .env.production\ � assess then whether to broaden the ignore pattern before any production env file is written.

Ruling: createServerClient() in src/lib/supabase/server.ts is now async (Next 16 async cookies()) and the supabase factory import is aliased. Consumers in Tasks 5/6/7 must await createServerClient(). This deviates from the plan's sync signature due to the Next 16 toolchain ruling. Also: middleware.ts triggers Next 16 deprecation warning (favors proxy.ts) � keep middleware for now, revisit in final review. Cost if wrong: a missed await is a build error, caught immediately.

Task 4: complete (commits 4c2c60a..7149a29, review clean)
Task 4: minor (deferred): trailing newlines missing on 4 new source files; /login?next= drops query string (plan-mandated); dashboardAccess ?? [] dead-safe; middleware.ts Next 16 deprecation (decision deferred).

Task 5: complete (commits 7149a29..adfa60b: d70c688 impl + adfa60b middleware-loop hotfix, review clean)
Ruling: task-5 deviated from plan text: (1) client.ts factory new file with aliased createSupabaseBrowserClient import (name-shadowing fix); (2) login page needs export const dynamic = \"force-dynamic\" for next build; (3) middleware fixed to exempt /login from redirect (infinite loop). All three approved inline.
Task 5: minor (deferred): /login exemption uses startsWith (slightly broad); React.FormEvent uses UMD-global style without import.

Task 6: complete (commits adfa60b..390ee16, review APPROVED-BUT-FIXED)
Task 6: fix commit 390ee16 corrects (1) acquirer/merchant dead allowlist entries (/acquirer /funding /merchant removed; real list endpoints used; viewer now allowed /dashboard), (2) claraFetch now prepends /api/v1 (CLARA_API_URL is a bare base URL per plan). Ruling: BFF-facing paths are prefix-free (/dashboard, /clearing/records); /api/v1 is a claraFetch internal detail. Exact-match allowlist intentionally excludes dynamic sub-routes (/cards/{ref}, /merchants/{id}, /disputes/{id}) � detail views deferred.
Task 6: minor (deferred): unused getEnv import in route.ts; trailing newline EOF on route.ts/clara.ts.

Task 7: complete (commits 390ee16..157c32a: 840d198 impl + 157c32a loop fix, review APPROVED-BUT-FIXED)
Task 7: fix 157c32a breaks /<->/login ping-pong for role-less authenticated users: root 3-state (no user->/login, role->home, no role->/overview) and login fallback /->/overview; (app) layout keeps role-less->notFound. ruling: do NOT grant role-less users viewer implicitly at this stage.
Task 7: minor (deferred): overview page doesn't await server client (no check needed, fine); trailing newlines missing on new files; force-dynamic added to (app)/layout and root page (justified). Nav links to /cards /funding /tokens etc. 404 until Task 9 (sanctioned).

Task 8: complete (commits 157c32a..0f711eb, review clean)
Task 8: incidental fix to Task 3 defect: fmtMinor produced \"- 50.00 EUR\" (stray space); new Task 8 money test exposed it; fixed to sign-concatenated \"-50.00 EUR\" (template-only change, math untouched). NOTE: Task 3 reviewer incorrectly verified the buggy output as correct � test caught what review missed.
Task 8: minor (deferred): ops re-fetches /dashboard (permitted); STAN-as-key fine for demo.

Task 9: complete (commits 0f711eb..15397cc: cca1bfb BFF query-forward fix + 15397cc dashboards, review clean)
Task 9: resolving note: brief Step 6/9 gates were wrong for merchant (DASHBOARD_ACCESS.merchant lacks /merchants); implementer gated /merchants to scheme_operator|acquirer and linked merchant home to /funding+/disputes. No dangling nav: all DASHBOARD_ACCESS paths resolve + have TITLES.
Task 9: minor (deferred, fold into final cleanup): disputes responseDue zero-time renders '1 Jan 0001' (seed always sets it; demo unaffected); funding page exposes merchant reserve data to merchant role (brief-mandated, demo-ok); unused getEnv import in route.ts (Task 6 deferred item).

Task 10: implemented, commit 7b4413a (base 15397cc..7b4413a). REVIEW PENDING - review dispatch was cancelled at session end; re-review task-10 before continuing. task-10-report.md + task-10-review-package.md on disk.
Task 10: implementer deviations to verify in review: override adds postgres: ports: !reset [] (compose-local PG never binds 5432); one-shot runs use --no-deps; teardown also stops dependency postgres container (not volumes); bash via Git Bash (WSL broken).
Task 10: REVIEW COMPLETE (re-review by controller) - clean. Files match brief: seed-docker.ps1 (PS5.1, !reset [], --no-deps, teardown), seed.sh (bash twin), create-users.mjs (idempotent admin API), demo-profile.sql (doc manifest), README runbook, db:seed/db:users scripts. Deviations verified OK (postgres ports reset; --no-deps; postgres teardown stop; Git Bash for bash -n). No live services touched (per brief).

Task 11 (deploy) + final whole-branch review: in progress (controller).
Task 11: prep done (no credentials required): web/vercel.json created (framework nextjs, build/install commands); package.json deploy:prod/deploy:preview scripts added; build passes (next build, 17 routes + Proxy); tsc clean; eslint clean (removed unused getEnv import in api/data/[...path]/route.ts). NOTE: vercel.json lives inside web/ (deploy runs from web/, so web/ is the Vercel project root). Deferred: middleware->proxy Next 16 migration (build warning only, works in prod).
Task 11: NO cloud credentials present (VERCEL_TOKEN/SUPABASE_ACCESS_TOKEN/RAILWAY_TOKEN/env vars all unset; no ~/.vercel, ~/.supabase, no Railway CLI). Actual cloud push (vercel --prod, railway up, supabase db push, db:users, db:seed against live Supabase) requires the repo owner's interactive CLI logins / tokens -> NOT performed by controller. Handing off explicit runbook.
Task 11: DEPLOY BLOCKER flagged: internal/adminapi is unauthenticated + wildcard CORS (Access-Control-Allow-Origin:*). BFF sends Authorization: Bearer <user.id> but adminapi ignores it. Fine for demo-by-design (BFF is the auth boundary); but if the Railway instance is publicly reachable, anyone can read all network data unauthenticated. Recommend restricting access until/or documenting as demo-only. Not changed (out of scope for web deploy; architecture decision for owner).

-- END OF SESSION 2026-08-31 (controller: check+complete+deploy) --
