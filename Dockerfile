FROM golang:1.15-alpine as builder

RUN apk update && \
  apk add --no-cache git && \
  apk add busybox-extras

WORKDIR /app
COPY . .

# RUN go mod download
# RUN go build -o main .

CMD ["go", "run", "main.go"]

# FROM alpine
# WORKDIR /app
# COPY --from=builder /app/main .
# CMD ["./main", "serveHttpUsers"]