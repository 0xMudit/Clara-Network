// src/lib/adminapi.ts
import "server-only";
import { cookies } from "next/headers";
import { getAppUrl } from "./env";
import { getMockForPath } from "./mock-data";

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
    const status = e instanceof AdminError ? e.status : undefined;

    // A real 401/403 means the session expired or the role lacks access —
    // never mask that with mock data, or we'd hide a genuine security issue.
    if (status === 401 || status === 403) {
      return { ok: false, error: "You don't have access to this data." };
    }

    // Network errors and upstream (5xx) failures: fall back to realistic
    // mock data so the demo stays alive when the Go admin API is unreachable.
    // The mock layer only knows a handful of endpoints; anything else falls
    // through to the generic error.
    const mock = getMockForPath(path);
    if (mock !== null) {
      console.warn(
        `[adminapi] using mock data for "${path}" (${status ?? "network error"})`
      );
      return { ok: true, data: mock as T };
    }

    const msg =
      e instanceof AdminError
        ? "The Admin API is temporarily unavailable. Try again shortly."
        : "Something went wrong loading this data.";
    return { ok: false, error: msg };
  }
}