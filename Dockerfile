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

ARG USER_ID
ARG GROUP_ID

WORKDIR /app

# Installing runtime dependencies
RUN apk update --no-cache && \
    apk upgrade --no-cache && \
    apk add --no-cache bind-tools ca-certificates chromium

COPY --from=builder /build/yelaa.txt .
COPY --from=builder /build/Yelaa .

RUN echo "gid: ${GROUP_ID} -- uid: ${USER_ID}"
RUN addgroup --gid $GROUP_ID -S yelaa_user && \
    adduser -S --uid $USER_ID -G yelaa_user yelaa_user && \
    chown -R yelaa_user: /app/Yelaa
USER yelaa_user

ENTRYPOINT [ "/app/Yelaa" ]
