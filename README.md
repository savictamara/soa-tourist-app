# soa-tourist-app

SOA tourist application

## Structure

- `frontend/tourist-app-ui` — Angular frontend
- `gateway/` — Nginx API Gateway
- `services/stakeholders-service` — .NET Web API (auth, users), PostgreSQL
- `services/blog-service` — .NET Web API (blogs, comments), PostgreSQL
- `services/tour-service` — Go REST API (tours, key points, reviews), MongoDB
- `services/follower-service` — Java Spring Boot (followers, recommendations), Neo4j

## Architecture

All traffic goes through the API Gateway. Services are not directly accessible from the host.

```
Browser → localhost:8080 (API Gateway)
              ├── /api/auth      → stakeholders-service
              ├── /api/users     → stakeholders-service
              ├── /api/blogs     → blog-service
              ├── /api/tours     → tour-service
              └── /api/followers → follower-service
```

## Ports

| Service | Host port | Notes |
|---|---|---|
| Frontend (Angular) | `4200` | |
| API Gateway (Nginx) | `8080` | Single entry point for all API calls |
| Stakeholders DB (PostgreSQL) | `5433` | |
| Blog DB (PostgreSQL) | `5434` | |
| Tour DB (MongoDB) | `27017` | |
| Follower DB (Neo4j browser) | `7474` | |
| Follower DB (Neo4j bolt) | `7687` | |

Backend services have no exposed host ports — accessible only within the Docker network via the gateway.

## Docker Compose run

Start everything:

```bash
docker compose up --build
```

Stop all containers:

```bash
docker compose down
```

Seed databases from mounted SQL files:

```bash
docker exec -it stakeholders-db psql -U postgres -d stakeholders_db -f /tmp/stakeholders-init.sql
docker exec -it blog-db psql -U postgres -d blog_db -f /tmp/blog-init.sql
```
