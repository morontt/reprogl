FROM debian:bookworm-slim

ENV TZ "Europe/Moscow"

VOLUME ["/var/www/app"]

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR "/var/www/app"

ENTRYPOINT ["./reprogl"]
