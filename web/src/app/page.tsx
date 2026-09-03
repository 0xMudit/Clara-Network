import { redirect } from "next/navigation";
import { createServerClient } from "@/lib/supabase/server";
import { roleFromAppMetadata, HOME_BY_ROLE } from "@/lib/roles";
import { LandingHeader } from "@/components/landing/header";
import { Hero } from "@/components/landing/hero";

export const dynamic = "force-dynamic";

export default async function LandingPage() {
  const supabase = await createServerClient();
  const { data } = await supabase.auth.getUser();
  // Already signed in → straight to their dashboard.
  if (data.user) {
    const role = roleFromAppMetadata(data.user.app_metadata);
    redirect(role ? HOME_BY_ROLE[role] : "/overview");
  }

  return (
    <div className="relative min-h-svh overflow-hidden">
      {/* Ambient glows */}
      <div aria-hidden className="pointer-events-none absolute inset-0 -z-10">
        <div className="absolute -top-32 left-1/2 h-[480px] w-[720px] -translate-x-1/2 rounded-full bg-violet-600/20 blur-[140px]" />
        <div className="absolute bottom-0 left-0 h-72 w-72 rounded-full bg-sky-600/10 blur-[120px]" />
        <div className="absolute right-0 top-1/3 h-72 w-72 rounded-full bg-emerald-600/10 blur-[120px]" />
      </div>

      <LandingHeader />
      <main className="pb-24">
        <Hero />
      </main>

      <footer className="relative border-t bg-card/40 py-8">
        <div className="mx-auto flex max-w-6xl flex-col items-center justify-between gap-3 px-4 text-center text-xs text-muted-foreground/60 sm:flex-row sm:px-6 lg:px-8">
          <span>Clara Network · Open-source payment infrastructure</span>
          <span className="flex items-center gap-1.5">
            Built for county &amp; bank engineers
          </span>
        </div>
      </footer>
    </div>
  );
}