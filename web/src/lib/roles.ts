export type Role = "scheme_operator" | "issuer" | "acquirer" | "merchant" | "viewer";
export const ROLE_LABEL: Record<Role, string> = {
  scheme_operator: "Scheme Operator", issuer: "Issuer", acquirer: "Acquirer",
  merchant: "Merchant", viewer: "Viewer (HR)",
};
export function roleFromAppMetadata(meta?: Record<string, unknown>): Role | null {
  const r = meta?.role;
  return r === "scheme_operator" || r === "issuer" || r === "acquirer" || r === "merchant" || r === "viewer" ? r : null;
}
export const HOME_BY_ROLE: Record<Role, string> = {
  scheme_operator: "/ops", issuer: "/issuer", acquirer: "/acquirer",
  merchant: "/merchant", viewer: "/overview",
};
export const DASHBOARD_ACCESS: Record<Role, string[]> = {
  scheme_operator: ["/ops", "/transactions", "/clearing", "/settlement", "/ledger", "/cards", "/merchants", "/disputes"],
  issuer: ["/issuer", "/cards", "/tokens"],
  acquirer: ["/acquirer", "/merchants", "/funding", "/disputes"],
  merchant: ["/merchant", "/funding", "/disputes"],
  viewer: ["/overview"],
};
export function dashboardAccess(role: Role): string[] { return DASHBOARD_ACCESS[role] ?? []; }