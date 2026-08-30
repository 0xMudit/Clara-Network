"use client";
import { useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { createBrowserClient } from "@/lib/supabase/client";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { roleFromAppMetadata, HOME_BY_ROLE } from "@/lib/roles";

export function LoginForm() {
  const router = useRouter();
  const params = useSearchParams();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError(null);
    const supabase = createBrowserClient();
    const { data, error } = await supabase.auth.signInWithPassword({ email, password });
    if (error) { setError(error.message); setLoading(false); return; }
    const role = roleFromAppMetadata(data.session?.user.app_metadata);
    const next = params.get("next");
    router.push(next && next.startsWith("/") ? next : (role ? HOME_BY_ROLE[role] : "/overview"));
    router.refresh();
  }

  return (
    <Card className="w-full max-w-sm">
      <CardHeader>
        <CardTitle>Clara Network</CardTitle>
        <CardDescription>Scheme operator console — sign in to continue</CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={onSubmit} className="grid gap-4">
          <div className="grid gap-2"><Label htmlFor="email">Work email</Label>
            <Input id="email" type="email" value={email} onChange={e => setEmail(e.target.value)} required autoComplete="email" /></div>
          <div className="grid gap-2"><Label htmlFor="password">Password</Label>
            <Input id="password" type="password" value={password} onChange={e => setPassword(e.target.value)} required autoComplete="current-password" /></div>
          {error && <p className="text-sm text-destructive">{error}</p>}
          <Button type="submit" disabled={loading}>{loading ? "Signing in…" : "Sign in"}</Button>
        </form>
      </CardContent>
    </Card>
  );
}
