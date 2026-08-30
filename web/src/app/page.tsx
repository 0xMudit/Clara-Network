import { redirect } from "next/navigation";
import { createServerClient } from "@/lib/supabase/server";
import { roleFromAppMetadata, HOME_BY_ROLE } from "@/lib/roles";

export const dynamic = "force-dynamic";

export default async function Home() {
  const supabase = await createServerClient();
  const { data } = await supabase.auth.getUser();
  if (!data.user) redirect("/login");
  const role = roleFromAppMetadata(data.user.app_metadata);
  redirect(role ? HOME_BY_ROLE[role] : "/overview");
}