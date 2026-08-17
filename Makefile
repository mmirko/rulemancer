MINIRULEMANCER_LOCATION := $(shell pwd)/../minirulemancer
include local.mk

all: rulemancer

.PHONY: rulemancer
rulemancer:
	@make -C pkg --no-print-directory all
	@go build
	@./rulemancer build

# If the ROLE is rulemancer then the target populate_minirulemancer can be 
.PHONY: populate_minirulemancer
populate_minirulemancer:
ifeq ($(ROLE), rulemancer)
#	Set the role
	@echo "ROLE=minirulemancer" > $(MINIRULEMANCER_LOCATION)/local.mk
# 	Custom readme
	@cp -a README-MINIRULEMANCER.md $(MINIRULEMANCER_LOCATION)/README.md
# 	Populate directories
	@for d in core; do \
		mkdir -p $(MINIRULEMANCER_LOCATION)/$$d; \
	done
# 	root directory files
	@for file in LICENSE .gitignore minilogo.png Makefile install-clips.sh wrapper.c ; do \
		cp $$file $(MINIRULEMANCER_LOCATION)/$$file; \
	done
else
	$(error "populate is meant to be run in the main repo")
endif
.PHONY: clean
clean:
	@rm -f ./rulemancer
	@rm -rf ./interface
	@make -C docs --no-print-directory clean
	@make -C pkg --no-print-directory clean
