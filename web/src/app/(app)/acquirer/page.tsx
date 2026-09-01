import { notFound } from "next/navigation";
import Link from "next/link";
import { createServerClient } from "@/lib/supabase/server";
import { roleFromAppMetadata } from "@/lib/roles";
import { tryFetchAdmin } from "@/lib/adminapi";
import { fmtCount } from "@/lib/format";
import { sparkFromValue } from "@/lib/spark-data";
import { StatCard } from "@/components/cards/stat-card";
import type { DashboardSummary } from "@/types/admin";
import { PageHeader, PageStack, SectionLabel } from "@/components/page-shell";
import { DataError } from "@/components/states";
import { Store, AlertTriangle, ArrowRight } from "lucide-react";

export default async function AcquirerPage() {
  const supabase = await createServerClient();
  const { data } = await supabase.auth.getUser();
  if (roleFromAppMetadata(data.user?.app_metadata) !== "acquirer") notFound();
  const dashboard = await tryFetchAdmin<DashboardSummary>("/dashboard");

  return (
    <PageStack>
      <PageHeader
        title="Acquirer"
        tone="emerald"
        description="Merchant management, funding, and dispute resolution."
      />

      {!dashboard.ok ? (
        <DataError message={dashboard.error} />
      ) : (
        <div>
          <SectionLabel>Portfolio</SectionLabel>
          <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
            <StatCard
              title="Merchants"
              value={fmtCount(dashboard.data.merchants)}
              hint="Active merchant accounts"
              icon={<Store className="size-5" />}
              accent="success"
              sparkData={sparkFromValue(dashboard.data.merchants)}
            />
            <StatCard
              title="Disputes"
              value={fmtCount(dashboard.data.disputes)}
              hint="Open & in-progress"
              icon={<AlertTriangle className="size-5" />}
              accent={dashboard.data.disputes > 0 ? "warning" : "default"}
              sparkData={sparkFromValue(dashboard.data.disputes, 10, 0.25)}
            />
          </div>
        </div>
      )}

      <div className="flex flex-wrap gap-3">
        <Link
          href="/merchants"
          className="group inline-flex items-center gap-2 rounded-xl border bg-card px-4 py-3 text-sm font-medium text-muted-foreground transition-all hover:border-border hover:bg-muted/50 hover:text-foreground hover:shadow-sm"
        >
          View merchants
          <ArrowRight className="size-4 transition-transform group-hover:translate-x-0.5" />
        </Link>
        <Link
          href="/disputes"
          className="group inline-flex items-center gap-2 rounded-xl border bg-card px-4 py-3 text-sm font-medium text-muted-foreground transition-all hover:border-border hover:bg-muted/50 hover:text-foreground hover:shadow-sm"
        >
          View disputes
          <ArrowRight className="size-4 transition-transform group-hover:translate-x-0.5" />
        </Link>
      </div>
    </PageStack>
  );
}