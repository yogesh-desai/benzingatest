# syntax=docker/dockerfile:1

#
# Build
#
FROM golang:1.16-alpine AS build

# Add necessary things
RUN apk add git
RUN apk add xz
RUN apk add binutils

# Create app directory and copy contents to it
RUN mkdir /app
COPY . /app
# Set app as working directory
WORKDIR /app

# Install UPX - This reduces the size of binary/artifacts
ADD https://github.com/upx/upx/releases/download/v3.96/upx-3.96-amd64_linux.tar.xz /usr/local
RUN xz -d -c /usr/local/upx-3.96-amd64_linux.tar.xz | \
    tar -xOf - upx-3.96-amd64_linux/upx > /bin/upx && \
    chmod a+x /bin/upx

# go mod download
RUN go mod download

# Build the go app and save binary with name mta
RUN CGO_ENABLED=0 go build -o mta .

# strip and compress the binary
RUN strip --strip-unneeded mta
RUN upx mta

#########################################################
#
# Deploy
#

# Taking bare minimum image.
FROM scratch

# Set working directory
WORKDIR /

# Copy necessary files
COPY --from=build /app/mta /mta
COPY --from=build /app/files/* ./files/
COPY --from=build /app/files/config/ ./files/config/

# Expose port
EXPOSE 9000

# Set entrypoint
CMD [ "./mta" ]
