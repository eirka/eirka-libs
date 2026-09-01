# eirka-libs: library module, no binary. Run `make check` before opening a PR.
# staticcheck is broken against go1.27 export data (v4 > v2), so lint = gofmt + go vet.
.PHONY: build lint preflight test check

build:
	go build ./...

lint:
	@out="$$(gofmt -l $$(go list -f '{{.Dir}}' ./...))"; \
	if [ -n "$$out" ]; then echo "gofmt: unformatted files:" >&2; echo "$$out" >&2; exit 1; fi
	go vet ./...

# redis/ tests exec `redis-server -` via github.com/stvp/tempredis and panic if it is missing.
preflight:
	@command -v redis-server >/dev/null 2>&1 || { \
	  echo "ERROR: redis-server not on PATH; redis/ tests exec it via tempredis (brew install redis)" >&2; exit 1; }
	@echo "preflight ok: $$(command -v redis-server)"

test: preflight
	go test -count=1 ./...

check: build lint test
