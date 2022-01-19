
.PHONY: all
all: install compile

.PHONY: help
help:
	@echo "\033[34mYelaa targets:\033[0m"
	@perl -nle'print $& if m{^[a-zA-Z_-\d]+:.*?## .*$$}' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-22s\033[0m %s\n", $$1, $$2}'

.PHONY: install
install: ## install necessary dependencies
	go install -v github.com/projectdiscovery/nuclei/v2/cmd/nuclei@latest
	# go install -v github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest
	# go install -v github.com/projectdiscovery/dnsx/cmd/dnsx@latest
	# go install -v github.com/projectdiscovery/httpx/cmd/httpx@latest

.PHONY: compile
compile: ## compile the project
	go build -o Yelaa main.go
