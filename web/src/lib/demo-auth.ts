// src/lib/demo-auth.ts
// Shared demo auth data + sign-in behavior for the landing role picker and the
// /login form. Keeps DEMO_PASSWORD, per-role credentials and the sign-in flow
// in ONE place so the two entry points can never drift apart.
"use client";

import type { AppRouterInstance } from "next/dist/shared/lib/app-router-context.shared-runtime";
import {
  Landmark,
  CreditCard,
  Store,
  ShoppingBag,
  Eye,
  type LucideIcon,
} from "lucide-react";
import { createBrowserClient } from "@/lib/supabase/client";
import { roleFromAppMetadata, HOME_BY_ROLE, type Role } from "@/lib/roles";

export const DEMO_PASSWORD = "ClaraDemo!2026";

export const DEMO_EMAILS: Record<Role, string> = {
  scheme_operator: "scheme-operator@clara.demo",
  issuer: "issuer@clara.demo",
  acquirer: "acquirer@clara.demo",
  merchant: "merchant@clara.demo",
  viewer: "viewer@clara.demo",
};

export interface RoleCardInfo {
  role: Role;
  label: string;
  email: string;
  blurb: string;
  Icon: LucideIcon;
  color: string;
  ring: string;
  bg: string;
  stat: { label: string; value: string };
}

export const ROLE_CARDS: RoleCardInfo[] = [
  {
    role: "scheme_operator",
    label: "Scheme Operator",
    email: DEMO_EMAILS.scheme_operator,
    blurb: "Full network view — clearing, settlement, ledger, disputes",
    Icon: Landmark,
    color: "text-violet-500",
    ring: "ring-violet-500/30",
    bg: "bg-violet-500/10",
    stat: { label: "Throughput", value: "1.2M/d" },
  },
  {
    role: "issuer",
    label: "Issuer",
    email: DEMO_EMAILS.issuer,
    blurb: "Cards, tokens and BIN ranges",
    Icon: CreditCard,
    color: "text-sky-500",
    ring: "ring-sky-500/30",
    bg: "bg-sky-500/10",
    stat: { label: "Cards issued", value: "8.4k" },
  },
  {
    role: "acquirer",
    label: "Acquirer",
    email: DEMO_EMAILS.acquirer,
    blurb: "Merchants, funding and disputes",
    Icon: Store,
    color: "text-emerald-500",
    ring: "ring-emerald-500/30",
    bg: "bg-emerald-500/10",
    stat: { label: "Merchants", value: "156" },
  },
  {
    role: "merchant",
    label: "Merchant",
    email: DEMO_EMAILS.merchant,
    blurb: "Funding lines and disputes",
    Icon: ShoppingBag,
    color: "text-amber-500",
    ring: "ring-amber-500/30",
    bg: "bg-amber-500/10",
    stat: { label: "Settled today", value: "$23.4k" },
  },
  {
    role: "viewer",
    label: "Viewer",
    email: DEMO_EMAILS.viewer,
    blurb: "Read-only overview, just the highlights",
    Icon: Eye,
    color: "text-rose-500",
    ring: "ring-rose-500/30",
    bg: "bg-rose-500/10",
    stat: { label: "Read-only", value: "live" },
  },
];

export type SignInOutcome = {
  ok: boolean;
  error?: string;
};

/**
 * Sign in as a demo role. Handles the Supabase password sign-in and routes to
 * the role's dashboard on success. `nextPath` (from a ?next= param) wins over
 * the role default when it is a safe, same-origin path.
 */
export async function signInAs(
  role: Role,
  router: AppRouterInstance,
  nextPath?: string | null,
): Promise<SignInOutcome> {
  const supabase = createBrowserClient();
  const { data, error } = await supabase.auth.signInWithPassword({
    email: DEMO_EMAILS[role],
    password: DEMO_PASSWORD,
  });
  if (error) return { ok: false, error: error.message };

  const signedInRole = roleFromAppMetadata(data.session?.user.app_metadata);
  const target =
    nextPath && nextPath.startsWith("/")
      ? nextPath
      : signedInRole
        ? HOME_BY_ROLE[signedInRole]
        : "/overview";
  router.push(target);
  router.refresh();
  return { ok: true };
}
