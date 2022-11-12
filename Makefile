POSTGRESQL_HOST = 0.0.0.0
POSTGRESQL_PORT = 8603
REDIS_HOST = 0.0.0.0
REDIS_PORT = 8601

.PHONY: generate
generate:
	sqlc generate

.PHONY: mock
mock:
	@echo "[SERVICE] mock post-feed..."
	mockgen -destination ./service/postfeed/mock/poster.go -package mock_postfeed -source ./service/postfeed/interface.go
	@echo "[PROVIDER] mock post-store..."
	mockgen -destination ./provider/poststore/mock/mocked.go -package mock_poststore -source ./provider/poststore/querier.go
	@echo "[PROVIDER] mock ad-store..."
	mockgen -destination ./provider/adstore/mock/mocked.go -package mock_adstore -source ./provider/adstore/interface.go
	@echo "[PROVIDER] mock feed-store..."
	mockgen -destination ./provider/feedstore/mock/mocked.go -package mock_feedstore -source ./provider/feedstore/interface.go

.PHONY: swagger
swagger:
	@echo "[SWAGGER] Generating swagger documentation..."
	swag i --parseDependency --parseInternal --parseDepth 1 -g ./server/rest/handler/router.go -o ./docs/

.PHONY: install
install:
	@echo "installing linter..."
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s latest

.PHONY: services_up
services_up:
	docker-compose -f docker-services.yml up --bu -d

.PHONY: api_up
api_up:
	docker-compose -f docker-compose.yml up --bu -d

.PHONY: unit_test
unit_test:
	@echo "[SERVER] testing handlers..."
	go test -v ./server/rest/handler/*.go
	@echo "\n[SERVICE] testing service..."
	go test -v ./service/postfeed/*.go
	@echo "\n[PROVIDER] testing provider..."
	go test -v ./provider/adstore/*.go

.PHONY: integration_test
integration_test: env_echo_main
	@echo "[REDIS][AD-STORE] testing adstore..."
	POSTGRESQL_HOST=${POSTGRESQL_HOST} POSTGRESQL_PORT=${POSTGRESQL_PORT} REDIS_HOST=${REDIS_HOST} REDIS_PORT=${REDIS_PORT} go test -v tests/integration/provider/adstore/ads_test.go
	@echo "\n[REDIS][FEED-STORE] testing feedstore..."
	POSTGRESQL_HOST=${POSTGRESQL_HOST} POSTGRESQL_PORT=${POSTGRESQL_PORT} REDIS_HOST=${REDIS_HOST} REDIS_PORT=${REDIS_PORT} go test -v tests/integration/provider/feedstore/feed_test.go
	@echo "\n[POSTGRES][POST-STORE] testing poststore..."
	POSTGRESQL_HOST=${POSTGRESQL_HOST} POSTGRESQL_PORT=${POSTGRESQL_PORT} REDIS_HOST=${REDIS_HOST} REDIS_PORT=${REDIS_PORT} go test -v tests/integration/provider/poststore/poststore_test.go

.PHONY: env_echo_main
env_echo_main:
	@echo "[ENV] setting env variables..."
	@echo POSTGRESQL_HOST=$(POSTGRESQL_HOST) POSTGRESQL_PORT=$(POSTGRESQL_PORT) REDIS_HOST=$(REDIS_HOST) REDIS_PORT=$(REDIS_PORT)
