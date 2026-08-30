import { notFound } from "next/navigation";
import { createServerClient } from "@/lib/supabase/server";
import { DASHBOARD_ACCESS, roleFromAppMetadata } from "@/lib/roles";
import { fetchAdmin } from "@/lib/adminapi";
import { fmtMinor } from "@/lib/money/minor";
import { fmtTs } from "@/lib/date";
import { StatCard } from "@/components/cards/stat-card";

interface PrefundAccount {
  member: string;
  balance: number;
  cap: number;
}

interface SettlementInstruction {
  cycleId: string;
  msgId: string;
  member: string;
  amount: number;
  direction: string;
  currency: string;
  instruction: string;
  final: boolean;
}

export default async function SettlementPage() {
  const supabase = await createServerClient();
  const { data } = await supabase.auth.getUser();
  const role = roleFromAppMetadata(data.user?.app_metadata);
  if (!role || !DASHBOARD_ACCESS[role].includes("/settlement")) notFound();
  const prefunds = await fetchAdmin<{ items: PrefundAccount[] }>("/settlement/prefunds");
  const df = await fetchAdmin<{ balance: number }>("/settlement/default-fund");
  const cycles = await fetchAdmin<{ items: string[] }>("/clearing/cycles");
  const cycle = cycles.items[0];
  const instructions = cycle
    ? await fetchAdmin<{ items: SettlementInstruction[] }>(`/settlement/instructions?cycle=${encodeURIComponent(cycle)}`)
    : { items: [] as SettlementInstruction[] };
  return (
    <div className="grid gap-4">
      <h1 className="text-2xl font-semibold">Settlement</h1>
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        <StatCard title="Default fund" value={fmtMinor(df.balance, "EUR")} />
      </div>
      <h2 className="text-lg font-semibold">Prefund balances</h2>
      <div className="rounded-lg border">
        <table className="w-full text-sm">
          <thead><tr className="border-b text-left text-muted-foreground">
            <th className="px-3 py-2">Member</th><th className="px-3 py-2">Balance</th><th className="px-3 py-2">Cap</th>
          </tr></thead>
          <tbody>
            {prefunds.items.map(a => (
              <tr key={a.member} className="border-b last:border-0">
                <td className="px-3 py-2">{a.member}</td>
                <td className="px-3 py-2">{fmtMinor(a.balance, "EUR")}</td>
                <td className="px-3 py-2">{fmtMinor(a.cap, "EUR")}</td>
              </tr>
            ))}
          </tbody>
        </table>
        {prefunds.items.length === 0 && <p className="p-4 text-sm text-muted-foreground">No prefund accounts yet — run the seed (Task 10).</p>}
      </div>
      <h2 className="text-lg font-semibold">Latest instructions</h2>
      {cycle ? (
        <div className="rounded-lg border">
          <table className="w-full text-sm">
            <thead><tr className="border-b text-left text-muted-foreground">
              <th className="px-3 py-2">MsgId</th><th className="px-3 py-2">Member</th>
              <th className="px-3 py-2">Dir</th><th className="px-3 py-2">Amount</th>
              <th className="px-3 py-2">Final</th><th className="px-3 py-2">Time</th>
            </tr></thead>
            <tbody>
              {instructions.items.map(i => (
                <tr key={i.msgId} className="border-b last:border-0">
                  <td className="px-3 py-2 font-mono">{i.msgId}</td>
                  <td className="px-3 py-2">{i.member}</td>
                  <td className="px-3 py-2">{i.direction}</td>
                  <td className="px-3 py-2">{fmtMinor(i.amount, i.currency)}</td>
                  <td className="px-3 py-2">{i.final ? "✅" : "—"}</td>
                  <td className="px-3 py-2 text-muted-foreground">{fmtTs(i.instruction)}</td>
                </tr>
              ))}
            </tbody>
          </table>
          {instructions.items.length === 0 && <p className="p-4 text-sm text-muted-foreground">No settlement instructions yet — run the seed (Task 10).</p>}
        </div>
      ) : (
        <p className="text-sm text-muted-foreground">No settlement instructions yet — run the seed (Task 10).</p>
      )}
    </div>
  );
}