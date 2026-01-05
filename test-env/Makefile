.PHONY: all zap start stop down clean test-env-network-create test-env-network-destroy

all:
	@echo all

zap: down start

test-env-network-create:
	docker network inspect deframer-test-env-net >/dev/null 2>&1 || docker network create deframer-test-env-net

test-env-network-destroy:
	docker network rm deframer-test-env-net || true

start: test-env-network-create
	$(MAKE) -C freshrss $@
	$(MAKE) -C rssbridge $@
	$(MAKE) -C wordpress $@

stop:
	$(MAKE) -C freshrss $@
	$(MAKE) -C rssbridge $@
	$(MAKE) -C wordpress $@

down:
	$(MAKE) -C freshrss $@
	$(MAKE) -C rssbridge $@
	$(MAKE) -C wordpress $@
	$(MAKE) test-env-network-destroy

clean: down
	$(MAKE) -C freshrss $@
	$(MAKE) -C rssbridge $@
	$(MAKE) -C wordpress $@
