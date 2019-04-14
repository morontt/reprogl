FROM nginx:1.14

RUN mkdir /etc/nginx/ssl \
  && ln -sf /var/www/app/docker/nginx/fullchain.pem /etc/nginx/ssl/cert.pem \
  && ln -sf /var/www/app/docker/nginx/privkey.pem /etc/nginx/ssl/key.pem

COPY ./web.conf /etc/nginx/conf.d

EXPOSE 80
