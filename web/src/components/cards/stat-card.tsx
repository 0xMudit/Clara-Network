"use client";
import { Card, CardContent } from "@/components/ui/card";
import { cn } from "@/lib/utils";
import type { ReactNode } from "react";

interface StatCardProps {
  title: string;
  value: string;
  hint?: string;
  icon?: ReactNode;
  trend?: "up" | "down" | "neutral";
  trendValue?: string;
  accent?: "default" | "success" | "warning" | "danger" | "info";
}

const ACCENT_MAP = {
  default: {
    iconBg: "bg-primary/10",
    iconText: "text-primary",
    ring: "ring-primary/20",
  },
  success: {
    iconBg: "bg-emerald-500/10",
    iconText: "text-emerald-500",
    ring: "ring-emerald-500/20",
  },
  warning: {
    iconBg: "bg-amber-500/10",
    iconText: "text-amber-500",
    ring: "ring-amber-500/20",
  },
  danger: {
    iconBg: "bg-red-500/10",
    iconText: "text-red-500",
    ring: "ring-red-500/20",
  },
  info: {
    iconBg: "bg-sky-500/10",
    iconText: "text-sky-500",
    ring: "ring-sky-500/20",
  },
};

export function StatCard({
  title,
  value,
  hint,
  icon,
  trend,
  trendValue,
  accent = "default",
}: StatCardProps) {
  const a = ACCENT_MAP[accent];

  return (
    <Card
      className={cn(
        "group relative overflow-hidden transition-all duration-300",
        "hover:shadow-lg hover:shadow-primary/8 hover:-translate-y-0.5",
        "ring-1 ring-border/60",
        "hover:ring-primary/20"
      )}
    >
      {/* Gradient glow on hover */}
      <div className="pointer-events-none absolute inset-0 bg-gradient-to-br from-primary/[0.04] via-accent/[0.02] to-transparent opacity-0 transition-opacity duration-300 group-hover:opacity-100" />

      <CardContent className="relative">
        <div className="flex items-start justify-between gap-3">
          <div className="flex-1 min-w-0">
            <p className="text-xs font-medium uppercase tracking-wider text-muted-foreground/80">
              {title}
            </p>
            <p className="mt-1.5 text-2xl font-bold tracking-tight text-foreground">
              {value}
            </p>

            {/* Trend indicator */}
            {trend && trendValue && (
              <div className="mt-1.5 flex items-center gap-1">
                <span
                  className={cn(
                    "inline-flex items-center gap-0.5 rounded-full px-1.5 py-0.5 text-[10px] font-semibold",
                    trend === "up" &&
                      "bg-emerald-500/10 text-emerald-600 dark:text-emerald-400",
                    trend === "down" &&
                      "bg-red-500/10 text-red-600 dark:text-red-400",
                    trend === "neutral" &&
                      "bg-muted text-muted-foreground"
                  )}
                >
                  {trend === "up" && "↑"}
                  {trend === "down" && "↓"}
                  {trend === "neutral" && "→"}
                  {trendValue}
                </span>
              </div>
            )}

            {hint && (
              <p className="mt-1 text-xs text-muted-foreground/70 line-clamp-2">
                {hint}
              </p>
            )}
          </div>

          {/* Icon */}
          {icon && (
            <div
              className={cn(
                "flex size-10 shrink-0 items-center justify-center rounded-xl ring-1",
                a.iconBg,
                a.iconText,
                a.ring
              )}
            >
              {icon}
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
