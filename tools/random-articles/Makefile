UV ?= uv

.PHONY: venv sync run pretty lint typecheck format-check help

help:
	@printf "Targets:\n"
	@printf "  make venv                Create .venv with uv\n"
	@printf "  make sync                Install dependencies with uv\n"
	@printf "  make run                 Fetch articles JSON\n"
	@printf "  make pretty              Fetch articles JSON pretty-printed\n"
	@printf "  make lint                Run ruff, format, mypy\n"

venv:
	$(UV) venv

sync:
	$(UV) sync

run:
	@test -n "$(DSN)" || (printf "DSN is required\n" && exit 1)
	$(UV) run random_articles.py --dsn "$(DSN)"

pretty:
	@test -n "$(DSN)" || (printf "DSN is required\n" && exit 1)
	$(UV) run random_articles.py --dsn "$(DSN)" --pretty

format-check:
	$(UV) run --with ruff ruff format --check random_articles.py

typecheck:
	$(UV) run --with mypy mypy random_articles.py

lint:
	$(UV) run --with ruff ruff check random_articles.py
	$(MAKE) format-check
	$(MAKE) typecheck
