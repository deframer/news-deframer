FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

# We disable CGO for a static binary (easier for Alpine/Scratch)
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/ ./cmd/worker ./cmd/admin

FROM gcr.io/distroless/static-debian13

# This ensures GitHub Actions has permission to write the package.
ARG REPO_URL=""
LABEL org.opencontainers.image.source=${REPO_URL}

WORKDIR /

COPY --from=builder /out/worker /worker
COPY --from=builder /out/admin /usr/bin/admin

USER nonroot:nonroot

ENTRYPOINT ["/worker"]
