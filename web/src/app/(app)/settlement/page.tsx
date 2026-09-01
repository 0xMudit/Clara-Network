import { notFound } from "next/navigation";
import { createServerClient } from "@/lib/supabase/server";
import { DASHBOARD_ACCESS, roleFromAppMetadata } from "@/lib/roles";
import { tryFetchAdmin } from "@/lib/adminapi";
import { fmtMinor } from "@/lib/money/minor";
import { fmtTs } from "@/lib/date";
import { StatCard } from "@/components/cards/stat-card";
import { PageHeader, PageStack, SectionLabel } from "@/components/page-shell";
import { DataTable, type Column } from "@/components/data-table";
import { Badge, MonoChip } from "@/components/ui/badge";
import { DataError } from "@/components/states";
import { Coins } from "lucide-react";

export interface PrefundAccount {
  member: string;
  balance: number;
  cap: number;
}

export interface SettlementInstruction {
  cycleId: string;
  msgId: string;
  member: string;
  amount: number;
  direction: string;
  currency: string;
  instruction: string;
  final: boolean;
}

const prefundColumns: Column<PrefundAccount>[] = [
  {
    key: "member",
    header: "Member",
    render: (a) => <span className="font-medium">{a.member}</span>,
  },
  {
    key: "balance",
    header: "Balance",
    render: (a) => <span className="font-medium">{fmtMinor(a.balance, "EUR")}</span>,
  },
  {
    key: "cap",
    header: "Cap",
    render: (a) => <span className="text-muted-foreground">{fmtMinor(a.cap, "EUR")}</span>,
  },
];

const instructionColumns: Column<SettlementInstruction>[] = [
  {
    key: "msgId",
    header: "MsgId",
    render: (inst) => <MonoChip>{inst.msgId}</MonoChip>,
  },
  {
    key: "member",
    header: "Member",
    render: (inst) => <span className="font-medium">{inst.member}</span>,
  },
  {
    key: "direction",
    header: "Dir",
    render: (inst) => (
      <Badge tone={inst.direction === "credit" ? "success" : "info"}>
        {inst.direction}
      </Badge>
    ),
  },
  {
    key: "amount",
    header: "Amount",
    render: (inst) => (
      <span className="font-medium">{fmtMinor(inst.amount, inst.currency)}</span>
    ),
  },
  {
    key: "final",
    header: "Final",
    render: (inst) =>
      inst.final ? (
        <Badge tone="success">Final</Badge>
      ) : (
        <span className="text-muted-foreground/40">—</span>
      ),
  },
  {
    key: "instruction",
    header: "Time",
    render: (inst) => (
      <span className="text-xs text-muted-foreground">
        {fmtTs(inst.instruction)}
      </span>
    ),
  },
];

export default async function SettlementPage() {
  const supabase = await createServerClient();
  const { data } = await supabase.auth.getUser();
  const role = roleFromAppMetadata(data.user?.app_metadata);
  if (!role || !DASHBOARD_ACCESS[role].includes("/settlement")) notFound();

  const prefunds = await tryFetchAdmin<{ items: PrefundAccount[] }>(
    "/settlement/prefunds"
  );
  const df = await tryFetchAdmin<{ balance: number }>(
    "/settlement/default-fund"
  );
  const cycles = await tryFetchAdmin<{ items: string[] }>("/clearing/cycles");
  const cycle = cycles.ok ? cycles.data.items[0] : undefined;
  const instructions = cycle
    ? await tryFetchAdmin<{ items: SettlementInstruction[] }>(
        `/settlement/instructions?cycle=${encodeURIComponent(cycle)}`
      )
    : null;

  return (
    <PageStack>
      <PageHeader
        title="Settlement"
        description="Prefund balances, net positions, and settlement instructions."
      />

      <div>
        <SectionLabel>Default fund</SectionLabel>
        {!df.ok ? (
          <DataError message={df.error} />
        ) : (
          <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
            <StatCard
              title="Default fund balance"
              value={fmtMinor(df.data.balance, "EUR")}
              icon={<Coins className="size-5" />}
              accent="success"
            />
          </div>
        )}
      </div>

      <div>
        <SectionLabel>Prefund balances</SectionLabel>
        {!prefunds.ok ? (
          <DataError message={prefunds.error} />
        ) : (
          <DataTable
            columns={prefundColumns}
            rows={prefunds.data.items}
            getKey={(a) => a.member}
            emptyTitle="No prefund accounts yet"
            emptyHint="Member prefund accounts will appear here."
          />
        )}
      </div>

      <div>
        <SectionLabel>Latest instructions</SectionLabel>
        {instructions ? (
          !instructions.ok ? (
            <DataError message={instructions.error} />
          ) : (
            <DataTable
              columns={instructionColumns}
              rows={instructions.data.items}
              getKey={(inst) => inst.msgId}
              emptyTitle="No settlement instructions yet"
              emptyHint="Generated instructions for the latest cycle will appear here."
            />
          )
        ) : (
          <DataError
            title="No settlement instructions yet"
            message="Instructions are generated once a clearing cycle exists."
          />
        )}
      </div>
    </PageStack>
  );
}