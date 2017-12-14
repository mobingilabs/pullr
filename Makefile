# This included makefile should define the 'custom' target rule which is called here.
include $(INCLUDE_MAKEFILE)

.PHONY: release
release: custom 

up:
	docker build -f dockerfile.apiserver --rm -t pullrapiserver --build-arg awsrgn=ap-northeast-1 --build-arg awsid=$(APISERVER_ACCESS_KEY_ID) --build-arg awssec=$(APISERVER_SECRET_ACCESS_KEY) --build-arg version="local" .
	docker-compose up -d

down:
	docker-compose down
	docker system prune -f
