# Task 2 Report: shadcn/ui + Stripe theme + theme switcher

## Status: DONE_WITH_CONCERNS

## What was implemented

1. **shadcn/ui initialized** in `web/` (Tailwind v4 mode → plain CSS variables, no `tailwind.config.ts`).
   - `web/components.json` (style `base-nova`, base color `neutral`, cssVariables `true`, aliases `@/components`, `@/components/ui`, `@/lib`, `@/lib/utils`).
   - `web/src/lib/utils.ts` (`cn` helper).
   - `web/src/app/globals.css` rewritten with shadcn design tokens (`:root`, `.dark`) + `@layer base`.

2. **tweakcn theme system installed** (registry `nextjs/theme-system`):
   - `web/src/lib/themes-config.ts` — `themes[]` (incl. `stripe` entry), `allThemeValues` (`{name}-light`/`{name}-dark`), `sortedThemes`, `themeNames`, `DEFAULT_THEME`.
   - `web/src/components/providers/theme-provider.tsx` — uses `attribute="data-theme"`, `themes={allThemeValues}`, `enableSystem={false}`, `disableTransitionOnChange`.
   - `web/src/components/theme-switcher.tsx` — full dropdown (mode + color theme incl. Stripe).
   - `web/src/components/ui/{button,dropdown-menu,scroll-area}.tsx` (Base UI).
   - `web/src/styles/themes/*.css` (all themes incl. `stripe.css`).

3. **Stripe theme** (`theme-stripe` registry) — idempotently overwrote/confirmed `web/src/styles/themes/stripe.css` with `[data-theme="stripe-light"]` / `[data-theme="stripe-dark"]` (Stripe brand OKLCH tokens).

4. **Default to `stripe-dark`**:
   - `web/src/lib/themes-config.ts`: `DEFAULT_THEME = "stripe-dark"` (was `default-dark`) + added `export const defaultTheme = "stripe-dark"`.

5. **Theme provider + toggle** added:
   - `web/src/components/theme-provider.tsx` — re-exports the registry provider (kept registry output, which is what wires `data-theme`/`allThemeValues`).
   - `web/src/components/theme-toggle.tsx` — plan snippet adapted to tweakcn scheme: `resolvedTheme?.endsWith("-dark")`, toggles `{colorTheme}-light`/`-dark`, uses `dark:` variant classes.

6. **Root layout wired** (`web/src/app/layout.tsx`):
   - Wrapped `<body>` children in `<ThemeProvider>`; added `suppressHydrationWarning` on `<html>`; metadata export preserved (unchanged).

7. **Dark-mode integration fix** (controller ruling): changed `@custom-variant dark` in `globals.css` from `(&:is(.dark *))` to `(&:is([data-theme$="-dark"] *))` so `dark:` Tailwind classes respond to tweakcn's `data-theme` scheme. Also imported the theme CSS (`@import "../styles/themes/index.css";`) and fixed a self-referential `--font-sans` in `@theme inline` to reference the Geist variable.

## Verification (exact commands + output)

- `cmd /c "cd web && set npm_config_allow_scripts= && npx.cmd shadcn@latest init -y -b base -p nova -f --no-reinstall"` → **exit 0**
- `cmd /c "cd web && set npm_config_allow_scripts= && npx.cmd shadcn@latest add https://tweakcn-picker.vercel.app/r/nextjs/theme-system.json"` → **exit 0** (created 49 files)
- `cmd /c "cd web && set npm_config_allow_scripts= && npx.cmd shadcn@latest add -y https://tweakcn-picker.vercel.app/r/theme-stripe.json"` → **exit 0** (skipped identical stripe.css)
- `cmd /c "cd web && npx.cmd next build"` → **exit 0**, "Compiled successfully", TypeScript OK, 2 static routes generated.
- `cmd /c "cd web && npm.cmd run lint"` → **exit 0** after fixing the registry theme-switcher's `react-hooks/set-state-in-effect` (applied `useSyncExternalStore` mounting guard, behavior preserved).
- Confirmed `stripe-dark`/`stripe-light` tokens present in production CSS via a recursive scan of `.next` (Turbopack emits CSS under `static/chunks/*.css`).

## Environment adaptation (npm install failure)

The brief/AGENTS said `web/.npmrc` has `allow-scripts=true`. The newer **shadcn v4 CLI + npm 11** combination causes shadcn's internal `npm install` to fail with `EALLOWSCRIPTS` (`--allow-scripts is not allowed in project-scoped installs`) because npx propagates `npm_config_allow_scripts=true` into the environment, which npm's `resolveAllowScripts` treats as the CLI layer.

**Workaround (no repo config change):** prefixing `set npm_config_allow_scripts=` before the shadcn invocation prevents the injected env var, letting shadcn's npm install succeed. `.npmrc` itself was left untouched.

**CLI mapping:** shadcn v4 removed `--base-color`; component library is now `-b/--base` (base | radix | aria) and design preset is `-p/--preset` (nova = Lucide/Geist, aligned with brand). Selected Base UI + Nova.

## Files changed (beyond the plan's enumerated list — all registry/CLI output)

Plan enumerated: `components.json`, `theme-provider.tsx`, `theme-toggle.tsx`, `globals.css`, `layout.tsx`, `package.json`, plus expected `themes-config.ts` and `lib/utils.ts`.

Registry/CLI added beyond that:
- `src/components/ui/{button,dropdown-menu,scroll-area}.tsx`
- `src/components/theme-switcher.tsx`
- `src/components/providers/theme-provider.tsx`
- `src/styles/themes/*.css` (43 theme files incl. `index.css`, `stripe.css`)
- `package-lock.json`
- deps added to `package.json`: `@base-ui/react`, `class-variance-authority`, `clsx`, `lucide-react`, `next-themes`, `shadcn`, `tailwind-merge`, `tw-animate-css`

Hand-edits (beyond the plan's listed modify set):
- `src/styles/themes/index.css` — removed dangling `@import "./default.css";` (registry shipped the import but never shipped `default.css`, which broke the production build). Fix = drop the dangling import; `default` theme stays in themes-config data. No OKLCH values hand-written.
- `src/components/theme-switcher.tsx` — minimal lint fix (mounted guard via `useSyncExternalStore` instead of `setState` in effect), preserving exact behavior.
- `src/lib/themes-config.ts` — also added `defaultTheme` export per brief Step 3.

## Self-review findings

- **Dark default works**: `DEFAULT_THEME = "stripe-dark"` → next-themes sets `data-theme="stripe-dark"` → Stripe dark tokens apply.
- **Stripe selectable in both modes**: `stripe-light`/`stripe-dark` are in `allThemeValues` (themes list passed to provider) + the theme-switcher lets users pick Stripe and Light/Dark.
- **`dark:` Tailwind selector** now responds to tweakcn `data-theme` scheme via the custom-variant change.
- **Toggle** dark icons render correctly under the `{name}-dark` scheme.
- **Build + lint + typecheck** all pass.
- Confirmed no secrets/config leaks; only `web/` was committed.

## Issues / concerns

1. **CI file out of scope**: `.github/workflows/ci.yml` was modified by a **concurrent process** (not me) during this session (added Go backend lint/docker-smoke/coverage jobs). It remains **unstaged** and was excluded from my commit. Flagging so the plan/controller is aware it's pending.
2. **Registry quirk — `default.css`**: index.css imported a file the registry never shipped; removed the import. The "Default" color theme in the switcher still renders (falls back to `:root` neutral tokens) but has no dedicated `[data-theme="default-*"]` block. Not part of requirements (Stripe is default + required selectable).
3. **Registry prop difference vs plan**: the registry's `ThemeProvider` hardcodes `attribute="data-theme"` (required for Stripe) and doesn't accept props, so Step 5's literal `attribute="class"` was intentionally not used — it would break the Stripe theme.
4. **Font**: Stripe theme CSS sets `--font-sans` to a system-ui stack; Geist variable restored for non-Stripe/neutral contexts. Brand ("Geist/system-ui") satisfied.
5. **Interactive CLI**: shadcn v4 init/add prompt on TTY and don't reliably read piped stdin; resolved by using non-interactive flags (`-b`, `-p`, `--no-reinstall`, `-y`).

## Commit

- `5bd97ac` feat(web): add shadcn ui with stripe theme (57 web/ files)
