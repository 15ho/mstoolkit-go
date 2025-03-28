PROJECTDIR:=$(shell pwd)
GOMODULEPREFIX:=github.com/15ho/mstoolkit-go

.PHONY: test
# go test
test:
	@go work edit -json | jq -r '.Use[].DiskPath'  | xargs -I{} go test -v {}

.PHONY: create
# create a new go's module
create:
	@echo "Creating a new go's module..."
	@read -p "Enter module name: " moduleName; \
		mkdir -p $(PROJECTDIR)/$$moduleName; \
		cd $(PROJECTDIR)/$$moduleName; \
		go mod init $(GOMODULEPREFIX)$$moduleName; \
		cd $(PROJECTDIR); \
		go work use ./$$moduleName


# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help