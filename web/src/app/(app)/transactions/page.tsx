import { notFound } from "next/navigation";
import Link from "next/link";
import { createServerClient } from "@/lib/supabase/server";
import { roleFromAppMetadata } from "@/lib/roles";
import { tryFetchAdmin } from "@/lib/adminapi";
import { fmtTs } from "@/lib/date";
import { fmtCount } from "@/lib/format";
import type { Page } from "@/types/admin";
import { PageHeader, PageStack, SectionLabel } from "@/components/page-shell";
import { DataTable, type Column } from "@/components/data-table";
import { MonoChip } from "@/components/ui/badge";
import { DataError } from "@/components/states";
import { ArrowRight } from "lucide-react";

export interface AuditEvent {
  stan: string;
  mti: string;
  pan: string;
  amount: string;
  responseCode: string;
  destination: string;
  createdAt: string;
}

const columns: Column<AuditEvent>[] = [
  { key: "stan", header: "STAN", render: (t) => <MonoChip>{t.stan}</MonoChip> },
  { key: "mti", header: "MTI", render: (t) => <MonoChip>{t.mti}</MonoChip> },
  {
    key: "pan",
    header: "PAN",
    render: (t) => <span className="font-mono text-xs">{t.pan}</span>,
  },
  {
    key: "amount",
    header: "Amount",
    render: (t) => <span className="font-medium">{t.amount}</span>,
  },
  {
    key: "responseCode",
    header: "Resp",
    render: (t) => <MonoChip>{t.responseCode}</MonoChip>,
  },
  {
    key: "destination",
    header: "Dest",
    render: (t) => <span className="font-mono text-xs">{t.destination}</span>,
  },
  {
    key: "createdAt",
    header: "Time",
    render: (t) => (
      <span className="text-xs text-muted-foreground">{fmtTs(t.createdAt)}</span>
    ),
  },
];

export default async function TransactionsPage() {
  const supabase = await createServerClient();
  const { data } = await supabase.auth.getUser();
  if (roleFromAppMetadata(data.user?.app_metadata) !== "scheme_operator")
    notFound();
  const page = await tryFetchAdmin<Page<AuditEvent>>("/transactions?limit=50");

  return (
    <PageStack>
      <PageHeader
        title="Transactions"
        description="Authorization requests flowing through the switch — every STAN, response, and destination."
        actions={
          <Link
            href="/ops"
            className="group inline-flex items-center gap-2 rounded-xl border bg-card px-4 py-2.5 text-sm font-medium text-muted-foreground transition-all hover:border-border hover:bg-muted/50 hover:text-foreground hover:shadow-sm"
          >
            Back to operations
            <ArrowRight className="size-4 transition-transform group-hover:translate-x-0.5" />
          </Link>
        }
      />

      <div>
        <SectionLabel>Switch audit log</SectionLabel>
        {!page.ok ? (
          <DataError message={page.error} />
        ) : (
          <DataTable
            columns={columns}
            rows={page.data.items}
            getKey={(t) => t.stan}
            emptyTitle="No transactions yet"
            emptyHint="The switch records every authorization here once traffic starts flowing."
            footer={
              page.data.total > 0
                ? `Showing ${page.data.items.length} of ${fmtCount(page.data.total)} transactions`
                : undefined
            }
          />
        )}
      </div>
    </PageStack>
  );
}