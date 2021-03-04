FROM golang:1.15.8-alpine as builder

WORKDIR /go/src/github.com/jfreeland/trace
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/trace-server cmd/traceserver/main.go

FROM golang:1.15.8-alpine
# SSL certs
#COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy our static executable
COPY --from=builder /go/bin/trace-server /usr/local/bin/trace-server

# Run the binary
CMD ["/usr/local/bin/trace-server"]
