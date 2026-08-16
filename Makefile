GO       ?= go
COMPOSE  ?= docker compose

.PHONY: build test vet run-switch run-issuer run-acquirer run-clearing run-ledger run-cardsvc run-card-sim run-acquiring run-disputes run-hsm run-resilience compose-up compose-down compose-logs clean

build:
	$(GO) build ./...

test:
	$(GO) test ./...

vet:
	$(GO) vet ./...

run-switch:
	$(GO) run ./cmd/switch

run-issuer:
	$(GO) run ./cmd/issuer-sim

run-acquirer:
	$(GO) run ./cmd/acquirer-sim

run-clearing:
	$(GO) run ./cmd/clearing-sim

run-ledger:
	$(GO) run ./cmd/ledger-sim

run-cardsvc:
	$(GO) run ./cmd/cardsvc

run-card-sim:
	$(GO) run ./cmd/card-sim

run-acquiring:
	$(GO) run ./cmd/acquiring-sim

run-disputes:
	$(GO) run ./cmd/disputes-sim

run-hsm:
	$(GO) run ./cmd/hsm-sim

run-resilience:
	$(GO) run ./cmd/resilience-sim

compose-up:
	$(COMPOSE) -f deploy/docker-compose.yml up --build -d

compose-down:
	$(COMPOSE) -f deploy/docker-compose.yml down

compose-logs:
	$(COMPOSE) -f deploy/docker-compose.yml logs -f

clean:
	$(GO) clean ./...
