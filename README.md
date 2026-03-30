# soa-tourist-app

KT1 SOA tourist application monorepo scaffold.

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
- Stakeholders service: `8081`
- Blog service: `8082`

## Notes

- Both backend services are separate runnable starter projects.
- Each backend service includes `Controllers`, `Models`, `DTOs`, `Services`, and `Repositories`.
- Sample health endpoints are available at `/api/health`.
- PostgreSQL preparation is limited to configuration placeholders and connection strings. Business logic is intentionally not implemented yet.
- Angular was generated with `--skip-install`.
- The current environment did not have the `.NET SDK` installed, so the Web API starter files were created manually.
