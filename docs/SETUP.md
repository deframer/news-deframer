# Setup

## LLM requirements

- **WARNING**: This might cost money.
- Get an OpenAPI account, Grok, Google Gemini
- Alternative: Install a Local LLM (e.g. LM-Studio) with an OpenAPI API (this will also cost money - in terms of your power bill)
- **Note**: Using commercial LLM APIs (OpenAI, Grok, Google Gemini) usually incurs usage costs.
- **Cloud**: Obtain an API key from your chosen provider.
- **Local**: Alternatively, you can host a local LLM (e.g., using LM-Studio) that exposes an OpenAI-compatible API. While this avoids API fees, it requires at least an Apple M1 or a lower end RTX with 16GB of VRAM.

## Docker Compose

1. Install Docker and Docker Compose.
2. Download the `docker-compose.yml` file.
3. Download `env.example` and save it as `.env`.
4. Edit `.env` to configure your API keys and preferences.
5. Start the services:
   ```bash
   docker compose pull
   docker compose up -d
   ```

## Handling Feeds

- Run `docker compose exec service admin -h`
- Add the example feeds `cat feed-example.json | docker compose exec -T  worker admin feed import`
- List `docker compose exec worker admin feed list`
- Force a Sync `docker compose exec worker admin feed sync-all`
You can manage feeds using the `admin` CLI tool inside the running container.

## RSS Proxy Usage

Once feeds are configured, you can access the proxied versions via the service.

- **URL Pattern**: `http://<your-server-ip>:<port>/rss?url=<upstream-feed-url>`
- **Example**: `http://192.168.1.1:8080/rss?url=https://my-fancy-newssite.com/rss`
- You might need to use a password. Some RSS reader require `https` connections. Use ngingx proxy manager or similar for a certificate.

**Note**: The requested feed URL must be registered in the system (see 'Handling Feeds') before it can be proxied. Allow some time for the background worker to download and process items.

## Browser Plugin

You can install the plugin by loading it as an unpacked extension.

### Option 1: Use a Pre-built Release

1.  Navigate to the project's GitHub Releases page.
2.  Download the `extension.zip` file from the latest release.
3.  Unzip the file's contents into a permanent folder on your computer. **Do not delete this folder**, as the browser loads the extension directly from these files.
4.  Open your Chrome-based browser (e.g., Chrome, Brave, Edge) and go to `chrome://extensions`.
5.  Enable **Developer mode** using the toggle in the top-right corner.
6.  Click **Load unpacked** and select the folder where you unzipped the extension.

### Option 2: Build from Source

1.  Clone this repository.
2.  Navigate into the `browser-plugin` directory.
3.  Install dependencies: `npm ci`.
4.  Build the project: `npm run build`. This creates the necessary files in the `dist/host` directory.
5.  Follow steps 4-6 from "Option 1", but select the `dist/host` directory when you "Load unpacked".
