.PHONY: all deps test validate lint

all: deps test validate

deps:
	go get -t ./...
	go get -u github.com/golang/lint/golint

test:
	go test -tags experimental -race -cover ./...

validate: lint
	go vet ./...
	test -z "$(gofmt -s -l . | tee /dev/stderr)"

lint:
	out="$$(golint ./...)"; \
	if [ -n "$$(golint ./...)" ]; then \
		echo "$$out"; \
		exit 1; \
	fi

.PHONY: bump
bump: ## Incriment version patch number
	@echo " > Bumping VERSION"
	@hack/bump/version -p $(shell cat VERSION) > VERSION
	@git commit -am "bumping version to $(VERSION)"
	@git push

.PHONY: release
release: bump ## Create a new release from the VERSION
	@echo " > Creating Release"
	@hack/make/release $(shell cat VERSION)

.PHONY: re_release
re_release: ## Create a new release from the VERSION
	@echo " > Recreating Release"
	@hack/make/release ${VERSION}