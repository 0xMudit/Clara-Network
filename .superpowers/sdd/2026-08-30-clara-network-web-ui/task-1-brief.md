# Task 1 Brief: Scaffold Next.js app + git branch

## Task Description (from plan)

Create Next.js app scaffold at repo root under `web/`, on feature branch `feat/web-ui` (already created by controller).

**Files:**
- Create: `web/` (full `create-next-app` output) at repo root
- Modify: `web/.gitignore`, `.gitignore`

**Interfaces:**
- Consumes: nothing.
- Produces: runnable Next.js app under `web/`; `web/src` layout (App Router, `--src-dir`).

## Steps

- **Step 1: Create feature branch** — ALREADY DONE by controller (branch `feat/web-ui` exists and is checked out). Skip this step. Instead verify you are on it: `git branch --show-current` should print `feat/web-ui`.

- **Step 2: Scaffold the app**

```bash
npx.cmd create-next-app@latest web --typescript --tailwind --eslint --app --src-dir --import-alias "@/*" --use-npm --yes
```

Expected: `web/` created; `npm run build` works inside `web/`.

IMPORTANT (controller ruling / environment): This Windows PowerShell shell has ExecutionPolicy that blocks `.ps1` shims — never call `npm` or `npx` alone; always use `npm.cmd` and `npx.cmd`. Node is v24.19.0, npm 11.17.0.

- **Step 3: Verify build**

```bash
cmd /c "cd web && npx.cmd tsc --noEmit"
cmd /c "cd web && npx.cmd next build"
```

Expected: both exit 0. NOTE: `next build` may fail before Tailwind/globals are configured if create-next-app output is mid-generate — that's fine; the real gate is `tsc --noEmit` and that create-next-app finished. If `next build` fails for reasons unrelated to task scope, report it in Concerns rather than fixing unrelated scaffold internals.

- **Step 4: Commit**

```bash
git add web/ .gitignore
git commit -m "feat(web): scaffold next.js app"
```

## Global Constraints (binding)

- Node `v24.19.0`; npm via `npm.cmd` (PowerShell ExecutionPolicy blocks `.ps1` shims — never call `npm` alone, always `npm.cmd` / `npx.cmd`).
- No Go on PATH in this shell. Docker `29.7.2` available.
- Brand: **Stripe theme**, dark mode default — but that's Task 2, not here.
- Commit messages in the repo's style use conventional prefixes (e.g. `feat:`, `docs:`).

## Controller Rulings for this task

- Branch creation step is already satisfied by the controller; do not create or switch branches.
- Absolute paths: work in `C:\Users\mudit\Videos\openSourceProjects\Clara-Network`. The `web/` dir must be created at repo root (sibling of `internal/`, `docs/`, `cmd/`, `deploy/`).
- create-next-app may ask interactive questions even with `--yes`; if it prompts, answer with the flags in the command (ts, tailwind, eslint, app router, src dir, `@/*` alias, npm). Try `--yes` first.