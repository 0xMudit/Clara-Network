import { notFound } from "next/navigation";
import { createServerClient } from "@/lib/supabase/server";
import { roleFromAppMetadata } from "@/lib/roles";

export default async function MerchantPage() {
  const supabase = await createServerClient();
  const { data } = await supabase.auth.getUser();
  if (roleFromAppMetadata(data.user?.app_metadata) !== "merchant") notFound();
  return <h1 className="text-2xl font-semibold">Merchant</h1>;
}