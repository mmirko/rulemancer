SHELL := /bin/bash
.PHONY: all
all:
	@re2c --lang go -o game/lex.go game/lex.go.re -i --api simple

.PHONY: clean
clean:
	@rm -f game/lex.go
