CC=go build
EXECUTABLES=bin/ch01-triangles bin/ch03-drawcommands bin/ch03-primitive_restart

default: $(EXECUTABLES)

bin/ch01-triangles: mkbin 01/triangles/main.go
	$(CC) -o $@ 01/triangles/main.go

bin/ch03-drawcommands: mkbin 03/drawcommands/main.go 
	$(CC) -o $@ 03/drawcommands/main.go

bin/ch03-primitive_restart: mkbin 03/primitive_restart/main.go 
	$(CC) -o $@ 03/primitive_restart/main.go

.PHONY: clean
clean:
	rm -f bin/*

.PHONY: mkbin
mkbin:
	mkdir -p bin


