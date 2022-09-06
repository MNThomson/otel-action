ARG GO_VERSION=1.19.0

###########################
# Build executable binary #
###########################
FROM golang:${GO_VERSION}-alpine AS builder

# Git is required for fetching the dependencies
RUN apk update && \
    apk add --no-cache git ca-certificates && \
    update-ca-certificates

# Create appuser
ENV USER=otel
ENV UID=10001
RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"

WORKDIR /src

# Fetch dependencies
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY *.go Makefile ./

# Build the binary
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -ldflags="-w -s" -o /otel-action

#########################
# Build the small image #
#########################
FROM scratch
LABEL maintainer="MNThomson"

# Setup SSL certs
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Setup user
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
USER otel:otel

# Copy the static executable
COPY --from=builder /otel-action /otel-action

# Run the binary
ENTRYPOINT ["/otel-action"]
