BINARY_NAME = clipr
USER_ENV_FILE = .env
USER_CREDENTIALS = .credentials
COMMIT_HASH = $(shell git rev-parse --short HEAD)
PACKAGE_CLIENT = "github.com/Denis-Kuso/rss_collector/client"
VERSION = "0.1.0"
# Linker flags
LD_FLAGS = "-s -w -X '$(PACKAGE_CLIENT)/cmd.version=$(VERSION)' -X '$(PACKAGE_CLIENT)/cmd.commit=$(COMMIT_HASH)'"
STATIC = 0

build:
	CGO_ENABLED=$(STATIC) go build -ldflags=$(LD_FLAGS) -o $(BINARY_NAME)

env:
	@echo "Creating env file: ${USER_ENV_FILE}\n"
	@touch ./${USER_ENV_FILE}
	@echo "Creating credentials file: ${USER_CREDENTIALS}\n"
	@touch ./${USER_CREDENTIALS}

clean:
	go clean
	rm ${BINARY_NAME}
