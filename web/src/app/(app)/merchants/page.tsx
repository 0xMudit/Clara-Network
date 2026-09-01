import { notFound } from "next/navigation";
import { createServerClient } from "@/lib/supabase/server";
import { DASHBOARD_ACCESS, roleFromAppMetadata } from "@/lib/roles";
import { tryFetchAdmin } from "@/lib/adminapi";
import { fmtMinor } from "@/lib/money/minor";
import { fmtCount } from "@/lib/format";
import type { Page } from "@/types/admin";
import { PageHeader, PageStack, SectionLabel } from "@/components/page-shell";
import { DataTable, type Column } from "@/components/data-table";
import { Badge, MonoChip } from "@/components/ui/badge";
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

const riskTone = (r: string) =>
  r === "low" ? "success" : r === "high" ? "danger" : "warning";
const statusTone = (s: string) =>
  s === "approved" ? "success" : s === "declined" ? "danger" : "neutral";

const columns: Column<Merchant>[] = [
  {
    key: "name",
    header: "Name",
    render: (m) => <span className="font-medium">{m.name}</span>,
  },
  {
    key: "dba",
    header: "DBA",
    render: (m) => <span className="text-muted-foreground">{m.dba}</span>,
  },
  {
    key: "mccs",
    header: "MCCs",
    render: (m) => (
      <div className="flex flex-wrap gap-1">
        {m.mccs.map((mcc) => (
          <MonoChip key={mcc}>{mcc}</MonoChip>
        ))}
      </div>
    ),
  },
  {
    key: "riskTier",
    header: "Risk",
    render: (m) => <Badge tone={riskTone(m.riskTier)}>{m.riskTier}</Badge>,
  },
  {
    key: "status",
    header: "Status",
    render: (m) => <Badge tone={statusTone(m.status)}>{m.status}</Badge>,
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

export default async function MerchantsPage() {
  const supabase = await createServerClient();
  const { data } = await supabase.auth.getUser();
  const role = roleFromAppMetadata(data.user?.app_metadata);
  if (!role || !DASHBOARD_ACCESS[role].includes("/merchants")) notFound();
  const page = await tryFetchAdmin<Page<Merchant>>("/merchants?limit=50");

  return (
    <PageStack>
      <PageHeader
        title="Merchants"
        tone="emerald"
        description="Boarded merchants with risk tiers, MCC codes, and funding status."
      />
      <div>
        <SectionLabel>Merchant portfolio</SectionLabel>
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