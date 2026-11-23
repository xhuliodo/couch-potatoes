# Build stage. The Go toolchain targets the platform of this stage, which is
# the build platform requested by --platform, so no GOARCH is pinned.
FROM golang:1.27-alpine AS test_and_build

# Set necessary environment variables needed for our image
ENV CGO_ENABLED=0

# Move to working directory /test_and_build
WORKDIR /test_and_build

# Copy and download dependencies using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Command to run unit testing
RUN go test ./...

# Build the application
RUN go build -o api .

# production environment
FROM alpine:3.24

# Create the runtime user and the writable directories it owns. The auth
# cache holds the tenant JWKS and is only readable by the runtime user.
RUN addgroup -S -g 10001 app \
    && adduser -S -D -H -u 10001 -G app app \
    && mkdir -p /app /logs /auth \
    && chown -R app:app /app /logs /auth \
    && chmod 0700 /auth /logs

WORKDIR /app
COPY --from=test_and_build /test_and_build/api .

USER app

# Export necessary port
EXPOSE 4001

# Command to run when starting the container
CMD ["./api"]
