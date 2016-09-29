.PHONY: test-cover
PACKAGES = $(shell find ./ -type d -not -path '*/\.*')

test-cover:
	echo "mode: count" > coverage-all.out
	$(foreach pkg,$(PACKAGES),\
		go test -coverprofile=coverage.out -covermode=count $(pkg);\
		tail -n +2 coverage.out >> coverage-all.out;)
