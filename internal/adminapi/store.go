// Package adminapi implements the Clara Network admin REST API: a read-only
// aggregator over the shared PostgreSQL schema that serves dashboard data to
// the frontend.
package adminapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Store is a read-only data layer over the shared Clara Network PostgreSQL
// schema. It never writes; all mutations go through the owning services.
type Store struct {
	Pool *pgxpool.Pool
}

// DashboardSummary is an aggregated overview for the landing page.
type DashboardSummary struct {
	Transactions    int64 `json:"transactions"`
	ClearingRecords int64 `json:"clearingRecords"`
	Merchants       int64 `json:"merchants"`
	Disputes        int64 `json:"disputes"`
	Cards           int64 `json:"cards"`
	Tokens          int64 `json:"tokens"`
}

// DashboardSummary returns aggregate counts across every major table.
func (s *Store) DashboardSummary(ctx context.Context) (*DashboardSummary, error) {
	var d DashboardSummary
	type result struct {
		field *int64
		val   int64
		err   error
	}
	tables := []struct {
		field *int64
		name  string
	}{
		{&d.Transactions, "switch_transactions"},
		{&d.ClearingRecords, "clearing_records"},
		{&d.Merchants, "merchants"},
		{&d.Disputes, "disputes"},
		{&d.Cards, "cards"},
		{&d.Tokens, "tokens"},
	}
	ch := make(chan result, len(tables))
	for _, t := range tables {
		go func(field *int64, name string) {
			n, err := s.count(ctx, name)
			ch <- result{field, n, err}
		}(t.field, t.name)
	}
	for range tables {
		r := <-ch
		if r.err != nil {
			return nil, r.err
		}
		*r.field = r.val
	}
	return &d, nil
}

func (s *Store) count(ctx context.Context, table string) (int64, error) {
	var n int64
	err := s.Pool.QueryRow(ctx, fmt.Sprintf("SELECT count(*) FROM %s", table)).Scan(&n)
	return n, err
}

// Page is a paginated list envelope.
type Page[T any] struct {
	Items []T   `json:"items"`
	Total int64 `json:"total"`
}

func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, ",")
}
