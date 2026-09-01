import { notFound } from "next/navigation";
import { createServerClient } from "@/lib/supabase/server";
import { DASHBOARD_ACCESS, roleFromAppMetadata } from "@/lib/roles";
import { tryFetchAdmin } from "@/lib/adminapi";
import { fmtCount } from "@/lib/format";
import type { Page } from "@/types/admin";
import { PageHeader, PageStack, SectionLabel } from "@/components/page-shell";
import { DataTable, type Column } from "@/components/data-table";
import { Badge, MonoChip } from "@/components/ui/badge";
import { DataError } from "@/components/states";

export interface Card {
  ref: string;
  panHash: string;
  panMask: string;
  bin: string;
  expiry: string;
  status: string;
  product: string;
  lastAtc: number;
}

const statusTone = (s: string) =>
  s === "active"
    ? "success"
    : s === "blocked"
      ? "danger"
      : "neutral";

const columns: Column<Card>[] = [
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
  {
    key: "bin",
    header: "BIN",
    render: (c) => <MonoChip>{c.bin}</MonoChip>,
  },
  { key: "product", header: "Product", render: (c) => c.product },
  {
    key: "status",
    header: "Status",
    render: (c) => <Badge tone={statusTone(c.status)}>{c.status}</Badge>,
  },
  { key: "expiry", header: "Expiry", render: (c) => c.expiry },
  {
    key: "lastAtc",
    header: "Last ATC",
    render: (c) => <span className="font-mono text-xs">{c.lastAtc}</span>,
  },
];

export default async function CardsPage() {
  const supabase = await createServerClient();
  const { data } = await supabase.auth.getUser();
  const role = roleFromAppMetadata(data.user?.app_metadata);
  if (!role || !DASHBOARD_ACCESS[role].includes("/cards")) notFound();
  const page = await tryFetchAdmin<Page<Card>>("/cards?limit=50");

  return (
    <PageStack>
      <PageHeader
        title="Cards"
        tone="sky"
        description="Issued cards with BIN ranges, status, and application cryptogram data."
      />
      <div>
        <SectionLabel>Issued cards</SectionLabel>
        {!page.ok ? (
          <DataError message={page.error} />
        ) : (
          <DataTable
            columns={columns}
            rows={page.data.items}
            getKey={(c) => c.ref}
            emptyTitle="No cards yet"
            emptyHint="Issued cards will appear here."
            footer={
              page.data.total > 0
                ? `Showing ${page.data.items.length} of ${fmtCount(page.data.total)} cards`
                : undefined
            }
          />
        )}
      </div>
    </PageStack>
  );
}