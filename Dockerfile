FROM golang:1.13-alpine AS builder

# Update the container with latest packages and get git
#   as it's a depency of go mods
RUN apk update && apk add git

# Build the folder structuer and working directory for the image where
#   we'll be copying over our files/code
WORKDIR /go/src/github.com/abelgoodwin1988/iata-finder
ADD . .
ENV GO111MODULE=on
RUN go build -mod=vendor -o ./iata-finder
RUN ls -a

# Build a leaner image with just the binary
FROM alpine:latest
# Creating folder structure which will container/run the binary. Match the build machine
#   so it's easier to debug when things go wrong
WORKDIR /go/src/github.com/abelgoodwin1988/iata-finder
# Get certs so we can perform network req's
RUN apk update && apk add ca-certificates
# Bring the built binary over
COPY --from=builder /go/src/github.com/abelgoodwin1988/iata-finder/iata-finder .
COPY --from=builder /go/src/github.com/abelgoodwin1988/iata-finder/configs/ ./configs/
# Create empty CSV's so we dont crash. We should make the code wait for the files or do some
#   kind of async channel running; don't serve the API until the dataCollector has run once
RUN mkdir ./assets/ && touch assets/airports.csv && touch assets/airlines.csv
ENTRYPOINT ./iata-finder

# Expose the standard gRPC port
EXPOSE 50051
