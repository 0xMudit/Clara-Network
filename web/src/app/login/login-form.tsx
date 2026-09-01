"use client";
import { useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { createBrowserClient } from "@/lib/supabase/client";
import { cn } from "@/lib/utils";
import {
  Landmark,
  CreditCard,
  Store,
  ShoppingBag,
  Eye,
  Loader2,
  ArrowRight,
  Shield,
  Zap,
  Globe,
} from "lucide-react";
import { roleFromAppMetadata, HOME_BY_ROLE, type Role } from "@/lib/roles";

const DEMO_PASSWORD = "ClaraDemo!2026";

const DEMO_ROLES: {
  role: Role;
  label: string;
  email: string;
  blurb: string;
  Icon: typeof Landmark;
  color: string;
  ring: string;
  bg: string;
}[] = [
  {
    role: "scheme_operator",
    label: "Scheme Operator",
    email: "scheme-operator@clara.demo",
    blurb: "Full network view — clearing, settlement, ledger, disputes",
    Icon: Landmark,
    color: "text-violet-500",
    ring: "ring-violet-500/30",
    bg: "bg-violet-500/10",
  },
  {
    role: "issuer",
    label: "Issuer",
    email: "issuer@clara.demo",
    blurb: "Cards, tokens and BIN ranges",
    Icon: CreditCard,
    color: "text-sky-500",
    ring: "ring-sky-500/30",
    bg: "bg-sky-500/10",
  },
  {
    role: "acquirer",
    label: "Acquirer",
    email: "acquirer@clara.demo",
    blurb: "Merchants, funding and disputes",
    Icon: Store,
    color: "text-emerald-500",
    ring: "ring-emerald-500/30",
    bg: "bg-emerald-500/10",
  },
  {
    role: "merchant",
    label: "Merchant",
    email: "merchant@clara.demo",
    blurb: "Funding lines and disputes",
    Icon: ShoppingBag,
    color: "text-amber-500",
    ring: "ring-amber-500/30",
    bg: "bg-amber-500/10",
  },
  {
    role: "viewer",
    label: "Viewer",
    email: "viewer@clara.demo",
    blurb: "Read-only overview, just the highlights",
    Icon: Eye,
    color: "text-rose-500",
    ring: "ring-rose-500/30",
    bg: "bg-rose-500/10",
  },
];

export function LoginForm() {
  const router = useRouter();
  const params = useSearchParams();
  const [selected, setSelected] = useState<(typeof DEMO_ROLES)[number] | null>(
    null
  );
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function signInAs(entry: (typeof DEMO_ROLES)[number]) {
    if (loading) return;
    setSelected(entry);
    setLoading(true);
    setError(null);
    const supabase = createBrowserClient();
    const { data, error } = await supabase.auth.signInWithPassword({
      email: entry.email,
      password: DEMO_PASSWORD,
    });
    if (error) {
      setError(error.message);
      setLoading(false);
      return;
    }
    const role = roleFromAppMetadata(data.session?.user.app_metadata);
    const next = params.get("next");
    router.push(
      next && next.startsWith("/")
        ? next
        : role
          ? HOME_BY_ROLE[role]
          : "/overview"
    );
    router.refresh();
  }

  return (
    <div className="flex min-h-svh">
      {/* Left panel — hero branding */}
      <div className="relative hidden w-1/2 overflow-hidden lg:flex lg:flex-col lg:justify-between bg-gradient-to-br from-[oklch(0.22_0.05_250)] via-[oklch(0.18_0.06_260)] to-[oklch(0.14_0.08_280)] p-10 text-white">
        {/* Decorative orbs */}
        <div className="pointer-events-none absolute -top-32 -left-32 size-96 rounded-full bg-violet-500/20 blur-[100px]" />
        <div className="pointer-events-none absolute -bottom-32 -right-32 size-96 rounded-full bg-sky-500/20 blur-[100px]" />
        <div className="pointer-events-none absolute top-1/2 left-1/2 size-64 -translate-x-1/2 -translate-y-1/2 rounded-full bg-primary/10 blur-[80px]" />

        {/* Grid pattern overlay */}
        <div
          className="pointer-events-none absolute inset-0 opacity-[0.03]"
          style={{
            backgroundImage:
              "linear-gradient(rgba(255,255,255,.1) 1px, transparent 1px), linear-gradient(90deg, rgba(255,255,255,.1) 1px, transparent 1px)",
            backgroundSize: "40px 40px",
          }}
        />

        {/* Top: Logo */}
        <div className="relative z-10 flex items-center gap-3">
          <div className="flex size-10 items-center justify-center rounded-xl bg-white/10 ring-1 ring-white/20 backdrop-blur-sm">
            <Landmark className="size-5" />
          </div>
          <span className="text-lg font-semibold tracking-tight">
            Clara Network
          </span>
        </div>

        {/* Center: Value proposition */}
        <div className="relative z-10 max-w-md">
          <h1 className="text-4xl font-bold leading-tight tracking-tight">
            The payment
            <br />
            network you
            <br />
            <span className="bg-gradient-to-r from-violet-400 via-sky-400 to-emerald-400 bg-clip-text text-transparent">
              build yourself.
            </span>
          </h1>
          <p className="mt-5 text-base leading-relaxed text-white/60">
            Open-source Mastercard / Visa-style card switching, clearing,
            settlement, and issuing — fork it, tweak it, ship your own
            network.
          </p>

          {/* Feature pills */}
          <div className="mt-8 flex flex-wrap gap-2.5">
            {[
              {
                Icon: Shield,
                text: "ISO 8583 + 20022",
              },
              { Icon: Zap, text: "Real-time auth" },
              { Icon: Globe, text: "Self-hosted" },
            ].map(({ Icon, text }) => (
              <span
                key={text}
                className="inline-flex items-center gap-1.5 rounded-full bg-white/10 px-3 py-1.5 text-xs font-medium ring-1 ring-white/10 backdrop-blur-sm"
              >
                <Icon className="size-3.5" />
                {text}
              </span>
            ))}
          </div>
        </div>

        {/* Bottom: Subtle attribution */}
        <p className="relative z-10 text-xs text-white/30">
          Open source · Built for county & bank engineers
        </p>
      </div>

      {/* Right panel — sign-in */}
      <div className="flex flex-1 flex-col items-center justify-center px-6 py-12 sm:px-12 lg:px-16">
        {/* Mobile logo */}
        <div className="mb-8 flex items-center gap-2.5 lg:hidden">
          <div className="flex size-9 items-center justify-center rounded-xl bg-primary/10 ring-1 ring-primary/20">
            <Landmark className="size-4.5 text-primary" />
          </div>
          <span className="text-lg font-semibold tracking-tight">
            Clara Network
          </span>
        </div>

        <div className="w-full max-w-md">
          {/* Header */}
          <div className="mb-8">
            <h2 className="text-2xl font-bold tracking-tight">
              Welcome back
            </h2>
            <p className="mt-1.5 text-sm text-muted-foreground">
              Pick a role and sign in instantly — no password needed for demo.
            </p>
          </div>

          {/* Role cards */}
          <div className="grid gap-2.5">
            {DEMO_ROLES.map((entry) => {
              const isActive = selected?.role === entry.role;
              return (
                <button
                  key={entry.role}
                  type="button"
                  disabled={loading}
                  onClick={() => signInAs(entry)}
                  className={cn(
                    "group/card flex w-full items-center gap-4 rounded-xl border px-4 py-3.5 text-left transition-all duration-200",
                    "hover:border-border hover:bg-muted/50 hover:shadow-sm",
                    "focus-visible:border-ring focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring/50",
                    "disabled:cursor-wait disabled:opacity-60",
                    isActive
                      ? "border-primary/40 bg-primary/5 shadow-sm ring-1 ring-primary/20"
                      : "border-border/60 bg-card"
                  )}
                >
                  {/* Icon */}
                  <div
                    className={cn(
                      "flex size-10 shrink-0 items-center justify-center rounded-xl ring-1 transition-all duration-200",
                      entry.bg,
                      entry.color,
                      entry.ring
                    )}
                  >
                    <entry.Icon className="size-5" />
                  </div>

                  {/* Text */}
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2">
                      <span className="text-sm font-semibold text-foreground">
                        {entry.label}
                      </span>
                      {isActive && loading && (
                        <Loader2 className="size-3.5 animate-spin text-primary" />
                      )}
                    </div>
                    <p className="mt-0.5 text-xs text-muted-foreground line-clamp-1">
                      {entry.blurb}
                    </p>
                  </div>

                  {/* Arrow */}
                  <ArrowRight
                    className={cn(
                      "size-4 shrink-0 text-muted-foreground/40 transition-all duration-200",
                      "group-hover/card:translate-x-0.5 group-hover/card:text-foreground/60",
                      isActive && "text-primary"
                    )}
                  />
                </button>
              );
            })}
          </div>

          {/* Error */}
          {error && (
            <div className="mt-4 rounded-lg border border-destructive/30 bg-destructive/5 px-4 py-3 text-sm text-destructive">
              {error}
            </div>
          )}

          {/* Footer */}
          <p className="mt-6 text-center text-xs text-muted-foreground/60">
            Demo accounts — selecting any role signs you straight in.
          </p>
        </div>
      </div>
    </div>
  );
}
