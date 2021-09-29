#####################################
#   STEP 1 build executable binary  #
#####################################
FROM golang:alpine AS builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go



#####################################
#   STEP 2 build a small image      #
#####################################
FROM alpine:latest

# Copy our static executable.
COPY --from=builder /app/main /app/main
COPY certs  /app/certs
COPY img /app/img
EXPOSE 8080
WORKDIR /app
ENTRYPOINT ["/app/main"]



