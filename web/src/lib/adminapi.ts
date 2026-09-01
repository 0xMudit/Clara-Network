// src/lib/adminapi.ts
import "server-only";
import { cookies } from "next/headers";
import { getAppUrl } from "./env";

export class AdminError extends Error {
  constructor(
    message: string,
    readonly status?: number,
  ) {
    super(message);
    this.name = "AdminError";
  }
}

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
    throw new AdminError(`adminapi ${path}: ${msg}`, res.status);
  }
  return res.json() as Promise<T>;
}

export type AdminResult<T> =
  | { ok: true; data: T }
  | { ok: false; error: string };

/** Non-throwing fetch: returns data or a user-facing error message. */
export async function tryFetchAdmin<T>(path: string): Promise<AdminResult<T>> {
  try {
    const data = await fetchAdmin<T>(path);
    return { ok: true, data };
  } catch (e) {
    const msg =
      e instanceof AdminError
        ? e.status === 401 || e.status === 403
          ? "You don't have access to this data."
          : "The Admin API is temporarily unavailable. Try again shortly."
        : "Something went wrong loading this data.";
    return { ok: false, error: msg };
  }
}