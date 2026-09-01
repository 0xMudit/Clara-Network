import { cn } from "@/lib/utils";
import type { ReactNode } from "react";

export interface Column<T> {
  key: string;
  header: ReactNode;
  className?: string;
  headerClassName?: string;
  render: (row: T) => ReactNode;
}

/**
 * Reusable table that normalizes the chrome across every dashboard list:
 * a clean container, sticky-ish header, zebra + hover rows, and a consistent
 * empty state. Each page supplies columns + rows.
 */
export function DataTable<T>({
  columns,
  rows,
  getKey,
  emptyTitle = "No data yet",
  emptyHint,
  footer,
  className,
}: {
  columns: Column<T>[];
  rows: T[];
  getKey: (row: T, index: number) => string;
  emptyTitle?: string;
  emptyHint?: string;
  footer?: ReactNode;
  className?: string;
}) {
  return (
    <div
      className={cn(
        "overflow-hidden rounded-xl border bg-card ring-1 ring-border/50",
        className
      )}
    >
      <div className="overflow-x-auto">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b bg-muted/40 text-left text-xs font-semibold uppercase tracking-wider text-muted-foreground/70">
              {columns.map((c) => (
                <th key={c.key} className={cn("px-4 py-3", c.headerClassName)}>
                  {c.header}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {rows.map((row, i) => (
              <tr
                key={getKey(row, i)}
                className={cn(
                  "border-b transition-colors last:border-0 hover:bg-muted/30",
                  i % 2 === 1 && "bg-muted/10"
                )}
              >
                {columns.map((c) => (
                  <td key={c.key} className={cn("px-4 py-2.5", c.className)}>
                    {c.render(row)}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {rows.length === 0 && (
        <div className="flex flex-col items-center justify-center py-12 text-center">
          <p className="text-sm font-medium text-muted-foreground">
            {emptyTitle}
          </p>
          {emptyHint && (
            <p className="mt-1 text-xs text-muted-foreground/60">{emptyHint}</p>
          )}
        </div>
      )}

      {footer && (
        <div className="border-t bg-muted/40 px-4 py-2.5 text-xs text-muted-foreground/70">
          {footer}
        </div>
      )}
    </div>
  );
}
