FROM golang:1.15.6
WORKDIR /go/src/github.com/chrisdoc/homewizard-p1-prometheus/
COPY go.mod .
COPY go.sum .
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/chrisdoc/homewizard-p1-prometheus/app .
CMD ["./app"]
