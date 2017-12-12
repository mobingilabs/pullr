VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || cat $(CURDIR)/.version 2> /dev/null || echo v0)
BASE = $(CURDIR)
MODULE = pullrapiserver
# BLDVER = version:$(VERSION);build:$(shell date -Ins)
BLDVER = module:$(MODULE),version:$(VERSION),build:$(shell date +"%Y%m%d.%H%M%S.%N.%z")

.PHONY: all $(MODULE)
all: version $(MODULE)

$(MODULE):| $(BASE)
	@go build -v -o $(BASE)/bin/$@

$(BASE):
	@mkdir -p $(dir $@)

# docker builds

.PHONY: custom docker local on off

# The single rule that is called by our root Makefile during CI builds.
custom: docker deploy
local: __local prune

docker:
	@docker build -f dockerfile.apiserver --rm -t $(MODULE):$(TRAVIS_COMMIT) --build-arg awsrgn=ap-northeast-1 --build-arg awsid=$(APISERVER_ACCESS_KEY_ID) --build-arg awssec=$(APISERVER_SECRET_ACCESS_KEY) --build-arg version="$(BLDVER)" .;

deploy:
	@export NAME=$(MODULE); \
	chmod +x $(BASE)/cmd/apiserver/deploy.sh; \
	$(BASE)/cmd/apiserver/deploy.sh;

__local:
	@docker build -f dockerfile.apiserver --rm -t $(MODULE) --build-arg awsrgn=ap-northeast-1 --build-arg awsid=$(APISERVER_ACCESS_KEY_ID) --build-arg awssec=$(APISERVER_SECRET_ACCESS_KEY) --build-arg version="$(BLDVER)" .;

# docker run containers

on: __on prune
off: __off prune

__on:
	@docker run --rm -d -p 8080:8080 --name $(MODULE) $(MODULE)

__off:
	@docker rm -f $(MODULE)

# misc
prune:
	@docker system prune -f

.PHONY: clean version list
clean:
	@rm -rfv bin; \
	docker rmi $(docker images --filter "dangling=true" -q --no-trunc); \
	exit 0

version:
	@echo "Version: $(VERSION)"

# From https://stackoverflow.com/questions/4219255/how-do-you-get-the-list-of-targets-in-a-makefile
list:
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | xargs
