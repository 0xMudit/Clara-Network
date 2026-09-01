/**
 * Generate a plausible sparkline array from a current value.
 *
 * Produces 8–12 points that visually trend toward the current value,
 * with slight randomness to feel organic. Used for stat cards where
 * real time-series data isn't available from the API yet.
 */
export function sparkFromValue(
  current: number,
  points = 10,
  volatility = 0.15
): number[] {
  if (current <= 0) return Array(points).fill(0);

  const result: number[] = [];
  // Start somewhere between 40-80% of the current value
  let v = current * (0.4 + Math.random() * 0.4);
  const step = (current - v) / (points - 1);

  for (let i = 0; i < points; i++) {
    // Drift toward current value with noise
    const noise = (Math.random() - 0.5) * current * volatility;
    v += step + noise;
    v = Math.max(0, v);
    result.push(Math.round(v));
  }

  // Ensure the last point is close to the actual value
  result[points - 1] = current;
  return result;
}
