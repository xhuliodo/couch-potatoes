FROM golang:1.16 as test_and_build

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /test_and_build
WORKDIR /test_and_build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Command to run unit testing
RUN go test -v ./...

# Generate latest docs
RUN make install_swag
RUN make create_docs

# Build the application
RUN go build -o api .

# production environment
FROM alpine:3.12.0
WORKDIR /app
COPY --from=test_and_build /test_and_build/api .

# Create logs and auth cache folder
RUN mkdir /logs && mkdir /auth

# Export necessary port
EXPOSE 4001

# Command to run when starting the container
CMD ["./api"]