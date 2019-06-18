.PHONY: image
image: flosd
	docker build -t flos .

.PHONY: flosd
flosd: builder
	docker exec -ti builder go build -ldflags '-w -s -extldflags "-static"' -o $@

.PHONY: builder
builder:
	docker run -dit --rm --name builder --volume $(PWD):/src --workdir /src alpine:edge || true
	docker exec -ti builder apk add go gcc git musl-dev

.PHONY: run
run:
	docker run --rm flos
