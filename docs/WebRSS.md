For a self-hosted RSS reader that runs in Docker and is fully accessible via a browser, **FreshRSS** and **Miniflux** are the top maintained recommendations for 2025.

Here is a breakdown of the best options, their maintenance status, and quick Docker Compose configurations to get you started.

### 1. FreshRSS (Best All-Rounder)
FreshRSS is the most balanced choice. It is actively maintained, feature-rich, and has a responsive web interface that works great on both desktop and mobile browsers.

*   **Best for:** Users who want a modern interface, themes, and extensions (e.g., to scrape full text from websites that only provide summaries).
*   **Maintenance:** Very active (Frequent releases in late 2024/2025).
*   **Key Features:**
    *   Mobile-friendly web UI (PWA support).
    *   Supports multiple users.
    *   Compatible with "Google Reader" API clients (Reeder, NetNewsWire, etc.).
    *   Lightweight resource usage.

**Quick Start (Docker Compose):**
Save this as `docker-compose.yml` and run `docker compose up -d`.

```yaml
services:
  freshrss:
    image: lscr.io/linuxserver/freshrss:latest
    container_name: freshrss
    environment:
      - PUID=1000
      - PGID=1000
      - TZ=Europe/Berlin
    volumes:
      - ./data:/config
    ports:
      - 8080:80
    restart: unless-stopped
```
*Access it at `http://localhost:8080`.*


```bash
# docker run -d \
#   --name=freshrss \
#   -e PUID=1000 \
#   -e PGID=1000 \
#   -e TZ=Europe/Berlin \
#   -p 8080:80 \
#   -v $(pwd)/data:/config \
#   --restart unless-stopped \
#   lscr.io/linuxserver/freshrss:latest

# https://github.com/FreshRSS/FreshRSS/blob/edge/Docker/README.md
docker run -d --restart unless-stopped --log-opt max-size=10m \
  -p 8080:80 \
  -e TZ=Europe/Paris \
  -e 'CRON_MIN=1,31' \
  -v freshrss_data:/var/www/FreshRSS/data \
  -v freshrss_extensions:/var/www/FreshRSS/extensions \
  --name freshrss \
  freshrss/freshrss
```

---

### 2. Miniflux (Best Minimalist Choice)
Miniflux is an "opinionated," minimalist reader written in Go. It strips away all the clutter (no social sharing, no fancy dashboards) and focuses purely on reading text.

*   **Best for:** Users who want speed, simplicity, and a distraction-free reading environment.
*   **Maintenance:** Extremely active; developer updates it frequently.
*   **Key Features:**
    *   Incredibly fast and lightweight.
    *   Clean, simple text-focused UI.
    *   Support for Fever and Miniflux APIs (for third-party apps).
    *   Built-in "Fetch original content" feature to grab full articles automatically.

**Quick Start (Docker Compose):**
Miniflux requires a PostgreSQL database.

```yaml
services:
  miniflux:
    image: miniflux/miniflux:latest
    container_name: miniflux
    ports:
      - "8081:8080"
    environment:
      - DATABASE_URL=postgres://miniflux:secret@db/miniflux?sslmode=disable
      - RUN_MIGRATIONS=1
      - CREATE_ADMIN=1
      - ADMIN_USERNAME=admin
      - ADMIN_PASSWORD=password
    depends_on:
      - db
    restart: unless-stopped

  db:
    image: postgres:15
    container_name: miniflux-db
    environment:
      - POSTGRES_USER=miniflux
      - POSTGRES_PASSWORD=secret
      - POSTGRES_DB=miniflux
    volumes:
      - miniflux-db:/var/lib/postgresql/data
    restart: unless-stopped

volumes:
  miniflux-db:
```
*Access it at `http://localhost:8081`. Login with `admin` / `password`.*

---

### 3. Tiny Tiny RSS (TT-RSS) (Best for Power Users)
TT-RSS is the "Swiss Army Knife" of RSS readers. It has been around for a long time and has a massive plugin ecosystem.

*   **Best for:** Power users who want to script filters, install complex plugins, or heavily customize the backend logic.
*   **Maintenance:** Continuous (Rolling releases via Git).
*   **Caveat:** The official setup is more complex than the others (often involves 3-4 containers). It is widely considered "feature-complete" but the developer can be strict about support.
*   **Deployment:** The official Docker setup is slightly involved. A popular community fork/image by **LinuxServer.io** (similar to the FreshRSS example) is often easier for beginners to deploy.

### 4. CommaFeed (Google Reader Nostalgia)
If you miss the old Google Reader interface from ~2013, CommaFeed is built specifically to mimic that look and feel.

*   **Best for:** Users who want a classic, familiar "Inbox" style interface.
*   **Maintenance:** Good (Steady updates).
*   **Key Features:** Bloat-free, supports RedHat/Postgres/MySQL/H2, very simple to set up.

### Summary Comparison

| Feature | FreshRSS | Miniflux | TT-RSS |
| :--- | :--- | :--- | :--- |
| **Interface** | Modern, Themable | Text-Only, Minimalist | Functional, Dense |
| **Speed** | Fast | Very Fast | Average |
| **Mobile Web** | Excellent | Good | Okay (Apps preferred) |
| **Database** | SQLite (Default), PgSQL | PostgreSQL Only | PostgreSQL Only |
| **Difficulty** | Easy | Easy | Medium |

**Recommendation:** Start with **FreshRSS**. It offers the easiest "out of the box" experience with SQLite (no separate database container needed) and looks great on a browser immediately.