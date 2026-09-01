"use client";

import { useEffect, useRef } from "react";
import { useRouter } from "next/navigation";

const POLL_INTERVAL_MS = 30_000;

/**
 * Invisible client component that triggers a full server-side re-render
 * of the current route every 30 seconds. This causes all server components
 * to re-fetch their data (dashboard stats, tables, etc.) and the UI to
 * update seamlessly — no manual state management required.
 *
 * The refresh is skipped when the tab is hidden (Page Visibility API)
 * to avoid wasteful fetches.
 */
export function RefreshProvider() {
  const router = useRouter();
  const timerRef = useRef<ReturnType<typeof setInterval> | null>(null);

  useEffect(() => {
    function startPolling() {
      timerRef.current = setInterval(() => {
        // Don't refresh if the tab is hidden
        if (document.hidden) return;
        router.refresh();
      }, POLL_INTERVAL_MS);
    }

    function handleVisibility() {
      if (document.hidden) {
        // Pause polling when tab is hidden
        if (timerRef.current) clearInterval(timerRef.current);
        timerRef.current = null;
      } else {
        // Resume polling + immediate refresh when tab becomes visible
        router.refresh();
        startPolling();
      }
    }

    startPolling();
    document.addEventListener("visibilitychange", handleVisibility);

    return () => {
      if (timerRef.current) clearInterval(timerRef.current);
      document.removeEventListener("visibilitychange", handleVisibility);
    };
  }, [router]);

  return null;
}
