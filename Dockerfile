FROM golang:alpine AS builder
RUN apk add git g++ gcc

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build
COPY . .
RUN go mod download
RUN go build -o main .
WORKDIR /dist
RUN cp /build/main .

FROM alpine
COPY --from=builder /dist/main .
ENTRYPOINT ["/main"]
