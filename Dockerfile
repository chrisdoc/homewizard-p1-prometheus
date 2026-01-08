FROM golang:1.23 as builder
WORKDIR /go/src/github.com/chrisdoc/homewizard-p1-prometheus/

# Copy go mod files first to cache dependencies layer
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy source code
COPY . .

# Build with cache mount for build cache
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -o /go/src/github.com/chrisdoc/homewizard-p1-prometheus/app

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /go/src/github.com/chrisdoc/homewizard-p1-prometheus/app .
CMD ["./app"]
