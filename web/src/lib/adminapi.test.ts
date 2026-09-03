import { describe, it, expect, beforeAll, beforeEach, afterAll, vi } from "vitest";

// The real `server-only` package throws outside a server component; neutralize it.
vi.mock("server-only", () => ({}));
vi.mock("next/headers", () => ({
  cookies: vi.fn(async () => ({
    toString: (): string => "",
  })),
}));
vi.mock("./env", () => ({
  getAppUrl: () => "https://clara-demo.app",
}));
vi.mock("./mock-data", async (importOriginal) => {
  const actual = await importOriginal<typeof import("./mock-data")>();
  return { ...actual, getMockForPath: vi.fn(actual.getMockForPath) };
});

import { tryFetchAdmin, AdminError } from "./adminapi";
import { getMockForPath } from "./mock-data";
import type { DashboardSummary } from "@/types/admin";

const mockGet = vi.mocked(getMockForPath);

describe("tryFetchAdmin", () => {
  const originalFetch = globalThis.fetch;
  let warnSpy: ReturnType<typeof vi.spyOn>;

  beforeAll(() => {
    // Silence the fallback warning in test output.
    warnSpy = vi.spyOn(console, "warn").mockImplementation(() => {});
  });
  afterAll(() => {
    globalThis.fetch = originalFetch;
    warnSpy.mockRestore();
    vi.useRealTimers();
  });
  beforeEach(() => {
    mockGet.mockClear();
  });

  function respond(status: number, body?: unknown) {
    globalThis.fetch = vi.fn().mockResolvedValue(
      new Response(body !== undefined ? JSON.stringify(body) : "", {
        status,
        headers: { "content-type": "application/json" },
      }),
    ) as unknown as typeof fetch;
  }

  it("returns data when the upstream succeeds (200)", async () => {
    respond(200, { transactions: 500, clearingRecords: 360, merchants: 160, disputes: 5, cards: 8300, tokens: 3200 });
    const got = await tryFetchAdmin<DashboardSummary>("/dashboard");
    expect(got.ok).toBe(true);
    if (got.ok) expect(got.data.transactions).toBe(500);
    expect(mockGet).not.toHaveBeenCalled();
  });

  it("returns access-denied and NO mock on 401", async () => {
    respond(401, { error: "Unauthorized" });
    const got = await tryFetchAdmin<DashboardSummary>("/dashboard");
    expect(got.ok).toBe(false);
    if (!got.ok) expect(got.error).toContain("access");
    expect(mockGet).not.toHaveBeenCalled();
  });

  it("returns access-denied and NO mock on 403", async () => {
    respond(403, { error: "Forbidden" });
    const got = await tryFetchAdmin<DashboardSummary>("/dashboard");
    expect(got.ok).toBe(false);
    if (!got.ok) expect(got.error).toContain("access");
    expect(mockGet).not.toHaveBeenCalled();
  });

  it("falls back to mock data on 500", async () => {
    respond(500, { error: "upstream failure" });
    const got = await tryFetchAdmin<DashboardSummary>("/dashboard");
    expect(got.ok).toBe(true);
    if (got.ok) expect(isFinite(got.data.transactions)).toBe(true);
    expect(mockGet).toHaveBeenCalledWith("/dashboard");
  });

  it("falls back to mock data on 502", async () => {
    respond(502, { error: "bad gateway" });
    const got = await tryFetchAdmin<DashboardSummary>("/dashboard");
    expect(got.ok).toBe(true);
    if (got.ok) expect(isFinite(got.data.clearingRecords)).toBe(true);
  });

  it("falls back to mock data on a network error (fetch rejects)", async () => {
    globalThis.fetch = vi.fn().mockRejectedValue(new TypeError("fetch failed")) as unknown as typeof fetch;
    const got = await tryFetchAdmin<DashboardSummary>("/dashboard");
    expect(got.ok).toBe(true);
    if (got.ok) expect(isFinite(got.data.transactions)).toBe(true);
  });

  it("returns a generic error when the upstream fails on an UNMOCKED path", async () => {
    respond(500, { error: "boom" });
    const got = await tryFetchAdmin("/unmocked/endpoint");
    expect(got.ok).toBe(false);
    if (!got.ok) expect(got.error).toMatch(/unavailable|temporarily/i);
  });

  it("returns a generic error on a network error for an unmocked path", async () => {
    globalThis.fetch = vi.fn().mockRejectedValue(new TypeError("fetch failed")) as unknown as typeof fetch;
    const got = await tryFetchAdmin("/unmocked/endpoint");
    expect(got.ok).toBe(false);
    if (!got.ok) expect(got.error).toMatch(/something went wrong/i);
  });

  it("AdminError carries the upstream status", () => {
    expect(AdminError).toBeInstanceOf(Function);
  });
});