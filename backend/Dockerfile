FROM golang:1.23-alpine AS builder

COPY . /github.com/merynayr/PingerVK/backend/source/
WORKDIR /github.com/merynayr/PingerVK/backend/source/

RUN go mod download
RUN go build -o ./bin/httpserver cmd/main.go

FROM alpine:latest


WORKDIR /root/

COPY --from=builder /github.com/merynayr/PingerVK/backend/source/bin/httpserver .

COPY .env .

CMD ["./httpserver", "-config-path=.env"]