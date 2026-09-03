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
import { Badge, MonoChip } from "@/components/ui/badge";
import { DataError } from "@/components/states";
import { Layers, Coins, ArrowRight, CreditCard } from "lucide-react";
import type { Card } from "../cards/page";

const cardColumns: Column<Card>[] = [
  {
    key: "ref",
    header: "Ref",
    render: (c) => <MonoChip>{c.ref}</MonoChip>,
  },
  {
    key: "pan",
    header: "PAN",
    render: (c) => <span className="font-mono text-xs">{c.panMask}</span>,
  },
  { key: "product", header: "Product", render: (c) => c.product },
  {
    key: "status",
    header: "Status",
    render: (c) => (
      <Badge
        tone={c.status === "active" ? "success" : c.status === "blocked" ? "danger" : "neutral"}
      >
        {c.status}
      </Badge>
    ),
  },
];

export default async function IssuerPage() {
  const supabase = await createServerClient();
  const { data } = await supabase.auth.getUser();
  if (roleFromAppMetadata(data.user?.app_metadata) !== "issuer") notFound();

  const dashboard = await tryFetchAdmin<DashboardSummary>("/dashboard");
  const cards = await tryFetchAdmin<Page<Card>>("/cards?limit=6");
  const series = await tryFetchAdmin<{ items: SeriesPoint[] }>("/dashboard/series?days=14");
  const seriesData: SeriesPointDatum[] = series.ok
    ? series.data.items.map((p) => ({ date: p.date, count: p.count }))
    : [];

  return (
    <PageStack>
      <PageHeader
        title="Issuer"
        tone="sky"
        description="Your card portfolio, token vault, and authorization activity across the network."
      />

      {!dashboard.ok ? (
        <DataError message={dashboard.error} />
      ) : (
        <div>
          <SectionLabel>Portfolio</SectionLabel>
          <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
            <StatCard
              title="Cards issued"
              value={fmtCount(dashboard.data.cards)}
              hint="Cards active across your BIN ranges"
              icon={<Layers className="size-5" />}
              accent="info"
            />
            <StatCard
              title="Tokens provisioned"
              value={fmtCount(dashboard.data.tokens)}
              hint="Digital tokens for mobile & contactless"
              icon={<Coins className="size-5" />}
              accent="default"
            />
            <StatCard
              title="Network transactions"
              value={fmtCount(dashboard.data.transactions)}
              hint="Authorizations routed via the switch"
              icon={<CreditCard className="size-5" />}
              accent="success"
            />
          </div>

          <div className="mt-8 grid gap-6 lg:grid-cols-2">
            <div>
              <SectionLabel>Authorization volume</SectionLabel>
              {!series.ok ? (
                <DataError message={series.error} />
              ) : seriesData.length === 0 ? (
                <p className="rounded-xl border border-dashed bg-card/40 px-4 py-8 text-center text-sm text-muted-foreground">
                  No authorization history yet — run the issuer demo to start moving volume.
                </p>
              ) : (
                <div className="rounded-2xl border bg-card p-5">
                  <TransactionChart data={seriesData} />
                </div>
              )}
            </div>

            <div>
              <SectionLabel>Recently issued cards</SectionLabel>
              {!cards.ok ? (
                <DataError message={cards.error} />
              ) : (
                <DataTable
                  columns={cardColumns}
                  rows={cards.data.items}
                  getKey={(c) => c.ref}
                  emptyTitle="No cards yet"
                  emptyHint="Issued cards will appear here."
                />
              )}
              <Link
                href="/cards"
                className="group mt-4 inline-flex items-center gap-2 rounded-xl border bg-card px-4 py-3 text-sm font-medium text-muted-foreground transition-all hover:border-border hover:bg-muted/50 hover:text-foreground hover:shadow-sm"
              >
                View all cards
                <ArrowRight className="size-4 transition-transform group-hover:translate-x-0.5" />
              </Link>
            </div>
          </div>
        </div>
      )}
    </PageStack>
  );
}
