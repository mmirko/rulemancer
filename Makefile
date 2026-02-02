
all: rulemancer

.PHONY: rulemancer
rulemancer:
	@make -C pkg --no-print-directory all
	@go build

.PHONY: clean
clean:
	@rm -f ./rulemancer
