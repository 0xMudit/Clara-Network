import { tryFetchAdmin } from "@/lib/adminapi";
import { fmtCount } from "@/lib/format";
import { sparkFromValue } from "@/lib/spark-data";
import { StatCard } from "@/components/cards/stat-card";
import { TransactionChart, type SeriesPointDatum } from "@/components/transaction-chart";
import type { DashboardSummary, SeriesPoint } from "@/types/admin";
import { PageHeader, PageStack, SectionLabel } from "@/components/page-shell";
import { DataError } from "@/components/states";
import {
  CreditCard,
  FileText,
  Store,
  AlertTriangle,
  Layers,
  Coins,
} from "lucide-react";

export default async function OverviewPage() {
  const dashboard = await tryFetchAdmin<DashboardSummary>("/dashboard");
  const series = await tryFetchAdmin<{ items: SeriesPoint[] }>("/dashboard/series?days=14");
  const seriesData: SeriesPointDatum[] = series.ok
    ? series.data.items.map((p) => ({ date: p.date, count: p.count }))
    : [];

  return (
    <PageStack>
      <PageHeader
        title="Network overview"
        description="Real-time snapshot of your Clara payment network — transactions, clearing, merchants, and more."
      />

      {!dashboard.ok ? (
        <DataError message={dashboard.error} />
      ) : (
        <div>
          <SectionLabel>Key metrics</SectionLabel>
          <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
            <StatCard
              title="Transactions"
              value={fmtCount(dashboard.data.transactions)}
              hint="Authorizations processed via the switch"
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
              hint="Open & in-progress disputes"
              icon={<AlertTriangle className="size-5" />}
              accent={
                dashboard.data.disputes > 0 ? "warning" : "default"
              }
              sparkData={sparkFromValue(dashboard.data.disputes, 10, 0.25)}
            />
            <StatCard
              title="Cards"
              value={fmtCount(dashboard.data.cards)}
              hint="Issued cards in the network"
              icon={<Layers className="size-5" />}
              accent="default"
              sparkData={sparkFromValue(dashboard.data.cards)}
            />
            <StatCard
              title="Tokens"
              value={fmtCount(dashboard.data.tokens)}
              hint="Digital tokens for contactless & mobile"
              icon={<Coins className="size-5" />}
              accent="default"
              sparkData={sparkFromValue(dashboard.data.tokens)}
            />
          </div>

          <div className="mt-8">
            <SectionLabel>Transactions over time</SectionLabel>
            {!series.ok ? (
              <DataError message={series.error} />
            ) : seriesData.length === 0 ? (
              <p className="rounded-xl border border-dashed bg-card/40 px-4 py-8 text-center text-sm text-muted-foreground">
                No transaction history yet — run the demos to start moving volume through the switch.
              </p>
            ) : (
              <TransactionChart data={seriesData} />
            )}
          </div>
        </div>
      )}
    </PageStack>
  );
}
