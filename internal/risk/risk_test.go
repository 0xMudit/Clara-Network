package risk

import (
	"context"
	"testing"
	"time"
)

const testRules = `{"rules":[
	{"name":"card-velocity","kind":"velocity-card","limit":2,"window":60,"code":"59","enabled":true},
	{"name":"merchant-velocity","kind":"velocity-merchant","limit":10,"window":60,"code":"59","enabled":true},
	{"name":"large-amount","kind":"amount","limit":500000,"code":"58","enabled":true}
]}`

func mustEngine(t *testing.T) *Engine {
	t.Helper()
	e, err := FromConfig([]byte(testRules), NewMemoryStore())
	if err != nil {
		t.Fatalf("FromConfig: %v", err)
	}
	return e
}

func TestCardVelocity(t *testing.T) {
	e := mustEngine(t)
	ctx := context.Background()
	for i := 0; i < 2; i++ {
		d, err := e.Evaluate(ctx, "4000001234567890", "TST00001", 1000)
		if err != nil {
			t.Fatalf("evaluate: %v", err)
		}
		if !d.Allow {
			t.Fatalf("request %d should be allowed", i+1)
		}
	}
	d, err := e.Evaluate(ctx, "4000001234567890", "TST00001", 1000)
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if d.Allow {
		t.Fatal("third request should be declined")
	}
	if d.Code != "59" {
		t.Fatalf("code = %q want 59", d.Code)
	}
	// A different card is unaffected.
	d, err = e.Evaluate(ctx, "4000011234567890", "TST00001", 1000)
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if !d.Allow {
		t.Fatal("different card should be allowed")
	}
}

func TestAmountLimit(t *testing.T) {
	e := mustEngine(t)
	ctx := context.Background()
	if d, _ := e.Evaluate(ctx, "4000001234567890", "TST00001", 500001); d.Allow {
		t.Fatal("large amount should be declined")
	}
	if d, _ := e.Evaluate(ctx, "4000001234567890", "TST00001", 500000); !d.Allow {
		t.Fatal("amount at limit should be allowed")
	}
}

func TestWindowReset(t *testing.T) {
	e := New(NewMemoryStore(), []Rule{
		{Name: "card-velocity", Kind: KindCardVelocity, Limit: 1, Window: 1, Code: "59", Enabled: true},
	})
	ctx := context.Background()
	if d, _ := e.Evaluate(ctx, "4000001234567890", "T", 100); !d.Allow {
		t.Fatal("first request should be allowed")
	}
	if d, _ := e.Evaluate(ctx, "4000001234567890", "T", 100); d.Allow {
		t.Fatal("second request should be declined")
	}
	time.Sleep(1100 * time.Millisecond)
	if d, _ := e.Evaluate(ctx, "4000001234567890", "T", 100); !d.Allow {
		t.Fatal("request after window reset should be allowed")
	}
}

func TestDisabledRuleIgnored(t *testing.T) {
	e := New(NewMemoryStore(), []Rule{
		{Name: "card-velocity", Kind: KindCardVelocity, Limit: 0, Window: 60, Code: "59", Enabled: false},
	})
	if d, _ := e.Evaluate(context.Background(), "4000001234567890", "T", 100); !d.Allow {
		t.Fatal("disabled rule must not decline")
	}
}

func TestNilEngineAllows(t *testing.T) {
	var e *Engine
	if d, _ := e.Evaluate(context.Background(), "x", "y", 1); !d.Allow {
		t.Fatal("nil engine must allow")
	}
}
