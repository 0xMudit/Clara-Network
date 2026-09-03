// src/components/landing/role-cards.tsx
// Client role picker for the landing page: one click signs in with a demo role
// and routes to that role's dashboard (honoring a ?next= deep link). Shares
// credentials and sign-in logic with /login via demo-auth.
"use client";

import { useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import Link from "next/link";
import { cn } from "@/lib/utils";
import { ArrowRight, Loader2, Sparkles } from "lucide-react";
import { ROLE_CARDS, signInAs, type RoleCardInfo } from "@/lib/demo-auth";

export function RoleCards() {
  const router = useRouter();
  const params = useSearchParams();
  const [selected, setSelected] = useState<RoleCardInfo | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const nextPath = params.get("next");

  async function handleSignIn(entry: RoleCardInfo) {
    if (loading) return;
    setSelected(entry);
    setLoading(true);
    setError(null);
    const result = await signInAs(entry.role, router, nextPath);
    if (!result.ok) {
      setError(result.error ?? "Sign in failed.");
      setLoading(false);
    }
  }

  return (
    <div>
      <div className="mb-6 text-center">
        <h2 className="inline-flex items-center gap-2 text-2xl font-bold tracking-tight">
          <Sparkles className="size-5 text-primary" />
          Explore it live — no signup
        </h2>
        <p className="mt-1.5 text-sm text-muted-foreground">
          Pick a role to enter the dashboard with pre-seeded demo data.
        </p>
      </div>

      <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
        {ROLE_CARDS.map((entry) => {
          const isActive = selected?.role === entry.role;
          return (
            <button
              key={entry.role}
              type="button"
              disabled={loading}
              onClick={() => handleSignIn(entry)}
              className={cn(
                "group/card flex flex-col rounded-2xl border bg-card p-5 text-left transition-all duration-200",
                "hover:-translate-y-0.5 hover:border-border hover:shadow-lg hover:shadow-primary/5",
                "focus-visible:border-ring focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring/50",
                "disabled:cursor-wait disabled:opacity-70",
                isActive
                  ? "border-primary/40 bg-primary/5 shadow-lg shadow-primary/10 ring-1 ring-primary/20"
                  : "border-border/60",
              )}
            >
              <div className="flex w-full items-start justify-between">
                <div
                  className={cn(
                    "flex size-11 shrink-0 items-center justify-center rounded-xl ring-1",
                    entry.bg,
                    entry.color,
                    entry.ring,
                  )}
                >
                  <entry.Icon className="size-5" />
                </div>
                <div className="text-right">
                  <div className="text-xs font-medium text-muted-foreground">
                    {entry.stat.label}
                  </div>
                  <div className={cn("text-sm font-semibold", entry.color)}>
                    {entry.stat.value}
                  </div>
                </div>
              </div>

              <div className="mt-4 flex-1">
                <div className="flex items-center gap-2">
                  <span className="text-base font-semibold">{entry.label}</span>
                  {isActive && loading && (
                    <Loader2 className="size-4 animate-spin text-primary" />
                  )}
                </div>
                <p className="mt-1 text-sm text-muted-foreground">
                  {entry.blurb}
                </p>
              </div>

              <div className="mt-4 inline-flex items-center gap-1.5 text-sm font-medium text-primary">
                Enter dashboard
                <ArrowRight className="size-4 transition-transform group-hover/card:translate-x-0.5" />
              </div>
            </button>
          );
        })}
      </div>

      {error && (
        <div className="mt-4 rounded-lg border border-destructive/30 bg-destructive/5 px-4 py-3 text-center text-sm text-destructive">
          {error}
        </div>
      )}

      <p className="mt-6 text-center text-xs text-muted-foreground/70">
        Prefer credentials?{" "}
        <Link href="/login" className="font-medium text-foreground underline underline-offset-2 hover:text-primary">
          Sign in directly
        </Link>{" "}
        with a demo account.
      </p>
    </div>
  );
}