# syntax=docker/dockerfile:1
FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/switch ./cmd/switch \
 && CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/issuer-sim ./cmd/issuer-sim \
 && CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/acquirer-sim ./cmd/acquirer-sim \
 && CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/clearing-sim ./cmd/clearing-sim \
 && CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/ledger-sim ./cmd/ledger-sim

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
COPY --from=build /out/switch /usr/local/bin/switch
COPY --from=build /out/issuer-sim /usr/local/bin/issuer-sim
COPY --from=build /out/acquirer-sim /usr/local/bin/acquirer-sim
COPY --from=build /out/clearing-sim /usr/local/bin/clearing-sim
COPY --from=build /out/ledger-sim /usr/local/bin/ledger-sim
ENTRYPOINT []
