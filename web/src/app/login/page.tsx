import { Suspense } from "react";
import { redirect } from "next/navigation";
import { createServerClient } from "@/lib/supabase/server";
import { LoginForm } from "./login-form";
import { roleFromAppMetadata, HOME_BY_ROLE } from "@/lib/roles";

export const dynamic = "force-dynamic";

export default async function LoginPage() {
  const supabase = await createServerClient();
  const { data } = await supabase.auth.getUser();
  if (data.user) {
    const role = roleFromAppMetadata(data.user.app_metadata);
    redirect(role ? HOME_BY_ROLE[role] : "/overview");
  }
  return (
    <main className="min-h-svh">
      <Suspense>
        <LoginForm />
      </Suspense>
    </main>
  );
}
