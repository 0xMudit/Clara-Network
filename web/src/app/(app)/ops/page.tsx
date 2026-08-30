import { notFound } from "next/navigation";
import Link from "next/link";
import { createServerClient } from "@/lib/supabase/server";
import { roleFromAppMetadata } from "@/lib/roles";
import { fetchAdmin } from "@/lib/adminapi";
import { StatCard } from "@/components/cards/stat-card";
import type { DashboardSummary } from "@/types/admin";

export default async function OpsPage() {
  const supabase = await createServerClient();
  const { data } = await supabase.auth.getUser();
  if (roleFromAppMetadata(data.user?.app_metadata) !== "scheme_operator") notFound();
  const d = await fetchAdmin<DashboardSummary>("/dashboard");
  return (
    <div className="grid gap-4">
      <h1 className="text-2xl font-semibold">Operations</h1>
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        <StatCard title="Transactions today" value={d.transactions.toLocaleString()} hint="Authorizations via the switch" />
        <StatCard title="Clearing records" value={d.clearingRecords.toLocaleString()} hint="Captured this settlement window" />
        <StatCard title="Merchants onboarded" value={d.merchants.toLocaleString()} />
      </div>
      <p className="text-sm text-muted-foreground">
        <Link href="/transactions" className="underline underline-offset-4 hover:text-foreground">View transaction log →</Link>
      </p>
    </div>
  );
}
