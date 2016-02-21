deps:
	git submodule update --init --recursive

build:
	docker run --rm \
		-v $$(pwd):/dnsovh \
		-e GOBIN=/go/bin/ \
		-e CGO_ENABLED=0 \
		-e GOPATH=/dnsovh:/dnsovh/vendor \
		-w /dnsovh \
		golang:1.6.0 \
		go build

test:
	./go-ovh-dns -f rt-config.json -c show -z blurb.space | jq .
