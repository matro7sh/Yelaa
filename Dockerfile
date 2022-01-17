FROM golang:1 as build

# Install nuclei
RUN go install -v github.com/projectdiscovery/nuclei/v2/cmd/nuclei@latest

WORKDIR /app

COPY . ./

# Build Yelaa
RUN go build -o /yelaa

# Use image that contains google-chrome binary
FROM chromedp/headless-shell:latest

# Install dumb-init
# dumb-init is used to prevent zombie process when running google-chrome
RUN export DEBIAN_FRONTEND=noninteractive \
  && apt-get update \
  && apt-get install -y --no-install-recommends \
    dumb-init \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*

# Copy the Yelaa binary from the previous image "build" into /usr/local/bin
COPY --from=build /yelaa /usr/local/bin
COPY --from=build /go/bin/nuclei /usr/local/bin

WORKDIR /data
VOLUME ["/data"]

ENTRYPOINT ["dumb-init", "--"]
CMD ["yelaa", "-h"]