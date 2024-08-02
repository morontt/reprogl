# Reprogl

![license](https://img.shields.io/github/license/morontt/reprogl)

Source code of new version my [blog](https://xelbot.com/)

This project is part of my personal blog system. Admin panel and information about the database (container, migrations,
etc.) are located [here](https://github.com/morontt/zend-blog-3-backend).

## Start project

Copy config file from dist and build application

```bash
  cp .env{.dist,} # edit the settings if necessary
  cp app.ini{.dist,}

  docker compose build
  docker compose up

  docker exec gopher bash -c "yarn install"

  ./assets.sh
  ./build.sh
```

Start docker-compose and application

```bash
  docker compose up
```
