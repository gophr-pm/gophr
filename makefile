.PHONY: test-cover-html
PACKAGES = $(shell find ./ -type d -not -path '*/\.*')

test-cover-html:
	echo "mode: count" > coverage-all.out
	$(foreach pkg,$(PACKAGES),\
		go test -coverprofile=coverage.out -covermode=count $(pkg);\
		tail -n +2 coverage.out >> coverage-all.out;)
	$HOME/gopath/bin/goveralls -coverprofile=coverage-all.out -service=travis-ci -repotoken $COVERALLS_TOKEN
