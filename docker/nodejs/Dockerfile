FROM node:16-bullseye-slim

ENV TZ "Europe/Moscow"

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN npm install npm@latest -g \
    && npm install -g grunt-cli@1.4.3 \
    ;

COPY ./.bashrc /root/.bashrc
COPY ./.bashrc /home/node/.bashrc

WORKDIR /var/www/app
