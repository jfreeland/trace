# Build our binary in the api_builder
FROM golang:1.15.8-alpine as api_builder
WORKDIR /go/src/github.com/jfreeland/trace
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/trace-server cmd/traceserver/main.go

FROM node:current-alpine3.13 as ui_builder
WORKDIR /src/ui
COPY ./ui .
RUN npm install && \
    npm run build

FROM scratch
# SSL certs
COPY --from=api_builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy our static executable
COPY --from=api_builder /go/bin/trace-server /trace-server
COPY --from=ui_builder /src/ui/public /ui

# Run the binary
CMD ["/trace-server"]
