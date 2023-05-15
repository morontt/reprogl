FROM debian:bookworm-slim

ENV TZ "Europe/Moscow"

VOLUME ["/var/www/app"]

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN apt-get update; \
    apt-get install -y --no-install-recommends \
    ca-certificates \
    ; \
    rm -rf /var/lib/apt/lists/*

WORKDIR "/var/www/app"

ENTRYPOINT ["./reprogl"]
