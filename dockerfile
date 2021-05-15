#Based on golang 1.16 builder
FROM golang:1.16 as builder

LABEL description "Custom image for the discordbot"

ADD . /

WORKDIR /

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o bot /bot/bot.go

FROM scratch

ENV BOT_TOKEN=ODM0MDY2NTk1NzAxMTk0ODEy.YH7fDA.S9q-1osgGnvAXLQ4_UxkAx_MxdQ

WORKDIR /

COPY --from=builder / .
RUN mv /bot/bot /main

CMD ["/main"]
