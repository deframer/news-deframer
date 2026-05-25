PYTHON ?= python3
UV ?= uv

.PHONY: venv sync validate run dry-run apply join-feed lint help

help:
	@printf "Targets:\n"
	@printf "  make venv                Create .venv with uv\n"
	@printf "  make sync                Install dependencies with uv\n"
	@printf "  make validate            Compare source/destination schema\n"
	@printf "  make dry-run             Run merge in dry-run mode\n"
	@printf "  make apply               Run merge with --apply\n"
	@printf "  make join-feed           Run join for one FEED_URL\n"
	@printf "  make lint                Syntax check script\n"

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

lint:
	$(PYTHON) -m py_compile db_joiner.py
