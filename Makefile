run:
	docker-compose up --build

tests-phase-1:
	docker-compose \
		-f ./docker-compose.phase.1.yaml \
		-f ./docker-compose.yaml \
		up --build
