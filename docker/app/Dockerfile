FROM alpine:3.9

VOLUME ["/var/www/app"]

RUN apk add --no-cache libc6-compat tzdata \
    && cp /usr/share/zoneinfo/Europe/Kiev /etc/localtime \
    && echo "Europe/Kiev" >  /etc/timezone

WORKDIR "/var/www/app"

ENTRYPOINT ["./reprogl"]
