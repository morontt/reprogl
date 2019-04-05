FROM nginx:1.14

COPY ./web.conf /etc/nginx/conf.d
COPY ./fullchain.pem /etc/nginx/fullchain.pem
COPY ./privkey.pem /etc/nginx/privkey.pem

EXPOSE 80
