import { notFound } from "next/navigation";
import { createServerClient } from "@/lib/supabase/server";
import { DASHBOARD_ACCESS, roleFromAppMetadata } from "@/lib/roles";
import { tryFetchAdmin } from "@/lib/adminapi";
import { fmtMinor } from "@/lib/money/minor";
import { fmtCount } from "@/lib/format";
import type { Page } from "@/types/admin";
import { PageHeader, PageStack, SectionLabel } from "@/components/page-shell";
import { DataTable, type Column } from "@/components/data-table";
import { Badge } from "@/components/ui/badge";
import { DataError } from "@/components/states";

export interface Merchant {
  id: string;
  name: string;
  dba: string;
  taxId: string;
  mccs: string[];
  status: string;
  riskTier: string;
  reserveRateBps: number;
  fundingDelayDays: number;
  transactionLimit: number;
  reserveBalance: number;
  volume: number;
  declineReason?: string;
  approvedAt: string;
}

const columns: Column<Merchant>[] = [
  {
    key: "name",
    header: "Name",
    render: (m) => <span className="font-medium">{m.name}</span>,
  },
  {
    key: "delay",
    header: "Delay (days)",
    render: (m) => (
      <Badge tone="info">
        {m.fundingDelayDays}d
      </Badge>
    ),
  },
  {
    key: "txnLimit",
    header: "Txn limit",
    render: (m) => (
      <span className="font-medium">{fmtMinor(m.transactionLimit)}</span>
    ),
  },
  {
    key: "reserve",
    header: "Reserve",
    render: (m) => (
      <span className="font-medium">{fmtMinor(m.reserveBalance)}</span>
    ),
  },
  {
    key: "volume",
    header: "Volume",
    render: (m) => <span className="font-medium">{fmtMinor(m.volume)}</span>,
  },
];

export default async function FundingPage() {
  const supabase = await createServerClient();
  const { data } = await supabase.auth.getUser();
  const role = roleFromAppMetadata(data.user?.app_metadata);
  if (!role || !DASHBOARD_ACCESS[role].includes("/funding")) notFound();
  const page = await tryFetchAdmin<Page<Merchant>>("/merchants?limit=50");

  return (
    <PageStack>
      <PageHeader
        title="Acquirer funding"
        tone="emerald"
        description="Merchant funding schedules, transaction limits, and reserve balances."
      />
      <div>
        <SectionLabel>Funding profiles</SectionLabel>
        {!page.ok ? (
          <DataError message={page.error} />
        ) : (
          <DataTable
            columns={columns}
            rows={page.data.items}
            getKey={(m) => m.id}
            emptyTitle="No merchants yet"
            emptyHint="Boarded merchants will appear here."
            footer={
              page.data.total > 0
                ? `Showing ${page.data.items.length} of ${fmtCount(page.data.total)} merchants`
                : undefined
            }
          />
        )}
      </div>
    </PageStack>
  );
}