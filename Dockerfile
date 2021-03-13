FROM alpine:3.5

LABEL maintainer="justlaputa@gmail.com"

WORKDIR /app

RUN apk --no-cache add ca-certificates && update-ca-certificates

COPY ./pt-search .

ENTRYPOINT ["/app/pt-search"]
