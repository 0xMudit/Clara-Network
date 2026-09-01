"use client";

import { useMemo } from "react";
import { cn } from "@/lib/utils";

interface SparklineProps {
  /** Array of numeric data points to plot */
  data: number[];
  /** Width of the SVG viewport */
  width?: number;
  /** Height of the SVG viewport */
  height?: number;
  /** Color accent — maps to the same palette as StatCard */
  accent?: "default" | "success" | "warning" | "danger" | "info";
  /** Whether to show the gradient fill under the line */
  fill?: boolean;
  /** Extra class on the wrapper */
  className?: string;
}

const STROKE_MAP = {
  default: "var(--color-primary)",
  success: "#10b981",
  warning: "#f59e0b",
  danger: "#ef4444",
  info: "#0ea5e9",
};

const FILL_ID_MAP = {
  default: "spark-grad-primary",
  success: "spark-grad-success",
  warning: "spark-grad-warning",
  danger: "spark-grad-danger",
  info: "spark-grad-info",
};

const GRADIENT_STOPS = {
  default: ["var(--color-primary)", "transparent"],
  success: ["#10b981", "transparent"],
  warning: ["#f59e0b", "transparent"],
  danger: ["#ef4444", "transparent"],
  info: ["#0ea5e9", "transparent"],
};

/**
 * Generates a smooth SVG path through the data points using
 * cubic bezier curves (Catmull-Rom → Bezier conversion).
 */
function buildSmoothPath(
  points: { x: number; y: number }[],
  closed = false
): string {
  if (points.length < 2) return "";

  const segments: string[] = [`M ${points[0].x},${points[0].y}`];

  for (let i = 0; i < points.length - 1; i++) {
    const p0 = points[Math.max(0, i - 1)];
    const p1 = points[i];
    const p2 = points[i + 1];
    const p3 = points[Math.min(points.length - 1, i + 2)];

    // Catmull-Rom to cubic bezier control points
    const tension = 0.3;
    const cp1x = p1.x + ((p2.x - p0.x) * tension);
    const cp1y = p1.y + ((p2.y - p0.y) * tension);
    const cp2x = p2.x - ((p3.x - p1.x) * tension);
    const cp2y = p2.y - ((p3.y - p1.y) * tension);

    segments.push(
      `C ${cp1x},${cp1y} ${cp2x},${cp2y} ${p2.x},${p2.y}`
    );
  }

  if (closed) {
    segments.push("Z");
  }

  return segments.join(" ");
}

export function Sparkline({
  data,
  width = 120,
  height = 40,
  accent = "default",
  fill = true,
  className,
}: SparklineProps) {
  const { strokePath, fillPath, gradientId } = useMemo(() => {
    if (data.length < 2) {
      return { strokePath: "", fillPath: "", gradientId: FILL_ID_MAP[accent] };
    }

    const padY = 4;
    const padX = 2;
    const w = width - padX * 2;
    const h = height - padY * 2;

    const min = Math.min(...data);
    const max = Math.max(...data);
    const range = max - min || 1;

    const points = data.map((v, i) => ({
      x: padX + (i / (data.length - 1)) * w,
      y: padY + h - ((v - min) / range) * h,
    }));

    const strokePath = buildSmoothPath(points);
    const gId = FILL_ID_MAP[accent];

    // Build a closed fill path: line along top, then straight down bottom edges
    const fillPath =
      strokePath +
      ` L ${points[points.length - 1].x},${height} L ${points[0].x},${height} Z`;

    return { strokePath, fillPath, gradientId: gId };
  }, [data, width, height, accent]);

  const strokeColor = STROKE_MAP[accent];
  const [gradFrom, gradTo] = GRADIENT_STOPS[accent];

  if (data.length < 2) return null;

  return (
    <svg
      viewBox={`0 0 ${width} ${height}`}
      className={cn("block overflow-visible", className)}
      aria-hidden="true"
    >
      <defs>
        <linearGradient id={gradientId} x1="0" y1="0" x2="0" y2="1">
          <stop offset="0%" stopColor={gradFrom} stopOpacity="0.2" />
          <stop offset="100%" stopColor={gradTo} stopOpacity="0" />
        </linearGradient>
      </defs>

      {/* Gradient fill under the line */}
      {fill && (
        <path
          d={fillPath}
          fill={`url(#${gradientId})`}
          className="spark-fill"
        />
      )}

      {/* The line itself */}
      <path
        d={strokePath}
        fill="none"
        stroke={strokeColor}
        strokeWidth="1.5"
        strokeLinecap="round"
        strokeLinejoin="round"
        className="spark-line"
      />
    </svg>
  );
}
