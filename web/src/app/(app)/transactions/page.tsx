// src/app/(app)/transactions/page.tsx
import { fetchAdmin } from "@/lib/adminapi";
import { fmtTs } from "@/lib/date";
import type { Page } from "@/types/admin";

interface AuditEvent {
  stan: string;
  mti: string;
  pan: string;
  amount: string;
  responseCode: string;
  destination: string;
  createdAt: string;
}

export default async function TransactionsPage() {
  const page = await fetchAdmin<Page<AuditEvent>>("/transactions?limit=50");
  return (
    <div className="grid gap-4">
      <h1 className="text-2xl font-semibold">Transactions</h1>
      <div className="rounded-lg border">
        <table className="w-full text-sm">
          <thead><tr className="border-b text-left text-muted-foreground">
            <th className="px-3 py-2">STAN</th><th className="px-3 py-2">MTI</th><th className="px-3 py-2">PAN</th>
            <th className="px-3 py-2">Amount</th><th className="px-3 py-2">Resp</th><th className="px-3 py-2">Dest</th>
            <th className="px-3 py-2">Time</th>
          </tr></thead>
          <tbody>
            {page.items.map(t => (
              <tr key={t.stan} className="border-b last:border-0">
                <td className="px-3 py-2 font-mono">{t.stan}</td>
                <td className="px-3 py-2 font-mono">{t.mti}</td>
                <td className="px-3 py-2 font-mono">{t.pan}</td>
                <td className="px-3 py-2">{t.amount}</td>
                <td className="px-3 py-2 font-mono">{t.responseCode}</td>
                <td className="px-3 py-2 font-mono">{t.destination}</td>
                <td className="px-3 py-2 text-muted-foreground">{fmtTs(t.createdAt)}</td>
              </tr>
            ))}
          </tbody>
        </table>
        {page.items.length === 0 && <p className="p-4 text-sm text-muted-foreground">No transactions yet — run the seed (Task 10).</p>}
      </div>
    </div>
  );
}