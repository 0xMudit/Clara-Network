"use client";

import { useTheme } from "next-themes";
import { Moon, Sun } from "lucide-react";
import { useEffect, useState } from "react";

export function ThemeToggle() {
  const { resolvedTheme, setTheme } = useTheme();
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    const id = window.setTimeout(() => setMounted(true), 0);
    return () => window.clearTimeout(id);
  }, []);

  // Avoid hydration mismatch: render a static placeholder until mounted.
  if (!mounted) {
    return (
      <button
        className="flex size-8 items-center justify-center rounded-lg text-muted-foreground"
        aria-label="Toggle theme"
        disabled
      >
        <Sun className="size-4" />
      </button>
    );
  }

  const isDark = resolvedTheme === "dark";
  return (
    <button
      className="flex size-8 items-center justify-center rounded-lg text-muted-foreground hover:bg-muted hover:text-foreground"
      onClick={() => setTheme(isDark ? "light" : "dark")}
      aria-label="Toggle theme"
      title={isDark ? "Switch to light mode" : "Switch to dark mode"}
    >
      {isDark ? <Sun className="size-4" /> : <Moon className="size-4" />}
    </button>
  );
}
