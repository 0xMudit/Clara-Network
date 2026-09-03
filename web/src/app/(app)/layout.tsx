import Nav from "@/components/nav";
import { RefreshProvider } from "@/components/providers/refresh-provider";
import { createServerClient } from "@/lib/supabase/server";
import { roleFromAppMetadata } from "@/lib/roles";
import { notFound } from "next/navigation";

export const dynamic = "force-dynamic";

export default async function AppLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const supabase = await createServerClient();
  const { data } = await supabase.auth.getUser();
  const role = roleFromAppMetadata(data.user?.app_metadata);
  if (!role) notFound();
  return (
    <div className="min-h-svh flex flex-col">
      <RefreshProvider />
      <Nav role={role} />
      <main className="mx-auto w-full max-w-6xl flex-1 px-4 py-8 sm:px-6 lg:px-8">
        {children}
      </main>

      {/* Footer */}
      <footer className="border-t bg-muted/30 py-4">
        <div className="mx-auto max-w-6xl px-4 text-center text-xs text-muted-foreground/50">
          Clara Network · Open-source payment infrastructure
        </div>
      </footer>
    </div>
  );
}
