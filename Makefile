.PHONY: check
check: ## check
	@go generate ./...
	@go fmt ./...
	@go vet ./...
	go test ./...

deploy:
	gcloud functions deploy SlashFunction --runtime go111 --trigger-http

describe:
	gcloud functions describe SlashFunction

.PHONY: help
help:	## show this help
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'