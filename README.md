# soa-tourist-app

KT1 SOA tourist application monorepo.

## Structure

- `frontend/tourist-app-ui` Angular starter app
- `services/stakeholders-service` independent .NET Web API
- `services/blog-service` independent .NET Web API

## Local ports

- Frontend: `4200`
- Stakeholders service: `5001`
- Blog service: `5002`

## Docker Compose ports

- Frontend: `4200`
- Stakeholders service: `5001`
- Blog service: `5002`
- Stakeholders database (PostgreSQL): `5433`
- Blog database (PostgreSQL): `5434`

## Docker Compose run

Start infrastructure and services:

```bash
docker compose up --build
```

Run migrations in dedicated migration containers:

```bash
docker compose -f docker-compose-migration.yml up --build
```

Note: application services do not run migrations automatically on startup.

Seed databases from mounted SQL files:

```bash
docker exec -it stakeholders-db psql -U postgres -d stakeholders_db -f /tmp/stakeholders-init.sql
docker exec -it blog-db psql -U postgres -d blog_db -f /tmp/blog-init.sql
```

Stop all containers:

```bash
docker compose down
```
