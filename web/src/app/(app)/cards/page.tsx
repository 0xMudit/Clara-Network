import { notFound } from "next/navigation";
import { createServerClient } from "@/lib/supabase/server";
import { DASHBOARD_ACCESS, roleFromAppMetadata } from "@/lib/roles";
import { fetchAdmin } from "@/lib/adminapi";
import type { Page } from "@/types/admin";

interface Card {
  ref: string;
  panHash: string;
  panMask: string;
  bin: string;
  expiry: string;
  status: string;
  product: string;
  lastAtc: number;
}

export default async function CardsPage() {
  const supabase = await createServerClient();
  const { data } = await supabase.auth.getUser();
  const role = roleFromAppMetadata(data.user?.app_metadata);
  if (!role || !DASHBOARD_ACCESS[role].includes("/cards")) notFound();
  const page = await fetchAdmin<Page<Card>>("/cards?limit=50");
  return (
    <div className="grid gap-4">
      <h1 className="text-2xl font-semibold">Cards</h1>
      <div className="rounded-lg border">
        <table className="w-full text-sm">
          <thead><tr className="border-b text-left text-muted-foreground">
            <th className="px-3 py-2">Ref</th><th className="px-3 py-2">PAN</th><th className="px-3 py-2">BIN</th>
            <th className="px-3 py-2">Product</th><th className="px-3 py-2">Status</th>
            <th className="px-3 py-2">Expiry</th><th className="px-3 py-2">Last ATC</th>
          </tr></thead>
          <tbody>
            {page.items.map(c => (
              <tr key={c.ref} className="border-b last:border-0">
                <td className="px-3 py-2 font-mono">{c.ref}</td>
                <td className="px-3 py-2 font-mono">{c.panMask}</td>
                <td className="px-3 py-2 font-mono">{c.bin}</td>
                <td className="px-3 py-2">{c.product}</td>
                <td className="px-3 py-2">{c.status}</td>
                <td className="px-3 py-2">{c.expiry}</td>
                <td className="px-3 py-2">{c.lastAtc}</td>
              </tr>
            ))}
          </tbody>
        </table>
        {page.items.length === 0 && <p className="p-4 text-sm text-muted-foreground">No cards yet — run the seed (Task 10).</p>}
      </div>
    </div>
  );
}