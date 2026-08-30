import { notFound } from "next/navigation";
import { createServerClient } from "@/lib/supabase/server";
import { DASHBOARD_ACCESS, roleFromAppMetadata } from "@/lib/roles";
import { fetchAdmin } from "@/lib/adminapi";
import { fmtMinor } from "@/lib/money/minor";

interface LedgerAccount {
  id: string;
  type: string;
  balance: number;
}

export default async function LedgerPage() {
  const supabase = await createServerClient();
  const { data } = await supabase.auth.getUser();
  const role = roleFromAppMetadata(data.user?.app_metadata);
  if (!role || !DASHBOARD_ACCESS[role].includes("/ledger")) notFound();
  const accounts = await fetchAdmin<{ items: LedgerAccount[] }>("/ledger/accounts");
  return (
    <div className="grid gap-4">
      <h1 className="text-2xl font-semibold">Ledger</h1>
      <div className="rounded-lg border">
        <table className="w-full text-sm">
          <thead><tr className="border-b text-left text-muted-foreground">
            <th className="px-3 py-2">Account</th><th className="px-3 py-2">Type</th><th className="px-3 py-2">Balance</th>
          </tr></thead>
          <tbody>
            {accounts.items.map(a => (
              <tr key={a.id} className="border-b last:border-0">
                <td className="px-3 py-2 font-mono">{a.id}</td>
                <td className="px-3 py-2">{a.type}</td>
                <td className="px-3 py-2">{fmtMinor(a.balance)}</td>
              </tr>
            ))}
          </tbody>
        </table>
        {accounts.items.length === 0 && <p className="p-4 text-sm text-muted-foreground">No ledger accounts yet — run the seed (Task 10).</p>}
      </div>
    </div>
  );
}