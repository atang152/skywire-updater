skywire-version = testing

lint: ## Run linters. Use make install-linters first.
#	vendorcheck ./...
	golangci-lint run -c .golangci.yml ./...
	# The govet version in golangci-lint is out of date and has spurious warnings, run it separately
	go vet -all ./...

install-linters: ## Install linters
	go get -u github.com/FiloSottile/vendorcheck
	# For some reason this install method is not recommended, see https://github.com/golangci/golangci-lint#install
	# However, they suggest `curl ... | bash` which we should not do
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

format: ## Formats the code. Must have goimports installed (use make install-linters).
	goimports -w -local github.com/watercompany/skywire-services ./...

test: ## Run tests for net
	@mkdir -p coverage/
	go test -coverpkg="github.com/watercompany/skywire-services/..." -coverprofile=coverage/go-test-cmd.coverage.out -timeout=5m ./... -race
	#go test -coverpkg="github.com/watercompany/skywire-messaging/..." -coverprofile=coverage/go-test-cmd.coverage.out -timeout=5m ./pkg/... -race
	# race testing is taking around one minute and a half in my computer