UV ?= uv

.PHONY: venv sync validate run dry-run apply join-feed lint typecheck format-check help

help:
	@printf "Targets:\n"
	@printf "  make venv                Create .venv with uv\n"
	@printf "  make sync                Install dependencies with uv\n"
	@printf "  make validate            Compare source/destination schema\n"
	@printf "  make dry-run             Run merge in dry-run mode\n"
	@printf "  make apply               Run merge with --apply\n"
	@printf "  make join-feed           Run join for one FEED_URL\n"
	@printf "  make lint                Run ruff, format, mypy\n"

venv:
	$(UV) venv

sync:
	$(UV) sync

validate:
	@test -n "$(SOURCE_DSN)" || (printf "SOURCE_DSN is required\n" && exit 1)
	@test -n "$(DEST_DSN)" || (printf "DEST_DSN is required\n" && exit 1)
	$(UV) run db_joiner.py --source-dsn "$(SOURCE_DSN)" --dest-dsn "$(DEST_DSN)" --validate-schema

dry-run:
	@test -n "$(SOURCE_DSN)" || (printf "SOURCE_DSN is required\n" && exit 1)
	@test -n "$(DEST_DSN)" || (printf "DEST_DSN is required\n" && exit 1)
	$(UV) run db_joiner.py --source-dsn "$(SOURCE_DSN)" --dest-dsn "$(DEST_DSN)"

apply:
	@test -n "$(SOURCE_DSN)" || (printf "SOURCE_DSN is required\n" && exit 1)
	@test -n "$(DEST_DSN)" || (printf "DEST_DSN is required\n" && exit 1)
	$(UV) run db_joiner.py --source-dsn "$(SOURCE_DSN)" --dest-dsn "$(DEST_DSN)" --apply

join-feed:
	@test -n "$(SOURCE_DSN)" || (printf "SOURCE_DSN is required\n" && exit 1)
	@test -n "$(DEST_DSN)" || (printf "DEST_DSN is required\n" && exit 1)
	@test -n "$(FEED_URL)" || (printf "FEED_URL is required\n" && exit 1)
	$(UV) run db_joiner.py --source-dsn "$(SOURCE_DSN)" --dest-dsn "$(DEST_DSN)" --join --feed-url "$(FEED_URL)"

format-check:
	$(UV) run --with ruff ruff format --check db_joiner.py

typecheck:
	$(UV) run --with mypy mypy db_joiner.py

lint:
	$(UV) run --with ruff ruff check db_joiner.py
	$(MAKE) format-check
	$(MAKE) typecheck
