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
import { Layers, Coins, ArrowRight } from "lucide-react";

export default async function IssuerPage() {
  const supabase = await createServerClient();
  const { data } = await supabase.auth.getUser();
  if (roleFromAppMetadata(data.user?.app_metadata) !== "issuer") notFound();
  const dashboard = await tryFetchAdmin<DashboardSummary>("/dashboard");

  return (
    <PageStack>
      <PageHeader
        title="Issuer"
        tone="sky"
        description="Manage cards, tokens, and BIN ranges for your issuing portfolio."
      />

      {!dashboard.ok ? (
        <DataError message={dashboard.error} />
      ) : (
        <div>
          <SectionLabel>Portfolio</SectionLabel>
          <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
            <StatCard
              title="Cards"
              value={fmtCount(dashboard.data.cards)}
              hint="Active cards in portfolio"
              icon={<Layers className="size-5" />}
              accent="info"
              sparkData={sparkFromValue(dashboard.data.cards)}
            />
            <StatCard
              title="Tokens"
              value={fmtCount(dashboard.data.tokens)}
              hint="Digital tokens for mobile & contactless"
              icon={<Coins className="size-5" />}
              accent="default"
              sparkData={sparkFromValue(dashboard.data.tokens)}
            />
          </div>
        </div>
      )}

      <div className="flex flex-wrap gap-3">
        <Link
          href="/cards"
          className="group inline-flex items-center gap-2 rounded-xl border bg-card px-4 py-3 text-sm font-medium text-muted-foreground transition-all hover:border-border hover:bg-muted/50 hover:text-foreground hover:shadow-sm"
        >
          View cards
          <ArrowRight className="size-4 transition-transform group-hover:translate-x-0.5" />
        </Link>
        <Link
          href="/tokens"
          className="group inline-flex items-center gap-2 rounded-xl border bg-card px-4 py-3 text-sm font-medium text-muted-foreground transition-all hover:border-border hover:bg-muted/50 hover:text-foreground hover:shadow-sm"
        >
          View tokens
          <ArrowRight className="size-4 transition-transform group-hover:translate-x-0.5" />
        </Link>
      </div>
    </PageStack>
  );
}