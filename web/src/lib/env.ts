const required = ["SUPABASE_URL", "SUPABASE_ANON_KEY", "CLARA_API_URL"] as const;
export function getEnv(): Record<typeof required[number], string> {
  const out = {} as Record<typeof required[number], string>;
  for (const k of required) {
    const v = process.env[k];
    if (!v) throw new Error(`missing env ${k}`);
    out[k] = v;
  }
  return out;
}

// Absolute origin of this app, used for internal server-side fetch calls
// (Node's fetch rejects relative URLs in the serverless runtime).
export function getAppUrl(): string {
  if (process.env.NEXT_PUBLIC_APP_URL) {
    return process.env.NEXT_PUBLIC_APP_URL.replace(/\/+$/, "");
  }
  if (process.env.VERCEL_PROJECT_PRODUCTION_URL) {
    return `https://${process.env.VERCEL_PROJECT_PRODUCTION_URL}`;
  }
  if (process.env.VERCEL_URL) {
    return `https://${process.env.VERCEL_URL}`;
  }
  return "http://localhost:3000";
}
