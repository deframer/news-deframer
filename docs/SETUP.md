# Setup

## LLM
## LLM Configuration

- **WARNING**: This might cost money.
- Get an OpenAPI account, Grok, Google Gemini
- Alternative: Install a Local LLM (e.g. LM-Studio) with an OpenAPI API (this will also cost money - in terms of your power bill)
- **Note**: Using commercial LLM APIs (OpenAI, Grok, Google Gemini) usually incurs usage costs.
- **Cloud**: Obtain an API key from your chosen provider.
- **Local**: Alternatively, you can host a local LLM (e.g., using LM-Studio) that exposes an OpenAI-compatible API. While this avoids API fees, it requires significant local hardware resources (GPU/Power).

## Docker Compose

- Install Docker
- Download the [docker compose](../docker/docker-compose.yml) file
- Download the [env-example](../docker/env.example) and save it as `.env`
- Edit the configuration
- Run `docker compose up -d`
1. Install Docker and Docker Compose.
2. Download the `docker-compose.yml` file.
3. Download `env.example` and save it as `.env`.
4. Edit `.env` to configure your API keys and preferences.
5. Start the services:
   ```bash
   docker compose up -d
   ```

## Handling Feeds

- Run `docker compose exec service admin -h`
- Add the example feeds `cat feed-example.json | docker compose exec -T  worker /admin feed import`
- List `docker compose exec worker /admin feed list`
- Force a Sync `docker compose exec worker /admin feed sync-all`
You can manage feeds using the `admin` CLI tool inside the running container.

## RSS Proxy
- **View Help**: `docker compose exec service admin -h`
- **Import Example Feeds**: `cat feed-example.json | docker compose exec -T worker admin feed import`
- **List Feeds**: `docker compose exec worker admin feed list`
- **Force Sync**: `docker compose exec worker admin feed sync-all`

- Pattern: `http://<your-server-ip>:<your-server-port>/rss?url=<url-of-your-newsfeed>`
- Example: `http://192.168.1.1:8080/rss?url=https://my-fancy-newssite/rss`
- Of course the news feed that you want to proxy needs to be part of your feeds.
- Give the AI / worker some time to download and progress it.
## RSS Proxy Usage

## Webbrowser Plugin
Once feeds are configured, you can access the proxied versions via the service.

- TBD
- **URL Pattern**: `http://<your-server-ip>:<port>/rss?url=<upstream-feed-url>`
- **Example**: `http://192.168.1.1:8080/rss?url=https://my-fancy-newssite.com/rss`

**Note**: The requested feed URL must be registered in the system (see 'Handling Feeds') before it can be proxied. Allow some time for the background worker to download and process items.

## Browser Plugin

- *Coming Soon*
