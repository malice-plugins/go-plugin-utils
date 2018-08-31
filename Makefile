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

_release-patch:
	@echo "version = \"`cat ${VERSION}/__init__.py | awk -F '("|")' '{ print($$2)}' | awk -F. '{$$NF = $$NF + 1;} 1' | sed 's/ /./g'`\"" > VERSION
release-patch: _release-patch git-release build upload current-version

.PHONY: release
release: ## Create a new release from the VERSION
	@echo " > Creating Release"
	@hack/bump/version -p ${shell cat VERSION} > VERSION
	@hack/make/release ${shell cat VERSION}

.PHONY: re_release
re_release: ## Create a new release from the VERSION
	@echo " > Recreating Release"
	@hack/make/release ${VERSION}