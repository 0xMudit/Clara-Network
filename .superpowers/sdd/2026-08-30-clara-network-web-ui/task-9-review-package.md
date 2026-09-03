BASE 0f711eb

15397cc feat(web): settlement/ledger/cards/tokens/merchant/dispute dashboards
cca1bfb fix(web): forward query strings in BFF proxy

=== STAT ===
 web/src/app/(app)/acquirer/page.tsx     | 20 +++++++-
 web/src/app/(app)/cards/page.tsx        | 52 +++++++++++++++++++
 web/src/app/(app)/clearing/page.tsx     | 86 +++++++++++++++++++++++++++++++
 web/src/app/(app)/disputes/page.tsx     | 69 +++++++++++++++++++++++++
 web/src/app/(app)/funding/page.tsx      | 56 ++++++++++++++++++++
 web/src/app/(app)/issuer/page.tsx       | 20 +++++++-
 web/src/app/(app)/ledger/page.tsx       | 41 +++++++++++++++
 web/src/app/(app)/merchant/page.tsx     | 20 +++++++-
 web/src/app/(app)/merchants/page.tsx    | 59 +++++++++++++++++++++
 web/src/app/(app)/settlement/page.tsx   | 91 +++++++++++++++++++++++++++++++++
 web/src/app/(app)/tokens/page.tsx       | 50 ++++++++++++++++++
 web/src/app/api/data/[...path]/route.ts |  7 +--
 12 files changed, 565 insertions(+), 6 deletions(-)
diff --git a/web/src/app/(app)/acquirer/page.tsx b/web/src/app/(app)/acquirer/page.tsx
index 2226234..6d9c922 100644
--- a/web/src/app/(app)/acquirer/page.tsx
+++ b/web/src/app/(app)/acquirer/page.tsx
@@ -1,10 +1,28 @@
 import { notFound } from "next/navigation";
+import Link from "next/link";
 import { createServerClient } from "@/lib/supabase/server";
 import { roleFromAppMetadata } from "@/lib/roles";
+import { fetchAdmin } from "@/lib/adminapi";
+import { StatCard } from "@/components/cards/stat-card";
+import type { DashboardSummary } from "@/types/admin";
 
 export default async function AcquirerPage() {
   const supabase = await createServerClient();
   const { data } = await supabase.auth.getUser();
   if (roleFromAppMetadata(data.user?.app_metadata) !== "acquirer") notFound();
-  return <h1 className="text-2xl font-semibold">Acquirer</h1>;
+  const d = await fetchAdmin<DashboardSummary>("/dashboard");
+  return (
+    <div className="grid gap-4">
+      <h1 className="text-2xl font-semibold">Acquirer</h1>
+      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
+        <StatCard title="Merchants" value={d.merchants.toLocaleString()} />
+        <StatCard title="Disputes" value={d.disputes.toLocaleString()} />
+      </div>
+      <p className="text-sm text-muted-foreground">
+        <Link href="/merchants" className="underline underline-offset-4 hover:text-foreground">View merchants ÔåÆ</Link>
+        {" ┬À "}
+        <Link href="/disputes" className="underline underline-offset-4 hover:text-foreground">View disputes ÔåÆ</Link>
+      </p>
+    </div>
+  );
 }
\ No newline at end of file
diff --git a/web/src/app/(app)/cards/page.tsx b/web/src/app/(app)/cards/page.tsx
new file mode 100644
index 0000000..38e6601
--- /dev/null
+++ b/web/src/app/(app)/cards/page.tsx
@@ -0,0 +1,52 @@
+import { notFound } from "next/navigation";
+import { createServerClient } from "@/lib/supabase/server";
+import { DASHBOARD_ACCESS, roleFromAppMetadata } from "@/lib/roles";
+import { fetchAdmin } from "@/lib/adminapi";
+import type { Page } from "@/types/admin";
+
+interface Card {
+  ref: string;
+  panHash: string;
+  panMask: string;
+  bin: string;
+  expiry: string;
+  status: string;
+  product: string;
+  lastAtc: number;
+}
+
+export default async function CardsPage() {
+  const supabase = await createServerClient();
+  const { data } = await supabase.auth.getUser();
+  const role = roleFromAppMetadata(data.user?.app_metadata);
+  if (!role || !DASHBOARD_ACCESS[role].includes("/cards")) notFound();
+  const page = await fetchAdmin<Page<Card>>("/cards?limit=50");
+  return (
+    <div className="grid gap-4">
+      <h1 className="text-2xl font-semibold">Cards</h1>
+      <div className="rounded-lg border">
+        <table className="w-full text-sm">
+          <thead><tr className="border-b text-left text-muted-foreground">
+            <th className="px-3 py-2">Ref</th><th className="px-3 py-2">PAN</th><th className="px-3 py-2">BIN</th>
+            <th className="px-3 py-2">Product</th><th className="px-3 py-2">Status</th>
+            <th className="px-3 py-2">Expiry</th><th className="px-3 py-2">Last ATC</th>
+          </tr></thead>
+          <tbody>
+            {page.items.map(c => (
+              <tr key={c.ref} className="border-b last:border-0">
+                <td className="px-3 py-2 font-mono">{c.ref}</td>
+                <td className="px-3 py-2 font-mono">{c.panMask}</td>
+                <td className="px-3 py-2 font-mono">{c.bin}</td>
+                <td className="px-3 py-2">{c.product}</td>
+                <td className="px-3 py-2">{c.status}</td>
+                <td className="px-3 py-2">{c.expiry}</td>
+                <td className="px-3 py-2">{c.lastAtc}</td>
+              </tr>
+            ))}
+          </tbody>
+        </table>
+        {page.items.length === 0 && <p className="p-4 text-sm text-muted-foreground">No cards yet ÔÇö run the seed (Task 10).</p>}
+      </div>
+    </div>
+  );
+}
\ No newline at end of file
diff --git a/web/src/app/(app)/clearing/page.tsx b/web/src/app/(app)/clearing/page.tsx
new file mode 100644
index 0000000..7ebfd3b
--- /dev/null
+++ b/web/src/app/(app)/clearing/page.tsx
@@ -0,0 +1,86 @@
+import { notFound } from "next/navigation";
+import { createServerClient } from "@/lib/supabase/server";
+import { DASHBOARD_ACCESS, roleFromAppMetadata } from "@/lib/roles";
+import { fetchAdmin } from "@/lib/adminapi";
+import { fmtMinor } from "@/lib/money/minor";
+
+interface ClearingRecord {
+  cycleId: string;
+  stan: string;
+  mti: string;
+  sender: string;
+  receiver: string;
+  amountMinor: number;
+  interchange: number;
+  currency: string;
+  refId: string;
+}
+
+interface NetPosition {
+  cycleId: string;
+  member: string;
+  net: number;
+}
+
+export default async function ClearingPage() {
+  const supabase = await createServerClient();
+  const { data } = await supabase.auth.getUser();
+  const role = roleFromAppMetadata(data.user?.app_metadata);
+  if (!role || !DASHBOARD_ACCESS[role].includes("/clearing")) notFound();
+  const cycles = await fetchAdmin<{ items: string[] }>("/clearing/cycles");
+  if (cycles.items.length === 0) {
+    return (
+      <div className="grid gap-4">
+        <h1 className="text-2xl font-semibold">Clearing</h1>
+        <p className="text-sm text-muted-foreground">No clearing cycles yet ÔÇö run the seed (Task 10).</p>
+      </div>
+    );
+  }
+  const cycle = cycles.items[0];
+  const records = await fetchAdmin<{ items: ClearingRecord[] }>(`/clearing/records?cycle=${encodeURIComponent(cycle)}`);
+  const positions = await fetchAdmin<{ items: NetPosition[] }>(`/clearing/positions?cycle=${encodeURIComponent(cycle)}`);
+  return (
+    <div className="grid gap-4">
+      <h1 className="text-2xl font-semibold">Clearing</h1>
+      <p className="text-sm text-muted-foreground">Cycle {cycle}</p>
+      <div className="rounded-lg border">
+        <table className="w-full text-sm">
+          <thead><tr className="border-b text-left text-muted-foreground">
+            <th className="px-3 py-2">STAN</th><th className="px-3 py-2">MTI</th>
+            <th className="px-3 py-2">Sender</th><th className="px-3 py-2">Receiver</th>
+            <th className="px-3 py-2">Amount</th><th className="px-3 py-2">Interchange</th>
+          </tr></thead>
+          <tbody>
+            {records.items.map(r => (
+              <tr key={r.stan} className="border-b last:border-0">
+                <td className="px-3 py-2 font-mono">{r.stan}</td>
+                <td className="px-3 py-2 font-mono">{r.mti}</td>
+                <td className="px-3 py-2">{r.sender}</td>
+                <td className="px-3 py-2">{r.receiver}</td>
+                <td className="px-3 py-2">{fmtMinor(r.amountMinor, r.currency)}</td>
+                <td className="px-3 py-2">{fmtMinor(r.interchange, r.currency)}</td>
+              </tr>
+            ))}
+          </tbody>
+        </table>
+        {records.items.length === 0 && <p className="p-4 text-sm text-muted-foreground">No clearing records yet ÔÇö run the seed (Task 10).</p>}
+      </div>
+      <div className="rounded-lg border">
+        <table className="w-full text-sm">
+          <thead><tr className="border-b text-left text-muted-foreground">
+            <th className="px-3 py-2">Member</th><th className="px-3 py-2">Net</th>
+          </tr></thead>
+          <tbody>
+            {positions.items.map(p => (
+              <tr key={p.member} className="border-b last:border-0">
+                <td className="px-3 py-2">{p.member}</td>
+                <td className="px-3 py-2">{fmtMinor(p.net)}</td>
+              </tr>
+            ))}
+          </tbody>
+        </table>
+        {positions.items.length === 0 && <p className="p-4 text-sm text-muted-foreground">No net positions yet ÔÇö run the seed (Task 10).</p>}
+      </div>
+    </div>
+  );
+}
\ No newline at end of file
diff --git a/web/src/app/(app)/disputes/page.tsx b/web/src/app/(app)/disputes/page.tsx
new file mode 100644
index 0000000..9c95c9b
--- /dev/null
+++ b/web/src/app/(app)/disputes/page.tsx
@@ -0,0 +1,69 @@
+import { notFound } from "next/navigation";
+import { createServerClient } from "@/lib/supabase/server";
+import { DASHBOARD_ACCESS, roleFromAppMetadata } from "@/lib/roles";
+import { fetchAdmin } from "@/lib/adminapi";
+import { fmtMinor } from "@/lib/money/minor";
+import { fmtTs } from "@/lib/date";
+import type { Page } from "@/types/admin";
+
+interface Dispute {
+  id: string;
+  refId: string;
+  merchantId: string;
+  cardholder: string;
+  amountMinor: number;
+  currency: string;
+  reasonCode: string;
+  category: string;
+  stage: string;
+  status: string;
+  filedAt: string;
+  responseDue: string;
+  respondedAt?: string;
+  escalatedAt?: string;
+  decision?: string;
+  winner?: string;
+  decisionAt?: string;
+  disputeFee: number;
+  arbitrationFee: number;
+  note?: string;
+}
+
+export default async function DisputesPage() {
+  const supabase = await createServerClient();
+  const { data } = await supabase.auth.getUser();
+  const role = roleFromAppMetadata(data.user?.app_metadata);
+  if (!role || !DASHBOARD_ACCESS[role].includes("/disputes")) notFound();
+  const page = await fetchAdmin<Page<Dispute>>("/disputes?limit=50");
+  return (
+    <div className="grid gap-4">
+      <h1 className="text-2xl font-semibold">Disputes</h1>
+      <div className="rounded-lg border">
+        <table className="w-full text-sm">
+          <thead><tr className="border-b text-left text-muted-foreground">
+            <th className="px-3 py-2">ID</th><th className="px-3 py-2">Cardholder</th>
+            <th className="px-3 py-2">Reason</th><th className="px-3 py-2">Category</th>
+            <th className="px-3 py-2">Stage</th><th className="px-3 py-2">Status</th>
+            <th className="px-3 py-2">Amount</th><th className="px-3 py-2">Filed</th><th className="px-3 py-2">Due</th>
+          </tr></thead>
+          <tbody>
+            {page.items.map(d => (
+              <tr key={d.id} className="border-b last:border-0">
+                <td className="px-3 py-2 font-mono">{d.id}</td>
+                <td className="px-3 py-2">{d.cardholder}</td>
+                <td className="px-3 py-2">{d.reasonCode}</td>
+                <td className="px-3 py-2">{d.category}</td>
+                <td className="px-3 py-2">{d.stage}</td>
+                <td className="px-3 py-2">{d.status}</td>
+                <td className="px-3 py-2">{fmtMinor(d.amountMinor, d.currency)}</td>
+                <td className="px-3 py-2 text-muted-foreground">{fmtTs(d.filedAt)}</td>
+                <td className="px-3 py-2 text-muted-foreground">{d.responseDue ? fmtTs(d.responseDue) : "ÔÇö"}</td>
+              </tr>
+            ))}
+          </tbody>
+        </table>
+        {page.items.length === 0 && <p className="p-4 text-sm text-muted-foreground">No disputes yet ÔÇö run the seed (Task 10).</p>}
+      </div>
+    </div>
+  );
+}
\ No newline at end of file
diff --git a/web/src/app/(app)/funding/page.tsx b/web/src/app/(app)/funding/page.tsx
new file mode 100644
index 0000000..0fde7aa
--- /dev/null
+++ b/web/src/app/(app)/funding/page.tsx
@@ -0,0 +1,56 @@
+import { notFound } from "next/navigation";
+import { createServerClient } from "@/lib/supabase/server";
+import { DASHBOARD_ACCESS, roleFromAppMetadata } from "@/lib/roles";
+import { fetchAdmin } from "@/lib/adminapi";
+import { fmtMinor } from "@/lib/money/minor";
+import type { Page } from "@/types/admin";
+
+interface Merchant {
+  id: string;
+  name: string;
+  dba: string;
+  taxId: string;
+  mccs: string[];
+  status: string;
+  riskTier: string;
+  reserveRateBps: number;
+  fundingDelayDays: number;
+  transactionLimit: number;
+  reserveBalance: number;
+  volume: number;
+  declineReason?: string;
+  approvedAt: string;
+}
+
+export default async function FundingPage() {
+  const supabase = await createServerClient();
+  const { data } = await supabase.auth.getUser();
+  const role = roleFromAppMetadata(data.user?.app_metadata);
+  if (!role || !DASHBOARD_ACCESS[role].includes("/funding")) notFound();
+  const page = await fetchAdmin<Page<Merchant>>("/merchants?limit=50");
+  return (
+    <div className="grid gap-4">
+      <h1 className="text-2xl font-semibold">Acquirer funding</h1>
+      <div className="rounded-lg border">
+        <table className="w-full text-sm">
+          <thead><tr className="border-b text-left text-muted-foreground">
+            <th className="px-3 py-2">Name</th><th className="px-3 py-2">Delay (days)</th>
+            <th className="px-3 py-2">Txn limit</th><th className="px-3 py-2">Reserve</th><th className="px-3 py-2">Volume</th>
+          </tr></thead>
+          <tbody>
+            {page.items.map(m => (
+              <tr key={m.id} className="border-b last:border-0">
+                <td className="px-3 py-2">{m.name}</td>
+                <td className="px-3 py-2">{m.fundingDelayDays}</td>
+                <td className="px-3 py-2">{fmtMinor(m.transactionLimit)}</td>
+                <td className="px-3 py-2">{fmtMinor(m.reserveBalance)}</td>
+                <td className="px-3 py-2">{fmtMinor(m.volume)}</td>
+              </tr>
+            ))}
+          </tbody>
+        </table>
+        {page.items.length === 0 && <p className="p-4 text-sm text-muted-foreground">No merchants yet ÔÇö run the seed (Task 10).</p>}
+      </div>
+    </div>
+  );
+}
\ No newline at end of file
diff --git a/web/src/app/(app)/issuer/page.tsx b/web/src/app/(app)/issuer/page.tsx
index 92557fb..ee3b609 100644
--- a/web/src/app/(app)/issuer/page.tsx
+++ b/web/src/app/(app)/issuer/page.tsx
@@ -1,10 +1,28 @@
 import { notFound } from "next/navigation";
+import Link from "next/link";
 import { createServerClient } from "@/lib/supabase/server";
 import { roleFromAppMetadata } from "@/lib/roles";
+import { fetchAdmin } from "@/lib/adminapi";
+import { StatCard } from "@/components/cards/stat-card";
+import type { DashboardSummary } from "@/types/admin";
 
 export default async function IssuerPage() {
   const supabase = await createServerClient();
   const { data } = await supabase.auth.getUser();
   if (roleFromAppMetadata(data.user?.app_metadata) !== "issuer") notFound();
-  return <h1 className="text-2xl font-semibold">Issuer</h1>;
+  const d = await fetchAdmin<DashboardSummary>("/dashboard");
+  return (
+    <div className="grid gap-4">
+      <h1 className="text-2xl font-semibold">Issuer</h1>
+      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
+        <StatCard title="Cards" value={d.cards.toLocaleString()} />
+        <StatCard title="Tokens" value={d.tokens.toLocaleString()} />
+      </div>
+      <p className="text-sm text-muted-foreground">
+        <Link href="/cards" className="underline underline-offset-4 hover:text-foreground">View cards ÔåÆ</Link>
+        {" ┬À "}
+        <Link href="/tokens" className="underline underline-offset-4 hover:text-foreground">View tokens ÔåÆ</Link>
+      </p>
+    </div>
+  );
 }
\ No newline at end of file
diff --git a/web/src/app/(app)/ledger/page.tsx b/web/src/app/(app)/ledger/page.tsx
new file mode 100644
index 0000000..ff2d71c
--- /dev/null
+++ b/web/src/app/(app)/ledger/page.tsx
@@ -0,0 +1,41 @@
+import { notFound } from "next/navigation";
+import { createServerClient } from "@/lib/supabase/server";
+import { DASHBOARD_ACCESS, roleFromAppMetadata } from "@/lib/roles";
+import { fetchAdmin } from "@/lib/adminapi";
+import { fmtMinor } from "@/lib/money/minor";
+
+interface LedgerAccount {
+  id: string;
+  type: string;
+  balance: number;
+}
+
+export default async function LedgerPage() {
+  const supabase = await createServerClient();
+  const { data } = await supabase.auth.getUser();
+  const role = roleFromAppMetadata(data.user?.app_metadata);
+  if (!role || !DASHBOARD_ACCESS[role].includes("/ledger")) notFound();
+  const accounts = await fetchAdmin<{ items: LedgerAccount[] }>("/ledger/accounts");
+  return (
+    <div className="grid gap-4">
+      <h1 className="text-2xl font-semibold">Ledger</h1>
+      <div className="rounded-lg border">
+        <table className="w-full text-sm">
+          <thead><tr className="border-b text-left text-muted-foreground">
+            <th className="px-3 py-2">Account</th><th className="px-3 py-2">Type</th><th className="px-3 py-2">Balance</th>
+          </tr></thead>
+          <tbody>
+            {accounts.items.map(a => (
+              <tr key={a.id} className="border-b last:border-0">
+                <td className="px-3 py-2 font-mono">{a.id}</td>
+                <td className="px-3 py-2">{a.type}</td>
+                <td className="px-3 py-2">{fmtMinor(a.balance)}</td>
+              </tr>
+            ))}
+          </tbody>
+        </table>
+        {accounts.items.length === 0 && <p className="p-4 text-sm text-muted-foreground">No ledger accounts yet ÔÇö run the seed (Task 10).</p>}
+      </div>
+    </div>
+  );
+}
\ No newline at end of file
diff --git a/web/src/app/(app)/merchant/page.tsx b/web/src/app/(app)/merchant/page.tsx
index fca4855..37fc166 100644
--- a/web/src/app/(app)/merchant/page.tsx
+++ b/web/src/app/(app)/merchant/page.tsx
@@ -1,10 +1,28 @@
 import { notFound } from "next/navigation";
+import Link from "next/link";
 import { createServerClient } from "@/lib/supabase/server";
 import { roleFromAppMetadata } from "@/lib/roles";
+import { fetchAdmin } from "@/lib/adminapi";
+import { StatCard } from "@/components/cards/stat-card";
+import type { DashboardSummary } from "@/types/admin";
 
 export default async function MerchantPage() {
   const supabase = await createServerClient();
   const { data } = await supabase.auth.getUser();
   if (roleFromAppMetadata(data.user?.app_metadata) !== "merchant") notFound();
-  return <h1 className="text-2xl font-semibold">Merchant</h1>;
+  const d = await fetchAdmin<DashboardSummary>("/dashboard");
+  return (
+    <div className="grid gap-4">
+      <h1 className="text-2xl font-semibold">Merchant</h1>
+      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
+        <StatCard title="Merchants" value={d.merchants.toLocaleString()} />
+        <StatCard title="Disputes" value={d.disputes.toLocaleString()} />
+      </div>
+      <p className="text-sm text-muted-foreground">
+        <Link href="/funding" className="underline underline-offset-4 hover:text-foreground">View funding ÔåÆ</Link>
+        {" ┬À "}
+        <Link href="/disputes" className="underline underline-offset-4 hover:text-foreground">View disputes ÔåÆ</Link>
+      </p>
+    </div>
+  );
 }
\ No newline at end of file
diff --git a/web/src/app/(app)/merchants/page.tsx b/web/src/app/(app)/merchants/page.tsx
new file mode 100644
index 0000000..02a7718
--- /dev/null
+++ b/web/src/app/(app)/merchants/page.tsx
@@ -0,0 +1,59 @@
+import { notFound } from "next/navigation";
+import { createServerClient } from "@/lib/supabase/server";
+import { DASHBOARD_ACCESS, roleFromAppMetadata } from "@/lib/roles";
+import { fetchAdmin } from "@/lib/adminapi";
+import { fmtMinor } from "@/lib/money/minor";
+import type { Page } from "@/types/admin";
+
+interface Merchant {
+  id: string;
+  name: string;
+  dba: string;
+  taxId: string;
+  mccs: string[];
+  status: string;
+  riskTier: string;
+  reserveRateBps: number;
+  fundingDelayDays: number;
+  transactionLimit: number;
+  reserveBalance: number;
+  volume: number;
+  declineReason?: string;
+  approvedAt: string;
+}
+
+export default async function MerchantsPage() {
+  const supabase = await createServerClient();
+  const { data } = await supabase.auth.getUser();
+  const role = roleFromAppMetadata(data.user?.app_metadata);
+  if (!role || !DASHBOARD_ACCESS[role].includes("/merchants")) notFound();
+  const page = await fetchAdmin<Page<Merchant>>("/merchants?limit=50");
+  return (
+    <div className="grid gap-4">
+      <h1 className="text-2xl font-semibold">Merchants</h1>
+      <div className="rounded-lg border">
+        <table className="w-full text-sm">
+          <thead><tr className="border-b text-left text-muted-foreground">
+            <th className="px-3 py-2">Name</th><th className="px-3 py-2">DBA</th><th className="px-3 py-2">MCCs</th>
+            <th className="px-3 py-2">Risk</th><th className="px-3 py-2">Status</th>
+            <th className="px-3 py-2">Reserve</th><th className="px-3 py-2">Volume</th>
+          </tr></thead>
+          <tbody>
+            {page.items.map(m => (
+              <tr key={m.id} className="border-b last:border-0">
+                <td className="px-3 py-2">{m.name}</td>
+                <td className="px-3 py-2">{m.dba}</td>
+                <td className="px-3 py-2">{m.mccs.join(", ")}</td>
+                <td className="px-3 py-2">{m.riskTier}</td>
+                <td className="px-3 py-2">{m.status}</td>
+                <td className="px-3 py-2">{fmtMinor(m.reserveBalance)}</td>
+                <td className="px-3 py-2">{fmtMinor(m.volume)}</td>
+              </tr>
+            ))}
+          </tbody>
+        </table>
+        {page.items.length === 0 && <p className="p-4 text-sm text-muted-foreground">No merchants yet ÔÇö run the seed (Task 10).</p>}
+      </div>
+    </div>
+  );
+}
\ No newline at end of file
diff --git a/web/src/app/(app)/settlement/page.tsx b/web/src/app/(app)/settlement/page.tsx
new file mode 100644
index 0000000..4a2e4f4
--- /dev/null
+++ b/web/src/app/(app)/settlement/page.tsx
@@ -0,0 +1,91 @@
+import { notFound } from "next/navigation";
+import { createServerClient } from "@/lib/supabase/server";
+import { DASHBOARD_ACCESS, roleFromAppMetadata } from "@/lib/roles";
+import { fetchAdmin } from "@/lib/adminapi";
+import { fmtMinor } from "@/lib/money/minor";
+import { fmtTs } from "@/lib/date";
+import { StatCard } from "@/components/cards/stat-card";
+
+interface PrefundAccount {
+  member: string;
+  balance: number;
+  cap: number;
+}
+
+interface SettlementInstruction {
+  cycleId: string;
+  msgId: string;
+  member: string;
+  amount: number;
+  direction: string;
+  currency: string;
+  instruction: string;
+  final: boolean;
+}
+
+export default async function SettlementPage() {
+  const supabase = await createServerClient();
+  const { data } = await supabase.auth.getUser();
+  const role = roleFromAppMetadata(data.user?.app_metadata);
+  if (!role || !DASHBOARD_ACCESS[role].includes("/settlement")) notFound();
+  const prefunds = await fetchAdmin<{ items: PrefundAccount[] }>("/settlement/prefunds");
+  const df = await fetchAdmin<{ balance: number }>("/settlement/default-fund");
+  const cycles = await fetchAdmin<{ items: string[] }>("/clearing/cycles");
+  const cycle = cycles.items[0];
+  const instructions = cycle
+    ? await fetchAdmin<{ items: SettlementInstruction[] }>(`/settlement/instructions?cycle=${encodeURIComponent(cycle)}`)
+    : { items: [] as SettlementInstruction[] };
+  return (
+    <div className="grid gap-4">
+      <h1 className="text-2xl font-semibold">Settlement</h1>
+      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
+        <StatCard title="Default fund" value={fmtMinor(df.balance, "EUR")} />
+      </div>
+      <h2 className="text-lg font-semibold">Prefund balances</h2>
+      <div className="rounded-lg border">
+        <table className="w-full text-sm">
+          <thead><tr className="border-b text-left text-muted-foreground">
+            <th className="px-3 py-2">Member</th><th className="px-3 py-2">Balance</th><th className="px-3 py-2">Cap</th>
+          </tr></thead>
+          <tbody>
+            {prefunds.items.map(a => (
+              <tr key={a.member} className="border-b last:border-0">
+                <td className="px-3 py-2">{a.member}</td>
+                <td className="px-3 py-2">{fmtMinor(a.balance, "EUR")}</td>
+                <td className="px-3 py-2">{fmtMinor(a.cap, "EUR")}</td>
+              </tr>
+            ))}
+          </tbody>
+        </table>
+        {prefunds.items.length === 0 && <p className="p-4 text-sm text-muted-foreground">No prefund accounts yet ÔÇö run the seed (Task 10).</p>}
+      </div>
+      <h2 className="text-lg font-semibold">Latest instructions</h2>
+      {cycle ? (
+        <div className="rounded-lg border">
+          <table className="w-full text-sm">
+            <thead><tr className="border-b text-left text-muted-foreground">
+              <th className="px-3 py-2">MsgId</th><th className="px-3 py-2">Member</th>
+              <th className="px-3 py-2">Dir</th><th className="px-3 py-2">Amount</th>
+              <th className="px-3 py-2">Final</th><th className="px-3 py-2">Time</th>
+            </tr></thead>
+            <tbody>
+              {instructions.items.map(i => (
+                <tr key={i.msgId} className="border-b last:border-0">
+                  <td className="px-3 py-2 font-mono">{i.msgId}</td>
+                  <td className="px-3 py-2">{i.member}</td>
+                  <td className="px-3 py-2">{i.direction}</td>
+                  <td className="px-3 py-2">{fmtMinor(i.amount, i.currency)}</td>
+                  <td className="px-3 py-2">{i.final ? "Ô£à" : "ÔÇö"}</td>
+                  <td className="px-3 py-2 text-muted-foreground">{fmtTs(i.instruction)}</td>
+                </tr>
+              ))}
+            </tbody>
+          </table>
+          {instructions.items.length === 0 && <p className="p-4 text-sm text-muted-foreground">No settlement instructions yet ÔÇö run the seed (Task 10).</p>}
+        </div>
+      ) : (
+        <p className="text-sm text-muted-foreground">No settlement instructions yet ÔÇö run the seed (Task 10).</p>
+      )}
+    </div>
+  );
+}
\ No newline at end of file
diff --git a/web/src/app/(app)/tokens/page.tsx b/web/src/app/(app)/tokens/page.tsx
new file mode 100644
index 0000000..a252f40
--- /dev/null
+++ b/web/src/app/(app)/tokens/page.tsx
@@ -0,0 +1,50 @@
+import { notFound } from "next/navigation";
+import { createServerClient } from "@/lib/supabase/server";
+import { DASHBOARD_ACCESS, roleFromAppMetadata } from "@/lib/roles";
+import { fetchAdmin } from "@/lib/adminapi";
+import { fmtTs } from "@/lib/date";
+import type { Page } from "@/types/admin";
+
+interface Token {
+  token: string;
+  par: string;
+  status: string;
+  bin: string;
+  requestor: string;
+  deviceId: string;
+  createdAt: string;
+}
+
+export default async function TokensPage() {
+  const supabase = await createServerClient();
+  const { data } = await supabase.auth.getUser();
+  const role = roleFromAppMetadata(data.user?.app_metadata);
+  if (!role || !DASHBOARD_ACCESS[role].includes("/tokens")) notFound();
+  const page = await fetchAdmin<Page<Token>>("/tokens?limit=50");
+  return (
+    <div className="grid gap-4">
+      <h1 className="text-2xl font-semibold">Tokens</h1>
+      <div className="rounded-lg border">
+        <table className="w-full text-sm">
+          <thead><tr className="border-b text-left text-muted-foreground">
+            <th className="px-3 py-2">Token</th><th className="px-3 py-2">PAR</th><th className="px-3 py-2">BIN</th>
+            <th className="px-3 py-2">Requestor</th><th className="px-3 py-2">Status</th><th className="px-3 py-2">Created</th>
+          </tr></thead>
+          <tbody>
+            {page.items.map(t => (
+              <tr key={t.token} className="border-b last:border-0">
+                <td className="px-3 py-2 font-mono">{t.token}</td>
+                <td className="px-3 py-2 font-mono">{t.par}</td>
+                <td className="px-3 py-2 font-mono">{t.bin}</td>
+                <td className="px-3 py-2">{t.requestor}</td>
+                <td className="px-3 py-2">{t.status}</td>
+                <td className="px-3 py-2 text-muted-foreground">{fmtTs(t.createdAt)}</td>
+              </tr>
+            ))}
+          </tbody>
+        </table>
+        {page.items.length === 0 && <p className="p-4 text-sm text-muted-foreground">No tokens yet ÔÇö run the seed (Task 10).</p>}
+      </div>
+    </div>
+  );
+}
\ No newline at end of file
diff --git a/web/src/app/api/data/[...path]/route.ts b/web/src/app/api/data/[...path]/route.ts
index fa10fa8..5cf2304 100644
--- a/web/src/app/api/data/[...path]/route.ts
+++ b/web/src/app/api/data/[...path]/route.ts
@@ -23,21 +23,22 @@ const ALLOWED_ROUTES: Record<string, string[]> = {
 
 export async function GET(
   req: Request,
   { params }: { params: Promise<{ path: string[] }> },
 ): Promise<Res<unknown>> {
   const supabase = await createServerClient();
   const { data, error } = await supabase.auth.getUser();
   if (error || !data.user) return NextResponse.json({ error: "unauthorized" }, { status: 401 });
   const role = data.user.app_metadata.role as string;
   const { path } = await params;
-  const target = "/" + path.join("/");
+  const url = new URL(req.url);
+  const base = "/" + path.join("/");
   const allowed = ALLOWED_ROUTES[role] ?? [];
-  if (!allowed.includes(target)) return NextResponse.json({ error: "forbidden" }, { status: 403 });
+  if (!allowed.includes(base)) return NextResponse.json({ error: "forbidden" }, { status: 403 });
   try {
-    const body = await claraFetch(target, data.user.id, { cache: "no-store" });
+    const body = await claraFetch(base + url.search, data.user.id, { cache: "no-store" });
     return NextResponse.json(body);
   } catch (e) {
     const status = e instanceof Error && /^clara \S+ 4\d\d/.test(e.message) ? 502 : 500;
     return NextResponse.json({ error: "upstream error", detail: e instanceof Error ? e.message : String(e) }, { status });
   }
 }
\ No newline at end of file
