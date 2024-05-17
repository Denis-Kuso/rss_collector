BINARY_NAME=clirs
USER_ENV_FILE=.env
USER_CREDENTIALS=.credentials
 
all: build test

build:
	go build -o ${BINARY_NAME} main.go
	@echo "Creating env file: ${USER_ENV_FILE}\n"
	@touch ./${USER_ENV_FILE}
	@echo "Creating credentials file: ${USER_CREDENTIALS}\n"
	@touch ./${USER_CREDENTIALS}
 
run:
	go build -o ${BINARY_NAME} main.go
	./${BINARY_NAME}

test:
	@echo "testing app...\n"
 
clean:
	go clean
	rm ${BINARY_NAME} ${USER_ENV_FILE} ${USER_CREDENTIALS} 
