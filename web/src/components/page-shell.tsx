import { cn } from "@/lib/utils";
import type { ReactNode } from "react";

type Tone =
  | "default"
  | "violet"
  | "sky"
  | "emerald"
  | "amber"
  | "rose";

const TONE_BG: Record<Tone, string> = {
  default: "from-primary/[0.06]",
  violet: "from-violet-500/[0.06]",
  sky: "from-sky-500/[0.06]",
  emerald: "from-emerald-500/[0.06]",
  amber: "from-amber-500/[0.06]",
  rose: "from-rose-500/[0.06]",
};

const TONE_ORB: Record<Tone, string> = {
  default: "bg-primary/[0.08]",
  violet: "bg-violet-500/[0.08]",
  sky: "bg-sky-500/[0.08]",
  emerald: "bg-emerald-500/[0.08]",
  amber: "bg-amber-500/[0.08]",
  rose: "bg-rose-500/[0.08]",
};

/**
 * Consistent page hero card (title + description) with the option of a
 * right-aligned actions slot. Use on every dashboard page so the chrome
 * stays uniform across the app.
 */
export function PageHeader({
  title,
  description,
  tone = "default",
  actions,
  className,
}: {
  title: string;
  description?: ReactNode;
  tone?: Tone;
  actions?: ReactNode;
  className?: string;
}) {
  return (
    <div
      className={cn(
        "relative overflow-hidden rounded-2xl border bg-gradient-to-br via-transparent to-transparent p-6 sm:p-8",
        TONE_BG[tone],
        className
      )}
    >
      <div
        className={cn(
          "pointer-events-none absolute -top-20 -right-20 size-40 rounded-full blur-[60px]",
          TONE_ORB[tone]
        )}
      />
      <div className="relative flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div className="max-w-xl">
          <h1 className="text-2xl font-bold tracking-tight sm:text-3xl">
            {title}
          </h1>
          {description && (
            <p className="mt-1.5 text-sm text-muted-foreground">
              {description}
            </p>
          )}
        </div>
        {actions && (
          <div className="flex flex-wrap items-center gap-2">{actions}</div>
        )}
      </div>
    </div>
  );
}

/**
 * Section label used to group content on a page (e.g. "Today's activity").
 */
export function SectionLabel({
  children,
  className,
}: {
  children: ReactNode;
  className?: string;
}) {
  return (
    <h2
      className={cn(
        "mb-3 text-xs font-semibold uppercase tracking-wider text-muted-foreground/70",
        className
      )}
    >
      {children}
    </h2>
  );
}

/**
 * Standard page vertical rhythm.
 */
export function PageStack({
  children,
  className,
}: {
  children: ReactNode;
  className?: string;
}) {
  return <div className={cn("grid gap-6", className)}>{children}</div>;
}
