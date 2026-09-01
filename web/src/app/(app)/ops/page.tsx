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
import { CreditCard, FileText, Store, ArrowRight } from "lucide-react";

export default async function OpsPage() {
  const supabase = await createServerClient();
  const { data } = await supabase.auth.getUser();
  if (roleFromAppMetadata(data.user?.app_metadata) !== "scheme_operator")
    notFound();
  const dashboard = await tryFetchAdmin<DashboardSummary>("/dashboard");

  return (
    <PageStack>
      <PageHeader
        title="Operations"
        tone="violet"
        description="Full network visibility — transactions, clearing, and merchant activity at a glance."
      />

      {!dashboard.ok ? (
        <DataError message={dashboard.error} />
      ) : (
        <div>
          <SectionLabel>Today&apos;s activity</SectionLabel>
          <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
            <StatCard
              title="Transactions today"
              value={fmtCount(dashboard.data.transactions)}
              hint="Authorizations via the switch"
              icon={<CreditCard className="size-5" />}
              accent="info"
              sparkData={sparkFromValue(dashboard.data.transactions)}
            />
            <StatCard
              title="Clearing records"
              value={fmtCount(dashboard.data.clearingRecords)}
              hint="Captured this settlement window"
              icon={<FileText className="size-5" />}
              accent="default"
              sparkData={sparkFromValue(dashboard.data.clearingRecords)}
            />
            <StatCard
              title="Merchants onboarded"
              value={fmtCount(dashboard.data.merchants)}
              icon={<Store className="size-5" />}
              accent="success"
              sparkData={sparkFromValue(dashboard.data.merchants)}
            />
          </div>
        </div>
      )}

      <Link
        href="/transactions"
        className="group inline-flex items-center gap-2 rounded-xl border bg-card px-4 py-3 text-sm font-medium text-muted-foreground transition-all hover:border-border hover:bg-muted/50 hover:text-foreground hover:shadow-sm"
      >
        View transaction log
        <ArrowRight className="size-4 transition-transform group-hover:translate-x-0.5" />
      </Link>
    </PageStack>
  );
}
