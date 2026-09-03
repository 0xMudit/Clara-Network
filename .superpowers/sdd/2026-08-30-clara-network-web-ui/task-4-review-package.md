BASE 4c2c60a

7149a29 feat(web): supabase auth clients, roles, auth middleware

=== STAT ===
 web/package-lock.json              | 128 +++++++++++++++++++++++++++++++++++++
 web/package.json                   |   2 +
 web/src/lib/roles.ts               |  21 ++++++
 web/src/lib/supabase/middleware.ts |  26 ++++++++
 web/src/lib/supabase/server.ts     |  11 ++++
 web/src/middleware.ts              |   7 ++
 6 files changed, 195 insertions(+)
diff --git a/web/package-lock.json b/web/package-lock.json
index 0604189..9904ca6 100644
--- a/web/package-lock.json
+++ b/web/package-lock.json
@@ -2,20 +2,22 @@
   "name": "web",
   "version": "0.1.0",
   "lockfileVersion": 3,
   "requires": true,
   "packages": {
     "": {
       "name": "web",
       "version": "0.1.0",
       "dependencies": {
         "@base-ui/react": "^1.7.0",
+        "@supabase/ssr": "^0.12.5",
+        "@supabase/supabase-js": "^2.112.4",
         "class-variance-authority": "^0.7.1",
         "clsx": "^2.1.1",
         "lucide-react": "^1.37.0",
         "next": "16.3.3",
         "next-themes": "^0.4.6",
         "react": "19.2.8",
         "react-dom": "19.2.8",
         "shadcn": "^4.19.0",
         "tailwind-merge": "^3.6.0",
         "tw-animate-css": "^1.4.0"
@@ -1877,20 +1879,137 @@
       "resolved": "https://registry.npmjs.org/@sindresorhus/merge-streams/-/merge-streams-4.0.0.tgz",
       "integrity": "sha512-tlqY9xq5ukxTUZBmoOp+m61cqwQD5pHJtFY3Mn8CA8ps6yghLH/Hw8UPdqg4OLmFW3IFlcXnQNmo/dh8HzXYIQ==",
       "license": "MIT",
       "engines": {
         "node": ">=18"
       },
       "funding": {
         "url": "https://github.com/sponsors/sindresorhus"
       }
     },
+    "node_modules/@supabase/auth-js": {
+      "version": "2.112.4",
+      "resolved": "https://registry.npmjs.org/@supabase/auth-js/-/auth-js-2.112.4.tgz",
+      "integrity": "sha512-z8DesgwLzKM5PiT0yNmJU8VJyh1zAhYi+20Z7drdJQLXg/wWW4yGt/un+He5ERYUo94Vz66t5aeyr1DIDemI5A==",
+      "license": "MIT",
+      "dependencies": {
+        "tslib": "2.8.1"
+      },
+      "engines": {
+        "node": ">=22.0.0"
+      }
+    },
+    "node_modules/@supabase/functions-js": {
+      "version": "2.112.4",
+      "resolved": "https://registry.npmjs.org/@supabase/functions-js/-/functions-js-2.112.4.tgz",
+      "integrity": "sha512-DQ0aVH8wSQAccVqNoEkec62qCu2QRNyoGN53RqsVZ1k6F1zq4/v8scrlR6LNT2RJmT97apiTmORijPVhErCS2g==",
+      "license": "MIT",
+      "dependencies": {
+        "tslib": "2.8.1"
+      },
+      "engines": {
+        "node": ">=22.0.0"
+      }
+    },
+    "node_modules/@supabase/phoenix": {
+      "version": "0.4.5",
+      "resolved": "https://registry.npmjs.org/@supabase/phoenix/-/phoenix-0.4.5.tgz",
+      "integrity": "sha512-aAn9H9ovVyeApKy11OWOrrOGq8DV68yWeH4ud2lN9fzn4aO8Zb5GLL9m1pUg9nLqIcT+ZDfAcsZe0E/nqdv2lw==",
+      "license": "MIT"
+    },
+    "node_modules/@supabase/postgrest-js": {
+      "version": "2.112.4",
+      "resolved": "https://registry.npmjs.org/@supabase/postgrest-js/-/postgrest-js-2.112.4.tgz",
+      "integrity": "sha512-uaubtPSeg2TR4wrtfQoQWgkTAe+a0qWX2KhmwvTfNl5mGN9+U7owiJt6abk3o/V6O899PSRD1yzxs5RlF4xTug==",
+      "license": "MIT",
+      "dependencies": {
+        "tslib": "2.8.1"
+      },
+      "engines": {
+        "node": ">=22.0.0"
+      }
+    },
+    "node_modules/@supabase/realtime-js": {
+      "version": "2.112.4",
+      "resolved": "https://registry.npmjs.org/@supabase/realtime-js/-/realtime-js-2.112.4.tgz",
+      "integrity": "sha512-vZ+j079SKrM0Xiq7MJCvQKLDpaH2kfKfLY68xuQE1sqsCsMmx1CyrDBJHsxZ3cX01VOs5SI9igmoZAF3BmdZxw==",
+      "license": "MIT",
+      "dependencies": {
+        "@supabase/phoenix": "0.4.5",
+        "tslib": "2.8.1"
+      },
+      "engines": {
+        "node": ">=22.0.0"
+      }
+    },
+    "node_modules/@supabase/ssr": {
+      "version": "0.12.5",
+      "resolved": "https://registry.npmjs.org/@supabase/ssr/-/ssr-0.12.5.tgz",
+      "integrity": "sha512-0GllaAtHe7FHs6tlSg1CwkFw0a7aai1Shs9rEqSs/HaVkJzPk35C9r0Z8WUtKJ3eY6G/Y5y47YvCro3LynrAJg==",
+      "license": "MIT",
+      "dependencies": {
+        "cookie": "^1.0.2"
+      },
+      "peerDependencies": {
+        "@supabase/supabase-js": "^2.112.4"
+      }
+    },
+    "node_modules/@supabase/ssr/node_modules/cookie": {
+      "version": "1.1.1",
+      "resolved": "https://registry.npmjs.org/cookie/-/cookie-1.1.1.tgz",
+      "integrity": "sha512-ei8Aos7ja0weRpFzJnEA9UHJ/7XQmqglbRwnf2ATjcB9Wq874VKH9kfjjirM6UhU2/E5fFYadylyhFldcqSidQ==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=18"
+      },
+      "funding": {
+        "type": "opencollective",
+        "url": "https://opencollective.com/express"
+      }
+    },
+    "node_modules/@supabase/storage-js": {
+      "version": "2.112.4",
+      "resolved": "https://registry.npmjs.org/@supabase/storage-js/-/storage-js-2.112.4.tgz",
+      "integrity": "sha512-lQ0JemuTlMIXVKgSci1qez8yPnM5hyDngeAfEBjZS2Om4D+Cus0EE5BE6glFobrxdyii1OF4UzWfF0zcQgDq5A==",
+      "license": "MIT",
+      "dependencies": {
+        "iceberg-js": "^0.8.1",
+        "tslib": "2.8.1"
+      },
+      "engines": {
+        "node": ">=22.0.0"
+      }
+    },
+    "node_modules/@supabase/supabase-js": {
+      "version": "2.112.4",
+      "resolved": "https://registry.npmjs.org/@supabase/supabase-js/-/supabase-js-2.112.4.tgz",
+      "integrity": "sha512-UiCX1udlFY1fQQrO7Z3GU7obQsju0w5Vk9mOOwalfo/+Gy+tahWVenSSuu5E/GTy/q//HxvGv2IrCdW66/61kw==",
+      "license": "MIT",
+      "dependencies": {
+        "@supabase/auth-js": "2.112.4",
+        "@supabase/functions-js": "2.112.4",
+        "@supabase/postgrest-js": "2.112.4",
+        "@supabase/realtime-js": "2.112.4",
+        "@supabase/storage-js": "2.112.4"
+      },
+      "engines": {
+        "node": ">=22.0.0"
+      },
+      "peerDependencies": {
+        "@opentelemetry/api": ">=1.0.0"
+      },
+      "peerDependenciesMeta": {
+        "@opentelemetry/api": {
+          "optional": true
+        }
+      }
+    },
     "node_modules/@swc/helpers": {
       "version": "0.5.23",
       "license": "Apache-2.0",
       "dependencies": {
         "tslib": "^2.8.0"
       }
     },
     "node_modules/@tailwindcss/node": {
       "version": "4.3.3",
       "dev": true,
@@ -5297,20 +5416,29 @@
     },
     "node_modules/human-signals": {
       "version": "8.0.1",
       "resolved": "https://registry.npmjs.org/human-signals/-/human-signals-8.0.1.tgz",
       "integrity": "sha512-eKCa6bwnJhvxj14kZk5NCPc6Hb6BdsU9DZcOnmQKSnO1VKrfV0zCvtttPZUsBvjmNDn8rpcJfpwSYnHBjc95MQ==",
       "license": "Apache-2.0",
       "engines": {
         "node": ">=18.18.0"
       }
     },
+    "node_modules/iceberg-js": {
+      "version": "0.8.1",
+      "resolved": "https://registry.npmjs.org/iceberg-js/-/iceberg-js-0.8.1.tgz",
+      "integrity": "sha512-1dhVQZXhcHje7798IVM+xoo/1ZdVfzOMIc8/rgVSijRK38EDqOJoGula9N/8ZI5RD8QTxNQtK/Gozpr+qUqRRA==",
+      "license": "MIT",
+      "engines": {
+        "node": ">=20.0.0"
+      }
+    },
     "node_modules/iconv-lite": {
       "version": "0.7.3",
       "resolved": "https://registry.npmjs.org/iconv-lite/-/iconv-lite-0.7.3.tgz",
       "integrity": "sha512-IKXpvIzjnC9XTAUbVBcMfGS0EPaIXtW6v+zr+RRp+hqULEpo0owZax6wyRwPOJbWbzjYspQwusTsfVr0ifh4uQ==",
       "license": "MIT",
       "dependencies": {
         "safer-buffer": ">= 2.1.2 < 3.0.0"
       },
       "engines": {
         "node": ">=0.10.0"
diff --git a/web/package.json b/web/package.json
index 455e64c..ab512ab 100644
--- a/web/package.json
+++ b/web/package.json
@@ -3,20 +3,22 @@
   "version": "0.1.0",
   "private": true,
   "scripts": {
     "dev": "next dev",
     "build": "next build",
     "start": "next start",
     "lint": "eslint"
   },
   "dependencies": {
     "@base-ui/react": "^1.7.0",
+    "@supabase/ssr": "^0.12.5",
+    "@supabase/supabase-js": "^2.112.4",
     "class-variance-authority": "^0.7.1",
     "clsx": "^2.1.1",
     "lucide-react": "^1.37.0",
     "next": "16.3.3",
     "next-themes": "^0.4.6",
     "react": "19.2.8",
     "react-dom": "19.2.8",
     "shadcn": "^4.19.0",
     "tailwind-merge": "^3.6.0",
     "tw-animate-css": "^1.4.0"
diff --git a/web/src/lib/roles.ts b/web/src/lib/roles.ts
new file mode 100644
index 0000000..3c2524b
--- /dev/null
+++ b/web/src/lib/roles.ts
@@ -0,0 +1,21 @@
+export type Role = "scheme_operator" | "issuer" | "acquirer" | "merchant" | "viewer";
+export const ROLE_LABEL: Record<Role, string> = {
+  scheme_operator: "Scheme Operator", issuer: "Issuer", acquirer: "Acquirer",
+  merchant: "Merchant", viewer: "Viewer (HR)",
+};
+export function roleFromAppMetadata(meta?: Record<string, unknown>): Role | null {
+  const r = meta?.role;
+  return r === "scheme_operator" || r === "issuer" || r === "acquirer" || r === "merchant" || r === "viewer" ? r : null;
+}
+export const HOME_BY_ROLE: Record<Role, string> = {
+  scheme_operator: "/ops", issuer: "/issuer", acquirer: "/acquirer",
+  merchant: "/merchant", viewer: "/overview",
+};
+export const DASHBOARD_ACCESS: Record<Role, string[]> = {
+  scheme_operator: ["/ops", "/transactions", "/clearing", "/settlement", "/ledger", "/cards", "/merchants", "/disputes"],
+  issuer: ["/issuer", "/cards", "/tokens"],
+  acquirer: ["/acquirer", "/merchants", "/funding", "/disputes"],
+  merchant: ["/merchant", "/funding", "/disputes"],
+  viewer: ["/overview"],
+};
+export function dashboardAccess(role: Role): string[] { return DASHBOARD_ACCESS[role] ?? []; }
\ No newline at end of file
diff --git a/web/src/lib/supabase/middleware.ts b/web/src/lib/supabase/middleware.ts
new file mode 100644
index 0000000..a781cf5
--- /dev/null
+++ b/web/src/lib/supabase/middleware.ts
@@ -0,0 +1,26 @@
+import { createServerClient } from "@supabase/ssr";
+import { NextResponse, type NextRequest } from "next/server";
+
+export async function updateSession(request: NextRequest) {
+  // eslint-disable-next-line @typescript-eslint/no-require-imports
+  const { SUPABASE_URL, SUPABASE_ANON_KEY } = require("../env").getEnv();
+  let response = NextResponse.next({ request });
+  const supabase = createServerClient(SUPABASE_URL, SUPABASE_ANON_KEY, {
+    cookies: {
+      getAll: () => request.cookies.getAll(),
+      setAll: (cookiesToSet) => {
+        cookiesToSet.forEach(({ name, value }) => request.cookies.set(name, value));
+        response = NextResponse.next({ request });
+        cookiesToSet.forEach(({ name, value }) => response.cookies.set(name, value));
+      },
+    },
+  });
+  const { data } = await supabase.auth.getUser();
+  if (!data.user) {
+    const url = request.nextUrl.clone();
+    url.pathname = "/login";
+    url.searchParams.set("next", request.nextUrl.pathname);
+    return NextResponse.redirect(url);
+  }
+  return response;
+}
\ No newline at end of file
diff --git a/web/src/lib/supabase/server.ts b/web/src/lib/supabase/server.ts
new file mode 100644
index 0000000..6d368ca
--- /dev/null
+++ b/web/src/lib/supabase/server.ts
@@ -0,0 +1,11 @@
+import { createServerClient as createSupabaseServerClient } from "@supabase/ssr";
+import { cookies } from "next/headers";
+
+export async function createServerClient() {
+  // eslint-disable-next-line @typescript-eslint/no-require-imports
+  const { SUPABASE_URL, SUPABASE_ANON_KEY } = require("../env").getEnv();
+  const cookieStore = await cookies();
+  return createSupabaseServerClient(SUPABASE_URL, SUPABASE_ANON_KEY, {
+    cookies: { getAll: () => cookieStore.getAll(), setAll: () => {} },
+  });
+}
\ No newline at end of file
diff --git a/web/src/middleware.ts b/web/src/middleware.ts
new file mode 100644
index 0000000..26c3bc7
--- /dev/null
+++ b/web/src/middleware.ts
@@ -0,0 +1,7 @@
+import { updateSession } from "@/lib/supabase/middleware";
+import { type NextRequest } from "next/server";
+
+export async function middleware(request: NextRequest) {
+  return updateSession(request);
+}
+export const config = { matcher: ["/((?!_next/static|_next/image|favicon.ico|.*\\.(?:svg|png|jpg|jpeg|gif|webp)$).*)"] };
\ No newline at end of file
