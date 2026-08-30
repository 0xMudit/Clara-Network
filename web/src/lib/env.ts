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
