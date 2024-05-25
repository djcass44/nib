.PHONY: build

build:
	NIB_DEFAULT_BASE_IMAGE=harbor.dcas.dev/ghcr.io/djcass44/nib/srv-debug go run main.go build ../demo/ --save /tmp/test.tar --v=10
	docker load < /tmp/test.tar