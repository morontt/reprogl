FROM node:16-bullseye-slim

RUN npm install npm@latest -g
RUN npm install -g grunt-cli@1.4.3

COPY ./.bashrc /root/.bashrc
COPY ./.bashrc /home/node/.bashrc

WORKDIR /var/www/app
