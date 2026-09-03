// src/components/landing/hero.tsx
// Marketing hero for the public landing page. Renders the headline, the
// open-source pitch, the "live right now" status strip, and the role-picker
// CTA. The role cards themselves live in role-cards.tsx (client) so the sign-in
// flow can run here on the same page.
import Link from "next/link";
import { ArrowRight, ShieldCheck, Activity } from "lucide-react";
import { GithubIcon } from "./github-icon";
import { RoleCards } from "./role-cards";

export function Hero() {
  return (
    <section className="relative px-4 pt-28 sm:px-6 sm:pt-32 lg:px-8">
      {/* Status strip */}
      <div className="mb-8 flex justify-center">
        <div className="inline-flex items-center gap-2 rounded-full border border-emerald-500/30 bg-emerald-500/10 px-3 py-1 text-xs font-medium text-emerald-400">
          <span className="relative flex size-2">
            <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-emerald-400 opacity-75" />
            <span className="relative inline-flex size-2 rounded-full bg-emerald-400" />
          </span>
          Live right now · real ISO 8583 traffic
        </div>
      </div>

      {/* Headline */}
      <div className="mx-auto max-w-3xl text-center">
        <h1 className="text-4xl font-bold tracking-tight text-balance sm:text-6xl">
          Open-source card payment network,{" "}
          <span className="bg-gradient-to-r from-violet-400 via-sky-400 to-emerald-400 bg-clip-text text-transparent">
            build it yourself.
          </span>
        </h1>
        <p className="mx-auto mt-5 max-w-2xl text-pretty text-base leading-relaxed text-muted-foreground sm:text-lg">
          Run your own payment network end-to-end — switching, clearing,
          settlement, card issuing and acquiring —{" "}
          <span className="text-foreground">without the license fees</span>{" "}
          or opaque vendor lock-in.
        </p>
        <p className="mx-auto mt-3 max-w-2xl text-pretty text-sm leading-relaxed text-muted-foreground/80">
          A complete Mastercard / Visa-style stack — ISO 8583 switching,
          clearing &amp; settlement, double-entry ledger, issuing, acquiring,
          disputes and instant payments. Fork it, tweak it, ship your own
          network.
        </p>
      </div>

      {/* Primary CTA */}
      <div className="mt-8 flex flex-wrap items-center justify-center gap-3">
        <Link
          href="#roles"
          className="group inline-flex items-center gap-2 rounded-xl bg-primary px-5 py-2.5 text-sm font-semibold text-primary-foreground shadow-lg shadow-primary/25 transition-all hover:-translate-y-0.5 hover:shadow-xl hover:shadow-primary/30"
        >
          Explore the live dashboard
          <ArrowRight className="size-4 transition-transform group-hover:translate-x-0.5" />
        </Link>
        <a
          href="https://github.com/0xMudit/Clara-Network"
          target="_blank"
          rel="noopener noreferrer"
          className="inline-flex items-center gap-2 rounded-xl border px-5 py-2.5 text-sm font-semibold transition-all hover:-translate-y-0.5 hover:bg-muted/50"
        >
          <GithubIcon className="size-4" />
          View source
        </a>
      </div>

      {/* Feature ticks */}
      <div className="mt-6 flex flex-wrap items-center justify-center gap-x-6 gap-y-2 text-xs text-muted-foreground">
        {[
          { Icon: ShieldCheck, text: "Zero cost identity / auth" },
          { Icon: Activity, text: "Real-time authorization" },
          { Icon: GithubIcon, text: "MIT-licensed Go + Next.js" },
        ].map(({ Icon, text }) => (
          <span key={text} className="inline-flex items-center gap-1.5">
            <Icon className="size-3.5 text-primary" />
            {text}
          </span>
        ))}
      </div>

      {/* Live demo + role picker */}
      <div id="roles" className="mx-auto mt-14 max-w-5xl scroll-mt-24">
        <RoleCards />
      </div>
    </section>
  );
}