REVISION := $(shell git describe)
REVISION += unknown
REVISION := $(word 1, $(REVISION))

BINARIES := zk-ex
BINARIES_LINUX := $(addsuffix -linux,${BINARIES})
BINARIES_LINUX := $(addsuffix -darwin,${BINARIES})

.PHONY: release
release: ${BINARIES}

.PHONY: all
all: ${BINARIES_LINUX} ${BINARIES_RASP} ${BINARIES_DARWIN}

.PHONY: ${BINARIES}
${BINARIES}:
	go get ./...
	go build -ldflags "-X main.GitRevision=$(REVISION)" -o $(BINARIES)

.PHONY: linux
linux: ${BINARIES_LINUX}

.PHONY: ${BINARIES_LINUX}
${BINARIES_LINUX}:
	GOOS=linux go get ./...
	GOOS=linux go build -ldflags "-X main.GitRevision=$(REVISION)" -o release/$(BINARIES_LINUX)-$(REVISION)

.PHONY: darwin
linux: ${BINARIES_DARWIN}

.PHONY: ${BINARIES_DARWIN}
${BINARIES_DARWIN}:
	GOOS=darwin go get ./...
	GOOS=darwin go build -ldflags "-X main.GitRevision=$(REVISION)" -o release/$(BINARIES_LINUX)-$(REVISION)

.PHONY: clean
clean:
	rm -rf release/*
