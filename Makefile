deps:
	git submodule update --init --recursive

build:
	docker run --rm \
		-v $$(pwd):/go/src/github.com/thbkrkr/go-ovh-dns \
		-e GOBIN=/go/bin/ \
		-e CGO_ENABLED=0 \
		-e GOPATH=/go \
		-w /go/src/github.com/thbkrkr/go-ovh-dns \
		golang:1.6.0 \
		go build

test:
	./go-ovh-dns -f rt-config.json -c show -z blurb.space | jq .
