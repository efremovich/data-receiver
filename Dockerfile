FROM --platform=linux/amd64 golang:1.22-alpine as builder

WORKDIR /build
COPY . .
RUN  go mod download
RUN  GOOS=linux GOARCH=amd64 go build -v -o ./package-receiver cmd/app/main.go

FROM --platform=linux/amd64 alpine:3.19.1
RUN apk add --no-cache gcompat
WORKDIR /app
COPY --from=builder /build/package-receiver ./
COPY --from=builder /build/docs/swagger ./docs/swagger
COPY --from=builder /build/config/vars.env ./config/

EXPOSE 8080

CMD ["/app/package-receiver"]
