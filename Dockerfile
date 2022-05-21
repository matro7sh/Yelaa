FROM golang:1.18-alpine as builder

WORKDIR /root
RUN apk update --no-cache && \
    apk upgrade --no-cache && \
    apk add --no-cache \
    make \
    build-base

COPY go.mod .
RUN go mod download

COPY . .
RUN make


FROM golang:1.18-alpine

ARG USER_ID
ARG GROUP_ID

WORKDIR /app

# Installing runtime dependencies
RUN apk update --no-cache && \
    apk upgrade --no-cache && \
    apk add --no-cache bind-tools ca-certificates chromium

COPY --from=builder /root/yelaa.txt .
COPY --from=builder /root/Yelaa .

RUN addgroup --gid $GROUP_ID -S yelaa_user && \
    adduser --uid $USER_ID -S -G yelaa_user yelaa_user && \
    chown -R yelaa_user: /app
USER yelaa_user

ENTRYPOINT [ "/app/Yelaa" ]
