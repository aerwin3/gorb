CC=go build
EXES=bin/ch01-triangles bin/ch03-drawcommands bin/ch03-primitive-restart bin/ch04-gouraud
DEPS=$(shell find ./utils -path '*doc.go' -prune -o -name '*.go' -print)

default : $(EXES)

bin/ch01-triangles: /bin 01/triangles/main.go $(DEPS)
	$(CC) -o $@ 01/triangles/main.go

bin/ch03-drawcommands: /bin 03/drawcommands/main.go $(DEPS)
	$(CC) -o $@ 03/drawcommands/main.go

bin/ch03-primitive-restart: /bin 03/primitive-restart/main.go $(DEPS)
	$(CC) -o $@ 03/primitive-restart/main.go

bin/ch04-gouraud: /bin 04/gouraud/main.go $(DEPS)
	$(CC) -o $@ 04/gouraud/main.go

/bin:
	mkdir -p bin

.PHONY: clean
clean:
	rm -f $(EXES)
