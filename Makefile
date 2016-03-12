CC=go build
EXES=bin/ch01-triangles bin/ch03-drawcommands bin/ch03-primitive-restart bin/ch04-gouraud bin/ch04-shadowmap

default : $(EXES)

bin/ch01-triangles: /bin 01/triangles/main.go
	$(CC) -o $@ 01/triangles/main.go

bin/ch03-drawcommands: /bin 03/drawcommands/main.go
	$(CC) -o $@ 03/drawcommands/main.go

bin/ch03-primitive-restart: /bin 03/primitive-restart/main.go
	$(CC) -o $@ 03/primitive-restart/main.go

bin/ch04-gouraud: /bin 04/gouraud/main.go
	$(CC) -o $@ 04/gouraud/main.go

bin/ch04-shadowmap: /bin 04/shadowmap/main.go
	$(CC) -o $@ 04/shadowmap/main.go

/bin:
	mkdir -p bin

.PHONY: clean
clean:
	rm -f $(EXES)
