# Task 2 Brief: shadcn/ui + Stripe theme + theme switcher

## Task Description (from plan)

Init shadcn/ui in the scaffolded `web/` app, install the tweakcn theme system + Stripe theme, add a theme provider + toggle, default to `stripe-dark`.

**Files:**
- Create: `web/components.json` (via init), `web/src/components/theme-provider.tsx`, `web/src/components/theme-toggle.tsx`, `web/src/app/globals.css` (updated)
- Modify: `web/src/app/layout.tsx`, `web/package.json` (deps via CLI)

**Interfaces:**
- Consumes: Task 1 scaffold (Next 16.3.3, React 19, Tailwind v4, `@/*` alias, src-dir).
- Produces: `<ThemeProvider>` and `<ThemeToggle>` components usable by all layouts; `--primary/*` OKLCH tokens in `globals.css`; `web/src/lib/themes-config.ts` with a Stripe entry.

## Steps

- **Step 1: Init shadcn**

```bash
cmd /c "cd web && npx.cmd shadcn@latest init --base-color neutral --yes"
```

Expected: `components.json`, `src/lib/utils.ts`, `components/ui` namespace configured. In Tailwind v4 mode shadcn writes plain CSS variables into `globals.css` (no `tailwind.config.ts`). Prompts may still appear despite `--yes`; answer defaults (style: "new-york", base color: neutral, CSS variables: yes).

- **Step 2: Add tweakcn theme system + Stripe theme**

```bash
cmd /c "cd web && npx.cmd shadcn@latest add https://tweakcn-picker.vercel.app/r/nextjs/theme-system.json"
cmd /c "cd web && npx.cmd shadcn@latest add https://tweakcn-picker.vercel.app/r/theme-stripe.json"
```

Expected: `src/lib/themes-config.ts` with a `{name,title,colors,fontSans}` entry for `stripe` (light+dark), and theme classes/CSV wiring for `next-themes`. If the registry command fails or prompts interactively, check what was added under `src/lib/` and `src/app/globals.css` and adapt: the REQUIREMENT is Stripe theme present as a selectable theme in both modes.

- **Step 3: Default to stripe-dark**

In `src/lib/themes-config.ts`: `export const defaultTheme = "stripe-dark";` (add if the registry file lacks a default export; keep the file's existing shape).

- **Step 4: Add theme-provider + toggle**

`src/components/theme-provider.tsx`:
```tsx
"use client";
import * as React from "react";
import { ThemeProvider as NextThemesProvider } from "next-themes";

export function ThemeProvider({ children, ...props }: React.ComponentProps<typeof NextThemesProvider>) {
  return <NextThemesProvider {...props}>{children}</NextThemesProvider>;
}
```

`src/components/theme-toggle.tsx`:
```tsx
"use client";
import { Moon, Sun } from "lucide-react";
import { useTheme } from "next-themes";
import { Button } from "@/components/ui/button";

export function ThemeToggle() {
  const { resolvedTheme, setTheme } = useTheme();
  return (
    <Button variant="outline" size="icon" onClick={() => setTheme(resolvedTheme === "dark" ? "light" : "dark")}>
      <Sun className="h-4 w-4 dark:hidden" />
      <Moon className="hidden h-4 w-4 dark:block" />
    </Button>
  );
}
```

Add `next-themes` + `lucide-react` if not auto-added by the registry: `cmd /c "cd web && npm.cmd i next-themes lucide-react"`.

IMPORTANT — theme naming: the tweakcn theme system keys themes as `{name}-light` / `{name}-dark` (e.g. `stripe-dark`). If the toggle uses `resolvedTheme === "dark"` but next-themes resolves to `stripe-dark`, the dark icon branch never shows. Fix accordingly: check `resolvedTheme?.endsWith("-dark")` (preferred — matches tweakcn docs `isDark = theme?.endsWith("-dark")`) so the toggle and `dark:` Tailwind classes both work under the tweakcn scheme.

- **Step 5: Wire provider into root layout**

`src/app/layout.tsx`: wrap `<body>` children in `<ThemeProvider attribute="class" defaultTheme="stripe-dark" enableSystem={false}>`. Root html: `suppressHydrationWarning`.

- **Step 6: Verify**

```bash
cmd /c "cd web && npx.cmd next build"
```

Expected: exit 0.

- **Step 7: Commit**

```bash
git add web/
git commit -m "feat(web): add shadcn ui with stripe theme"
```

## Global Constraints (binding)

- Node v24.19.0; npm via `npm.cmd`/`npx.cmd` (never bare `npm`/`npx`). Windows PowerShell.
- Brand: **Stripe theme**, dark mode default, toggle to light. Geist/system-ui font.
- Do not break the existing scaffold (Keep the App Router layout; `web/src/app/layout.tsx` currently renders directly into body — preserve metadata export).
- No heavy additions beyond the task's deps (shadcn/ui + next-themes + lucide-react + the tweakcn registry output).

## Controller Rulings for this task

- Toolchain is Next 16.3.3 + Tailwind v4 (scaffold default); shadcn init will run in v4 mode (CSS vars, no config file). Plan assumed v4-style output already; proceed with v4 paths.
- The tweakcn theme system may generate files/classes the plan didn't enumerate (e.g. a `<ThemeProvider>` or theme classes in globals.css). That is expected registry output — keep it minimal and consistent; do NOT delete registry-provided files that the Stripe theme needs.
- If the tweakcn theme uses the `{name}-dark` scheme, adapt the toggle per the IMPORTANT note above so the dark icon state and Tailwind `dark:` selector behave under the actual resolved theme name.
- Exact numeric values (colors) come from the registry, not from the plan — do not hand-write OKLCH values.