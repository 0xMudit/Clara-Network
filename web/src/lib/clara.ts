// src/lib/clara.ts
import "server-only";
import { getEnv } from "./env";

export function claraFetch<T>(path: string, accessToken: string, init?: RequestInit): Promise<T> {
  const { CLARA_API_URL } = getEnv();
  return fetch(`${CLARA_API_URL}/api/v1${path}`, {
    headers: { Authorization: `Bearer ${accessToken}`, Accept: "application/json" },
    ...init,
  }).then(r => {
    if (!r.ok) throw new Error(`clara ${path}: ${r.status}`);
    return r.json() as Promise<T>;
  });
}