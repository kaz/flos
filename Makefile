.PHONY: flos
flos:
	GOOS=linux go build -ldflags '-w -s -extldflags "-static"' -o $@

.PHONY: flos.exe
flos.exe:
	GOOS=windows go build -ldflags '-w -s -extldflags "-static"' -o $@
