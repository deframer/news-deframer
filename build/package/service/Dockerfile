FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

# We disable CGO for a static binary (easier for Alpine/Scratch)
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /out/ ./cmd/service ./cmd/admin

FROM gcr.io/distroless/static-debian13

# This ensures GitHub Actions has permission to write the package.
ARG REPO_URL=""
LABEL org.opencontainers.image.source=${REPO_URL}

WORKDIR /

COPY --from=builder /out/service /usr/bin/service
COPY --from=builder /out/admin /usr/bin/admin

USER nonroot:nonroot

ENTRYPOINT ["service"]
