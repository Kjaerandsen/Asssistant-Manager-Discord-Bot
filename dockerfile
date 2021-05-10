#Based on golang 1.16 builder
FROM golang:1.16 as builder

LABEL description "Custom image for the discordbot"

ADD . /

WORKDIR /

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o bot /bot/bot.go

FROM scratch

#ENV BOT_TOKEN=ODM1NDcxODgxOTY0MjI0NTEz.YIP70g.Pj7Shx_nZprANAliX46VdvpLYfU

WORKDIR /

COPY --from=builder / .

CMD ["/bot/bot"]
