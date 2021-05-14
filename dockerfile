#Based on golang 1.16 builder
FROM golang:1.16 as builder

LABEL description "Custom image for the discordbot"

ADD . /

WORKDIR /

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o bot /bot/bot.go

FROM scratch

ENV BOT_TOKEN=ODM0MDE1NzE0MjAwNjQ5NzU4.YH6vqQ.XyxiV8tp0sQqDNkSBW5bT6Wobmg

WORKDIR /

COPY --from=builder / .
RUN mv /bot/bot /main

CMD ["/main"]
