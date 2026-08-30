import { notFound } from "next/navigation";
import Link from "next/link";
import { createServerClient } from "@/lib/supabase/server";
import { roleFromAppMetadata } from "@/lib/roles";
import { fetchAdmin } from "@/lib/adminapi";
import { StatCard } from "@/components/cards/stat-card";
import type { DashboardSummary } from "@/types/admin";

export default async function MerchantPage() {
  const supabase = await createServerClient();
  const { data } = await supabase.auth.getUser();
  if (roleFromAppMetadata(data.user?.app_metadata) !== "merchant") notFound();
  const d = await fetchAdmin<DashboardSummary>("/dashboard");
  return (
    <div className="grid gap-4">
      <h1 className="text-2xl font-semibold">Merchant</h1>
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        <StatCard title="Merchants" value={d.merchants.toLocaleString()} />
        <StatCard title="Disputes" value={d.disputes.toLocaleString()} />
      </div>
      <p className="text-sm text-muted-foreground">
        <Link href="/funding" className="underline underline-offset-4 hover:text-foreground">View funding →</Link>
        {" · "}
        <Link href="/disputes" className="underline underline-offset-4 hover:text-foreground">View disputes →</Link>
      </p>
    </div>
  );
}