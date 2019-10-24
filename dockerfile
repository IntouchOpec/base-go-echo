
# Start from golang:1.12-alpine base image
FROM golang:1.12-alpine

# The latest alpine images don't have some tools like (`git` and `bash`).
# Adding git, bash and openssh to the image
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

# Add Maintainer Info
LABEL maintainer="Rajeev Singh <rajeevhub@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

ENV GO111MODULE=on

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD ["./main"]

# FROM golang:1.12-alpine as builder

# # To fix go get and build with cgo
# RUN apk add --no-cache --virtual .build-deps \
#     bash \
#     gcc \
#     git \
#     musl-dev

# RUN mkdir build
# COPY . /build
# WORKDIR /build

# ADD conf /conf

# RUN go get
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o webserver .
# RUN adduser -S -D -H -h /build webserver
# USER webserver

# FROM scratch
# COPY --from=builder /build/webserver /app/
# WORKDIR /app
# EXPOSE 5000
# CMD ["./webserver"]
