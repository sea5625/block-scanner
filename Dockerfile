############################
# STEP 1 build executable binary
############################
FROM golang:1.12-alpine  AS builder

# Install dependencies and certification.
RUN apk add bash ca-certificates make git curl gcc g++ libc-dev yarn nodejs

# Add Maintainer Info
LABEL maintainer="ICONLOOP, Inc"

# Set the Current Working Directory inside the container
ADD  .  /go/src/motherbear
WORKDIR /go/src/motherbear

# Build
RUN  make build-linux

############################
# STEP 2 build a small image
############################
FROM alpine:latest AS product

# Use unicode
RUN locale-gen C.UTF-8 || true
ENV LANG=C.UTF-8

# We add the certificates to connectnode via  HTTPS.
RUN apk add ca-certificates

#this seems dumb, but the libc from the build stage is not the same as the alpine libc
#create a symlink to where it expects it since they are compatable. https://stackoverflow.com/a/35613430/3105368
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

# Copy our static executable.
WORKDIR  /
COPY --from=builder  /go/src/motherbear/isaac  .
RUN mkdir frontend
RUN mkdir frontend/public
COPY --from=builder  /go/src/motherbear/frontend/build/  ./frontend/build/

# This container exposes port 6553 to the outside world
EXPOSE 6553

# Volume for  configuration file.
VOLUME /config
VOLUME /data
VOLUME /log

# Release mode for gin.
ENV GIN_MODE=release

# Run the executable
CMD  ./isaac
