// src/app/(app)/overview/page.tsx
import { fetchAdmin } from "@/lib/adminapi";
import { StatCard } from "@/components/cards/stat-card";
import type { DashboardSummary } from "@/types/admin";

export default async function OverviewPage() {
  const d = await fetchAdmin<DashboardSummary>("/dashboard");
  return (
    <div className="grid gap-4">
      <h1 className="text-2xl font-semibold">Network overview</h1>
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        <StatCard title="Transactions" value={d.transactions.toLocaleString()} />
        <StatCard title="Clearing records" value={d.clearingRecords.toLocaleString()} />
        <StatCard title="Merchants" value={d.merchants.toLocaleString()} />
        <StatCard title="Disputes" value={d.disputes.toLocaleString()} />
        <StatCard title="Cards" value={d.cards.toLocaleString()} />
        <StatCard title="Tokens" value={d.tokens.toLocaleString()} />
      </div>
    </div>
  );
}