# Image-Based Deployment

The server does not need the backend or frontend source code. GitHub Actions builds images and publishes them to GitHub Container Registry (GHCR).

## Publish Images

Push the backend repository and frontend repository to their default branches. The workflows publish these images:

```text
ghcr.io/jienigui06091/task-manager-api
ghcr.io/jienigui06091/task-manager-postgres
ghcr.io/jienigui06091/task-manager-frontend
```

## Server Files

Create a deployment directory on the server containing only:

```text
deploy/
  docker-compose.yml
  .env
```

Copy `docker-compose.yml` and `.env.example` to that directory, then create `.env` from the example. Replace `POSTGRES_PASSWORD` and `JWT_SECRET` before starting.

If the GHCR packages are private, authenticate once with a GitHub personal access token that has `read:packages`:

```sh
docker login ghcr.io
```

## Start And Update

Start the full application without building source code:

```sh
docker compose up -d
```

Open `http://SERVER_IP:8080`. Change `APP_PORT` in `.env` if port 8080 is occupied.

To update to newly published image tags:

```sh
docker compose pull
docker compose up -d
```

`docker compose down` keeps the PostgreSQL data volume. `docker compose down -v` also removes the database data.
