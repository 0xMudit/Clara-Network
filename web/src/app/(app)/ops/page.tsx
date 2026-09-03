import { notFound } from "next/navigation";
import Link from "next/link";
import { createServerClient } from "@/lib/supabase/server";
import { roleFromAppMetadata } from "@/lib/roles";
import { tryFetchAdmin } from "@/lib/adminapi";
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
import { MonoChip } from "@/components/ui/badge";
import { DataError } from "@/components/states";
import { CreditCard, FileText, Store, ArrowRight } from "lucide-react";
import type { AuditEvent } from "../transactions/page";

const txnColumns: Column<AuditEvent>[] = [
  {
    key: "stan",
    header: "STAN",
    render: (t) => <MonoChip>{t.stan}</MonoChip>,
  },
  {
    key: "mti",
    header: "MTI",
    render: (t) => <MonoChip>{t.mti}</MonoChip>,
  },
  {
    key: "pan",
    header: "PAN",
    render: (t) => <span className="font-mono text-xs">{t.pan}</span>,
  },
  {
    key: "amount",
    header: "Amount",
    render: (t) => <span className="font-medium">{t.amount}</span>,
  },
  {
    key: "responseCode",
    header: "Resp",
    render: (t) => <MonoChip>{t.responseCode}</MonoChip>,
  },
];

export default async function OpsPage() {
  const supabase = await createServerClient();
  const { data } = await supabase.auth.getUser();
  if (roleFromAppMetadata(data.user?.app_metadata) !== "scheme_operator")
    notFound();

  const dashboard = await tryFetchAdmin<DashboardSummary>("/dashboard");
  const txns = await tryFetchAdmin<Page<AuditEvent>>("/transactions?limit=6");
  const series = await tryFetchAdmin<{ items: SeriesPoint[] }>("/dashboard/series?days=14");
  const seriesData: SeriesPointDatum[] = series.ok
    ? series.data.items.map((p) => ({ date: p.date, count: p.count }))
    : [];

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
          <SectionLabel>Network activity</SectionLabel>
          <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
            <StatCard
              title="Transactions today"
              value={fmtCount(dashboard.data.transactions)}
              hint="Authorizations via the switch"
              icon={<CreditCard className="size-5" />}
              accent="info"
            />
            <StatCard
              title="Clearing records"
              value={fmtCount(dashboard.data.clearingRecords)}
              hint="Captured this settlement window"
              icon={<FileText className="size-5" />}
              accent="default"
            />
            <StatCard
              title="Merchants onboarded"
              value={fmtCount(dashboard.data.merchants)}
              hint="Active merchant accounts"
              icon={<Store className="size-5" />}
              accent="success"
            />
          </div>

          <div className="mt-8 grid gap-6 lg:grid-cols-2">
            <div>
              <SectionLabel>Throughput</SectionLabel>
              {!series.ok ? (
                <DataError message={series.error} />
              ) : seriesData.length === 0 ? (
                <p className="rounded-xl border border-dashed bg-card/40 px-4 py-8 text-center text-sm text-muted-foreground">
                  No throughput history yet — run the demos to start moving volume through the switch.
                </p>
              ) : (
                <div className="rounded-2xl border bg-card p-5">
                  <TransactionChart data={seriesData} />
                </div>
              )}
            </div>

            <div>
              <SectionLabel>Latest authorizations</SectionLabel>
              {!txns.ok ? (
                <DataError message={txns.error} />
              ) : (
                <DataTable
                  columns={txnColumns}
                  rows={txns.data.items}
                  getKey={(t) => t.stan}
                  emptyTitle="No transactions yet"
                  emptyHint="The switch records every authorization here once traffic starts flowing."
                />
              )}
              <Link
                href="/transactions"
                className="group mt-4 inline-flex items-center gap-2 rounded-xl border bg-card px-4 py-3 text-sm font-medium text-muted-foreground transition-all hover:border-border hover:bg-muted/50 hover:text-foreground hover:shadow-sm"
              >
                View full transaction log
                <ArrowRight className="size-4 transition-transform group-hover:translate-x-0.5" />
              </Link>
            </div>
          </div>
        </div>
      )}
    </PageStack>
  );
}
