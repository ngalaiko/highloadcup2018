APP_NAME := highloadcup
BIN_DIR := ./bin
CMD_DIR := ./cmd/${APP_NAME}

BIN := ${BIN_DIR}/${APP_NAME}

build:
	go build -o ${BIN} ${CMD_DIR}/main.go

run-local: build
	${BIN} --data_path=./test_data/data/data.zip
