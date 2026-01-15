# SERVER_IMAGE=netshield-server
# SERVER_CONTAINER=netshield-server
# AGENT_EXE=agent/bin/shieldagent.exe

# execute:
# 	@echo "▶ Starting NetShield  (Git Bash)..."

# 	@echo "▶ Building server image"
# 	docker build -t $(SERVER_IMAGE) -f server/Dockerfile .

# 	@echo "▶ Starting server container"
# 	-docker rm -f $(SERVER_CONTAINER)
# 	docker run -d \
# 		--name $(SERVER_CONTAINER) \
# 		-p 8082:8082 \
# 		-p 50051:50051 \
# 		$(SERVER_IMAGE)

# 	@echo "▶ Starting agent (Windows native)"
# 	$(AGENT_EXE) &

# 	@echo "▶ Waiting for services to be healthy..."
# 	@until curl -s http://localhost:8082/health > /dev/null; do sleep 1; done
# 	@until curl -s http://localhost:9090/health > /dev/null; do sleep 1; done

# 	@echo "✅ NetShield is READY"
# 	@echo "Server  : http://localhost:8082/status"
# 	@echo "Agent   : http://localhost:9090/current"

# stop:
# 	@echo "⏹ Stopping NetShield"
# 	-docker rm -f $(SERVER_CONTAINER)




# =========================
# NetShield Makefile
# =========================

PROJECT_NAME=netshield
COMPOSE_FILE=deploy/docker-compose.yaml

SERVER_IMAGE=netshield-server
AGENT_EXE=agent/bin/shieldagent.exe

# -------------------------
# Default target
# -------------------------
.PHONY: execute
execute: build up agent wait
	@echo "✅ NetShield is READY"
	@echo "Server  : http://localhost:8082/status"
	@echo "pgAdmin : http://localhost:5050"
	@echo "Agent   : http://localhost:9090/current"

# -------------------------
# Build server image
# -------------------------
.PHONY: build
build:
	@echo "▶ Building NetShield server image"
	docker build -t $(SERVER_IMAGE) -f server/Dockerfile .

# -------------------------
# Start full stack (deploy)
# -------------------------
.PHONY: up
up:
	@echo "▶ Starting NetShield stack (docker-compose)"
	docker compose -f $(COMPOSE_FILE) up -d --build

# -------------------------
# Stop full stack
# -------------------------
.PHONY: down
down:
	@echo "⏹ Stopping NetShield stack"
	docker compose -f $(COMPOSE_FILE) down

# -------------------------
# Clean everything (⚠ destructive)
# -------------------------
.PHONY: clean
clean:
	@echo "🧹 Cleaning NetShield (containers, volumes, images)"
	docker compose -f $(COMPOSE_FILE) down -v
	-docker rmi $(SERVER_IMAGE)

# -------------------------
# Logs
# -------------------------
.PHONY: logs
logs:
	docker compose -f $(COMPOSE_FILE) logs -f

# -------------------------
# Shell into server
# -------------------------
.PHONY: shell
shell:
	docker exec -it netshield-server sh

# -------------------------
# Start Windows agent
# -------------------------
.PHONY: agent
agent:
	@echo "▶ Starting NetShield agent (Windows native)"
	$(AGENT_EXE) &

# -------------------------
# Health wait (server + agent)
# -------------------------
.PHONY: wait
wait:
	@echo "⏳ Waiting for server health..."
	@until curl -sf http://localhost:8082/health > /dev/null; do sleep 1; done

	@echo "⏳ Waiting for agent health..."
	@until curl -sf http://localhost:9090/health > /dev/null; do sleep 1; done

# -------------------------
# Restart server only
# -------------------------
.PHONY: restart-server
restart-server:
	docker compose -f $(COMPOSE_FILE) restart server
