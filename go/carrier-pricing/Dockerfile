FROM golang:1.13.4-alpine

# The latest alpine images don't have some tools like (`git` and `bash`).
# Adding git, bash and openssh to the image
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

LABEL maintainer="Alessio Savi <alessiosavibtc@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . /app

RUN go clean

# Build the Go app
RUN go build -o carrier-pricing .

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD ["./carrier-pricing"]