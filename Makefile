.PHONY: image
image: flos
	docker build -t flos .

.PHONY: flos
flos:
	GOOS=linux go build -ldflags '-w -s -extldflags "-static"' -o $@

.PHONY: run
run:
	docker run --rm flos
