import { Card, CardContent } from "@/components/ui/card";
import { AlertCircle } from "lucide-react";
import { cn } from "@/lib/utils";

/** Full-area loading skeleton, shown by route-level loading.tsx. */
export function PageSkeleton() {
  return (
    <div className="mx-auto w-full max-w-6xl flex-1 px-4 py-8 sm:px-6 lg:px-8">
      <div className="grid gap-6">
        {/* Hero skeleton */}
        <div className="h-40 animate-pulse rounded-2xl border bg-muted/40" />
        {/* Stat-card skeletons */}
        <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
          {Array.from({ length: 3 }).map((_, i) => (
            <div
              key={i}
              className="h-28 animate-pulse rounded-xl border bg-muted/30"
            />
          ))}
        </div>
        {/* Table skeleton */}
        <div className="h-64 animate-pulse rounded-xl border bg-muted/20" />
      </div>
    </div>
  );
}

/** Inline error card for a failed admin-api fetch. */
export function DataError({
  title = "Couldn't load data",
  message,
  className,
}: {
  title?: string;
  message?: string;
  className?: string;
}) {
  return (
    <Card className={cn("overflow-hidden", className)}>
      <CardContent className="flex flex-col items-center justify-center gap-3 py-12 text-center">
        <div className="flex size-9 items-center justify-center rounded-full bg-destructive/10 text-destructive">
          <AlertCircle className="size-4" />
        </div>
        <p className="text-sm font-medium text-foreground">{title}</p>
        {message && (
          <p className="max-w-md text-xs text-muted-foreground">{message}</p>
        )}
      </CardContent>
    </Card>
  );
}

/** Inline empty-state card with a subtle icon for empty data sets. */
export function EmptyState({
  title,
  hint,
  icon,
  className,
}: {
  title: string;
  hint?: string;
  icon?: React.ReactNode;
  className?: string;
}) {
  return (
    <Card className={cn("overflow-hidden", className)}>
      <CardContent className="flex flex-col items-center justify-center gap-3 py-12 text-center">
        {icon && (
          <div className="flex size-9 items-center justify-center rounded-full bg-primary/10 text-primary">
            {icon}
          </div>
        )}
        <p className="text-sm font-medium text-foreground">{title}</p>
        {hint && (
          <p className="max-w-md text-xs text-muted-foreground">{hint}</p>
        )}
      </CardContent>
    </Card>
  );
}
