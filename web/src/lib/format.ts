const intlCount = new Intl.NumberFormat("en-US");

/** Integer counts with thousands separators, e.g. 1234567 -> "1,234,567". */
export function fmtCount(n: number | null | undefined): string {
  if (n == null) return "—";
  return intlCount.format(Math.trunc(n));
}