import { cn } from "@/lib/utils";
import type { ReactNode } from "react";

export type BadgeTone =
  | "default"
  | "success"
  | "warning"
  | "danger"
  | "info"
  | "neutral";

const TONE_MAP: Record<BadgeTone, string> = {
  default: "bg-muted text-muted-foreground",
  neutral: "bg-muted/60 text-muted-foreground",
  success: "bg-emerald-500/10 text-emerald-600 dark:text-emerald-400",
  warning: "bg-amber-500/10 text-amber-600 dark:text-amber-400",
  danger: "bg-red-500/10 text-red-600 dark:text-red-400",
  info: "bg-sky-500/10 text-sky-600 dark:text-sky-400",
};

/**
 * Semantic status pill. Standardizes the status color language across tables
 * and cards so green/amber/red always mean the same thing.
 */
export function Badge({
  children,
  tone = "default",
  className,
  title,
}: {
  children: ReactNode;
  tone?: BadgeTone;
  className?: string;
  title?: string;
}) {
  return (
    <span
      title={title}
      className={cn(
        "inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-xs font-medium whitespace-nowrap",
        TONE_MAP[tone],
        className
      )}
    >
      {children}
    </span>
  );
}

/** Tiny mono chip for codes (STAN, MTI, BIN, reason codes, ids). */
export function MonoChip({
  children,
  className,
}: {
  children: ReactNode;
  className?: string;
}) {
  return (
    <span
      className={cn(
        "inline-flex rounded-md bg-muted px-1.5 py-0.5 font-mono text-xs text-foreground/90",
        className
      )}
    >
      {children}
    </span>
  );
}
