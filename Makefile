
TARGET 		= Yelaa
DYN_TARGET 	= DynYelaa
SRC 		= main.go

GROUP_ID 	= $$(id -g)
USER_ID 	= $$(id -u)

.PHONY: all
all: compile

.PHONY: help
help:
	@echo "\033[34mYelaa targets:\033[0m"
	@perl -nle'print $& if m{^[a-zA-Z_-\d]+:.*?## .*$$}' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-22s\033[0m %s\n", $$1, $$2}'

.PHONY: compile
compile: ## compile the project
	@go build -o $(TARGET) $(SRC)

.PHONY: docker
docker: ## Builds a docker image from source
	@docker build -t yelaa \
		--build-arg USER_ID=$(USER_ID) \
		--build-arg GROUP_ID=$(GROUP_ID) \
		.

.PHONY: dynamic
dynamic: ## Builds a dynamically linked binary (if you need to use Proxychains with it)
	@./scripts/compile_dyn.sh

.PHONY: clean
clean: ## Cleans up the project
	rm -f $(TARGET) $(DynYelaa)
