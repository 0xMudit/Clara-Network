import { notFound } from "next/navigation";
import Link from "next/link";
import { createServerClient } from "@/lib/supabase/server";
import { roleFromAppMetadata } from "@/lib/roles";
import { fetchAdmin } from "@/lib/adminapi";
import { StatCard } from "@/components/cards/stat-card";
import type { DashboardSummary } from "@/types/admin";

export default async function IssuerPage() {
  const supabase = await createServerClient();
  const { data } = await supabase.auth.getUser();
  if (roleFromAppMetadata(data.user?.app_metadata) !== "issuer") notFound();
  const d = await fetchAdmin<DashboardSummary>("/dashboard");
  return (
    <div className="grid gap-4">
      <h1 className="text-2xl font-semibold">Issuer</h1>
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        <StatCard title="Cards" value={d.cards.toLocaleString()} />
        <StatCard title="Tokens" value={d.tokens.toLocaleString()} />
      </div>
      <p className="text-sm text-muted-foreground">
        <Link href="/cards" className="underline underline-offset-4 hover:text-foreground">View cards →</Link>
        {" · "}
        <Link href="/tokens" className="underline underline-offset-4 hover:text-foreground">View tokens →</Link>
      </p>
    </div>
  );
}