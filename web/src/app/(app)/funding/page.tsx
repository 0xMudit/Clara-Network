import { notFound } from "next/navigation";
import { createServerClient } from "@/lib/supabase/server";
import { DASHBOARD_ACCESS, roleFromAppMetadata } from "@/lib/roles";
import { fetchAdmin } from "@/lib/adminapi";
import { fmtMinor } from "@/lib/money/minor";
import type { Page } from "@/types/admin";

interface Merchant {
  id: string;
  name: string;
  dba: string;
  taxId: string;
  mccs: string[];
  status: string;
  riskTier: string;
  reserveRateBps: number;
  fundingDelayDays: number;
  transactionLimit: number;
  reserveBalance: number;
  volume: number;
  declineReason?: string;
  approvedAt: string;
}

export default async function FundingPage() {
  const supabase = await createServerClient();
  const { data } = await supabase.auth.getUser();
  const role = roleFromAppMetadata(data.user?.app_metadata);
  if (!role || !DASHBOARD_ACCESS[role].includes("/funding")) notFound();
  const page = await fetchAdmin<Page<Merchant>>("/merchants?limit=50");
  return (
    <div className="grid gap-4">
      <h1 className="text-2xl font-semibold">Acquirer funding</h1>
      <div className="rounded-lg border">
        <table className="w-full text-sm">
          <thead><tr className="border-b text-left text-muted-foreground">
            <th className="px-3 py-2">Name</th><th className="px-3 py-2">Delay (days)</th>
            <th className="px-3 py-2">Txn limit</th><th className="px-3 py-2">Reserve</th><th className="px-3 py-2">Volume</th>
          </tr></thead>
          <tbody>
            {page.items.map(m => (
              <tr key={m.id} className="border-b last:border-0">
                <td className="px-3 py-2">{m.name}</td>
                <td className="px-3 py-2">{m.fundingDelayDays}</td>
                <td className="px-3 py-2">{fmtMinor(m.transactionLimit)}</td>
                <td className="px-3 py-2">{fmtMinor(m.reserveBalance)}</td>
                <td className="px-3 py-2">{fmtMinor(m.volume)}</td>
              </tr>
            ))}
          </tbody>
        </table>
        {page.items.length === 0 && <p className="p-4 text-sm text-muted-foreground">No merchants yet — run the seed (Task 10).</p>}
      </div>
    </div>
  );
}