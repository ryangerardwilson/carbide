FROM golang:1.26-bookworm AS build

WORKDIR /workspace
COPY db/go.mod db/go.sum ./db/
COPY api/go.mod api/go.sum ./api/

WORKDIR /workspace/api
RUN go mod download

WORKDIR /workspace
COPY db ./db
COPY api ./api

WORKDIR /workspace/api
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags "-s -w" -o /out/carbide-app .

FROM debian:bookworm-slim

RUN apt-get update \
  && apt-get install -y --no-install-recommends ca-certificates curl \
  && rm -rf /var/lib/apt/lists/*

COPY --from=build /out/carbide-app /usr/local/bin/carbide-app

EXPOSE 8080

HEALTHCHECK --interval=3s --timeout=3s --retries=20 \
  CMD curl -fsS http://localhost:8080/health >/dev/null || exit 1

CMD ["carbide-app"]
