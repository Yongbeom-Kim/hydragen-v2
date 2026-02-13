# AGENTS Guide for `hydragen-v2`

This document is for coding agents working in this repository.

## 1) Repository overview

`hydragen-v2` is a Docker Compose–orchestrated monorepo with multiple services:

- `client/` — React + TanStack Router frontend (Vite + pnpm).
- `server/` — Go backend API.
- `dataset-loader/` — Python data ingestion/loader jobs.
- `flyway/` — SQL migrations and Flyway config.
- `caddy/` — reverse-proxy/web entrypoint config.
- `infra/` — Terraform files for infrastructure.
- Root compose files and `Makefile` coordinate all services.

## 2) Critical rule: do **not** run services directly

Do **not** start services with ad-hoc local commands like:

- `go run ...`
- `pnpm dev`
- `python index.py`
- `docker compose ...` with custom file combinations that bypass project conventions

Instead, use the **Makefile targets** that define supported Docker Compose flows.

## 3) Standard development workflow

From repository root:

1. Ensure env file exists and is configured (`.env`).
2. Use Make targets for stack lifecycle.
3. Run checks/tests in the relevant service (prefer containerized commands when practical).
4. Stop stacks when done.

### Primary commands

- `make help` — list supported targets.
- `make dev` — start development stack (build + up).
- `make start` / `make start_detached` — production-style stack.
- `make start-cx23` / `make start_detached-cx23` — constrained-resource production profile.
- `make logs` / `make logs-cx23` — stream logs.
- `make stop` — tear down running stacks.

### Useful shell attach targets

- `make attach-postgres`
- `make attach-flyway`
- `make attach-dataset-loader`

## 4) Docker Compose file intent

- `docker-compose.base.yaml` — shared service definitions.
- `docker-compose.dev.yaml` — local development overrides (ports, bind mounts, dev targets).
- `docker-compose.prod.yaml` — production/runtime configuration.
- `docker-compose.cx23.yaml` — resource limits profile for smaller hosts.

Agents should not invent alternate orchestration paths unless explicitly requested.

## 5) Change guidelines by area

### Frontend (`client/`)

- Keep route structure aligned with TanStack Router file conventions.
- Keep API access centralized in feature API/query modules.
- For UI changes, include a short note about impacted routes/components.

### Backend (`server/`)

- Preserve separation between handlers, domain, db, and utility layers.
- Prefer small, focused handler changes with clear input/output behavior.

### Data loading (`dataset-loader/`)

- Keep loader logic deterministic and resumable when possible.
- Avoid hardcoding environment-specific paths/credentials.

### Database (`flyway/sql`)

- Add forward-only migration files with clear naming.
- Avoid editing old migrations that may have already been applied.

## 6) Validation expectations for agents

Before finalizing changes:

- Run the most relevant checks/tests for modified areas.
- If services are required, use Make/Docker Compose flows from this repo.
- Include concise evidence of what commands were run and their outcomes.

## 7) Safety and operational constraints

- Do not delete volumes/data unless explicitly asked.
- Do not run destructive database operations unless explicitly asked.
- Do not modify infrastructure state files casually (`infra/terraform.tfstate`).
- Keep edits scoped to the task.

## 8) If unsure

When uncertain about runtime behavior, default to:

1. `make help`
2. `make dev`
3. inspect logs with `make logs`

and proceed with minimal, reversible changes.
