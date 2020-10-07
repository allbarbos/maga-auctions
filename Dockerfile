FROM golang:1.13.10-alpine AS builder
RUN apk add --update --no-cache \
  build-base \
  upx
WORKDIR /go/src/github.com/maga-auctions/
COPY . .
RUN go build -o main ./api/cmd/main.go && upx main

FROM alpine:latest AS runtime
WORKDIR /root/
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/github.com/maga-auctions/main ./
CMD ["./main"]