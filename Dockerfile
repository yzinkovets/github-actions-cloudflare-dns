FROM golang:1.23 AS builder

WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o /bin/cf-dns .


FROM alpine:latest
COPY --from=builder /bin/cf-dns /bin/cf-dns
RUN apk --no-cache add ca-certificates
ENTRYPOINT ["/bin/cf-dns"]