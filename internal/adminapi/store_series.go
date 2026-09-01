package adminapi

import (
	"context"
	"time"
)

// SeriesPoint is a single bucket in a time series (e.g. transactions per day).
type SeriesPoint struct {
	Date  string `json:"date"`  // ISO date (YYYY-MM-DD) in the series' local timezone
	Count int64  `json:"count"`
}

// TransactionSeries returns the number of switch_transactions per calendar
// day for the last `days` days (including today), oldest first.
func (s *Store) TransactionSeries(ctx context.Context, days int) ([]SeriesPoint, error) {
	if days <= 0 || days > 366 {
		return []SeriesPoint{}, nil
	}

	start := time.Now().UTC().AddDate(0, 0, -(days - 1))
	rows, err := s.Pool.Query(ctx, `
		SELECT date_trunc('day', created_at)::date AS day, count(*) AS n
		FROM switch_transactions
		WHERE created_at >= $1
		GROUP BY day
		ORDER BY day ASC`, start)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	byDay := make(map[time.Time]int64, days)
	for rows.Next() {
		var day time.Time
		var n int64
		if err := rows.Scan(&day, &n); err != nil {
			return nil, err
		}
		byDay[day.UTC().Truncate(24*time.Hour)] = n
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Fill any gaps so the chart always renders a continuous line.
	out := make([]SeriesPoint, 0, days)
	for i := 0; i < days; i++ {
		day := start.AddDate(0, 0, i).UTC().Truncate(24 * time.Hour)
		out = append(out, SeriesPoint{
			Date:  day.Format("2006-01-02"),
			Count: byDay[day],
		})
	}
	return out, nil
}