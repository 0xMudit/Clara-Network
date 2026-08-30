export function fmtMinor(minor: number, currency = "EUR"): string {
  const sign = minor < 0 ? "-" : "";
  const abs = Math.abs(minor);
  const major = Math.floor(abs / 100);
  const frac = String(abs % 100).padStart(2, "0");
  const intl = new Intl.NumberFormat("en-US").format(major);
  return `${sign}${intl}.${frac} ${currency}`;
}
