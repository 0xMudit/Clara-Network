# Task 6 Fix Brief: correct admin-API mapping in BFF

Two defects in the Task 6 BFF were found in review + controller inspection:

## Defect 1 (Reviewer, Important): dead allowlist entries

`web/src/app/api/data/[...path]/route.ts` lines 18-19 allow `/acquirer`, `/funding`, `/merchant` — none of which exist on the admin API. Real endpoints (internal/adminapi/server.go:35-55) are served as `/api/v1/<path>` where `<path>` is one of:
- dashboard, transactions, clearing/cycles, clearing/records, clearing/positions, settlement/instructions, settlement/prefunds, settlement/default-fund, ledger/accounts, ledger/entries, cards, bin-ranges, tokens, merchants, disputes
- (dynamic sub-routes cards/{ref}, merchants/{id}, merchants/{id}/funding, disputes/{id}, disputes/overdue, disputes/{id}/ratio exist but are EXACT-matched below and not needed by the dashboards)

Fix the ALLOWED_ROUTES map to contain ONLY real, listable endpoints (BFF-facing paths, NO `/api/v1` prefix):

```ts
const ALLOWED_ROUTES: Record<string, string[]> = {
  scheme_operator: [
    "/dashboard", "/transactions", "/clearing/cycles", "/clearing/records",
    "/clearing/positions", "/settlement/instructions", "/settlement/prefunds",
    "/settlement/default-fund", "/ledger/accounts", "/ledger/entries",
    "/cards", "/bin-ranges", "/tokens", "/merchants", "/disputes",
  ],
  issuer: ["/dashboard", "/cards", "/tokens", "/bin-ranges"],
  acquirer: ["/dashboard", "/merchants", "/disputes"],
  merchant: ["/dashboard", "/merchants", "/disputes"],
  viewer: ["/dashboard"],
};
```

Note the controller addition of `viewer: ["/dashboard"]` (roles.ts already defines viewer; without an entry, viewer logs in to `/overview` and every data call 403s). The `/merchants/{id}` and `/disputes/{id}` detail routes are intentionally NOT allowed (exact-match allowlist; detail views are a later-task concern).

## Defect 2 (Controller, Important): missing /api/v1 prefix

`web/src/lib/clara.ts` fetches `${CLARA_API_URL}${path}`. But `CLARA_API_URL` is the admin API BASE URL (e.g. `http://<railway-host>` or `http://localhost:18083`) and the admin API serves under `/api/v1/...` — so every BFF request would 404 upstream. Fix claraFetch to prepend `/api/v1`:

```ts
return fetch(`${CLARA_API_URL}/api/v1${path}`, {
```

`path` is always an allowlisted string already starting with "/". The BFF-facing paths (`/dashboard` etc.) stay WITHOUT the prefix everywhere else (allowlist, UI calls, docs) — the prefix is purely a claraFetch internal detail.

## Constraints

- Only edit the two quoted files (route.ts + clara.ts). One commit: `fix(web): correct admin-API paths and role allowlist in BFF`.
- Windows PowerShell: `cmd /c "cd web && npx.cmd next build"` and `npx.cmd tsc --noEmit` must pass; commit ONLY these two web/ files (do not touch unrelated working-tree files from the other process).
- Do not spawn subagents.

## Report back
Status + commit SHA + verification (build/tsc) + confirmation that ONLY the two files are in the commit.