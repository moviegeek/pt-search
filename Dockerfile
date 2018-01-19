FROM alpine:3.5

LABEL maintainer="justlaputa@gmail.com"

RUN apk --no-cache add ca-certificates && update-ca-certificates

COPY ./pt-search /app/pt-search
COPY ./build /app/build

WORKDIR /app

ENV GOCOOKIES /app/cookies/.cookies
VOLUME /app/cookies

ENTRYPOINT ["/app/pt-search"]
