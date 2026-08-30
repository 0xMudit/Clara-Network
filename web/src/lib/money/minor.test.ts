// src/lib/money/minor.test.ts
import assert from "node:assert/strict";
import { test } from "node:test";
import { fmtMinor } from "./minor";

test("fmtMinor formats minor units", () => {
  assert.equal(fmtMinor(123456), "1,234.56 EUR");
  assert.equal(fmtMinor(0), "0.00 EUR");
  assert.equal(fmtMinor(-5000), "-50.00 EUR");
});