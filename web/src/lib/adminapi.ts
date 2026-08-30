// src/lib/adminapi.ts
import "server-only";

export async function fetchAdmin<T>(path: string): Promise<T> {
  const url = `/api/data${path}`;
  const res = await fetch(url, { cache: "no-store" });
  if (!res.ok) {
    const body = await res.json().catch(() => null);
    const msg = body && (body as { error?: string }).error ? (body as { error: string }).error : String(res.status);
    throw new Error(`adminapi ${path}: ${msg}`);
  }
  return res.json() as Promise<T>;
}