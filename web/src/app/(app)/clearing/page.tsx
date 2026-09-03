import { notFound } from "next/navigation";
import { createServerClient } from "@/lib/supabase/server";
import { DASHBOARD_ACCESS, roleFromAppMetadata } from "@/lib/roles";
import { tryFetchAdmin } from "@/lib/adminapi";
import { fmtMinor } from "@/lib/money/minor";
import { PageHeader, PageStack, SectionLabel } from "@/components/page-shell";
import { DataTable, type Column } from "@/components/data-table";
import { MonoChip } from "@/components/ui/badge";
import { DataError } from "@/components/states";

export interface ClearingRecord {
  cycleId: string;
  stan: string;
  mti: string;
  sender: string;
  receiver: string;
  amountMinor: number;
  interchange: number;
  currency: string;
  refId: string;
}

export interface NetPosition {
  cycleId: string;
  member: string;
  net: number;
}

const recordColumns: Column<ClearingRecord>[] = [
  {
    key: "stan",
    header: "STAN",
    render: (r) => <MonoChip>{r.stan}</MonoChip>,
  },
  {
    key: "mti",
    header: "MTI",
    render: (r) => <MonoChip>{r.mti}</MonoChip>,
  },
  {
    key: "sender",
    header: "Sender",
    render: (r) => <span className="font-medium">{r.sender}</span>,
  },
  {
    key: "receiver",
    header: "Receiver",
    render: (r) => <span className="font-medium">{r.receiver}</span>,
  },
  {
    key: "amount",
    header: "Amount",
    render: (r) => (
      <span className="font-medium">{fmtMinor(r.amountMinor, r.currency)}</span>
    ),
  },
  {
    key: "interchange",
    header: "Interchange",
    render: (r) => (
      <span className="text-muted-foreground">
        {fmtMinor(r.interchange, r.currency)}
      </span>
    ),
  },
];

const positionColumns: Column<NetPosition>[] = [
  {
    key: "member",
    header: "Member",
    render: (p) => <span className="font-medium">{p.member}</span>,
  },
  {
    key: "net",
    header: "Net",
    render: (p) => (
      <span className="font-medium">{fmtMinor(p.net)}</span>
    ),
  },
];

export default async function ClearingPage() {
  const supabase = await createServerClient();
  const { data } = await supabase.auth.getUser();
  const role = roleFromAppMetadata(data.user?.app_metadata);
  if (!role || !DASHBOARD_ACCESS[role].includes("/clearing")) notFound();

  const cycles = await tryFetchAdmin<{ items: string[] }>("/clearing/cycles");
  if (!cycles.ok) {
    return (
      <PageStack>
        <PageHeader title="Clearing" description="Captured transactions and interchange calculations." />
        <DataError message={cycles.error} />
      </PageStack>
    );
  }
  if (cycles.data.items.length === 0) {
    return (
      <PageStack>
        <PageHeader
          title="Clearing"
          description="Captured transactions and interchange calculations."
        />
        <DataError
          title="No clearing cycles yet"
          message="Clearing runs after a settlement window captures transactions."
        />
      </PageStack>
    );
  }

  const cycle = cycles.data.items[0];
  const records = await tryFetchAdmin<{ items: ClearingRecord[] }>(
    `/clearing/records?cycle=${encodeURIComponent(cycle)}`
  );
  const positions = await tryFetchAdmin<{ items: NetPosition[] }>(
    `/clearing/positions?cycle=${encodeURIComponent(cycle)}`
  );

  return (
    <PageStack>
      <PageHeader
        title="Clearing"
        description={
          <>
            Cycle{" "}
            <span className="font-mono text-foreground">{cycle}</span> —
            captured transactions and interchange calculations.
          </>
        }
      />

      <div>
        <SectionLabel>Clearing records</SectionLabel>
        {!records.ok ? (
          <DataError message={records.error} />
        ) : (
          <DataTable
            columns={recordColumns}
            rows={records.data.items}
            getKey={(r) => r.stan}
            emptyTitle="No clearing records yet"
            emptyHint="Captured transactions for this cycle will appear here."
          />
        )}
      </div>

      <div>
        <SectionLabel>Net positions</SectionLabel>
        {!positions.ok ? (
          <DataError message={positions.error} />
        ) : (
          <DataTable
            columns={positionColumns}
            rows={positions.data.items}
            getKey={(p) => p.member}
            emptyTitle="No net positions yet"
            emptyHint="Netting results for this cycle will appear here."
          />
        )}
      </div>
    </PageStack>
  );
}