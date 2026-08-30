BASE adfa60b

42ee415 feat(web): BFF data proxy, role-gated admin API access

=== STAT ===
 web/package-lock.json                   |  7 ++++++
 web/package.json                        |  1 +
 web/src/app/api/data/[...path]/route.ts | 41 +++++++++++++++++++++++++++++++++
 web/src/lib/clara.ts                    | 14 +++++++++++
 4 files changed, 63 insertions(+)
diff --git a/web/package-lock.json b/web/package-lock.json
index 9904ca6..f902622 100644
--- a/web/package-lock.json
+++ b/web/package-lock.json
@@ -11,20 +11,21 @@
         "@base-ui/react": "^1.7.0",
         "@supabase/ssr": "^0.12.5",
         "@supabase/supabase-js": "^2.112.4",
         "class-variance-authority": "^0.7.1",
         "clsx": "^2.1.1",
         "lucide-react": "^1.37.0",
         "next": "16.3.3",
         "next-themes": "^0.4.6",
         "react": "19.2.8",
         "react-dom": "19.2.8",
+        "server-only": "^0.0.1",
         "shadcn": "^4.19.0",
         "tailwind-merge": "^3.6.0",
         "tw-animate-css": "^1.4.0"
       },
       "devDependencies": {
         "@tailwindcss/postcss": "^4",
         "@types/node": "^20",
         "@types/react": "^19",
         "@types/react-dom": "^19",
         "eslint": "^9",
@@ -7918,20 +7919,26 @@
         "send": "^1.2.0"
       },
       "engines": {
         "node": ">= 18"
       },
       "funding": {
         "type": "opencollective",
         "url": "https://opencollective.com/express"
       }
     },
+    "node_modules/server-only": {
+      "version": "0.0.1",
+      "resolved": "https://registry.npmjs.org/server-only/-/server-only-0.0.1.tgz",
+      "integrity": "sha512-qepMx2JxAa5jjfzxG79yPPq+8BuFToHd1hm7kI+Z4zAq1ftQiP7HcxMhDDItrbtwVeLg/cY2JnKnrcFkmiswNA==",
+      "license": "MIT"
+    },
     "node_modules/set-function-length": {
       "version": "1.2.2",
       "dev": true,
       "license": "MIT",
       "dependencies": {
         "define-data-property": "^1.1.4",
         "es-errors": "^1.3.0",
         "function-bind": "^1.1.2",
         "get-intrinsic": "^1.2.4",
         "gopd": "^1.0.1",
diff --git a/web/package.json b/web/package.json
index ab512ab..fff09c8 100644
--- a/web/package.json
+++ b/web/package.json
@@ -12,20 +12,21 @@
     "@base-ui/react": "^1.7.0",
     "@supabase/ssr": "^0.12.5",
     "@supabase/supabase-js": "^2.112.4",
     "class-variance-authority": "^0.7.1",
     "clsx": "^2.1.1",
     "lucide-react": "^1.37.0",
     "next": "16.3.3",
     "next-themes": "^0.4.6",
     "react": "19.2.8",
     "react-dom": "19.2.8",
+    "server-only": "^0.0.1",
     "shadcn": "^4.19.0",
     "tailwind-merge": "^3.6.0",
     "tw-animate-css": "^1.4.0"
   },
   "devDependencies": {
     "@tailwindcss/postcss": "^4",
     "@types/node": "^20",
     "@types/react": "^19",
     "@types/react-dom": "^19",
     "eslint": "^9",
diff --git a/web/src/app/api/data/[...path]/route.ts b/web/src/app/api/data/[...path]/route.ts
new file mode 100644
index 0000000..6c51177
--- /dev/null
+++ b/web/src/app/api/data/[...path]/route.ts
@@ -0,0 +1,41 @@
+// src/app/api/data/[...path]/route.ts
+import { NextResponse } from "next/server";
+import { createServerClient } from "@/lib/supabase/server";
+import { getEnv } from "@/lib/env";
+import { claraFetch } from "@/lib/clara";
+
+export const dynamic = "force-dynamic";
+
+type Res<T> = NextResponse<T | { error: string }>;
+
+const ALLOWED_ROUTES: Record<string, string[]> = {
+  scheme_operator: [
+    "/dashboard", "/transactions", "/clearing/cycles", "/clearing/records",
+    "/clearing/positions", "/settlement/instructions", "/settlement/prefunds",
+    "/settlement/default-fund", "/ledger/accounts", "/ledger/entries",
+  ],
+  issuer: ["/dashboard", "/cards", "/tokens"],
+  acquirer: ["/dashboard", "/acquirer", "/merchants", "/funding", "/disputes"],
+  merchant: ["/dashboard", "/merchant", "/funding", "/disputes"],
+};
+
+export async function GET(
+  req: Request,
+  { params }: { params: Promise<{ path: string[] }> },
+): Promise<Res<unknown>> {
+  const supabase = await createServerClient();
+  const { data, error } = await supabase.auth.getUser();
+  if (error || !data.user) return NextResponse.json({ error: "unauthorized" }, { status: 401 });
+  const role = data.user.app_metadata.role as string;
+  const { path } = await params;
+  const target = "/" + path.join("/");
+  const allowed = ALLOWED_ROUTES[role] ?? [];
+  if (!allowed.includes(target)) return NextResponse.json({ error: "forbidden" }, { status: 403 });
+  try {
+    const body = await claraFetch(target, data.user.id, { cache: "no-store" });
+    return NextResponse.json(body);
+  } catch (e) {
+    const status = e instanceof Error && /^clara \S+ 4\d\d/.test(e.message) ? 502 : 500;
+    return NextResponse.json({ error: "upstream error", detail: e instanceof Error ? e.message : String(e) }, { status });
+  }
+}
\ No newline at end of file
diff --git a/web/src/lib/clara.ts b/web/src/lib/clara.ts
new file mode 100644
index 0000000..5565877
--- /dev/null
+++ b/web/src/lib/clara.ts
@@ -0,0 +1,14 @@
+// src/lib/clara.ts
+import "server-only";
+import { getEnv } from "./env";
+
+export function claraFetch<T>(path: string, accessToken: string, init?: RequestInit): Promise<T> {
+  const { CLARA_API_URL } = getEnv();
+  return fetch(`${CLARA_API_URL}${path}`, {
+    headers: { Authorization: `Bearer ${accessToken}`, Accept: "application/json" },
+    ...init,
+  }).then(r => {
+    if (!r.ok) throw new Error(`clara ${path}: ${r.status}`);
+    return r.json() as Promise<T>;
+  });
+}
\ No newline at end of file
