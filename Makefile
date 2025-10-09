# StreamForge - Makefile para gestión del ecosistema

.PHONY: help up down logs clean build test lint format

# Colores para output
GREEN=\033[0;32m
YELLOW=\033[1;33m
RED=\033[0;31m
NC=\033[0m # No Color

help: ## Mostrar ayuda
	@echo "$(GREEN)StreamForge - Comandos disponibles:$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}'

up: ## Levantar todo el ecosistema
	@echo "$(GREEN)🚀 Levantando StreamForge...$(NC)"
	docker-compose up -d
	@echo "$(GREEN)✅ StreamForge está corriendo!$(NC)"
	@echo "$(YELLOW)📊 Grafana: http://localhost:3000$(NC)"
	@echo "$(YELLOW)📈 Prometheus: http://localhost:9090$(NC)"
	@echo "$(YELLOW)🔍 Jaeger: http://localhost:16686$(NC)"

down: ## Parar todo el ecosistema
	@echo "$(RED)🛑 Parando StreamForge...$(NC)"
	docker-compose down

logs: ## Ver logs de todos los servicios
	docker-compose logs -f

clean: ## Limpiar todo (containers, volumes, networks)
	@echo "$(RED)🧹 Limpiando StreamForge...$(NC)"
	docker-compose down -v --remove-orphans
	docker system prune -f
	@echo "$(GREEN)✅ Limpieza completada!$(NC)"

build: ## Construir todas las imágenes
	@echo "$(GREEN)🔨 Construyendo imágenes...$(NC)"
	docker-compose build

test: ## Ejecutar tests de todos los proyectos
	@echo "$(GREEN)🧪 Ejecutando tests...$(NC)"
	@for project in projects/*/; do \
		if [ -f "$$project/Makefile" ]; then \
			echo "$(YELLOW)Testing $$project$(NC)"; \
			cd $$project && make test && cd ../..; \
		fi; \
	done

lint: ## Ejecutar linters en todos los proyectos
	@echo "$(GREEN)🔍 Ejecutando linters...$(NC)"
	@for project in projects/*/; do \
		if [ -f "$$project/Makefile" ]; then \
			echo "$(YELLOW)Linting $$project$(NC)"; \
			cd $$project && make lint && cd ../..; \
		fi; \
	done

format: ## Formatear código en todos los proyectos
	@echo "$(GREEN)✨ Formateando código...$(NC)"
	@for project in projects/*/; do \
		if [ -f "$$project/Makefile" ]; then \
			echo "$(YELLOW)Formatting $$project$(NC)"; \
			cd $$project && make format && cd ../..; \
		fi; \
	done

status: ## Mostrar estado de los servicios
	@echo "$(GREEN)📊 Estado de StreamForge:$(NC)"
	docker-compose ps

restart: ## Reiniciar todos los servicios
	@echo "$(YELLOW)🔄 Reiniciando StreamForge...$(NC)"
	docker-compose restart

# Comandos específicos por proyecto
up-project: ## Levantar proyecto específico (uso: make up-project PROJECT=event-bridge-kafka)
	@if [ -z "$(PROJECT)" ]; then \
		echo "$(RED)❌ Especifica el proyecto: make up-project PROJECT=event-bridge-kafka$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)🚀 Levantando $(PROJECT)...$(NC)"
	cd projects/$(PROJECT) && make up

down-project: ## Parar proyecto específico (uso: make down-project PROJECT=event-bridge-kafka)
	@if [ -z "$(PROJECT)" ]; then \
		echo "$(RED)❌ Especifica el proyecto: make down-project PROJECT=event-bridge-kafka$(NC)"; \
		exit 1; \
	fi
	@echo "$(RED)🛑 Parando $(PROJECT)...$(NC)"
	cd projects/$(PROJECT) && make down

# Comandos de desarrollo
dev: ## Modo desarrollo (hot reload)
	@echo "$(GREEN)🔥 Iniciando modo desarrollo...$(NC)"
	docker-compose -f docker-compose.dev.yml up

# Comandos de producción
prod: ## Despliegue en producción
	@echo "$(GREEN)🚀 Desplegando en producción...$(NC)"
	docker-compose -f docker-compose.prod.yml up -d

# Comandos de monitoreo
monitor: ## Abrir dashboards de monitoreo
	@echo "$(GREEN)📊 Abriendo dashboards...$(NC)"
	@echo "$(YELLOW)Grafana: http://localhost:3000$(NC)"
	@echo "$(YELLOW)Prometheus: http://localhost:9090$(NC)"
	@echo "$(YELLOW)Jaeger: http://localhost:16686$(NC)"
	@open http://localhost:3000 || echo "Abre manualmente: http://localhost:3000"

# Comandos de base de datos
db-migrate: ## Ejecutar migraciones de base de datos
	@echo "$(GREEN)🗄️ Ejecutando migraciones...$(NC)"
	@for project in projects/*/; do \
		if [ -f "$$project/migrations" ]; then \
			echo "$(YELLOW)Migrating $$project$(NC)"; \
			cd $$project && make migrate && cd ../..; \
		fi; \
	done

# Comandos de seguridad
security-scan: ## Escanear vulnerabilidades
	@echo "$(GREEN)🔒 Escaneando vulnerabilidades...$(NC)"
	docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \
		aquasec/trivy image streamforge/event-bridge-kafka:latest

# Comandos de backup
backup: ## Crear backup de datos
	@echo "$(GREEN)💾 Creando backup...$(NC)"
	mkdir -p backups/$(shell date +%Y%m%d_%H%M%S)
	docker-compose exec postgres pg_dump -U streamforge > backups/$(shell date +%Y%m%d_%H%M%S)/postgres.sql

# Comandos de logs
logs-project: ## Ver logs de proyecto específico (uso: make logs-project PROJECT=event-bridge-kafka)
	@if [ -z "$(PROJECT)" ]; then \
		echo "$(RED)❌ Especifica el proyecto: make logs-project PROJECT=event-bridge-kafka$(NC)"; \
		exit 1; \
	fi
	docker-compose logs -f $(PROJECT)
