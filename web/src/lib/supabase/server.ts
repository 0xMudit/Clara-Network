import { createServerClient as createSupabaseServerClient } from "@supabase/ssr";
import { cookies } from "next/headers";

export async function createServerClient() {
  // eslint-disable-next-line @typescript-eslint/no-require-imports
  const { SUPABASE_URL, SUPABASE_ANON_KEY } = require("../env").getEnv();
  const cookieStore = await cookies();
  return createSupabaseServerClient(SUPABASE_URL, SUPABASE_ANON_KEY, {
    cookies: { getAll: () => cookieStore.getAll(), setAll: () => {} },
  });
}