FROM golang:1.16 as builder
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -o pt-search ./cmd/pt-search

FROM alpine:3.5
WORKDIR /app
RUN apk --no-cache add ca-certificates && update-ca-certificates
COPY --from=builder /src/pt-search ./

ENTRYPOINT ["/app/pt-search"]
