# Start from golang base image
FROM golang:alpine as builder

# ENV GO111MODULE=on

# Add Maintainer info
LABEL maintainer="Alberto Llamedo"

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git \
&& rm -rf /var/lib/apt/lists/*


# Set the current working directory inside the container 
WORKDIR /monitoring-aggr-service

# Copy go mod and sum files 
COPY go.mod go.sum ./

# RUN go mod tidy

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed 
RUN go mod download 

# Copy the source from the current directory to the working Directory inside the container 

COPY . .
# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o aggregator .

# Start a new stage from scratch
FROM alpine:3.17
RUN apk --no-cache add ca-certificates

RUN addgroup -S icos && adduser -S icos -G icos -u 1001
USER icos


WORKDIR /home/icos/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder //monitoring-aggr-service/aggregator .
# COPY --from=builder //monitoring-aggr-service/.env .

# Expose port 8080 to the outside world
EXPOSE 8080

#Command to run the executable
CMD ["./aggregator"]
