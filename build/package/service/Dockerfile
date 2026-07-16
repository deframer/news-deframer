FROM --platform=$BUILDPLATFORM golang:1.26.5-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

RUN go install goa.design/goa/v3/cmd/goa@latest

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN goa gen github.com/deframer/news-deframer/pkg/design

# We disable CGO for a static binary (easier for Alpine/Scratch)
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /out/ ./cmd/service ./cmd/admin ./cmd/migration

FROM gcr.io/distroless/static-debian13

# This ensures GitHub Actions has permission to write the package.
ARG REPO_URL=""
LABEL org.opencontainers.image.source=${REPO_URL}

WORKDIR /

COPY --from=builder /out/service /usr/bin/service
COPY --from=builder /out/admin /usr/bin/admin
COPY --from=builder /out/migration /usr/bin/migration
COPY --from=builder /app/gen/http/openapi3.json /gen/http/openapi3.json

USER nonroot:nonroot

ENTRYPOINT ["service"]
