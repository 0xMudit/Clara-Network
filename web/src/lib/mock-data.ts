// src/lib/mock-data.ts
// Deterministic, realistic-looking mock data for the admin API. Used only as a
// fallback when the live Go admin API is unreachable (network error or 5xx) —
// a real 401/403 still returns "access denied". Dates are generated relative
// to now so the "last 14 days" throughput window always looks live.
import type { DashboardSummary, Page, SeriesPoint } from "@/types/admin";

// ---- Deterministic PRNG (mulberry32) seeded per-process so a single page load
// is stable but the data isn't literally the same forever. ----------------------------------------------------------
function mulberry32(seed: number) {
  return function () {
    seed |= 0;
    seed = (seed + 0x6d2b79f5) | 0;
    let t = Math.imul(seed ^ (seed >>> 15), 1 | seed);
    t = (t + Math.imul(t ^ (t >>> 7), 61 | t)) ^ t;
    return ((t ^ (t >>> 14)) >>> 0) / 4294967296;
  };
}

// Seeded with a value derived from the current day so the data "progresses"
// naturally day to day but is stable across requests within the same day.
function rngFor(): () => number {
  const now = new Date();
  const seedKey =
    now.getUTCFullYear() * 10000 +
    (now.getUTCMonth() + 1) * 100 +
    now.getUTCDate();
  return mulberry32(seedKey);
}

function isoDate(d: Date): string {
  return d.toISOString().slice(0, 10);
}

const RESPONSE_CODES = ["00", "00", "00", "00", "51", "14", "55"];
const MTIS = ["0100", "0100", "0100", "0200", "0200"];
const DESTINATIONS = ["issuer-a", "issuer-b", "issuer-c"];
const MERCHANTS = [
  "Acme Retail",
  "Harbor Market",
  "Sana Cafe",
  "Northstar Travel",
  "Blue Peak Gas",
  "Lumen Bookstore",
  "Crestline Hotel",
  "Fern & Field Grocery",
];

function pick<T>(rng: () => number, arr: T[]): T {
  return arr[Math.floor(rng() * arr.length)];
}

function randInt(rng: () => number, min: number, max: number): number {
  return Math.floor(rng() * (max - min + 1)) + min;
}

function maskPan(rng: () => number): string {
  // Realistic-ish PAN prefixes for demo issuers, last four random.
  const issuer = pick(rng, ["400000", "411111", "520000", "550000", "620000"]);
  const lastFour = String(randInt(rng, 1000, 9999));
  return `${issuer}********${lastFour}`;
}

function amountFor(rng: () => number): string {
  const whole = randInt(rng, 5, 1200);
  const cents = String(randInt(rng, 0, 99)).padStart(2, "0");
  return `${whole}.${cents}`;
}

function stanFor(rng: () => number): string {
  return String(randInt(rng, 100000, 999999));
}

// ---- DashboardSummary ---------------------------------------------------------------------------------------------
const BASE_DAILY = 450;

function dashboard() {
  const rng = rngFor();
  const txDelta = randInt(rng, -80, 140);
  const transactions = Math.max(400, BASE_DAILY + txDelta);
  return {
    transactions,
    clearingRecords: Math.round(transactions * 0.72),
    merchants: randInt(rng, 150, 165),
    disputes: randInt(rng, 4, 9),
    cards: randInt(rng, 8200, 8600),
    tokens: randInt(rng, 3100, 3400),
  } satisfies DashboardSummary;
}

// ---- Series (14 days) ---------------------------------------------------------------------------------------------
interface SeriesEnvelope {
  items: SeriesPoint[];
  // extra metadata the ops page doesn't need but keeps the envelope plausible
  total: number;
  start: string;
  end: string;
}

function series(): SeriesEnvelope {
  const rng = rngFor();
  const now = new Date();
  const items: SeriesPoint[] = [];
  for (let i = 13; i >= 0; i--) {
    const d = new Date(now);
    d.setDate(now.getDate() - i);
    const base = 80 + (i % 3) * 25;
    items.push({
      date: isoDate(d),
      count: Math.max(40, base + randInt(rng, -30, 90)),
    });
  }
  return {
    items,
    total: items.reduce((s, p) => s + p.count, 0),
    start: items[0].date,
    end: items[items.length - 1].date,
  };
}

// ---- Transactions (recent authorizations) -------------------------------------------------------------------------
interface TransactionRow {
  stan: string;
  mti: string;
  pan: string;
  amount: string;
  responseCode: string;
  destination: string;
  createdAt: string;
  merchant: string;
  brand: string;
}

function transactions(limit: number): Page<TransactionRow> {
  const rng = rngFor();
  const items: TransactionRow[] = [];
  const now = Date.now();
  for (let i = 0; i < limit; i++) {
    const createdAt = new Date(now - i * randInt(rng, 3, 14) * 60_000);
    const code = pick(rng, RESPONSE_CODES);
    items.push({
      stan: stanFor(rng),
      mti: pick(rng, MTIS),
      pan: maskPan(rng),
      amount: code === "00" ? amountFor(rng) : amountFor(rng),
      responseCode: code,
      destination: pick(rng, DESTINATIONS),
      createdAt: createdAt.toISOString(),
      merchant: pick(rng, MERCHANTS),
      brand: pick(rng, ["visa", "mastercard", "visa", "mastercard", "amex"]),
    });
  }
  return { items, total: limit };
}

// ---- Public path dispatcher ---------------------------------------------------------------------------------------
const MOCK_PATHS = new Map<string, () => unknown>([
  ["/dashboard", dashboard],
  ["/dashboard/series", series],
  ["/transactions", () => transactions(6)],
]);

/**
 * Return a realistic mock response for an admin API path, or null when no mock
 * exists for that path (the caller should then surface a normal error).
 *
 * Handles query strings: "/dashboard/series?days=14" and "/transactions?limit=6"
 * both match. Numbers embedded in the query are honored where they map to a
 * limit (transactions) and ignored elsewhere.
 */
export function getMockForPath(path: string): unknown | null {
  const base = path.split("?")[0];
  const url = new URL(path, "https://clara.local");

  const limitRaw = Number(url.searchParams.get("limit"));

  if (base === "/dashboard/series") {
    // series() always returns the trailing 14-day window; ?days= is honored
    // only when a caller passes it as a limit-style int.
    return series();
  }
  if (base === "/transactions") {
    const limit = Number.isInteger(limitRaw) && limitRaw > 0 ? limitRaw : 6;
    return transactions(limit);
  }
  const hit = MOCK_PATHS.get(base);
  return hit ? hit() : null;
}
