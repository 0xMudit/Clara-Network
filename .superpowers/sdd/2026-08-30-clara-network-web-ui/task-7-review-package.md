BASE 390ee16

840d198 feat(web): authenticated app shell with role nav

=== STAT ===
 web/src/app/(app)/acquirer/page.tsx | 10 +++++
 web/src/app/(app)/issuer/page.tsx   | 10 +++++
 web/src/app/(app)/layout.tsx        | 19 +++++++++
 web/src/app/(app)/merchant/page.tsx | 10 +++++
 web/src/app/(app)/ops/page.tsx      | 10 +++++
 web/src/app/(app)/overview/page.tsx |  3 ++
 web/src/app/page.tsx                | 79 ++++++-------------------------------
 web/src/components/nav.tsx          | 36 +++++++++++++++++
 8 files changed, 109 insertions(+), 68 deletions(-)
diff --git a/web/src/app/(app)/acquirer/page.tsx b/web/src/app/(app)/acquirer/page.tsx
new file mode 100644
index 0000000..2226234
--- /dev/null
+++ b/web/src/app/(app)/acquirer/page.tsx
@@ -0,0 +1,10 @@
+import { notFound } from "next/navigation";
+import { createServerClient } from "@/lib/supabase/server";
+import { roleFromAppMetadata } from "@/lib/roles";
+
+export default async function AcquirerPage() {
+  const supabase = await createServerClient();
+  const { data } = await supabase.auth.getUser();
+  if (roleFromAppMetadata(data.user?.app_metadata) !== "acquirer") notFound();
+  return <h1 className="text-2xl font-semibold">Acquirer</h1>;
+}
\ No newline at end of file
diff --git a/web/src/app/(app)/issuer/page.tsx b/web/src/app/(app)/issuer/page.tsx
new file mode 100644
index 0000000..92557fb
--- /dev/null
+++ b/web/src/app/(app)/issuer/page.tsx
@@ -0,0 +1,10 @@
+import { notFound } from "next/navigation";
+import { createServerClient } from "@/lib/supabase/server";
+import { roleFromAppMetadata } from "@/lib/roles";
+
+export default async function IssuerPage() {
+  const supabase = await createServerClient();
+  const { data } = await supabase.auth.getUser();
+  if (roleFromAppMetadata(data.user?.app_metadata) !== "issuer") notFound();
+  return <h1 className="text-2xl font-semibold">Issuer</h1>;
+}
\ No newline at end of file
diff --git a/web/src/app/(app)/layout.tsx b/web/src/app/(app)/layout.tsx
new file mode 100644
index 0000000..91357f5
--- /dev/null
+++ b/web/src/app/(app)/layout.tsx
@@ -0,0 +1,19 @@
+import Nav from "@/components/nav";
+import { createServerClient } from "@/lib/supabase/server";
+import { roleFromAppMetadata } from "@/lib/roles";
+import { notFound } from "next/navigation";
+
+export const dynamic = "force-dynamic";
+
+export default async function AppLayout({ children }: { children: React.ReactNode }) {
+  const supabase = await createServerClient();
+  const { data } = await supabase.auth.getUser();
+  const role = roleFromAppMetadata(data.user?.app_metadata);
+  if (!role) notFound();
+  return (
+    <div className="min-h-svh">
+      <Nav role={role} />
+      <main className="mx-auto max-w-6xl px-4 py-6">{children}</main>
+    </div>
+  );
+}
\ No newline at end of file
diff --git a/web/src/app/(app)/merchant/page.tsx b/web/src/app/(app)/merchant/page.tsx
new file mode 100644
index 0000000..fca4855
--- /dev/null
+++ b/web/src/app/(app)/merchant/page.tsx
@@ -0,0 +1,10 @@
+import { notFound } from "next/navigation";
+import { createServerClient } from "@/lib/supabase/server";
+import { roleFromAppMetadata } from "@/lib/roles";
+
+export default async function MerchantPage() {
+  const supabase = await createServerClient();
+  const { data } = await supabase.auth.getUser();
+  if (roleFromAppMetadata(data.user?.app_metadata) !== "merchant") notFound();
+  return <h1 className="text-2xl font-semibold">Merchant</h1>;
+}
\ No newline at end of file
diff --git a/web/src/app/(app)/ops/page.tsx b/web/src/app/(app)/ops/page.tsx
new file mode 100644
index 0000000..9f9c057
--- /dev/null
+++ b/web/src/app/(app)/ops/page.tsx
@@ -0,0 +1,10 @@
+import { notFound } from "next/navigation";
+import { createServerClient } from "@/lib/supabase/server";
+import { roleFromAppMetadata } from "@/lib/roles";
+
+export default async function OpsPage() {
+  const supabase = await createServerClient();
+  const { data } = await supabase.auth.getUser();
+  if (roleFromAppMetadata(data.user?.app_metadata) !== "scheme_operator") notFound();
+  return <h1 className="text-2xl font-semibold">Operations</h1>;
+}
\ No newline at end of file
diff --git a/web/src/app/(app)/overview/page.tsx b/web/src/app/(app)/overview/page.tsx
new file mode 100644
index 0000000..8d0ec73
--- /dev/null
+++ b/web/src/app/(app)/overview/page.tsx
@@ -0,0 +1,3 @@
+export default async function OverviewPage() {
+  return <h1 className="text-2xl font-semibold">Overview</h1>;
+}
\ No newline at end of file
diff --git a/web/src/app/page.tsx b/web/src/app/page.tsx
index c887311..98f7a0f 100644
--- a/web/src/app/page.tsx
+++ b/web/src/app/page.tsx
@@ -1,69 +1,12 @@
-import Image from "next/image";
+import { redirect } from "next/navigation";
+import { createServerClient } from "@/lib/supabase/server";
+import { roleFromAppMetadata, HOME_BY_ROLE } from "@/lib/roles";
 
-export default function Home() {
-  return (
-    <div className="flex flex-col flex-1 items-center justify-center bg-zinc-50 font-sans dark:bg-black">
-      <main className="flex flex-1 w-full max-w-3xl flex-col items-center justify-between py-32 px-16 bg-white dark:bg-black sm:items-start">
-        <Image
-          className="dark:invert h-5 w-[100px]"
-          src="/next.svg"
-          alt="Next.js logo"
-          width={100}
-          height={20}
-          priority
-        />
-        <div className="flex flex-col items-center gap-6 text-center sm:items-start sm:text-left">
-          <h1 className="max-w-xs text-3xl font-semibold leading-10 tracking-tight text-black dark:text-zinc-50">
-            To get started, edit the{" "}
-            <code className="rounded bg-black/[.06] px-1.5 py-0.5 font-mono text-[0.9em] dark:bg-white/[.08]">
-              page.tsx
-            </code>{" "}
-            file.
-          </h1>
-          <p className="max-w-md text-lg leading-8 text-zinc-600 dark:text-zinc-400">
-            Looking for a starting point or more instructions? Head over to{" "}
-            <a
-              href="https://vercel.com/templates?framework=next.js&utm_source=create-next-app&utm_medium=appdir-template-tw&utm_campaign=create-next-app"
-              className="font-medium text-zinc-950 dark:text-zinc-50"
-            >
-              Templates
-            </a>{" "}
-            or the{" "}
-            <a
-              href="https://nextjs.org/learn?utm_source=create-next-app&utm_medium=appdir-template-tw&utm_campaign=create-next-app"
-              className="font-medium text-zinc-950 dark:text-zinc-50"
-            >
-              Learning
-            </a>{" "}
-            center.
-          </p>
-        </div>
-        <div className="flex flex-col gap-4 text-base font-medium sm:flex-row">
-          <a
-            className="flex h-12 w-full items-center justify-center gap-2 rounded-full bg-foreground px-5 text-background transition-colors hover:bg-[#383838] dark:hover:bg-[#ccc] md:w-[158px]"
-            href="https://vercel.com/new?utm_source=create-next-app&utm_medium=appdir-template-tw&utm_campaign=create-next-app"
-            target="_blank"
-            rel="noopener noreferrer"
-          >
-            <Image
-              className="dark:invert h-[14px] w-4"
-              src="/vercel.svg"
-              alt="Vercel logomark"
-              width={16}
-              height={14}
-            />
-            Deploy Now
-          </a>
-          <a
-            className="flex h-12 w-full items-center justify-center rounded-full border border-solid border-black/[.08] px-5 transition-colors hover:border-transparent hover:bg-black/[.04] dark:border-white/[.145] dark:hover:bg-[#1a1a1a] md:w-[158px]"
-            href="https://nextjs.org/docs?utm_source=create-next-app&utm_medium=appdir-template-tw&utm_campaign=create-next-app"
-            target="_blank"
-            rel="noopener noreferrer"
-          >
-            Documentation
-          </a>
-        </div>
-      </main>
-    </div>
-  );
-}
+export const dynamic = "force-dynamic";
+
+export default async function Home() {
+  const supabase = await createServerClient();
+  const { data } = await supabase.auth.getUser();
+  const role = roleFromAppMetadata(data.user?.app_metadata);
+  redirect(role ? HOME_BY_ROLE[role] : "/login");
+}
\ No newline at end of file
diff --git a/web/src/components/nav.tsx b/web/src/components/nav.tsx
new file mode 100644
index 0000000..4b7c75f
--- /dev/null
+++ b/web/src/components/nav.tsx
@@ -0,0 +1,36 @@
+import { HOME_BY_ROLE, DASHBOARD_ACCESS, ROLE_LABEL, type Role } from "@/lib/roles";
+import { LogoutButton } from "@/components/logout-button";
+import { ThemeToggle } from "@/components/theme-toggle";
+import Link from "next/link";
+
+const TITLES: Record<string, string> = {
+  "/overview": "Overview", "/ops": "Operations", "/issuer": "Issuer", "/acquirer": "Acquirer",
+  "/merchant": "Merchant", "/transactions": "Transactions", "/clearing": "Clearing",
+  "/settlement": "Settlement", "/ledger": "Ledger", "/cards": "Cards", "/tokens": "Tokens",
+  "/merchants": "Merchants", "/funding": "Funding", "/disputes": "Disputes",
+};
+
+export default function Nav({ role }: { role: Role }) {
+  const items = DASHBOARD_ACCESS[role] ?? [];
+  return (
+    <header className="sticky top-0 z-40 border-b bg-background/95 backdrop-blur">
+      <div className="flex h-14 items-center justify-between px-4">
+        <div className="flex items-center gap-6">
+          <Link href={HOME_BY_ROLE[role]} className="font-semibold tracking-tight">Clara Network</Link>
+          <nav className="hidden items-center gap-1 md:flex">
+            {items.map(p => (
+              <Link key={p} href={p} className="rounded-md px-3 py-1.5 text-sm text-muted-foreground hover:bg-muted hover:text-foreground">
+                {TITLES[p] ?? p}
+              </Link>
+            ))}
+          </nav>
+        </div>
+        <div className="flex items-center gap-3">
+          <span className="text-sm text-muted-foreground">{ROLE_LABEL[role]}</span>
+          <ThemeToggle />
+          <LogoutButton />
+        </div>
+      </div>
+    </header>
+  );
+}
\ No newline at end of file
