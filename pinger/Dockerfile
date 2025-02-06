FROM golang:1.23-alpine AS builder

COPY . /github.com/merynayr/PingerVK/pinger/source/
WORKDIR /github.com/merynayr/PingerVK/pinger/source/

RUN go mod download
RUN go build -o ./bin/pinger cmd/main.go

FROM alpine:latest


WORKDIR /root/

COPY --from=builder /github.com/merynayr/PingerVK/pinger/source/bin/pinger .

COPY .env .

CMD ["./pinger",  "-config-path=.env"]