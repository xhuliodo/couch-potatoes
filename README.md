# Couch Potatoes

Couch Potatoes is a local full-stack movie recommender. You sign in, pick a few
favourite genres, rate fifteen movies by swiping, and the app suggests films two
ways: **content-based** (movies that share cast, directors, or writers with films
you liked) and **user-based** (movies liked by the users most similar to you).
Popular movies in your favourite genres give you something to rate during
onboarding.

Everything runs on your machine with Docker Compose; nothing is hosted.

## Services

| Service | Container | What it is | Local address |
| --- | --- | --- | --- |
| Client | `cp-client` | React single-page app served by nginx | <http://localhost:3000> |
| API | `cp-api` | Go REST API, Clean Architecture layout | <http://localhost:4001> |
| Database | `cp-db` | Neo4j 4.4 with APOC, seeds itself on first start | <http://localhost:7474> (browser), `bolt://localhost:7687` |
| Thumbnails | `cp-thumbnails` | nginx serving poster files from `thumbnails/img` | <http://localhost:4003> |

## Prerequisites

- Docker Desktop with Docker Compose v2 (`docker compose`, not the old
  `docker-compose` binary).
- Free ports **3000**, **4001**, **7474**, **7687**, and **4003**.
- At least **4 GB** of memory available to Docker. The first image build and the
  Neo4j seed are the memory-heavy parts.
- Network access for the first build (base images and Go/Node dependencies) and
  at API start (the Auth0 JWKS fetch).

All images build natively for the host architecture: **arm64 on Apple Silicon,
amd64 on Intel**. You do not need `--platform` or emulation.

## Quick start

```sh
cp .env.example .env
docker compose up -d --build
```

The first build takes several minutes. After the containers start, the database
seeds itself; on a current laptop that takes about **20 seconds** (observed
8–21 s), and the API does not start serving until the database is healthy.

Watch progress if you want:

```sh
docker compose ps
docker compose logs -f cp-db
```

When `cp-db` and `cp-api` are healthy and `cp-client` is running, open
**<http://localhost:3000>**.

If you prefer not to create a `.env` file, pass the example file directly. Every
value already has a working default:

```sh
docker compose --env-file .env.example up -d --build
```

Stop the stack with `docker compose down`; data stays in `db/data`.

## Ports and endpoints

| Port | Service | Notes |
| --- | --- | --- |
| `3000` | Client | App entry point. Keep this port: the Auth0 tenant allows exactly `http://localhost:3000`. |
| `4001` | API | REST base URL. Swagger UI at `/docs/`, liveness at `/healthcheck`, readiness at `/ready` (both verify database connectivity). |
| `7474` | Neo4j | Browser UI. Connect with `bolt://localhost:7687` and the credentials from `.env`. |
| `7687` | Neo4j | Bolt endpoint used by the API (inside Compose: `bolt://cp-db:7687`). |
| `4003` | Thumbnails | Static poster images named `thumbnail_<movieId>.jpg`. |

Quick checks:

```sh
curl http://localhost:4001/healthcheck
curl http://localhost:4001/ready
```

## Configuration

Configuration comes from `.env` in the repository root; `.env.example` documents
every variable. Real `.env` files are git-ignored.

| Variable | Default | Purpose |
| --- | --- | --- |
| `COMPOSE_PROJECT_NAME` | `couch_potatoes` | Compose project name; prefixes the network and labels the stack. |
| `NEO4J_USER` | `neo4j` | Database user. Compose builds `NEO4J_AUTH` from user and password and passes both to the API. |
| `NEO4J_PASSWORD` | `letmein` | Database password. Local development default; change user and password together. |
| `NEO4J_URI` | `bolt://cp-db:7687` | Bolt endpoint the API connects to. Keep the Compose service name unless the database runs elsewhere. |
| `API_LISTEN_PORT` | `:4001` | API listen address, including the leading colon. |
| `API_ALLOWED_ORIGINS` | `http://localhost:3000` | Comma-separated CORS allowlist. Must contain the client origin. |
| `JWKS_URL` | `https://dev-ps5dqqis.eu.auth0.com/.well-known/jwks.json` | Tenant JWKS used to validate Auth0 access tokens. |
| `JWKS_CACHE_PATH` | `/auth/jwks.json` | Writable JWKS cache path inside the API container. |
| `REACT_APP_AUTH0_DOMAIN` | `dev-ps5dqqis.eu.auth0.com` | Auth0 tenant domain baked into the client. |
| `REACT_APP_AUTH0_CLIENT_ID` | `DqbCvZtL8cn5plDla9TYlrJLhWIXpZtV` | Public SPA client id baked into the client. Not a secret. |
| `REACT_APP_REDIRECT_BASE` | `http://localhost:3000` | Base for the login callback (`/solo`) and logout return (`/`). |
| `REACT_APP_API_ENDPOINT` | `http://localhost:4001` | API base URL used by the browser. |
| `REACT_APP_THUMBNAILS_ENDPOINT` | `http://localhost:4003` | Thumbnail base URL used by the browser. |

`REACT_APP_*` values are **baked into the client bundle at image build time**.
Editing `.env` does not change a running client: rebuild it.

```sh
docker compose up -d --build cp-client
```

## Authentication

The app uses the existing Auth0 tenant. These are the settings it expects; the
tenant already allows them, so **no Auth0 dashboard change is needed**.

| Setting | Value |
| --- | --- |
| Domain | `dev-ps5dqqis.eu.auth0.com` |
| Application type | Single Page Application (public client, no secret) |
| Client ID | `DqbCvZtL8cn5plDla9TYlrJLhWIXpZtV` |
| Allowed Callback URLs | `http://localhost:3000/solo` |
| Allowed Logout URLs | `http://localhost:3000/` |
| Allowed Web Origins | `http://localhost:3000` |

The API validates tokens against the tenant's JWKS over HTTPS, so sign-in needs
network access even though everything else is local.

## Data and seeding

The seed is the small [MovieLens Latest Dataset](dataset/01-ml-latest-small/README.txt),
enriched with cast, director, and writer data. It creates roughly **34,000
nodes** and 100,600 rating relationships:

- 9,687 movies and 19 genres
- 610 dataset users and about 100,600 ratings
- about 23,700 people (actors, directors, writers)

Seeding runs inside `cp-db`. On first start there is no seed marker, so the
entrypoint imports the dataset and writes a `SeedMarker` node on success. Later
starts find the marker and skip the import, coming up quickly. If an import
fails the marker is left unset, so the next start re-runs the idempotent import
even when the graph already holds a partial import — the safe recovery path. The
graph persists in the `db/data` bind mount across `docker compose down` and
`docker compose up`.

Inspect it at <http://localhost:7474> with the Neo4j connection URL
`bolt://localhost:7687` and the user/password from `.env`.

To reset the database and seed again:

```sh
docker compose down
find db/data -mindepth 1 ! -name .gitkeep -delete
docker compose up -d
```

If the files are owned by the container user, run the removal in a throwaway
container instead:

```sh
docker compose down
docker run --rm -v "$PWD/db/data:/data" alpine:3.24 \
  sh -c 'find /data -mindepth 1 ! -name .gitkeep -exec rm -rf {} +'
docker compose up -d
```

## Posters (optional)

Poster files are not committed. Without them, the client shows styled
posterless cards, so the app works fine. To fetch real posters, use a free
[TMDB](https://www.themoviedb.org/) v3 API key:

```sh
TMDB_API_KEY=your-key ./scripts/backfill-posters.sh
```

The script needs `curl` and `jq`, maps MovieLens ids to TMDB ids, and writes
`thumbnails/img/thumbnail_<movieId>.jpg`. Expect about **9,700 lookups** and
**tens of minutes** at TMDB's rate limit. It is resumable and idempotent:
completed files are skipped, and downloads are moved into place atomically, so
an interrupted or rate-limited run can simply be re-run. The thumbnails
directory is a bind mount, so no restart is needed; reload the page.

## Manual acceptance checklist

Run this after the documented setup path. Sign-in is the only step that needs a
human.

- [ ] `cp .env.example .env && docker compose up -d --build` finishes, and
      `docker compose ps` shows the database and API healthy and the client
      running.
- [ ] <http://localhost:3000> loads and **Sign in** redirects to the Auth0
      tenant and back to the app.
- [ ] Onboarding asks for genres; selecting fewer than three is refused, and
      three or more can be saved.
- [ ] Swiping **Loved it** / **Not for me** records ratings; after 15 ratings the
      **Next** button enables and leads to the recommendations.
- [ ] The **Our suggestions for you** tab (content-based) shows movie cards
      after at least one liked movie.
- [ ] The **Other people also liked** tab (user-based) shows movie cards once
      your ratings overlap enough with other users.
- [ ] A movie can be added to **Watch later** and appears under
      **Movies waiting for you**; removing it makes it disappear.
      **Your collection** shows watchlist movies you have already rated.
- [ ] `docker compose down && docker compose up -d`, then sign in again: the
      watchlist is still there.
- [ ] After `./scripts/backfill-posters.sh`, posters render on the cards; before
      it, cards render without broken-image icons.

## Troubleshooting

### Apple Silicon

Use the **arm64** build of Docker Desktop. Every base image is multi-arch, so
the stack builds natively; do not add `--platform linux/amd64`, which silently
reintroduces emulation.

Confirm what you built:

```sh
docker image inspect --format '{{.Os}}/{{.Architecture}}' \
  couch-potatoes-db:local couch-potatoes-api:local \
  couch-potatoes-client:local couch-potatoes-thumbnails:local
```

Each line should say `linux/arm64`. If an image says `linux/amd64`, it came from
an older build cache: rebuild without cache (`docker compose build --no-cache`)
and, if needed, `docker builder prune`.

### Ports already in use

Compose fails to start if a port is taken. Find the process and stop it rather
than moving the client off `3000`, which would break the Auth0 allowlist:

```sh
lsof -nP -iTCP:3000 -sTCP:LISTEN
lsof -nP -iTCP:4001 -sTCP:LISTEN
```

### Seed did not finish or data looks wrong

Check `docker compose logs cp-db` for the import result. If the seed failed, the
marker is left unset and the next start retries. To start from a clean graph, use
the reset commands under [Data and seeding](#data-and-seeding).

### Sign-in loop, 401s, or CORS errors

These are configuration mismatches, not stack failures:

- The client must run at exactly `http://localhost:3000`; the tenant's callback,
  logout, and web-origin entries are fixed to that address.
- `REACT_APP_REDIRECT_BASE` must be `http://localhost:3000`, and
  `API_ALLOWED_ORIGINS` must contain `http://localhost:3000`.
- After changing any `REACT_APP_*` value, rebuild the client:
  `docker compose up -d --build cp-client`.

### API cannot fetch JWKS

Token validation needs outbound HTTPS to `dev-ps5dqqis.eu.auth0.com`. If the API
logs a JWKS fetch failure, check network access from Docker and retry. The cache
lives at `JWKS_CACHE_PATH` inside the container; keep that path writable.

## Development checks

Go tests, with a current Go toolchain (CI uses 1.27):

```sh
cd clean-api
go test ./...
```

End-to-end smoke check (builds every image, seeds, probes the API and client,
stops the database to prove `/ready` returns 503, restores it, then tears the
stack down):

```sh
./scripts/stack-smoke.sh
```

> The smoke check is destructive to local state: it removes the stack and clears
> `db/data` back to `.gitkeep` when it finishes. Run it when you do not mind
> re-seeding.

CI (`.github/workflows/ci.yml`) runs `go vet` and `go test`, builds the client
bundle, and runs the same stack smoke check on pull requests and pushes to
`master`; it publishes no images.

## Further reading

- [GLOSSARY.md](GLOSSARY.md) — domain nouns and rating vocabulary.
- [docs/adr/0001-rating-contract.md](docs/adr/0001-rating-contract.md) — why
  user ratings are binary while seed ratings stay decimal.
