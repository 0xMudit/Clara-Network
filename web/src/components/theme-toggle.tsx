"use client";

import { Moon, Sun } from "lucide-react";
import { useTheme } from "next-themes";

import { Button } from "@/components/ui/button";

// Tweakcn themes use the "{name}-dark" / "{name}-light" scheme, so next-themes
// resolves to values like "stripe-dark". Detect dark via the "-dark" suffix and
// toggle between the current color theme's light / dark variants.
export function ThemeToggle() {
  const { resolvedTheme, setTheme } = useTheme();

  const isDark = resolvedTheme?.endsWith("-dark");
  const colorTheme = resolvedTheme?.replace(/-(light|dark)$/, "") || "stripe";

  const toggle = () =>
    setTheme(`${colorTheme}${isDark ? "-light" : "-dark"}`);

  return (
    <Button
      variant="outline"
      size="icon"
      aria-label="Toggle theme"
      onClick={toggle}
    >
      <Sun className="h-4 w-4 dark:hidden" />
      <Moon className="hidden h-4 w-4 dark:block" />
    </Button>
  );
}
