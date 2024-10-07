FROM golang:1.22.5-alpine AS builder

WORKDIR /app

RUN apk --no-cache add build-base make git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make build

FROM scratch

WORKDIR /app

ARG IMMO_MONGO_URI

ENV IMMO_MONGO_URI=${IMMO_MONGO_URI}

COPY --from=builder /app/build/immortal .
COPY --from=builder /config/config.yml ./config.yml

ENTRYPOINT ["./immortal", "run", "./config.yml"]