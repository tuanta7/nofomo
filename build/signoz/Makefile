.PHONY: install-foundry signoz-up

install-foundry:
	# https://signoz.io/docs/install/docker/
	curl -fsSL https://signoz.io/foundry.sh | bash

signoz-up:
	# Open ports: 8080 (SigNoz UI), 4317 and 4318 (OTLP ingestion)
	echo "Need at least 4GB of memory allocated to Docker"
	docker info --format 'OS: {{.OperatingSystem}} | Memory: {{.MemTotal}} bytes | CPU: {{.NCPU}}'
	foundryctl cast -f casting.yaml