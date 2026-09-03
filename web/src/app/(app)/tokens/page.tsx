import { notFound } from "next/navigation";
import { createServerClient } from "@/lib/supabase/server";
import { DASHBOARD_ACCESS, roleFromAppMetadata } from "@/lib/roles";
import { tryFetchAdmin } from "@/lib/adminapi";
import { fmtTs } from "@/lib/date";
import { fmtCount } from "@/lib/format";
import type { Page } from "@/types/admin";
import { PageHeader, PageStack, SectionLabel } from "@/components/page-shell";
import { DataTable, type Column } from "@/components/data-table";
import { Badge, MonoChip } from "@/components/ui/badge";
import { DataError } from "@/components/states";

export interface Token {
  token: string;
  par: string;
  status: string;
  bin: string;
  requestor: string;
  deviceId: string;
  createdAt: string;
}

const columns: Column<Token>[] = [
  {
    key: "token",
    header: "Token",
    render: (t) => <span className="font-mono text-xs">{t.token}</span>,
  },
  {
    key: "par",
    header: "PAR",
    render: (t) => <MonoChip>{t.par}</MonoChip>,
  },
  {
    key: "bin",
    header: "BIN",
    render: (t) => <MonoChip>{t.bin}</MonoChip>,
  },
  {
    key: "requestor",
    header: "Requestor",
    render: (t) => <span className="font-medium">{t.requestor}</span>,
  },
  {
    key: "status",
    header: "Status",
    render: (t) => (
      <Badge tone={t.status === "active" ? "success" : "neutral"}>
        {t.status}
      </Badge>
    ),
  },
  {
    key: "createdAt",
    header: "Created",
    render: (t) => (
      <span className="text-xs text-muted-foreground">
        {fmtTs(t.createdAt)}
      </span>
    ),
  },
];

export default async function TokensPage() {
  const supabase = await createServerClient();
  const { data } = await supabase.auth.getUser();
  const role = roleFromAppMetadata(data.user?.app_metadata);
  if (!role || !DASHBOARD_ACCESS[role].includes("/tokens")) notFound();
  const page = await tryFetchAdmin<Page<Token>>("/tokens?limit=50");

  return (
    <PageStack>
      <PageHeader
        title="Tokens"
        tone="sky"
        description="Digital payment tokens for contactless, mobile, and in-app transactions."
      />
      <div>
        <SectionLabel>Token vault</SectionLabel>
        {!page.ok ? (
          <DataError message={page.error} />
        ) : (
          <DataTable
            columns={columns}
            rows={page.data.items}
            getKey={(t) => t.token}
            emptyTitle="No tokens yet"
            emptyHint="Provisioned tokens will appear here."
            footer={
              page.data.total > 0
                ? `Showing ${page.data.items.length} of ${fmtCount(page.data.total)} tokens`
                : undefined
            }
          />
        )}
      </div>
    </PageStack>
  );
}