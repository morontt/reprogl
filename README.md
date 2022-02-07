# Reprogl

Source code of new version my blog (coming soon)

This project is part of my personal blog system. Admin panel and information about the database (container, migrations,
etc.) are located [here](https://github.com/morontt/zend-blog-3-backend).

## Start project

Copy config file from dist and build application

```bash
  cp .env{.dist,} # edit the settings if necessary
  cp app.ini{.dist,}

  ./build.sh
```

Start docker-compose and application

```bash
  docker-compose up --build
```
