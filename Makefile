CC=go build
EXECUTABLES=bin/ch01-triangles

default: $(EXECUTABLES)

bin/ch01-triangles: mkbin 01/triangles/main.go 01/triangles/triangles.vert 01/triangles/triangles.frag
	$(CC) -o $@ 01/triangles/main.go

.PHONY: clean
clean:
	rm -f bin/*

.PHONY: mkbin
mkbin:
	mkdir -p bin


