BASE 7149a29

adfa60b fix(web): exempt /login from auth middleware redirect
d70c688 feat(web): login page, logout, demo role flows

=== STAT ===
 web/.env.local.example               |   3 +
 web/src/app/login/login-form.tsx     |  50 +++++++++++++++++
 web/src/app/login/page.tsx           |  21 +++++++
 web/src/components/logout-button.tsx |  15 +++++
 web/src/components/ui/card.tsx       | 103 +++++++++++++++++++++++++++++++++++
 web/src/components/ui/input.tsx      |  20 +++++++
 web/src/components/ui/label.tsx      |  20 +++++++
 web/src/lib/supabase/client.ts       |   9 +++
 web/src/lib/supabase/middleware.ts   |   1 +
 9 files changed, 242 insertions(+)
diff --git a/web/src/app/login/login-form.tsx b/web/src/app/login/login-form.tsx
new file mode 100644
index 0000000..ff4b502
--- /dev/null
+++ b/web/src/app/login/login-form.tsx
@@ -0,0 +1,50 @@
+"use client";
+import { useState } from "react";
+import { useRouter, useSearchParams } from "next/navigation";
+import { createBrowserClient } from "@/lib/supabase/client";
+import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
+import { Label } from "@/components/ui/label";
+import { Input } from "@/components/ui/input";
+import { Button } from "@/components/ui/button";
+import { roleFromAppMetadata, HOME_BY_ROLE } from "@/lib/roles";
+
+export function LoginForm() {
+  const router = useRouter();
+  const params = useSearchParams();
+  const [email, setEmail] = useState("");
+  const [password, setPassword] = useState("");
+  const [error, setError] = useState<string | null>(null);
+  const [loading, setLoading] = useState(false);
+
+  async function onSubmit(e: React.FormEvent) {
+    e.preventDefault();
+    setLoading(true);
+    setError(null);
+    const supabase = createBrowserClient();
+    const { data, error } = await supabase.auth.signInWithPassword({ email, password });
+    if (error) { setError(error.message); setLoading(false); return; }
+    const role = roleFromAppMetadata(data.session?.user.app_metadata);
+    const next = params.get("next");
+    router.push(next && next.startsWith("/") ? next : (role ? HOME_BY_ROLE[role] : "/overview"));
+    router.refresh();
+  }
+
+  return (
+    <Card className="w-full max-w-sm">
+      <CardHeader>
+        <CardTitle>Clara Network</CardTitle>
+        <CardDescription>Scheme operator console ÔÇö sign in to continue</CardDescription>
+      </CardHeader>
+      <CardContent>
+        <form onSubmit={onSubmit} className="grid gap-4">
+          <div className="grid gap-2"><Label htmlFor="email">Work email</Label>
+            <Input id="email" type="email" value={email} onChange={e => setEmail(e.target.value)} required autoComplete="email" /></div>
+          <div className="grid gap-2"><Label htmlFor="password">Password</Label>
+            <Input id="password" type="password" value={password} onChange={e => setPassword(e.target.value)} required autoComplete="current-password" /></div>
+          {error && <p className="text-sm text-destructive">{error}</p>}
+          <Button type="submit" disabled={loading}>{loading ? "Signing inÔÇª" : "Sign in"}</Button>
+        </form>
+      </CardContent>
+    </Card>
+  );
+}
diff --git a/web/src/app/login/page.tsx b/web/src/app/login/page.tsx
new file mode 100644
index 0000000..adaa84b
--- /dev/null
+++ b/web/src/app/login/page.tsx
@@ -0,0 +1,21 @@
+import { Suspense } from "react";
+import { redirect } from "next/navigation";
+import { createServerClient } from "@/lib/supabase/server";
+import { LoginForm } from "./login-form";
+import { roleFromAppMetadata, HOME_BY_ROLE } from "@/lib/roles";
+
+export const dynamic = "force-dynamic";
+
+export default async function LoginPage() {
+  const supabase = await createServerClient();
+  const { data } = await supabase.auth.getUser();
+  if (data.user) {
+    const role = roleFromAppMetadata(data.user.app_metadata);
+    redirect(role ? HOME_BY_ROLE[role] : "/");
+  }
+  return (
+    <main className="flex min-h-svh items-center justify-center p-6">
+      <Suspense><LoginForm /></Suspense>
+    </main>
+  );
+}
diff --git a/web/src/components/logout-button.tsx b/web/src/components/logout-button.tsx
new file mode 100644
index 0000000..c28cd7f
--- /dev/null
+++ b/web/src/components/logout-button.tsx
@@ -0,0 +1,15 @@
+"use client";
+import { useRouter } from "next/navigation";
+import { createBrowserClient } from "@/lib/supabase/client";
+import { Button } from "@/components/ui/button";
+
+export function LogoutButton() {
+  const router = useRouter();
+  async function logout() {
+    const supabase = createBrowserClient();
+    await supabase.auth.signOut();
+    router.replace("/login");
+    router.refresh();
+  }
+  return <Button variant="outline" onClick={logout}>Sign out</Button>;
+}
diff --git a/web/src/components/ui/card.tsx b/web/src/components/ui/card.tsx
new file mode 100644
index 0000000..4458dae
--- /dev/null
+++ b/web/src/components/ui/card.tsx
@@ -0,0 +1,103 @@
+import * as React from "react"
+
+import { cn } from "@/lib/utils"
+
+function Card({
+  className,
+  size = "default",
+  ...props
+}: React.ComponentProps<"div"> & { size?: "default" | "sm" }) {
+  return (
+    <div
+      data-slot="card"
+      data-size={size}
+      className={cn(
+        "group/card flex flex-col gap-(--card-spacing) overflow-hidden rounded-xl bg-card py-(--card-spacing) text-sm text-card-foreground ring-1 ring-foreground/10 [--card-spacing:--spacing(4)] has-data-[slot=card-footer]:pb-0 has-[>img:first-child]:pt-0 data-[size=sm]:[--card-spacing:--spacing(3)] data-[size=sm]:has-data-[slot=card-footer]:pb-0 *:[img:first-child]:rounded-t-xl *:[img:last-child]:rounded-b-xl",
+        className
+      )}
+      {...props}
+    />
+  )
+}
+
+function CardHeader({ className, ...props }: React.ComponentProps<"div">) {
+  return (
+    <div
+      data-slot="card-header"
+      className={cn(
+        "group/card-header @container/card-header grid auto-rows-min items-start gap-1 rounded-t-xl px-(--card-spacing) has-data-[slot=card-action]:grid-cols-[1fr_auto] has-data-[slot=card-description]:grid-rows-[auto_auto] [.border-b]:pb-(--card-spacing)",
+        className
+      )}
+      {...props}
+    />
+  )
+}
+
+function CardTitle({ className, ...props }: React.ComponentProps<"div">) {
+  return (
+    <div
+      data-slot="card-title"
+      className={cn(
+        "font-heading text-base leading-snug font-medium group-data-[size=sm]/card:text-sm",
+        className
+      )}
+      {...props}
+    />
+  )
+}
+
+function CardDescription({ className, ...props }: React.ComponentProps<"div">) {
+  return (
+    <div
+      data-slot="card-description"
+      className={cn("text-sm text-muted-foreground", className)}
+      {...props}
+    />
+  )
+}
+
+function CardAction({ className, ...props }: React.ComponentProps<"div">) {
+  return (
+    <div
+      data-slot="card-action"
+      className={cn(
+        "col-start-2 row-span-2 row-start-1 self-start justify-self-end",
+        className
+      )}
+      {...props}
+    />
+  )
+}
+
+function CardContent({ className, ...props }: React.ComponentProps<"div">) {
+  return (
+    <div
+      data-slot="card-content"
+      className={cn("px-(--card-spacing)", className)}
+      {...props}
+    />
+  )
+}
+
+function CardFooter({ className, ...props }: React.ComponentProps<"div">) {
+  return (
+    <div
+      data-slot="card-footer"
+      className={cn(
+        "flex items-center rounded-b-xl border-t bg-muted/50 p-(--card-spacing)",
+        className
+      )}
+      {...props}
+    />
+  )
+}
+
+export {
+  Card,
+  CardHeader,
+  CardFooter,
+  CardTitle,
+  CardAction,
+  CardDescription,
+  CardContent,
+}
diff --git a/web/src/components/ui/input.tsx b/web/src/components/ui/input.tsx
new file mode 100644
index 0000000..7d21bab
--- /dev/null
+++ b/web/src/components/ui/input.tsx
@@ -0,0 +1,20 @@
+import * as React from "react"
+import { Input as InputPrimitive } from "@base-ui/react/input"
+
+import { cn } from "@/lib/utils"
+
+function Input({ className, type, ...props }: React.ComponentProps<"input">) {
+  return (
+    <InputPrimitive
+      type={type}
+      data-slot="input"
+      className={cn(
+        "h-8 w-full min-w-0 rounded-lg border border-input bg-transparent px-2.5 py-1 text-base transition-colors outline-none file:inline-flex file:h-6 file:border-0 file:bg-transparent file:text-sm file:font-medium file:text-foreground placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 disabled:pointer-events-none disabled:cursor-not-allowed disabled:bg-input/50 disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-3 aria-invalid:ring-destructive/20 md:text-sm dark:bg-input/30 dark:disabled:bg-input/80 dark:aria-invalid:border-destructive/50 dark:aria-invalid:ring-destructive/40",
+        className
+      )}
+      {...props}
+    />
+  )
+}
+
+export { Input }
diff --git a/web/src/components/ui/label.tsx b/web/src/components/ui/label.tsx
new file mode 100644
index 0000000..74da65c
--- /dev/null
+++ b/web/src/components/ui/label.tsx
@@ -0,0 +1,20 @@
+"use client"
+
+import * as React from "react"
+
+import { cn } from "@/lib/utils"
+
+function Label({ className, ...props }: React.ComponentProps<"label">) {
+  return (
+    <label
+      data-slot="label"
+      className={cn(
+        "flex items-center gap-2 text-sm leading-none font-medium select-none group-data-[disabled=true]:pointer-events-none group-data-[disabled=true]:opacity-50 peer-disabled:cursor-not-allowed peer-disabled:opacity-50",
+        className
+      )}
+      {...props}
+    />
+  )
+}
+
+export { Label }
diff --git a/web/src/lib/supabase/client.ts b/web/src/lib/supabase/client.ts
new file mode 100644
index 0000000..089defa
--- /dev/null
+++ b/web/src/lib/supabase/client.ts
@@ -0,0 +1,9 @@
+"use client";
+import { createBrowserClient as createSupabaseBrowserClient } from "@supabase/ssr";
+
+export function createBrowserClient() {
+  return createSupabaseBrowserClient(
+    process.env.NEXT_PUBLIC_SUPABASE_URL!,
+    process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY!,
+  );
+}
diff --git a/web/src/lib/supabase/middleware.ts b/web/src/lib/supabase/middleware.ts
index a781cf5..fe20f07 100644
--- a/web/src/lib/supabase/middleware.ts
+++ b/web/src/lib/supabase/middleware.ts
@@ -9,18 +9,19 @@ export async function updateSession(request: NextRequest) {
     cookies: {
       getAll: () => request.cookies.getAll(),
       setAll: (cookiesToSet) => {
         cookiesToSet.forEach(({ name, value }) => request.cookies.set(name, value));
         response = NextResponse.next({ request });
         cookiesToSet.forEach(({ name, value }) => response.cookies.set(name, value));
       },
     },
   });
   const { data } = await supabase.auth.getUser();
+  if (request.nextUrl.pathname.startsWith("/login")) return response;
   if (!data.user) {
     const url = request.nextUrl.clone();
     url.pathname = "/login";
     url.searchParams.set("next", request.nextUrl.pathname);
     return NextResponse.redirect(url);
   }
   return response;
 }
\ No newline at end of file
diff --git a/web/.env.local.example b/web/.env.local.example
index ae7e332..41f54c0 100644
--- a/web/.env.local.example
+++ b/web/.env.local.example
@@ -1,5 +1,8 @@
 # Supabase project (project settings -> API)
 SUPABASE_URL=https://<ref>.supabase.co
 SUPABASE_ANON_KEY=<anon key>
 # Railway Go adminapi URL (or http://localhost:18083 during local dev)
 CLARA_API_URL=http://localhost:18083
+# Client-side public Supabase vars (next.js requires NEXT_PUBLIC_ prefix, values match SUPABASE_URL/SUPABASE_ANON_KEY)
+NEXT_PUBLIC_SUPABASE_URL=https://<ref>.supabase.co
+NEXT_PUBLIC_SUPABASE_ANON_KEY=<anon key>
diff --git a/web/src/app/login/login-form.tsx b/web/src/app/login/login-form.tsx
new file mode 100644
index 0000000..ff4b502
--- /dev/null
+++ b/web/src/app/login/login-form.tsx
@@ -0,0 +1,50 @@
+"use client";
+import { useState } from "react";
+import { useRouter, useSearchParams } from "next/navigation";
+import { createBrowserClient } from "@/lib/supabase/client";
+import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
+import { Label } from "@/components/ui/label";
+import { Input } from "@/components/ui/input";
+import { Button } from "@/components/ui/button";
+import { roleFromAppMetadata, HOME_BY_ROLE } from "@/lib/roles";
+
+export function LoginForm() {
+  const router = useRouter();
+  const params = useSearchParams();
+  const [email, setEmail] = useState("");
+  const [password, setPassword] = useState("");
+  const [error, setError] = useState<string | null>(null);
+  const [loading, setLoading] = useState(false);
+
+  async function onSubmit(e: React.FormEvent) {
+    e.preventDefault();
+    setLoading(true);
+    setError(null);
+    const supabase = createBrowserClient();
+    const { data, error } = await supabase.auth.signInWithPassword({ email, password });
+    if (error) { setError(error.message); setLoading(false); return; }
+    const role = roleFromAppMetadata(data.session?.user.app_metadata);
+    const next = params.get("next");
+    router.push(next && next.startsWith("/") ? next : (role ? HOME_BY_ROLE[role] : "/overview"));
+    router.refresh();
+  }
+
+  return (
+    <Card className="w-full max-w-sm">
+      <CardHeader>
+        <CardTitle>Clara Network</CardTitle>
+        <CardDescription>Scheme operator console ÔÇö sign in to continue</CardDescription>
+      </CardHeader>
+      <CardContent>
+        <form onSubmit={onSubmit} className="grid gap-4">
+          <div className="grid gap-2"><Label htmlFor="email">Work email</Label>
+            <Input id="email" type="email" value={email} onChange={e => setEmail(e.target.value)} required autoComplete="email" /></div>
+          <div className="grid gap-2"><Label htmlFor="password">Password</Label>
+            <Input id="password" type="password" value={password} onChange={e => setPassword(e.target.value)} required autoComplete="current-password" /></div>
+          {error && <p className="text-sm text-destructive">{error}</p>}
+          <Button type="submit" disabled={loading}>{loading ? "Signing inÔÇª" : "Sign in"}</Button>
+        </form>
+      </CardContent>
+    </Card>
+  );
+}
diff --git a/web/src/app/login/page.tsx b/web/src/app/login/page.tsx
new file mode 100644
index 0000000..adaa84b
--- /dev/null
+++ b/web/src/app/login/page.tsx
@@ -0,0 +1,21 @@
+import { Suspense } from "react";
+import { redirect } from "next/navigation";
+import { createServerClient } from "@/lib/supabase/server";
+import { LoginForm } from "./login-form";
+import { roleFromAppMetadata, HOME_BY_ROLE } from "@/lib/roles";
+
+export const dynamic = "force-dynamic";
+
+export default async function LoginPage() {
+  const supabase = await createServerClient();
+  const { data } = await supabase.auth.getUser();
+  if (data.user) {
+    const role = roleFromAppMetadata(data.user.app_metadata);
+    redirect(role ? HOME_BY_ROLE[role] : "/");
+  }
+  return (
+    <main className="flex min-h-svh items-center justify-center p-6">
+      <Suspense><LoginForm /></Suspense>
+    </main>
+  );
+}
diff --git a/web/src/components/logout-button.tsx b/web/src/components/logout-button.tsx
new file mode 100644
index 0000000..c28cd7f
--- /dev/null
+++ b/web/src/components/logout-button.tsx
@@ -0,0 +1,15 @@
+"use client";
+import { useRouter } from "next/navigation";
+import { createBrowserClient } from "@/lib/supabase/client";
+import { Button } from "@/components/ui/button";
+
+export function LogoutButton() {
+  const router = useRouter();
+  async function logout() {
+    const supabase = createBrowserClient();
+    await supabase.auth.signOut();
+    router.replace("/login");
+    router.refresh();
+  }
+  return <Button variant="outline" onClick={logout}>Sign out</Button>;
+}
diff --git a/web/src/components/ui/card.tsx b/web/src/components/ui/card.tsx
new file mode 100644
index 0000000..4458dae
--- /dev/null
+++ b/web/src/components/ui/card.tsx
@@ -0,0 +1,103 @@
+import * as React from "react"
+
+import { cn } from "@/lib/utils"
+
+function Card({
+  className,
+  size = "default",
+  ...props
+}: React.ComponentProps<"div"> & { size?: "default" | "sm" }) {
+  return (
+    <div
+      data-slot="card"
+      data-size={size}
+      className={cn(
+        "group/card flex flex-col gap-(--card-spacing) overflow-hidden rounded-xl bg-card py-(--card-spacing) text-sm text-card-foreground ring-1 ring-foreground/10 [--card-spacing:--spacing(4)] has-data-[slot=card-footer]:pb-0 has-[>img:first-child]:pt-0 data-[size=sm]:[--card-spacing:--spacing(3)] data-[size=sm]:has-data-[slot=card-footer]:pb-0 *:[img:first-child]:rounded-t-xl *:[img:last-child]:rounded-b-xl",
+        className
+      )}
+      {...props}
+    />
+  )
+}
+
+function CardHeader({ className, ...props }: React.ComponentProps<"div">) {
+  return (
+    <div
+      data-slot="card-header"
+      className={cn(
+        "group/card-header @container/card-header grid auto-rows-min items-start gap-1 rounded-t-xl px-(--card-spacing) has-data-[slot=card-action]:grid-cols-[1fr_auto] has-data-[slot=card-description]:grid-rows-[auto_auto] [.border-b]:pb-(--card-spacing)",
+        className
+      )}
+      {...props}
+    />
+  )
+}
+
+function CardTitle({ className, ...props }: React.ComponentProps<"div">) {
+  return (
+    <div
+      data-slot="card-title"
+      className={cn(
+        "font-heading text-base leading-snug font-medium group-data-[size=sm]/card:text-sm",
+        className
+      )}
+      {...props}
+    />
+  )
+}
+
+function CardDescription({ className, ...props }: React.ComponentProps<"div">) {
+  return (
+    <div
+      data-slot="card-description"
+      className={cn("text-sm text-muted-foreground", className)}
+      {...props}
+    />
+  )
+}
+
+function CardAction({ className, ...props }: React.ComponentProps<"div">) {
+  return (
+    <div
+      data-slot="card-action"
+      className={cn(
+        "col-start-2 row-span-2 row-start-1 self-start justify-self-end",
+        className
+      )}
+      {...props}
+    />
+  )
+}
+
+function CardContent({ className, ...props }: React.ComponentProps<"div">) {
+  return (
+    <div
+      data-slot="card-content"
+      className={cn("px-(--card-spacing)", className)}
+      {...props}
+    />
+  )
+}
+
+function CardFooter({ className, ...props }: React.ComponentProps<"div">) {
+  return (
+    <div
+      data-slot="card-footer"
+      className={cn(
+        "flex items-center rounded-b-xl border-t bg-muted/50 p-(--card-spacing)",
+        className
+      )}
+      {...props}
+    />
+  )
+}
+
+export {
+  Card,
+  CardHeader,
+  CardFooter,
+  CardTitle,
+  CardAction,
+  CardDescription,
+  CardContent,
+}
diff --git a/web/src/components/ui/input.tsx b/web/src/components/ui/input.tsx
new file mode 100644
index 0000000..7d21bab
--- /dev/null
+++ b/web/src/components/ui/input.tsx
@@ -0,0 +1,20 @@
+import * as React from "react"
+import { Input as InputPrimitive } from "@base-ui/react/input"
+
+import { cn } from "@/lib/utils"
+
+function Input({ className, type, ...props }: React.ComponentProps<"input">) {
+  return (
+    <InputPrimitive
+      type={type}
+      data-slot="input"
+      className={cn(
+        "h-8 w-full min-w-0 rounded-lg border border-input bg-transparent px-2.5 py-1 text-base transition-colors outline-none file:inline-flex file:h-6 file:border-0 file:bg-transparent file:text-sm file:font-medium file:text-foreground placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 disabled:pointer-events-none disabled:cursor-not-allowed disabled:bg-input/50 disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-3 aria-invalid:ring-destructive/20 md:text-sm dark:bg-input/30 dark:disabled:bg-input/80 dark:aria-invalid:border-destructive/50 dark:aria-invalid:ring-destructive/40",
+        className
+      )}
+      {...props}
+    />
+  )
+}
+
+export { Input }
diff --git a/web/src/components/ui/label.tsx b/web/src/components/ui/label.tsx
new file mode 100644
index 0000000..74da65c
--- /dev/null
+++ b/web/src/components/ui/label.tsx
@@ -0,0 +1,20 @@
+"use client"
+
+import * as React from "react"
+
+import { cn } from "@/lib/utils"
+
+function Label({ className, ...props }: React.ComponentProps<"label">) {
+  return (
+    <label
+      data-slot="label"
+      className={cn(
+        "flex items-center gap-2 text-sm leading-none font-medium select-none group-data-[disabled=true]:pointer-events-none group-data-[disabled=true]:opacity-50 peer-disabled:cursor-not-allowed peer-disabled:opacity-50",
+        className
+      )}
+      {...props}
+    />
+  )
+}
+
+export { Label }
diff --git a/web/src/lib/supabase/client.ts b/web/src/lib/supabase/client.ts
new file mode 100644
index 0000000..089defa
--- /dev/null
+++ b/web/src/lib/supabase/client.ts
@@ -0,0 +1,9 @@
+"use client";
+import { createBrowserClient as createSupabaseBrowserClient } from "@supabase/ssr";
+
+export function createBrowserClient() {
+  return createSupabaseBrowserClient(
+    process.env.NEXT_PUBLIC_SUPABASE_URL!,
+    process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY!,
+  );
+}
diff --git a/web/src/lib/supabase/middleware.ts b/web/src/lib/supabase/middleware.ts
index a781cf5..fe20f07 100644
--- a/web/src/lib/supabase/middleware.ts
+++ b/web/src/lib/supabase/middleware.ts
@@ -9,18 +9,19 @@ export async function updateSession(request: NextRequest) {
     cookies: {
       getAll: () => request.cookies.getAll(),
       setAll: (cookiesToSet) => {
         cookiesToSet.forEach(({ name, value }) => request.cookies.set(name, value));
         response = NextResponse.next({ request });
         cookiesToSet.forEach(({ name, value }) => response.cookies.set(name, value));
       },
     },
   });
   const { data } = await supabase.auth.getUser();
+  if (request.nextUrl.pathname.startsWith("/login")) return response;
   if (!data.user) {
     const url = request.nextUrl.clone();
     url.pathname = "/login";
     url.searchParams.set("next", request.nextUrl.pathname);
     return NextResponse.redirect(url);
   }
   return response;
 }
\ No newline at end of file
