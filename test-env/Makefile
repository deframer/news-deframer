.PHONY: all zap start stop down clean test-env-net-create test-env-net-destroy

all:
	@echo all

zap: down start

test-env-net-create:
	docker network inspect deframer-test-env-net >/dev/null 2>&1 || docker network create deframer-test-env-net

test-env-net-destroy:
	docker network rm deframer-test-env-net || true

start: test-env-net-create
	$(MAKE) -C freshrss $@
	$(MAKE) -C rssbridge $@
	$(MAKE) -C wordpress $@
	$(MAKE) -C rsshub $@

stop:
	$(MAKE) -C freshrss $@
	$(MAKE) -C rssbridge $@
	$(MAKE) -C wordpress $@
	$(MAKE) -C rsshub $@

down:
	$(MAKE) -C freshrss $@
	$(MAKE) -C rssbridge $@
	$(MAKE) -C wordpress $@
	$(MAKE) -C rsshub $@

clean: down
	$(MAKE) -C freshrss $@
	$(MAKE) -C rssbridge $@
	$(MAKE) -C wordpress $@
	$(MAKE) -C rsshub $@
	$(MAKE) test-env-net-destroy
