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
import { Store, AlertTriangle, Wallet, ArrowRight } from "lucide-react";
import type { Merchant } from "../merchants/page";

const riskTone = (r: string) =>
  r === "low" ? "success" : r === "high" ? "danger" : "warning";

const merchantColumns: Column<Merchant>[] = [
  {
    key: "name",
    header: "Merchant",
    render: (m) => <span className="font-medium">{m.name}</span>,
  },
  {
    key: "riskTier",
    header: "Risk",
    render: (m) => <Badge tone={riskTone(m.riskTier)}>{m.riskTier}</Badge>,
  },
  {
    key: "reserve",
    header: "Reserve",
    render: (m) => <span className="font-medium">{fmtMinor(m.reserveBalance)}</span>,
  },
  {
    key: "volume",
    header: "Volume",
    render: (m) => <span className="font-medium">{fmtMinor(m.volume)}</span>,
  },
];

export default async function AcquirerPage() {
  const supabase = await createServerClient();
  const { data } = await supabase.auth.getUser();
  if (roleFromAppMetadata(data.user?.app_metadata) !== "acquirer") notFound();

  const dashboard = await tryFetchAdmin<DashboardSummary>("/dashboard");
  const merchants = await tryFetchAdmin<Page<Merchant>>("/merchants?limit=6");
  const series = await tryFetchAdmin<{ items: SeriesPoint[] }>("/dashboard/series?days=14");
  const seriesData: SeriesPointDatum[] = series.ok
    ? series.data.items.map((p) => ({ date: p.date, count: p.count }))
    : [];

  return (
    <PageStack>
      <PageHeader
        title="Acquirer"
        tone="emerald"
        description="Your merchant portfolio, reserves, funding lines, and dispute exposure."
      />

      {!dashboard.ok ? (
        <DataError message={dashboard.error} />
      ) : (
        <div>
          <SectionLabel>Portfolio</SectionLabel>
          <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
            <StatCard
              title="Merchants onboarded"
              value={fmtCount(dashboard.data.merchants)}
              hint="Active merchant accounts"
              icon={<Store className="size-5" />}
              accent="success"
            />
            <StatCard
              title="Open disputes"
              value={fmtCount(dashboard.data.disputes)}
              hint="Chargebacks & representments in flight"
              icon={<AlertTriangle className="size-5" />}
              accent={dashboard.data.disputes > 0 ? "warning" : "default"}
            />
            <StatCard
              title="Network transactions"
              value={fmtCount(dashboard.data.transactions)}
              hint="Authorizations routed via the switch"
              icon={<Wallet className="size-5" />}
              accent="info"
            />
          </div>

          <div className="mt-8 grid gap-6 lg:grid-cols-2">
            <div>
              <SectionLabel>Acquiring volume</SectionLabel>
              {!series.ok ? (
                <DataError message={series.error} />
              ) : seriesData.length === 0 ? (
                <p className="rounded-xl border border-dashed bg-card/40 px-4 py-8 text-center text-sm text-muted-foreground">
                  No acquiring history yet — run the acquirer demo to start moving volume.
                </p>
              ) : (
                <div className="rounded-2xl border bg-card p-5">
                  <TransactionChart data={seriesData} />
                </div>
              )}
            </div>

            <div>
              <SectionLabel>Top merchants</SectionLabel>
              {!merchants.ok ? (
                <DataError message={merchants.error} />
              ) : (
                <DataTable
                  columns={merchantColumns}
                  rows={merchants.data.items}
                  getKey={(m) => m.id}
                  emptyTitle="No merchants yet"
                  emptyHint="Boarded merchants will appear here."
                />
              )}
              <div className="mt-4 flex flex-wrap gap-3">
                <Link
                  href="/merchants"
                  className="group inline-flex items-center gap-2 rounded-xl border bg-card px-4 py-3 text-sm font-medium text-muted-foreground transition-all hover:border-border hover:bg-muted/50 hover:text-foreground hover:shadow-sm"
                >
                  View merchants
                  <ArrowRight className="size-4 transition-transform group-hover:translate-x-0.5" />
                </Link>
                <Link
                  href="/funding"
                  className="group inline-flex items-center gap-2 rounded-xl border bg-card px-4 py-3 text-sm font-medium text-muted-foreground transition-all hover:border-border hover:bg-muted/50 hover:text-foreground hover:shadow-sm"
                >
                  Funding profiles
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
