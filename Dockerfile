##
## Step 1 - Get dependencies
##
FROM golang:1.17.6-alpine as builder

WORKDIR /build

RUN apk update --no-cache && \
    apk upgrade --no-cache && \
    apk add --no-cache \
    make \
    build-base

COPY . .

RUN make

##
## Step 2 - Build lean container
##
FROM golang:1.17.6-alpine

WORKDIR /app

# Installing runtime dependencies
RUN apk update --no-cache && \
    apk upgrade --no-cache && \
    apk add --no-cache bind-tools ca-certificates chromium

COPY --from=builder /build/yelaa.txt .
COPY --from=builder /build/Yelaa .

RUN adduser -D yelaa_user && chown -R yelaa_user: /app/Yelaa
USER yelaa_user

ENTRYPOINT [ "/app/Yelaa" ]
