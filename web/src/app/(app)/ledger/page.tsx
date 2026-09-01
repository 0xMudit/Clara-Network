import { notFound } from "next/navigation";
import { createServerClient } from "@/lib/supabase/server";
import { DASHBOARD_ACCESS, roleFromAppMetadata } from "@/lib/roles";
import { tryFetchAdmin } from "@/lib/adminapi";
import { fmtMinor } from "@/lib/money/minor";
import { PageHeader, PageStack, SectionLabel } from "@/components/page-shell";
import { DataTable, type Column } from "@/components/data-table";
import { Badge, MonoChip } from "@/components/ui/badge";
import { DataError } from "@/components/states";

export interface LedgerAccount {
  id: string;
  type: string;
  balance: number;
}

const columns: Column<LedgerAccount>[] = [
  {
    key: "id",
    header: "Account",
    render: (a) => <MonoChip>{a.id}</MonoChip>,
  },
  {
    key: "type",
    header: "Type",
    render: (a) => <Badge tone="neutral">{a.type}</Badge>,
  },
  {
    key: "balance",
    header: "Balance",
    render: (a) => <span className="font-medium">{fmtMinor(a.balance)}</span>,
  },
];

export default async function LedgerPage() {
  const supabase = await createServerClient();
  const { data } = await supabase.auth.getUser();
  const role = roleFromAppMetadata(data.user?.app_metadata);
  if (!role || !DASHBOARD_ACCESS[role].includes("/ledger")) notFound();
  const accounts = await tryFetchAdmin<{ items: LedgerAccount[] }>(
    "/ledger/accounts"
  );

  return (
    <PageStack>
      <PageHeader
        title="Ledger"
        description="Double-entry accounting — every transaction has a debit and a credit."
      />
      <div>
        <SectionLabel>Accounts</SectionLabel>
        {!accounts.ok ? (
          <DataError message={accounts.error} />
        ) : (
          <DataTable
            columns={columns}
            rows={accounts.data.items}
            getKey={(a) => a.id}
            emptyTitle="No ledger accounts yet"
            emptyHint="Double-entry accounts will appear here."
          />
        )}
      </div>
    </PageStack>
  );
}