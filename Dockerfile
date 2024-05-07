FROM golang:1.20 as builder
WORKDIR /go/src/github.com/chrisdoc/homewizard-p1-prometheus/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /go/src/github.com/chrisdoc/homewizard-p1-prometheus/app .
CMD ["./app"]
