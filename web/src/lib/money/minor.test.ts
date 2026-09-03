// src/lib/money/minor.test.ts
import { describe, it, expect } from "vitest";
import { fmtMinor } from "./minor";

describe("fmtMinor", () => {
  it("formats minor units with locale grouping", () => {
    expect(fmtMinor(123456)).toBe("1,234.56 EUR");
    expect(fmtMinor(0)).toBe("0.00 EUR");
    expect(fmtMinor(-5000)).toBe("-50.00 EUR");
  });
});