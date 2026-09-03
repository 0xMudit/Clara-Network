import { redirect } from "next/navigation";
import { createServerClient } from "@/lib/supabase/server";
import { roleFromAppMetadata, HOME_BY_ROLE } from "@/lib/roles";

export const dynamic = "force-dynamic";

export default async function LoginPage({
  searchParams,
}: {
  searchParams: Promise<{ next?: string }>;
}) {
  const supabase = await createServerClient();
  const { data } = await supabase.auth.getUser();
  // Logged in → straight to their dashboard. Not logged in → the landing page
  // at "/" is the single sign-in surface (hero + role picker). Preserve any
  // ?next= deep link so post-login we land where the visitor was headed.
  if (data.user) {
    const role = roleFromAppMetadata(data.user.app_metadata);
    redirect(role ? HOME_BY_ROLE[role] : "/overview");
  }
  const { next } = await searchParams;
  if (next && next.startsWith("/") && next !== "/") {
    redirect(`/?next=${encodeURIComponent(next)}`);
  }
  redirect("/");
}
