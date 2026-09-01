"use client";
import { useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { createBrowserClient } from "@/lib/supabase/client";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu, DropdownMenuContent, DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { ChevronDown, Landmark, CreditCard, Store, ShoppingBag, Eye, Sparkles } from "lucide-react";
import { roleFromAppMetadata, HOME_BY_ROLE, type Role } from "@/lib/roles";

// Demo credentials are provisioned by web/scripts/create-users.mjs (npm run db:users).
const DEMO_PASSWORD = "ClaraDemo!2026";

const DEMO_ROLES: { role: Role; label: string; email: string; blurb: string; Icon: typeof Landmark }[] = [
  { role: "scheme_operator", label: "Scheme Operator", email: "scheme-operator@clara.demo", blurb: "Full network view — clearing, settlement, ledger, disputes", Icon: Landmark },
  { role: "issuer", label: "Issuer", email: "issuer@clara.demo", blurb: "Cards, tokens and BIN ranges", Icon: CreditCard },
  { role: "acquirer", label: "Acquirer", email: "acquirer@clara.demo", blurb: "Merchants, funding and disputes", Icon: Store },
  { role: "merchant", label: "Merchant", email: "merchant@clara.demo", blurb: "Funding lines and disputes", Icon: ShoppingBag },
  { role: "viewer", label: "Viewer", email: "viewer@clara.demo", blurb: "Read-only overview, just the highlights", Icon: Eye },
];

export function LoginForm() {
  const router = useRouter();
  const params = useSearchParams();
  const [selected, setSelected] = useState<typeof DEMO_ROLES[number] | null>(null);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  function pickRole(entry: typeof DEMO_ROLES[number]) {
    setSelected(entry);
    setEmail(entry.email);
    setPassword(DEMO_PASSWORD);
    setError(null);
  }

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

  const SelectedIcon = selected?.Icon ?? Sparkles;

  return (
    <Card className="w-full max-w-md">
      <CardHeader className="items-center text-center">
        <CardTitle className="flex items-center gap-2">
          <Landmark className="size-5" /> Clara Network
        </CardTitle>
        <CardDescription>Pick a persona to auto-fill a demo account — then one click to sign in.</CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={onSubmit} className="grid gap-4">
          <div className="grid gap-2">
            <Label htmlFor="role">I am a…</Label>
            <DropdownMenu>
              <DropdownMenuTrigger
                render={<Button type="button" variant="outline" className="justify-between text-left font-normal" />}
              >
                <span className="flex items-center gap-2 truncate">
                  <SelectedIcon className="size-4 text-muted-foreground" />
                  {selected ? selected.label : "Choose your role"}
                </span>
                <ChevronDown className="size-4 opacity-60" />
              </DropdownMenuTrigger>
              <DropdownMenuContent className="w-72">
                {DEMO_ROLES.map((entry) => (
                  <DropdownMenuItem key={entry.role} onSelect={() => pickRole(entry)}>
                    <entry.Icon className="size-4" />
                    <span className="flex flex-col">
                      <span className="text-sm font-medium">{entry.label}</span>
                      <span className="text-xs text-muted-foreground">{entry.blurb}</span>
                    </span>
                  </DropdownMenuItem>
                ))}
              </DropdownMenuContent>
            </DropdownMenu>
          </div>

          {selected && (
            <p className="text-xs text-muted-foreground">
              Signing in as <span className="font-medium text-foreground">{selected.email}</span> — password
              auto-filled for the demo.
            </p>
          )}

          <div className="grid gap-2">
            <Label htmlFor="email">Work email</Label>
            <Input id="email" type="email" value={email} onChange={e => setEmail(e.target.value)} required autoComplete="email"
              placeholder="you@company.com" />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="password">Password</Label>
            <Input id="password" type="password" value={password} onChange={e => setPassword(e.target.value)} required autoComplete="current-password" />
          </div>

          {error && <p className="text-sm text-destructive">{error}</p>}
          <Button type="submit" disabled={loading || !email || !password}>
            {loading ? "Signing in…" : selected ? `Sign in as ${selected.label}` : "Sign in"}
          </Button>
        </form>
      </CardContent>
    </Card>
  );
}
