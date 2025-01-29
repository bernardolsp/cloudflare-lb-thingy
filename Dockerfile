FROM golang:1.23.2-alpine AS builder
WORKDIR /app
COPY . ./
RUN go build -o app .

FROM scratch
COPY --from=builder /app/app /app
ENTRYPOINT ["/app"]