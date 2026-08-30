import { notFound } from "next/navigation";
import { createServerClient } from "@/lib/supabase/server";
import { DASHBOARD_ACCESS, roleFromAppMetadata } from "@/lib/roles";
import { fetchAdmin } from "@/lib/adminapi";
import { fmtMinor } from "@/lib/money/minor";
import { fmtTs } from "@/lib/date";
import type { Page } from "@/types/admin";

interface Dispute {
  id: string;
  refId: string;
  merchantId: string;
  cardholder: string;
  amountMinor: number;
  currency: string;
  reasonCode: string;
  category: string;
  stage: string;
  status: string;
  filedAt: string;
  responseDue: string;
  respondedAt?: string;
  escalatedAt?: string;
  decision?: string;
  winner?: string;
  decisionAt?: string;
  disputeFee: number;
  arbitrationFee: number;
  note?: string;
}

export default async function DisputesPage() {
  const supabase = await createServerClient();
  const { data } = await supabase.auth.getUser();
  const role = roleFromAppMetadata(data.user?.app_metadata);
  if (!role || !DASHBOARD_ACCESS[role].includes("/disputes")) notFound();
  const page = await fetchAdmin<Page<Dispute>>("/disputes?limit=50");
  return (
    <div className="grid gap-4">
      <h1 className="text-2xl font-semibold">Disputes</h1>
      <div className="rounded-lg border">
        <table className="w-full text-sm">
          <thead><tr className="border-b text-left text-muted-foreground">
            <th className="px-3 py-2">ID</th><th className="px-3 py-2">Cardholder</th>
            <th className="px-3 py-2">Reason</th><th className="px-3 py-2">Category</th>
            <th className="px-3 py-2">Stage</th><th className="px-3 py-2">Status</th>
            <th className="px-3 py-2">Amount</th><th className="px-3 py-2">Filed</th><th className="px-3 py-2">Due</th>
          </tr></thead>
          <tbody>
            {page.items.map(d => (
              <tr key={d.id} className="border-b last:border-0">
                <td className="px-3 py-2 font-mono">{d.id}</td>
                <td className="px-3 py-2">{d.cardholder}</td>
                <td className="px-3 py-2">{d.reasonCode}</td>
                <td className="px-3 py-2">{d.category}</td>
                <td className="px-3 py-2">{d.stage}</td>
                <td className="px-3 py-2">{d.status}</td>
                <td className="px-3 py-2">{fmtMinor(d.amountMinor, d.currency)}</td>
                <td className="px-3 py-2 text-muted-foreground">{fmtTs(d.filedAt)}</td>
                <td className="px-3 py-2 text-muted-foreground">{d.responseDue ? fmtTs(d.responseDue) : "—"}</td>
              </tr>
            ))}
          </tbody>
        </table>
        {page.items.length === 0 && <p className="p-4 text-sm text-muted-foreground">No disputes yet — run the seed (Task 10).</p>}
      </div>
    </div>
  );
}