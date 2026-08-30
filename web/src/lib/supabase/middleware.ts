import { createServerClient } from "@supabase/ssr";
import { NextResponse, type NextRequest } from "next/server";

export async function updateSession(request: NextRequest) {
  // eslint-disable-next-line @typescript-eslint/no-require-imports
  const { SUPABASE_URL, SUPABASE_ANON_KEY } = require("../env").getEnv();
  let response = NextResponse.next({ request });
  const supabase = createServerClient(SUPABASE_URL, SUPABASE_ANON_KEY, {
    cookies: {
      getAll: () => request.cookies.getAll(),
      setAll: (cookiesToSet) => {
        cookiesToSet.forEach(({ name, value }) => request.cookies.set(name, value));
        response = NextResponse.next({ request });
        cookiesToSet.forEach(({ name, value }) => response.cookies.set(name, value));
      },
    },
  });
  const { data } = await supabase.auth.getUser();
  if (request.nextUrl.pathname.startsWith("/login")) return response;
  if (!data.user) {
    const url = request.nextUrl.clone();
    url.pathname = "/login";
    url.searchParams.set("next", request.nextUrl.pathname);
    return NextResponse.redirect(url);
  }
  return response;
}