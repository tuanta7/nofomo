setup-local:
	docker compose -f ./build/docker-compose.yml up -d

backtest:
	go run ./cmd/backtest