FROM golang:1.23 as builder
WORKDIR /go/src/github.com/chrisdoc/homewizard-p1-prometheus/

COPY . .

RUN go mod download
RUN go vet -v
RUN go test -v

RUN CGO_ENABLED=0 go build -o /go/src/github.com/chrisdoc/homewizard-p1-prometheus/app

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /go/src/github.com/chrisdoc/homewizard-p1-prometheus/app .
CMD ["./app"]
