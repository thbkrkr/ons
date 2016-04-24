build:
	docker run --rm \
		-v $$(pwd):/go/src/github.com/thbkrkr/ons \
		-e GOBIN=/go/bin/ \
		-e CGO_ENABLED=0 \
		-e GOPATH=/go \
		-w /go/src/github.com/thbkrkr/ons \
		golang:1.6.0 go build
