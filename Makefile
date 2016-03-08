CC=go build
EXES = $(shell find . -path ./utils -prune -o -name 'main.go' -print | sed 's/^\.\//ch/g' | sed 's/\/main.go$$//g' | sed 's/\//-/g' | sed 's/^/bin\//g')
DEPS=$(shell find ./utils -path '*doc.go' -prune -o -name '*.go' -print)

default : $(EXES)

# TODO: "{ch}/{example}/main.go" should be included as a dependencie
#       but for some reason the following does not work
#$(EXES) : /bin $(shell echo "$(subst -,/,$(subst bin/ch,,$@))" | sed 's/$$/\/main.go/g') $(DEPS)
$(EXES) : /bin $(DEPS)
	$(CC) -o $@ $(shell echo "$(subst -,/,$(subst bin/ch,,$@))" | sed 's/$$/\/main.go/g')

/bin:
	mkdir -p bin

.PHONY: clean
clean:
	rm -f $(EXES)
