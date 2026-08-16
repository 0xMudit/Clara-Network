GO       ?= go
COMPOSE  ?= docker compose

.PHONY: build test vet run-switch run-issuer run-acquirer compose-up compose-down compose-logs clean

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

compose-up:
	$(COMPOSE) -f deploy/docker-compose.yml up --build -d

compose-down:
	$(COMPOSE) -f deploy/docker-compose.yml down

compose-logs:
	$(COMPOSE) -f deploy/docker-compose.yml logs -f

clean:
	$(GO) clean ./...
