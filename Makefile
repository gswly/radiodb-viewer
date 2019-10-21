
GO_BASE_IMAGE = amd64/golang:1.12-alpine3.10

help:
	@echo "usage: make [action]"
	@echo ""
	@echo "available actions:"
	@echo ""
	@echo "  mod-tidy"
	@echo "  format"
	@echo "  dev"
	@echo "  prod"
	@echo ""

mod-tidy:
	echo "FROM $(GO_BASE_IMAGE) \n\
	RUN apk add git \n\
	" | docker build - -t rdbviewer-modtidy
	docker run --rm -it -v $(PWD):/s rdbviewer-modtidy \
	sh -c "cd /s/back && go get -m ./... && go mod tidy"

format:
	docker run --rm -it -v $(PWD):/s $(GO_BASE_IMAGE) \
	sh -c "cd /s && find . -type f -name '*.go' | xargs gofmt -l -w -s"

BUILD = docker build . \
    --build-arg BUILD_MODE=$(BUILD_MODE) \
    -t radiodb-viewer-$(BUILD_MODE)

dev:
	@command -v inotifywait >/dev/null 2>&1 || { echo "inotifywait is required"; exit 1; }
	docker run --rm -it -v radiodb:/out amd64/alpine:3.8 \
	sh -c "apk add curl && curl --compressed -o/out/radiodb.json https://radiodb.freeddns.org/dumpget"
	make dev-inner

dev-inner:
	docker kill radiodb-viewer-dev >/dev/null 2>&1 || exit 0

	$(eval BUILD_MODE = dev)
	$(BUILD)

	docker run --rm \
	--read-only \
	-v radiodb:/data \
	-p 7446:7446 \
	--name radiodb-viewer-dev \
	radiodb-viewer-dev &

	inotifywait -qre close_write ./*/
	make dev-inner

prod:
	$(eval BUILD_MODE = prod)
	$(BUILD)
