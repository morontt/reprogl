FROM alpine:3.9

VOLUME ["/var/www/app"]

RUN apk add --no-cache libc6-compat

WORKDIR "/var/www/app"

ENTRYPOINT ["./reprogl"]
