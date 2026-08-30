"use client";
import { useRouter } from "next/navigation";
import { createBrowserClient } from "@/lib/supabase/client";
import { Button } from "@/components/ui/button";

export function LogoutButton() {
  const router = useRouter();
  async function logout() {
    const supabase = createBrowserClient();
    await supabase.auth.signOut();
    router.replace("/login");
    router.refresh();
  }
  return <Button variant="outline" onClick={logout}>Sign out</Button>;
}
