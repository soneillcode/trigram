FROM golang:alpine AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY cmd ./cmd
COPY pkg ./pkg

# Build the application
RUN go build -ldflags="-w -s" -o main ./cmd

# Build a small image
FROM scratch

COPY --from=builder /build/main /

# Command to run when starting the container
ENTRYPOINT ["/main"]