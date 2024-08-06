FRONT_END_BINARY=frontApp
ACCCOUNT_BINARY=accountApp
LOGGER_BINARY=loggerServiceApp
MAIL_BINARY=mailerApp
USER_BINARY=userApp
LISTENER_BINARY=listenerApp

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build:  build_account build_user build_listener build_logger build_mail
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## build_account: builds the account binary as a linux executable
build_account:
	@echo "Building of account service started"
	cd ../account-service && env GOOS=linux CGO_ENABLED=0 go build -o ${ACCCOUNT_BINARY} ./cmd/api
	@echo "Building of account service finished\n"

## build_user: builds the user binary as a linux executable
build_user:
	@echo "Building of user service started"
	cd ../user-service && env GOOS=linux CGO_ENABLED=0 go build -o ${USER_BINARY} ./cmd/api
	@echo "Building of user service finished\n"

## build_listener: builds the listener binary as a linux executable
build_listener:
	@echo "Building of listener service started"
	cd ../listener-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LISTENER_BINARY} .
	@echo "Building of listener service finished\n"

## build_logger: builds the logger binary as a linux executable
build_logger:
	@echo "Building of logger service started"
	cd ../logger-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LOGGER_BINARY} ./cmd/api
	@echo "Building of logger service finished\n"

## build_mail: builds the mail binary as a linux executable
build_mail:
	@echo "Building of mail service started"
	cd ../mail-service && env GOOS=linux CGO_ENABLED=0 go build -o ${MAIL_BINARY} ./cmd/api
	@echo "Building of mail service finished\n"
