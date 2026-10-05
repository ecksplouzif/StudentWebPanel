# StudentWebPanel

A web panel for student groups: class schedule, homework deadlines, and announcements in one place.

**Stack:** Go (standard `net/http`), HTML templates, Supabase (Postgres + Google auth), Nginx, VPS.

Development progress is tracked in [docs/PLAN.md](docs/PLAN.md). At the moment the panel supports
Google sign-in, a profile page, logout, and the database schema for schedules; the schedule UI and
importer are in progress.

## Contents

- [Requirements](#requirements)
- [Quick start](#quick-start)
- [Configuration](#configuration)
- [Database migrations](#database-migrations)
- [Make targets](#make-targets)
- [Routes](#routes)
- [Schedule importer](#schedule-importer)
- [Deployment](#deployment)
- [Troubleshooting](#troubleshooting)
- [Project structure](#project-structure)

## Requirements

| Tool | Version | Used for |
| --- | --- | --- |
| [Go](https://go.dev/dl/) | 1.26+ | building and running the app |
| [GNU Make](https://www.gnu.org/software/make/) | any | `make run` / `build` / `lint` |
| A [Supabase](https://supabase.com) project | — | Google sign-in and the Postgres database |
| [goose](https://github.com/pressly/goose) | v3 | applying database migrations |
| [golangci-lint](https://golangci-lint.run/) | v2 | linting and formatting (optional) |
| [Docker](https://www.docker.com/) | any | local Postgres instead of Supabase's (optional) |

Install the Go tools:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

```bash
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
```

(On macOS `brew install goose golangci-lint` works too.) Make sure `$(go env GOPATH)/bin` is on your `PATH`.

## Quick start

### 1. Clone and configure

```bash
git clone https://github.com/sqprrr/StudentWebPanel.git
```

```bash
cd StudentWebPanel && cp .env.example .env
```

Fill in `.env` as described in [Configuration](#configuration).

### 2. Set up Supabase

1. Create a project at [supabase.com](https://supabase.com).
2. **Authentication → Sign In / Providers → Google**: enable it and paste the Client ID and secret
   of a Google Cloud OAuth client. In Google Cloud, add
   `https://<project-ref>.supabase.co/auth/v1/callback` as an authorized redirect URI.
3. **Authentication → URL Configuration → Redirect URLs**: add
   `http://localhost:8080/auth/callback` (and your production URL, e.g.
   `https://grouphub.pp.ua/auth/callback`).
4. **Project Settings → JWT Keys**: make sure the project signs tokens with an asymmetric key
   (ECC/RSA). The app verifies tokens against the project's public JWKS, which is empty for the
   legacy shared-secret (HS256) setup.
5. Copy the project URL, the anon/publishable API key, and the database connection string
   (**Connect** button at the top of the dashboard) into `.env`.

### 3. Apply migrations

```bash
goose -dir migrations postgres "$(grep '^CONECT_DATABASE_URL=' .env | cut -d= -f2-)" up
```

See [Database migrations](#database-migrations) for details and a local Postgres alternative.

### 4. Run

```bash
make run
```

Check that it's up:

```bash
curl http://localhost:8080/healthz
```

It should answer `{"message":"OK"}`. Open <http://localhost:8080/login> and sign in with Google;
you'll land on `/profile` showing your name.

> The session cookies are set with `Secure`. Chrome and Firefox accept them on
> `http://localhost`; if sign-in loops back to `/login` in another browser, try one of those.

## Configuration

All configuration comes from environment variables. `make` loads them from `.env` automatically
(`-include .env` + `export` in the [Makefile](Makefile)). Write values **without quotes** — `make`
would keep the quotes as part of the value.

### App

| Variable | Required | Description | Example |
| --- | --- | --- | --- |
| `PORT` | no | HTTP port (default `8080`) | `8080` |
| `PROJECT_URL` | **yes** | Supabase project URL; the app fetches `<PROJECT_URL>/auth/v1/.well-known/jwks.json` at startup to verify tokens | `https://abcd1234.supabase.co` |
| `PROVIDER_GOOGLE_URL` | **yes** | Supabase authorize URL for Google, ending in `code_challenge=` — the app appends the PKCE challenge | `https://abcd1234.supabase.co/auth/v1/authorize?provider=google&redirect_to=http://localhost:8080/auth/callback&code_challenge=` |
| `TOKEN_URL` | **yes** | Supabase endpoint that exchanges the auth code for a session | `https://abcd1234.supabase.co/auth/v1/token?grant_type=pkce` |
| `SUPABASE_API_KEY` | **yes** | Supabase anon/publishable key, sent as the `apikey` header | `eyJhbGciOi...` |
| `CONECT_DATABASE_URL` | **yes** | Postgres connection string (note the spelling: `CONECT`) | `postgres://postgres:password@db.abcd1234.supabase.co:5432/postgres` |

### Deployment (used only by `deploy/deploy.sh`)

| Variable | Description | Example |
| --- | --- | --- |
| `DEPLOY_PORT` | SSH port of the server | `22` |
| `SERVER_USER` | SSH user on the server | `deploy` |
| `SERVER_IP` | Server address | `203.0.113.10` |
| `HEALTHZ_URL` | Public health check URL polled after deploy | `https://grouphub.pp.ua/healthz` |

## Database migrations

Migrations live in [migrations/](migrations) and use [goose](https://github.com/pressly/goose)
annotations (`-- +goose Up` / `-- +goose Down`). They create the `groups`, `roles`, `users`,
`subjects`, `class_types`, `schedule` and `links` tables and seed the roles and class types.

The commands below assume the connection string is exported in your shell
(`export CONECT_DATABASE_URL=...`).

```bash
goose -dir migrations postgres "$CONECT_DATABASE_URL" up
```

```bash
goose -dir migrations postgres "$CONECT_DATABASE_URL" status
```

Roll back the last migration with `down`, or everything with `reset`.

### Local Postgres with Docker

If you'd rather not touch the Supabase database while developing, run Postgres locally. Auth
still goes through Supabase; only the user/schedule data is stored locally.

```bash
docker run -d --name studentwebpanel-db -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=webpanel -p 5432:5432 postgres:17
```

Then set `CONECT_DATABASE_URL=postgres://postgres:postgres@localhost:5432/webpanel?sslmode=disable`
in `.env` and apply the migrations as above.

## Make targets

| Command | What it does |
| --- | --- |
| `make run` | `go run .` with `.env` loaded |
| `make build` | Builds the `WebPanel` binary in the repo root (templates are embedded, so the binary is self-contained) |
| `make lint` | `golangci-lint run` with the linters from [.golangci.yml](.golangci.yml) |
| `make forms` | `golangci-lint fmt` — formats the code with `gofumpt` |

There are no tests yet; `go test ./...` compiles every package and `go vet ./...` should pass.

## Routes

| Path | Auth | Description |
| --- | --- | --- |
| `/healthz` | — | Health check, returns `{"message":"OK"}` |
| `/login` | — | Login page with the "Sign in with Google" button |
| `/auth` | — | Generates a PKCE verifier (stored in a cookie for 5 min) and redirects to Google via Supabase |
| `/auth/callback` | — | Exchanges the code for an access token, stores it in the `access_token` cookie, upserts the user into `users`, redirects to `/profile` |
| `/profile` | required | Shows the signed-in user's name and ID |
| `/logout` | — | Clears the session cookie and redirects to `/login` |
| everything else | required | JSON `404` |

Protected routes redirect to `/login` when the `access_token` cookie is missing or its JWT
doesn't verify against the Supabase JWKS. The session cookie currently lives for 5 minutes.

## Schedule importer

`cmd/import` is a work-in-progress CLI that reads the university's schedule export
(Windows-1251 CSV with `\r` line endings), decodes it, and prints the parsed rows. It currently
expects the file `TimeTable_04_10_2026.csv` in the working directory:

```bash
go run ./cmd/import
```

CSV files are git-ignored, so the export is never committed.

## Deployment

Production runs as a systemd service behind Nginx on a VPS. Deploy from your machine with:

```bash
./deploy/deploy.sh
```

The script reads the deployment variables from `.env`, then:

1. cross-compiles a `linux/amd64` binary (`StudentWebPanel.new`);
2. copies it to `/opt/StudentWebPanel/` over `scp`;
3. stops the `StudentWebPanel` service, keeps the previous binary as `StudentWebPanel.old`,
   swaps in the new one and starts the service;
4. waits 5 seconds and checks `HEALTHZ_URL` — the script fails if it doesn't return 2xx.

The server is expected to have:

- a systemd unit named `StudentWebPanel` running `/opt/StudentWebPanel/StudentWebPanel`, with the
  [app variables](#app) provided through its environment (e.g. `EnvironmentFile=`) — secrets live
  on the server, not in the repo;
- an SSH user that can write to `/opt/StudentWebPanel/` and run `sudo systemctl` without a password;
- Nginx proxying to the app on port 8080 — see [deploy/nginx.conf](deploy/nginx.conf).

### Rolling back

The previous binary is kept on the server. To restore it:

```bash
ssh -p "$DEPLOY_PORT" "$SERVER_USER@$SERVER_IP" "sudo systemctl stop StudentWebPanel && mv /opt/StudentWebPanel/StudentWebPanel.old /opt/StudentWebPanel/StudentWebPanel && sudo systemctl start StudentWebPanel"
```

## Troubleshooting

| Symptom | Cause / fix |
| --- | --- |
| `panic: runtime error: invalid memory address` right after `Failed get public key Supabase` | `PROJECT_URL` is empty or unreachable — the JWKS is fetched at startup. Check `.env`. |
| Sign-in always returns to `/login` | The JWT doesn't verify: the project uses the legacy HS256 secret (switch to asymmetric JWT keys), the cookie expired (5 min), or the browser rejects `Secure` cookies on `http://localhost`. Check the server log for `Failed parse JWT token`. |
| Supabase shows "redirect URL not allowed" | Add the exact `redirect_to` from `PROVIDER_GOOGLE_URL` to Supabase's Redirect URLs. |
| `Non-OK response from token endpoint` in the log | Wrong `TOKEN_URL` or `SUPABASE_API_KEY`, or the PKCE cookie expired (sign-in took longer than 5 minutes). |
| `Add user information` error in the log | Migrations weren't applied to the database in `CONECT_DATABASE_URL`. |
| `make: golangci-lint: No such file or directory` | Install golangci-lint (see [Requirements](#requirements)). |

## Project structure

```
.
├── main.go                  # entry point: config, routes, server startup
├── cmd/
│   └── import/              # schedule CSV importer (work in progress)
├── internal/
│   ├── handlers/            # HTTP handlers: login, auth, callback, profile, logout, healthz
│   └── middleware/          # JWT auth and request logging
├── web/                     # HTML templates (embedded into the binary)
├── migrations/              # goose SQL migrations
├── deploy/                  # deploy script and nginx config
├── docs/                    # development plan and decision log
├── Makefile                 # run, build, lint, fmt
├── .golangci.yml            # linter and formatter settings
└── .env.example             # required environment variables
```
