SHELL:=/bin/bash

## Development time

all:

python-env:
	mkdir -p .python-venv
	python3 -m venv .python-venv
	. .python-venv/bin/activate && \
		python3 -m pip install --quiet --upgrade pip && \
		pip install --quiet -r requirements.txt && \
		playwright install
	@echo always execute . .python-venv/bin/activate
