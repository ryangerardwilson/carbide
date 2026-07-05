# Environment

The docs app keeps secrets in the remote `.env` file created during deploy.
Secret values are never printed by `carbide doctor`, `deploy preview`, or
`deploy apply`.

Required production values:

- `CARBIDE_HTTP_PORT`
- `PUBLIC_URL`
- `POSTGRES_PASSWORD`
- `DATABASE_URL`

`POSTGRES_PASSWORD` and `DATABASE_URL` are framework-owned deploy secrets.
