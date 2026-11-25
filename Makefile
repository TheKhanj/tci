VERSION = $(shell git describe --tags --exact-match 2>/dev/null || echo -n dev)
LD_FLAGS = -X 'main.VERSION=$(VERSION)'
ifneq ($(VERSION),dev)
LD_FLAGS += -s -w
endif

all: tci

.PHONY: tci
tci:
	go build \
		-ldflags "$(LD_FLAGS)" \
		-o $@ \
		./cmd
