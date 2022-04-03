FROM golang:1.18.0-alpine3.15 AS builder

RUN apk --no-cache add ca-certificates git

WORKDIR /app
COPY  . .

RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/ssh

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app

COPY --from=builder /app/app .

CMD ["./app"]