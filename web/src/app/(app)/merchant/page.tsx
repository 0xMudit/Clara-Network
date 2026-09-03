import { notFound } from "next/navigation";
import Link from "next/link";
import { createServerClient } from "@/lib/supabase/server";
import { roleFromAppMetadata } from "@/lib/roles";
import { tryFetchAdmin } from "@/lib/adminapi";
import { fmtMinor } from "@/lib/money/minor";
import { fmtCount } from "@/lib/format";
import { StatCard } from "@/components/cards/stat-card";
import {
  TransactionChart,
  type SeriesPointDatum,
} from "@/components/transaction-chart";
import type { DashboardSummary, Page, SeriesPoint } from "@/types/admin";
import {
  PageHeader,
  PageStack,
  SectionLabel,
} from "@/components/page-shell";
import { DataTable, type Column } from "@/components/data-table";
import { Badge } from "@/components/ui/badge";
import { DataError } from "@/components/states";
import { Store, AlertTriangle, Landmark, ArrowRight } from "lucide-react";
import type { Merchant } from "../merchants/page";

const fundColumns: Column<Merchant>[] = [
  {
    key: "name",
    header: "Merchant",
    render: (m) => <span className="font-medium">{m.name}</span>,
  },
  {
    key: "delay",
    header: "Funding delay",
    render: (m) => <Badge tone="info">{m.fundingDelayDays}d</Badge>,
  },
  {
    key: "limit",
    header: "Txn limit",
    render: (m) => <span className="font-medium">{fmtMinor(m.transactionLimit)}</span>,
  },
  {
    key: "reserve",
    header: "Reserve",
    render: (m) => <span className="font-medium">{fmtMinor(m.reserveBalance)}</span>,
  },
];

export default async function MerchantPage() {
  const supabase = await createServerClient();
  const { data } = await supabase.auth.getUser();
  if (roleFromAppMetadata(data.user?.app_metadata) !== "merchant") notFound();

  const dashboard = await tryFetchAdmin<DashboardSummary>("/dashboard");
  const funding = await tryFetchAdmin<Page<Merchant>>("/merchants?limit=6");
  const series = await tryFetchAdmin<{ items: SeriesPoint[] }>("/dashboard/series?days=14");
  const seriesData: SeriesPointDatum[] = series.ok
    ? series.data.items.map((p) => ({ date: p.date, count: p.count }))
    : [];

  return (
    <PageStack>
      <PageHeader
        title="Merchant"
        tone="amber"
        description="Your funding schedule, reserves, transaction limits, and dispute status."
      />

      {!dashboard.ok ? (
        <DataError message={dashboard.error} />
      ) : (
        <div>
          <SectionLabel>Summary</SectionLabel>
          <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
            <StatCard
              title="Accounts"
              value={fmtCount(dashboard.data.merchants)}
              hint="Merchant accounts in your portfolio"
              icon={<Store className="size-5" />}
              accent="success"
            />
            <StatCard
              title="Open disputes"
              value={fmtCount(dashboard.data.disputes)}
              hint="Cases open against your accounts"
              icon={<AlertTriangle className="size-5" />}
              accent={dashboard.data.disputes > 0 ? "danger" : "default"}
            />
            <StatCard
              title="Network transactions"
              value={fmtCount(dashboard.data.transactions)}
              hint="Authorizations routed via the switch"
              icon={<Landmark className="size-5" />}
              accent="info"
            />
          </div>

          <div className="mt-8 grid gap-6 lg:grid-cols-2">
            <div>
              <SectionLabel>Transaction volume</SectionLabel>
              {!series.ok ? (
                <DataError message={series.error} />
              ) : seriesData.length === 0 ? (
                <p className="rounded-xl border border-dashed bg-card/40 px-4 py-8 text-center text-sm text-muted-foreground">
                  No transaction history yet — run the merchant demo to start moving volume.
                </p>
              ) : (
                <div className="rounded-2xl border bg-card p-5">
                  <TransactionChart data={seriesData} />
                </div>
              )}
            </div>

            <div>
              <SectionLabel>Funding profiles</SectionLabel>
              {!funding.ok ? (
                <DataError message={funding.error} />
              ) : (
                <DataTable
                  columns={fundColumns}
                  rows={funding.data.items}
                  getKey={(m) => m.id}
                  emptyTitle="No funding profiles yet"
                  emptyHint="Funding schedules will appear here."
                />
              )}
              <div className="mt-4 flex flex-wrap gap-3">
                <Link
                  href="/funding"
                  className="group inline-flex items-center gap-2 rounded-xl border bg-card px-4 py-3 text-sm font-medium text-muted-foreground transition-all hover:border-border hover:bg-muted/50 hover:text-foreground hover:shadow-sm"
                >
                  View funding
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
            </div>
          </div>
        </div>
      )}
    </PageStack>
  );
}
