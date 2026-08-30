import { notFound } from "next/navigation";
import { createServerClient } from "@/lib/supabase/server";
import { roleFromAppMetadata } from "@/lib/roles";

export default async function IssuerPage() {
  const supabase = await createServerClient();
  const { data } = await supabase.auth.getUser();
  if (roleFromAppMetadata(data.user?.app_metadata) !== "issuer") notFound();
  return <h1 className="text-2xl font-semibold">Issuer</h1>;
}