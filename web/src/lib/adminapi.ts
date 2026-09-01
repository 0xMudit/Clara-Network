// src/lib/adminapi.ts
import "server-only";
import { cookies } from "next/headers";
import { getAppUrl } from "./env";

export async function fetchAdmin<T>(path: string): Promise<T> {
  const cookieStore = await cookies();
  const url = `${getAppUrl()}/api/data${path}`;
  const res = await fetch(url, {
    cache: "no-store",
    headers: { Cookie: cookieStore.toString() },
  });
  if (!res.ok) {
    const body = await res.json().catch(() => null);
    const msg = body && (body as { error?: string }).error ? (body as { error: string }).error : String(res.status);
    throw new Error(`adminapi ${path}: ${msg}`);
  }
  return res.json() as Promise<T>;
}