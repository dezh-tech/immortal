# --------------------------
#! Stage 1: Build the Go binary
# --------------------------
FROM golang:1.23.3-alpine AS builder

WORKDIR /app

#* Install required build tools
RUN apk --no-cache add build-base git

#* Cache dependencies by copying go.mod and go.sum first
COPY go.mod go.sum ./
RUN go mod download

#* Copy the rest of the application source code
COPY . .

#* Build the Go binary using the provided Makefile
RUN make build

# --------------------------
#! Stage 2: Create the runtime image
# --------------------------
FROM alpine:latest

WORKDIR /app

#* Define build-time arguments and environment variables
ARG IMMO_MONGO_URI
ENV IMMO_MONGO_URI=${IMMO_MONGO_URI}

ARG IMMO_REDIS_URI
ENV IMMO_REDIS_URI=${IMMO_REDIS_URI}

#* Copy the compiled binary from the builder stage
COPY --from=builder /app/build/immortal .
COPY --from=builder /app/config/config.yml .

#* Expose necessary ports for the application

# websocket port
EXPOSE 7777

# grpc port
EXPOSE 50050

#* Set the entrypoint to run the application
ENTRYPOINT ["./immortal", "run", "./config.yml"]
