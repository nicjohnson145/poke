FROM golang:1.20-alpine AS builder
WORKDIR /src
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build

# Mabye scratch? :shrug:
FROM alpine:3.17.1
COPY --from=builder /src/echoserver /bin/echoserver
ENTRYPOINT ["/bin/echoserver"]
