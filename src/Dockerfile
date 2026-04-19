FROM golang:1.26-alpine
WORKDIR /app
COPY . .
RUN go build -o app
RUN mkdir -p /app/logs /app/config
CMD ["./app"]