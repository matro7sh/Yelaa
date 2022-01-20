##
## Step 1 - Get dependencies
##
FROM golang:1.16-alpine as builder

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
FROM golang:1.17.6-alpine

WORKDIR /app

# Installing runtime dependencies
RUN apk update --no-cache && \
    apk upgrade --no-cache && \
    apk add --no-cache bind-tools ca-certificates chromium && \
    go install -v github.com/projectdiscovery/nuclei/v2/cmd/nuclei@latest

COPY --from=builder /build/yelaa.txt .
COPY --from=builder /build/Yelaa .

RUN adduser -D yelaa_user && chown -R yelaa_user: /app/Yelaa
USER yelaa_user

# Example command:
# docker run -v $PWD:/mnt/ yelaa-local scan -t /mnt/target.txt
ENTRYPOINT [ "/app/Yelaa" ]
