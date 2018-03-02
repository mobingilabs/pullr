DIST:=bin
DOCKER_TAG_PREFIX:=mobingilabs/pullr-
DOCKER_REGISTRY:=DOCKER_REGISTRY
VERSION:=localdev

CMDS:=apisrv buildctl trackersvc
LINUX_CMDS:=$(addsuffix -linux,$(CMDS))
DOCKER_CMDS:=$(addsuffix -docker,$(CMDS))
PUSH_CMDS:=$(addsuffix -push,$(CMDS)) ui-push

.SECONDEXPANSION:

.PHONE: show
show:
	@echo Variables:
	@echo ============================================================
	@echo "DIST:              $(DIST)"
	@echo "VERSION:           $(VERSION)"
	@echo "CMDS:              $(CMDS)"
	@echo "DOCKER_TAG_PREFIX: $(DOCKER_TAG_PREFIX)"
	@echo
	@echo Services:
	@echo ============================================================
	@echo apisrv, buildctl, ui
	@echo
	@echo Tasks:
	@echo ============================================================
	@echo "helm               Prepares helm package"
	@echo "build              Build all services"
	@echo "build-linux        Build all services for linux"
	@echo "dev                Run pullr for development"
	@echo "docker             Build docker images for all services"
	@echo "                   arguments:"
	@echo "                     VERSION=<imageversion>"
	@echo "push               Pushes all the docker images to given registry"
	@echo "                   arguments:"
	@echo "                     DOCKER_REGISTRY=<registryurl>"
	@echo "                     VERSION=<imageversion>"
	@echo "[svcname]          Build only given service"
	@echo "[svcname]-linux    Build only given service for linux"
	@echo "[svcname]-docker   Build docker image only for given service"
	@echo "                   arguments:"
	@echo "                     VERSION=<imageversion>"
	@echo "[svcname]-push     Pushes docker image only for given service"
	@echo "                   arguments:"
	@echo "                     DOCKER_REGISTRY=<registryurl>"
	@echo "                     VERSION=<imageversion>"

.PHONY: build build-linux
build: $(CMDS) ui
build-linux: $(LINUX_CMDS) ui
docker: $(DOCKER_CMDS) ui-docker
push: $(PUSH_CMDS)

.PHONY: dev show-dev-help
dev:
	docker-compose -f docker-compose.yaml -f docker-compose.dev.yaml up -d
	@make show-dev-help
	@sh -c "cd ui; npm install && cd dist && parcel serve --https -d . ../src/index.html"

show-dev-help:
	@echo "If this is your first run please add pullr's names"
	@echo "$ sudo echo 127.0.0.1 pullr.dev reg.pullr.dev >> /etc/hosts"
	@echo "pullr is running at https://pullr.local"

.PHONY: push
pushcmd=$(patsubst %-push,%,$@)
$(PUSH_CMDS):
	docker tag $(DOCKER_TAG_PREFIX)$(pushcmd):$(VERSION) $(DOCKER_REGISTRY)/$(DOCKER_TAG_PREFIX)$(pushcmd):$(VERSION)
	docker push $(DOCKER_REGISTRY)/$(DOCKER_TAG_PREFIX)$(pushcmd):$(VERSION)

.PHONY: $(CMDS) $(LINUX_CMDS) $(DOCKER_CMDS)
cmd=$(word 1, $@)
$(CMDS):
	go build -o $(DIST)/$(cmd) ./cmd/$(cmd)

linuxcmd=$(patsubst %-linux,%,$@)
$(LINUX_CMDS):
	GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o $(DIST)/$(linuxcmd) ./cmd/$(linuxcmd)

dockercmd=$(patsubst %-docker,%,$@)
dockerdep=$(patsubst %-docker,%-linux,$@)
$(DOCKER_CMDS): % : $$(dockerdep)
	docker build -t $(DOCKER_TAG_PREFIX)$(dockercmd):$(VERSION) -f cmd/$(dockercmd)/Dockerfile .

.PHONY: ui ui-docker
ui:
	- mkdir -p ui/dist
	cd ui/dist; parcel build -d . ../src/index.html
ui-docker: ui
	docker build -t $(DOCKER_TAG_PREFIX)ui:$(VERSION) -f ui/Dockerfile .

