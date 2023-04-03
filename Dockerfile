# FOR GAE Flexible Environment Custom Runtime 
FROM golang:1.19

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# install build essentials
RUN apt-get update && \
    apt-get install -y wget build-essential pkg-config --no-install-recommends

# Move to working directory /build
WORKDIR /build

# Copy the code into the container\
COPY . .

# Copy and download dependency using go mod
RUN go mod download

# Build the application
RUN go build -o main .

# Export necessary port # default GCP App Engine Port
EXPOSE 8080

# Command to run when starting the container
CMD ["/build/main"]