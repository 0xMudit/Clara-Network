# Task 7 Fix Brief: break authenticated no-role redirect loop

## Bug (reviewer, Important)

`web/src/app/page.tsx` redirects a role-less authenticated user to `/login` (`redirect(role ? HOME_BY_ROLE[role] : "/login")`), while `web/src/app/login/page.tsx` sends role-less authenticated users back to `/` (`redirect(role ? HOME_BY_ROLE[role] : "/")`). With middleware allowing authenticated users through, this is an infinite `/` ↔ `/login` ping-pong for any session whose `app_metadata.role` is absent. Reachable via login-form's no-role fallback to `/overview` (which the `(app)` layout then `notFound()`s — the role-less state is supported in code).

## Fix (make the two fallbacks agree; always terminate next=`-less)

- `web/src/app/page.tsx`: three-state disambiguation —
  - `!data.user` → `redirect("/login")`
  - `user + role` → `redirect(HOME_BY_ROLE[role])`
  - `user + no role` → `redirect("/overview")` (terminal: the `(app)` layout `notFound()`s role-less users — an explicit no-access 404, NOT a loop)
- `web/src/app/login/page.tsx`: change the authenticated-user redirect fallback from `"/"` to `"/overview"` — `redirect(role ? HOME_BY_ROLE[role] : "/overview")`. Both paths must agree.

Do NOT change the `(app)` layout gate (role-less → notFound) and do NOT grant role-less strangers the viewer role as part of this fix — that is a separate decision. Demo users (Task 10) always get an explicit role; this merely turns a hang into a clear 404.

## Constraints

- Edit ONLY the two files above. One commit: `fix(web): break no-role redirect loop on root and login`.
- Verify: `cmd /c "cd web && npx.cmd tsc --noEmit"` and `cmd /c "cd web && npx.cmd next build"` pass.
- Commit ONLY these two web/ files. Leave unrelated working-tree files (`.github/**`, docs etc.) untouched.
- Windows: use `npx.cmd`, never bare `npx`. Do not spawn subagents.

## Report back (under 8 lines): Status, commit SHA, build/tsc result, only-two-files confirmation.