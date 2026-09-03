BASE 5bd97ac

4c2c60a feat(web): env config, types, money/date helpers

=== STAT ===
 web/.env.local.example     |  5 +++++
 web/.gitignore             |  3 ++-
 web/src/lib/date.ts        |  5 +++++
 web/src/lib/env.ts         | 10 ++++++++++
 web/src/lib/money/minor.ts |  8 ++++++++
 web/src/types/admin.ts     |  5 +++++
 6 files changed, 35 insertions(+), 1 deletion(-)
diff --git a/web/.env.local.example b/web/.env.local.example
new file mode 100644
index 0000000..ae7e332
--- /dev/null
+++ b/web/.env.local.example
@@ -0,0 +1,5 @@
+# Supabase project (project settings -> API)
+SUPABASE_URL=https://<ref>.supabase.co
+SUPABASE_ANON_KEY=<anon key>
+# Railway Go adminapi URL (or http://localhost:18083 during local dev)
+CLARA_API_URL=http://localhost:18083
diff --git a/web/.gitignore b/web/.gitignore
index 5ef6a52..ae0b095 100644
--- a/web/.gitignore
+++ b/web/.gitignore
@@ -24,18 +24,19 @@
 .DS_Store
 *.pem
 
 # debug
 npm-debug.log*
 yarn-debug.log*
 yarn-error.log*
 .pnpm-debug.log*
 
 # env files (can opt-in for committing if needed)
-.env*
+.env*.local
+!.env.local.example
 
 # vercel
 .vercel
 
 # typescript
 *.tsbuildinfo
 next-env.d.ts
diff --git a/web/src/lib/date.ts b/web/src/lib/date.ts
new file mode 100644
index 0000000..e623513
--- /dev/null
+++ b/web/src/lib/date.ts
@@ -0,0 +1,5 @@
+export function fmtTs(iso: string): string {
+  const d = new Date(iso);
+  if (Number.isNaN(d.getTime())) return iso;
+  return d.toLocaleString("en-GB", { dateStyle: "medium", timeStyle: "short" });
+}
diff --git a/web/src/lib/env.ts b/web/src/lib/env.ts
new file mode 100644
index 0000000..3db5243
--- /dev/null
+++ b/web/src/lib/env.ts
@@ -0,0 +1,10 @@
+const required = ["SUPABASE_URL", "SUPABASE_ANON_KEY", "CLARA_API_URL"] as const;
+export function getEnv(): Record<typeof required[number], string> {
+  const out = {} as Record<typeof required[number], string>;
+  for (const k of required) {
+    const v = process.env[k];
+    if (!v) throw new Error(`missing env ${k}`);
+    out[k] = v;
+  }
+  return out;
+}
diff --git a/web/src/lib/money/minor.ts b/web/src/lib/money/minor.ts
new file mode 100644
index 0000000..8dbcbd5
--- /dev/null
+++ b/web/src/lib/money/minor.ts
@@ -0,0 +1,8 @@
+export function fmtMinor(minor: number, currency = "EUR"): string {
+  const sign = minor < 0 ? "-" : "";
+  const abs = Math.abs(minor);
+  const major = Math.floor(abs / 100);
+  const frac = String(abs % 100).padStart(2, "0");
+  const intl = new Intl.NumberFormat("en-US").format(major);
+  return `${sign} ${intl}.${frac} ${currency}`.trim();
+}
diff --git a/web/src/types/admin.ts b/web/src/types/admin.ts
new file mode 100644
index 0000000..e052c33
--- /dev/null
+++ b/web/src/types/admin.ts
@@ -0,0 +1,5 @@
+export interface Page<T> { items: T[]; total: number }
+export interface DashboardSummary {
+  transactions: number; clearingRecords: number; merchants: number;
+  disputes: number; cards: number; tokens: number;
+}
