"use client";

import { HOME_BY_ROLE, DASHBOARD_ACCESS, ROLE_LABEL, type Role } from "@/lib/roles";
import { LogoutButton } from "@/components/logout-button";
import { ThemeToggle } from "@/components/theme-toggle";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { cn } from "@/lib/utils";
import { Landmark } from "lucide-react";
import { useState } from "react";

const TITLES: Record<string, string> = {
  "/overview": "Overview",
  "/ops": "Operations",
  "/issuer": "Issuer",
  "/acquirer": "Acquirer",
  "/merchant": "Merchant",
  "/transactions": "Transactions",
  "/clearing": "Clearing",
  "/settlement": "Settlement",
  "/ledger": "Ledger",
  "/cards": "Cards",
  "/tokens": "Tokens",
  "/merchants": "Merchants",
  "/funding": "Funding",
  "/disputes": "Disputes",
};

export default function Nav({ role }: { role: Role }) {
  const items = DASHBOARD_ACCESS[role] ?? [];
  const pathname = usePathname();
  const [mobileOpen, setMobileOpen] = useState(false);

  return (
    <header className="sticky top-0 z-40 border-b bg-background/80 backdrop-blur-xl supports-[backdrop-filter]:bg-background/60">
      <div className="mx-auto flex h-14 max-w-6xl items-center justify-between px-4">
        {/* Left: Logo + nav */}
        <div className="flex items-center gap-6">
          <Link
            href={HOME_BY_ROLE[role]}
            className="flex items-center gap-2.5 transition-opacity hover:opacity-80"
          >
            <div className="flex size-8 items-center justify-center rounded-lg bg-primary/10 ring-1 ring-primary/20">
              <Landmark className="size-4 text-primary" />
            </div>
            <span className="text-sm font-bold tracking-tight">
              Clara
            </span>
          </Link>

          {/* Desktop nav */}
          <nav className="hidden items-center gap-0.5 md:flex">
            {items.map((p) => {
              const isActive = pathname === p;
              return (
                <Link
                  key={p}
                  href={p}
                  className={cn(
                    "relative rounded-lg px-3 py-1.5 text-sm font-medium transition-colors",
                    isActive
                      ? "text-foreground"
                      : "text-muted-foreground hover:text-foreground hover:bg-muted/50"
                  )}
                >
                  {TITLES[p] ?? p}
                  {isActive && (
                    <span className="absolute inset-x-2 bottom-0 h-0.5 rounded-full bg-primary" />
                  )}
                </Link>
              );
            })}
          </nav>
        </div>

        {/* Right: Role badge, theme, logout */}
        <div className="flex items-center gap-2">
          {/* Role badge */}
          <span className="hidden rounded-full bg-muted/80 px-2.5 py-1 text-[11px] font-medium text-muted-foreground ring-1 ring-border/50 sm:inline-block">
            {ROLE_LABEL[role]}
          </span>
          <ThemeToggle />
          <LogoutButton />

          {/* Mobile hamburger */}
          <button
            className="flex size-8 items-center justify-center rounded-lg text-muted-foreground hover:bg-muted md:hidden"
            onClick={() => setMobileOpen(!mobileOpen)}
            aria-label="Toggle menu"
          >
            <svg
              className="size-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              {mobileOpen ? (
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M6 18L18 6M6 6l12 12"
                />
              ) : (
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M4 6h16M4 12h16M4 18h16"
                />
              )}
            </svg>
          </button>
        </div>
      </div>

      {/* Mobile nav dropdown */}
      {mobileOpen && (
        <nav className="border-t bg-background/95 px-4 py-2 backdrop-blur-xl md:hidden">
          {items.map((p) => {
            const isActive = pathname === p;
            return (
              <Link
                key={p}
                href={p}
                onClick={() => setMobileOpen(false)}
                className={cn(
                  "block rounded-lg px-3 py-2 text-sm font-medium transition-colors",
                  isActive
                    ? "bg-primary/10 text-primary"
                    : "text-muted-foreground hover:bg-muted hover:text-foreground"
                )}
              >
                {TITLES[p] ?? p}
              </Link>
            );
          })}
        </nav>
      )}
    </header>
  );
}
