FROM golang:1.17 as cli

# Setup src directory
RUN mkdir /src
COPY cli /src
WORKDIR /src

# Build boots
RUN go build ./cmd/boots

FROM ubuntu:latest

# Update
RUN apt-get update && apt-get upgrade

# Copy CLI tools
COPY --from=cli /src/boots /usr/bin/boots

# Setup work dir
RUN mkdir /lab
WORKDIR /lab