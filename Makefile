SERVER_IMAGE=netshield-server
SERVER_CONTAINER=netshield-server
AGENT_EXE=agent/bin/shieldagent.exe

execute:
	@echo "▶ Starting NetShield  (Git Bash)..."

	@echo "▶ Building server image"
	docker build -t $(SERVER_IMAGE) -f server/Dockerfile .

	@echo "▶ Starting server container"
	-docker rm -f $(SERVER_CONTAINER)
	docker run -d \
		--name $(SERVER_CONTAINER) \
		-p 8082:8082 \
		-p 50051:50051 \
		$(SERVER_IMAGE)

	@echo "▶ Starting agent (Windows native)"
	$(AGENT_EXE) &

	@echo "▶ Waiting for services to be healthy..."
	@until curl -s http://localhost:8082/health > /dev/null; do sleep 1; done
	@until curl -s http://localhost:9090/health > /dev/null; do sleep 1; done

	@echo "✅ NetShield is READY"
	@echo "Server  : http://localhost:8082/status"
	@echo "Agent   : http://localhost:9090/current"

stop:
	@echo "⏹ Stopping NetShield"
	-docker rm -f $(SERVER_CONTAINER)
