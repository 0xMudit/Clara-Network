import { notFound } from "next/navigation";
import { createServerClient } from "@/lib/supabase/server";
import { roleFromAppMetadata } from "@/lib/roles";

export default async function OpsPage() {
  const supabase = await createServerClient();
  const { data } = await supabase.auth.getUser();
  if (roleFromAppMetadata(data.user?.app_metadata) !== "scheme_operator") notFound();
  return <h1 className="text-2xl font-semibold">Operations</h1>;
}