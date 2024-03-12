# Build stage
FROM golang:1.22.1-alpine3.19 AS builder
WORKDIR /app
COPY . .
RUN go env -w  GOPROXY=https://goproxy.cn,direct && go build -o main main.go

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/app.env .

EXPOSE 8888
CMD [ "/app/main" ]