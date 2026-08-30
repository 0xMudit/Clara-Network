// src/app/api/data/[...path]/route.ts
import { NextResponse } from "next/server";
import { createServerClient } from "@/lib/supabase/server";
import { getEnv } from "@/lib/env";
import { claraFetch } from "@/lib/clara";

export const dynamic = "force-dynamic";

type Res<T> = NextResponse<T | { error: string }>;

const ALLOWED_ROUTES: Record<string, string[]> = {
  scheme_operator: [
    "/dashboard", "/transactions", "/clearing/cycles", "/clearing/records",
    "/clearing/positions", "/settlement/instructions", "/settlement/prefunds",
    "/settlement/default-fund", "/ledger/accounts", "/ledger/entries",
  ],
  issuer: ["/dashboard", "/cards", "/tokens"],
  acquirer: ["/dashboard", "/acquirer", "/merchants", "/funding", "/disputes"],
  merchant: ["/dashboard", "/merchant", "/funding", "/disputes"],
};

export async function GET(
  req: Request,
  { params }: { params: Promise<{ path: string[] }> },
): Promise<Res<unknown>> {
  const supabase = await createServerClient();
  const { data, error } = await supabase.auth.getUser();
  if (error || !data.user) return NextResponse.json({ error: "unauthorized" }, { status: 401 });
  const role = data.user.app_metadata.role as string;
  const { path } = await params;
  const target = "/" + path.join("/");
  const allowed = ALLOWED_ROUTES[role] ?? [];
  if (!allowed.includes(target)) return NextResponse.json({ error: "forbidden" }, { status: 403 });
  try {
    const body = await claraFetch(target, data.user.id, { cache: "no-store" });
    return NextResponse.json(body);
  } catch (e) {
    const status = e instanceof Error && /^clara \S+ 4\d\d/.test(e.message) ? 502 : 500;
    return NextResponse.json({ error: "upstream error", detail: e instanceof Error ? e.message : String(e) }, { status });
  }
}