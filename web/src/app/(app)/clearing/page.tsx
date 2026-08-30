import { notFound } from "next/navigation";
import { createServerClient } from "@/lib/supabase/server";
import { DASHBOARD_ACCESS, roleFromAppMetadata } from "@/lib/roles";
import { fetchAdmin } from "@/lib/adminapi";
import { fmtMinor } from "@/lib/money/minor";

interface ClearingRecord {
  cycleId: string;
  stan: string;
  mti: string;
  sender: string;
  receiver: string;
  amountMinor: number;
  interchange: number;
  currency: string;
  refId: string;
}

interface NetPosition {
  cycleId: string;
  member: string;
  net: number;
}

export default async function ClearingPage() {
  const supabase = await createServerClient();
  const { data } = await supabase.auth.getUser();
  const role = roleFromAppMetadata(data.user?.app_metadata);
  if (!role || !DASHBOARD_ACCESS[role].includes("/clearing")) notFound();
  const cycles = await fetchAdmin<{ items: string[] }>("/clearing/cycles");
  if (cycles.items.length === 0) {
    return (
      <div className="grid gap-4">
        <h1 className="text-2xl font-semibold">Clearing</h1>
        <p className="text-sm text-muted-foreground">No clearing cycles yet — run the seed (Task 10).</p>
      </div>
    );
  }
  const cycle = cycles.items[0];
  const records = await fetchAdmin<{ items: ClearingRecord[] }>(`/clearing/records?cycle=${encodeURIComponent(cycle)}`);
  const positions = await fetchAdmin<{ items: NetPosition[] }>(`/clearing/positions?cycle=${encodeURIComponent(cycle)}`);
  return (
    <div className="grid gap-4">
      <h1 className="text-2xl font-semibold">Clearing</h1>
      <p className="text-sm text-muted-foreground">Cycle {cycle}</p>
      <div className="rounded-lg border">
        <table className="w-full text-sm">
          <thead><tr className="border-b text-left text-muted-foreground">
            <th className="px-3 py-2">STAN</th><th className="px-3 py-2">MTI</th>
            <th className="px-3 py-2">Sender</th><th className="px-3 py-2">Receiver</th>
            <th className="px-3 py-2">Amount</th><th className="px-3 py-2">Interchange</th>
          </tr></thead>
          <tbody>
            {records.items.map(r => (
              <tr key={r.stan} className="border-b last:border-0">
                <td className="px-3 py-2 font-mono">{r.stan}</td>
                <td className="px-3 py-2 font-mono">{r.mti}</td>
                <td className="px-3 py-2">{r.sender}</td>
                <td className="px-3 py-2">{r.receiver}</td>
                <td className="px-3 py-2">{fmtMinor(r.amountMinor, r.currency)}</td>
                <td className="px-3 py-2">{fmtMinor(r.interchange, r.currency)}</td>
              </tr>
            ))}
          </tbody>
        </table>
        {records.items.length === 0 && <p className="p-4 text-sm text-muted-foreground">No clearing records yet — run the seed (Task 10).</p>}
      </div>
      <div className="rounded-lg border">
        <table className="w-full text-sm">
          <thead><tr className="border-b text-left text-muted-foreground">
            <th className="px-3 py-2">Member</th><th className="px-3 py-2">Net</th>
          </tr></thead>
          <tbody>
            {positions.items.map(p => (
              <tr key={p.member} className="border-b last:border-0">
                <td className="px-3 py-2">{p.member}</td>
                <td className="px-3 py-2">{fmtMinor(p.net)}</td>
              </tr>
            ))}
          </tbody>
        </table>
        {positions.items.length === 0 && <p className="p-4 text-sm text-muted-foreground">No net positions yet — run the seed (Task 10).</p>}
      </div>
    </div>
  );
}