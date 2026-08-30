import { notFound } from "next/navigation";
import { createServerClient } from "@/lib/supabase/server";
import { DASHBOARD_ACCESS, roleFromAppMetadata } from "@/lib/roles";
import { fetchAdmin } from "@/lib/adminapi";
import { fmtTs } from "@/lib/date";
import type { Page } from "@/types/admin";

interface Token {
  token: string;
  par: string;
  status: string;
  bin: string;
  requestor: string;
  deviceId: string;
  createdAt: string;
}

export default async function TokensPage() {
  const supabase = await createServerClient();
  const { data } = await supabase.auth.getUser();
  const role = roleFromAppMetadata(data.user?.app_metadata);
  if (!role || !DASHBOARD_ACCESS[role].includes("/tokens")) notFound();
  const page = await fetchAdmin<Page<Token>>("/tokens?limit=50");
  return (
    <div className="grid gap-4">
      <h1 className="text-2xl font-semibold">Tokens</h1>
      <div className="rounded-lg border">
        <table className="w-full text-sm">
          <thead><tr className="border-b text-left text-muted-foreground">
            <th className="px-3 py-2">Token</th><th className="px-3 py-2">PAR</th><th className="px-3 py-2">BIN</th>
            <th className="px-3 py-2">Requestor</th><th className="px-3 py-2">Status</th><th className="px-3 py-2">Created</th>
          </tr></thead>
          <tbody>
            {page.items.map(t => (
              <tr key={t.token} className="border-b last:border-0">
                <td className="px-3 py-2 font-mono">{t.token}</td>
                <td className="px-3 py-2 font-mono">{t.par}</td>
                <td className="px-3 py-2 font-mono">{t.bin}</td>
                <td className="px-3 py-2">{t.requestor}</td>
                <td className="px-3 py-2">{t.status}</td>
                <td className="px-3 py-2 text-muted-foreground">{fmtTs(t.createdAt)}</td>
              </tr>
            ))}
          </tbody>
        </table>
        {page.items.length === 0 && <p className="p-4 text-sm text-muted-foreground">No tokens yet — run the seed (Task 10).</p>}
      </div>
    </div>
  );
}