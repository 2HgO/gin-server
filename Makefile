export APP_ENV?=debug
export APP_PORT?=55099
export DB_NAME?=gin-server
export DB_URL?=mongodb://gin-server-db:27017/gin-server
export JWT_SECRET?=8hrh8f93h4urbf82h3i4j2
export GIN_MODE:=${APP_ENV}

.PHONY: run
run: stop-app
	@ \
	docker run -dt --restart=unless-stopped --name=gin-server -p ${APP_PORT}:${APP_PORT} --link gin-server-db \
	-e GIN_MODE=${GIN_MODE} \
	-e APP_ENV=${APP_ENV} \
	-e APP_PORT=${APP_PORT} \
	-e DB_URL=${DB_URL} \
	-e DB_NAME=${DB_NAME} \
	-e JWT_SECRET=${JWT_SECRET} \
	oodemwingie/gin-server:latest &> /dev/null
	@ echo 'App successfully started at http://localhost:${APP_PORT}'

.PHONY: start-db
start-db:
	@ docker start gin-server-db &> /dev/null || docker run -dit --rm --name gin-server-db mongo &> /dev/null
	@ echo 'gin-server-db started successfully'

.PHONY: stop-db
stop-db:
	@ docker stop gin-server-db &> /dev/null | true
	@ echo 'gin-server-db stopped successfully'

.PHONY: kill-db
kill-db:
	@ docker rm gin-server-db &> /dev/null | true
	@ echo 'gin-server-db killed successfully'

.PHONY: build
build:
	@ docker image rm oodemwingie/gin-server:latest &> /dev/null | true
	@ docker build -t oodemwingie/gin-server:latest .

.PHONY: stop
stop: stop-app
	@ echo 'gin-server stopped successfully'

stop-app:
	@ docker stop gin-server &> /dev/null | true
	@ docker rm gin-server &> /dev/null | true

.PHONY: test
test:
	DB_URL='mongodb://localhost:27017/gin-server' go test ./...