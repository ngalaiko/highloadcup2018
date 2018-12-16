APP_NAME := highloadcup
APP_DIR := ./app
BIN_DIR := ${APP_DIR}/bin
CMD_DIR := ${APP_DIR}/cmd/${APP_NAME}
DOCKER_REPO := stor.highloadcup.ru/accounts/cat_winner

BIN := ${BIN_DIR}/${APP_NAME}

local-build:
	go build -o ${BIN} ${CMD_DIR}/main.go

local-run: local-build
	${BIN} --data_path=./test_data/data/data.zip --profile_addr=:8080

local-profile-heap: 
	go tool pprof ${BIN} http://127.0.0.1:8080/debug/pprof/heap

local-profile-cpu: 
	go tool pprof ${BIN} http://127.0.0.1:8080/debug/pprof/profile?seconds=10

docker-run: docker-build
	cp -R ./test_data/data /tmp/data
	docker run --rm -it \
		-v /tmp/data:/tmp/data \
		-p 80:80 \
		${DOCKER_REPO} 

docker-build:
	docker build ${APP_DIR} -t ${DOCKER_REPO}

docker-push: docker-build
	docker push ${DOCKER_REPO}
