.PHONY: image
image: flos
	docker build -t flos .

.PHONY: flos
flos: builder
	docker exec -ti builder go build -ldflags '-w -s -extldflags "-static"' -o $@

.PHONY: builder
builder:
	docker run -dit --name builder --volume $(PWD):/src --workdir /src alpine:edge || true
	docker exec -ti builder apk add go gcc git musl-dev || true

.PHONY: run
run:
	docker run --rm flos
