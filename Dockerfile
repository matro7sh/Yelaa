##
## Step 1 - Get dependencies
##
FROM golang:1.16-alpine as otp-builder

WORKDIR /build

RUN apk update --no-cache && \
    apk upgrade --no-cache && \
    apk add --no-cache \
    make \
    build-base

COPY . .

RUN go mod download && \
    make compile

##
## Step 2 - Build lean container
##
FROM alpine:latest

WORKDIR /app

COPY --from=otp-builder /build/yelaa.txt .
COPY --from=otp-builder /build/Yelaa .

# Enabling run permissions
RUN adduser -D yelaa_user && \
    chown -R yelaa_user: /app/Yelaa
USER yelaa_user

# Example command:
# docker run -v $PWD:/mnt/ yelaa-local scan -t /mnt/target.txt
ENTRYPOINT [ "/app/Yelaa" ]
