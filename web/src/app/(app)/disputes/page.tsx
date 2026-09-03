import { notFound } from "next/navigation";
import { createServerClient } from "@/lib/supabase/server";
import { DASHBOARD_ACCESS, roleFromAppMetadata } from "@/lib/roles";
import { tryFetchAdmin } from "@/lib/adminapi";
import { fmtMinor } from "@/lib/money/minor";
import { fmtTs } from "@/lib/date";
import { fmtCount } from "@/lib/format";
import type { Page } from "@/types/admin";
import { PageHeader, PageStack, SectionLabel } from "@/components/page-shell";
import { DataTable, type Column } from "@/components/data-table";
import { Badge, MonoChip } from "@/components/ui/badge";
import { DataError } from "@/components/states";

export interface Dispute {
  id: string;
  refId: string;
  merchantId: string;
  cardholder: string;
  amountMinor: number;
  currency: string;
  reasonCode: string;
  category: string;
  stage: string;
  status: string;
  filedAt: string;
  responseDue: string;
  respondedAt?: string;
  escalatedAt?: string;
  decision?: string;
  winner?: string;
  decisionAt?: string;
  disputeFee: number;
  arbitrationFee: number;
  note?: string;
}

const stageTone = (s: string) =>
  s === "arbitration"
    ? "danger"
    : s === "representment"
      ? "warning"
      : "neutral";
const statusTone = (s: string) =>
  s === "won" ? "success" : s === "lost" ? "danger" : "neutral";

const columns: Column<Dispute>[] = [
  {
    key: "id",
    header: "ID",
    render: (d) => <MonoChip>{d.id}</MonoChip>,
  },
  {
    key: "cardholder",
    header: "Cardholder",
    render: (d) => <span className="font-medium">{d.cardholder}</span>,
  },
  {
    key: "reason",
    header: "Reason",
    render: (d) => <MonoChip>{d.reasonCode}</MonoChip>,
  },
  {
    key: "category",
    header: "Category",
    render: (d) => (
      <span className="text-xs text-muted-foreground">{d.category}</span>
    ),
  },
  {
    key: "stage",
    header: "Stage",
    render: (d) => <Badge tone={stageTone(d.stage)}>{d.stage}</Badge>,
  },
  {
    key: "status",
    header: "Status",
    render: (d) => <Badge tone={statusTone(d.status)}>{d.status}</Badge>,
  },
  {
    key: "amount",
    header: "Amount",
    render: (d) => (
      <span className="font-medium">{fmtMinor(d.amountMinor, d.currency)}</span>
    ),
  },
  {
    key: "filedAt",
    header: "Filed",
    render: (d) => (
      <span className="text-xs text-muted-foreground">{fmtTs(d.filedAt)}</span>
    ),
  },
  {
    key: "responseDue",
    header: "Due",
    render: (d) =>
      d.responseDue ? (
        <span className="text-xs text-muted-foreground">
          {fmtTs(d.responseDue)}
        </span>
      ) : (
        <span className="text-muted-foreground/40">—</span>
      ),
  },
];

export default async function DisputesPage() {
  const supabase = await createServerClient();
  const { data } = await supabase.auth.getUser();
  const role = roleFromAppMetadata(data.user?.app_metadata);
  if (!role || !DASHBOARD_ACCESS[role].includes("/disputes")) notFound();
  const page = await tryFetchAdmin<Page<Dispute>>("/disputes?limit=50");

  return (
    <PageStack>
      <PageHeader
        title="Disputes"
        tone="rose"
        description="Chargebacks, representments, and arbitration — track every case from filing to resolution."
      />
      <div>
        <SectionLabel>Dispute cases</SectionLabel>
        {!page.ok ? (
          <DataError message={page.error} />
        ) : (
          <DataTable
            columns={columns}
            rows={page.data.items}
            getKey={(d) => d.id}
            emptyTitle="No disputes yet"
            emptyHint="Filed disputes will appear here."
            footer={
              page.data.total > 0
                ? `Showing ${page.data.items.length} of ${fmtCount(page.data.total)} disputes`
                : undefined
            }
          />
        )}
      </div>
    </PageStack>
  );
}