import { describe, it, expect, vi, beforeEach } from "vitest";

vi.mock("@/lib/supabase/client", () => ({
  createBrowserClient: vi.fn(),
}));

import { createBrowserClient } from "@/lib/supabase/client";
import { signInAs, DEMO_PASSWORD, DEMO_EMAILS, ROLE_CARDS } from "./demo-auth";
import type { Role } from "@/lib/roles";

const mockCreateClient = vi.mocked(createBrowserClient);

function router() {
  return { push: vi.fn(), refresh: vi.fn() };
}

describe("demo-auth", () => {
  beforeEach(() => {
    mockCreateClient.mockReset();
  });

  it("exposes a single demo password and one email per role", () => {
    expect(DEMO_PASSWORD).toBe("ClaraDemo!2026");
    expect(Object.keys(DEMO_EMAILS).sort()).toEqual([
      "acquirer", "issuer", "merchant", "scheme_operator", "viewer",
    ]);
  });

  it("signs in with the demo credentials for the chosen role", async () => {
    const signInWithPassword = vi.fn(async () => ({
      data: { session: { user: { app_metadata: { role: "issuer" } } } },
      error: null,
    }));
    mockCreateClient.mockReturnValue({
      auth: { signInWithPassword },
    } as unknown as ReturnType<typeof createBrowserClient>);

    const r = router();
    const result = await signInAs("issuer", r as never);
    expect(result.ok).toBe(true);
    // signed in with issuer's email + the shared demo password
    expect(signInWithPassword).toHaveBeenCalledWith({
      email: DEMO_EMAILS.issuer,
      password: DEMO_PASSWORD,
    });
    expect(r.push).toHaveBeenCalledWith("/issuer");
    expect(r.refresh).toHaveBeenCalled();
  });

  it("returns an error and does NOT navigate when sign-in fails", async () => {
    const signInWithPassword = vi.fn(async () => ({
      data: { session: null },
      error: { message: "Invalid login credentials" },
    }));
    mockCreateClient.mockReturnValue({
      auth: { signInWithPassword },
    } as unknown as ReturnType<typeof createBrowserClient>);

    const r = router();
    const result = await signInAs("merchant", r as never);
    expect(result.ok).toBe(false);
    if (!result.ok) expect(result.error).toContain("Invalid login");
    expect(r.push).not.toHaveBeenCalled();
  });

  it("honors a safe ?next= deep link over the role default", async () => {
    const signInWithPassword = vi.fn(async () => ({
      data: { session: { user: { app_metadata: { role: "scheme_operator" } } } },
      error: null,
    }));
    mockCreateClient.mockReturnValue({
      auth: { signInWithPassword },
    } as unknown as ReturnType<typeof createBrowserClient>);

    const r = router();
    await signInAs("scheme_operator", r as never, "/transactions");
    expect(r.push).toHaveBeenCalledWith("/transactions");
  });

  it("ignores an unsafe next target (protocol-relative/external)", async () => {
    const signInWithPassword = vi.fn(async () => ({
      data: { session: { user: { app_metadata: { role: "viewer" } } } },
      error: null,
    }));
    mockCreateClient.mockReturnValue({
      auth: { signInWithPassword },
    } as unknown as ReturnType<typeof createBrowserClient>);

    const r = router();
    await signInAs("viewer", r as never, "https://evil.example");
    expect(r.push).toHaveBeenCalledWith("/overview");
  });

  it("defines exactly the five demo roles with live demo stats", () => {
    expect(ROLE_CARDS).toHaveLength(5);
    const roles: Role[] = ROLE_CARDS.map((c) => c.role);
    expect(roles).toContain("scheme_operator");
    expect(roles).toContain("issuer");
    expect(roles).toContain("acquirer");
    expect(roles).toContain("merchant");
    expect(roles).toContain("viewer");
    for (const c of ROLE_CARDS) {
      expect(c.email).toBe(DEMO_EMAILS[c.role]);
      expect(c.stat.value.length).toBeGreaterThan(0);
    }
  });
});