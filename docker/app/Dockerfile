FROM xelbot/docker-reprogl

ARG HOST_UID=1000

VOLUME ["/var/www/app"]

RUN usermod -u ${HOST_UID} www-data && groupmod -g ${HOST_UID} www-data

WORKDIR "/var/www/app"

ENTRYPOINT ["./reprogl"]
