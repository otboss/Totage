# TODO:
# WORK IN PROGRESS
# THIS FILE IS NOT COMPLETED AND MAY NOt work as expected

# Ensure this file runs using a clean shell environment and stops on any error
.SHELLFLAGS = -ec

# Automatically load and export your existing .env file 
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# ⚠️ UPDATE THIS: The name of your actual live running PostGIS Docker container
LIVE_CONTAINER_NAME = postgis_database

# Configuration for the temporary schema analyzer engine
SHADOW_CONTAINER = postgres_shadow_diff_engine
SHADOW_PORT = 54321

.DEFAULT_GOAL := help

.PHONY: help run

help:
	@echo "================================================================="
	@echo "📋 INTERACTIVE BOUNDARY-SAFE POSTGIS MIGRATION UTILITY"
	@echo "================================================================="
	@echo "Available commands:"
	@echo "  make -f migrations.mk run      Validate schemas, preview ALTERs, & sync safely"
	@echo "================================================================="

run:
	@echo "🔍 Step 1: Pre-validating queries and schemas with SQLC..."
	@sqlc generate -f ./sqlc.yml
	@echo "✅ Validation successful!"

	@echo "🏗️ Step 2: Spinning up ephemeral open-source schema diff environment..."
	@docker rm -f $(SHADOW_CONTAINER) > /dev/null 2>&1 || true
	@docker run --rm -d \
		--name $(SHADOW_CONTAINER) \
		-p $(SHADOW_PORT):5432 \
		-e POSTGRES_PASSWORD=postgres \
		postgis/postgis:17-3.5 > /dev/null
	@until docker exec $(SHADOW_CONTAINER) pg_isready -U postgres >/dev/null 2>&1; do sleep 0.5; done

	@echo "📦 Step 3: Loading domain layout files into analyzer engine..."
	# 1. Initialize PostGIS extension on the analyzer
	@docker exec -i $(SHADOW_CONTAINER) psql -U postgres -d postgres -c "CREATE EXTENSION IF NOT EXISTS postgis;"
	
	# 2. Extract and pre-create all schema namespaces exactly as defined in your code
	@grep -ih "CREATE SCHEMA" shared/domain/entities/*.schema.sql modules/*/domain/entities/*.schema.sql 2>/dev/null | \
		docker exec -i $(SHADOW_CONTAINER) psql -U postgres -d postgres || true

	# 3. Stream all layouts dynamically into your transaction block arrangement
	@( \
		echo "BEGIN;"; \
		cat "shared/domain/entities/shared.schema.sql" 2>/dev/null || true; \
		cat modules/user/domain/entities/*.schema.sql 2>/dev/null || true; \
		cat modules/payment/domain/entities/*.schema.sql 2>/dev/null || true; \
		cat modules/organization/domain/entities/*.schema.sql 2>/dev/null || true; \
		cat modules/logistic/domain/entities/*.schema.sql 2>/dev/null || true; \
		cat modules/order/domain/entities/*.schema.sql 2>/dev/null || true; \
		echo "COMMIT;" \
	) | docker exec -i $(SHADOW_CONTAINER) psql -U postgres -d postgres -v ON_ERROR_STOP=1

	@echo "⚙️ Step 4: Generating and reviewing ALTER statements natively..."
	$(eval SSL_MODE := $(if $(filter true,$(DATABASE_SSL_ENABLED)),require,disable))
	$(eval LIVE_DB_URL := postgres://$(DATABASE_USER):$(DATABASE_PASSWORD)@host.docker.internal:$(DATABASE_PORT)/$(DATABASE_NAME)?sslmode=$(SSL_MODE))
	$(eval SHADOW_DB_URL := postgres://postgres:postgres@host.docker.internal:$(SHADOW_PORT)/postgres?sslmode=disable)

	# Using the official public Python slim image (which is perfectly multi-arch compatible)
	# and installing migra on-the-fly inside an isolated runtime container
	@docker run --rm \
		--add-host=host.docker.internal:host-gateway \
		python:3.11-slim bash -c "pip install --quiet migra psycopg2-binary && migra --unsafe '$(LIVE_DB_URL)' '$(SHADOW_DB_URL)'" > migration_changes.sql || true

	@if [ -s migration_changes.sql ] && grep -q '[a-zA-Z]' migration_changes.sql; then \
		echo ""; \
		echo "📋 The following differences were detected between your code and the live database:"; \
		echo "-----------------------------------------------------------------"; \
		cat migration_changes.sql; \
		echo "-----------------------------------------------------------------"; \
		echo ""; \
		exec 3<&0; \
		read -p "⚠️ Do you want to apply these alterations to your live database? [y/N]: " CONFIRM <&3; \
		if [ "$$CONFIRM" = "y" ] || [ "$$CONFIRM" = "Y" ]; then \
			echo "🚀 Applying database alterations safely..."; \
			docker exec -i $(LIVE_CONTAINER_NAME) psql -U $(DATABASE_USER) -d $(DATABASE_NAME) < migration_changes.sql; \
			echo "🎉 Database synchronized successfully!"; \
		else \
			echo "❌ Migration cancelled. No changes were made to your live data."; \
		fi \
	else \
		echo "✅ Live database columns are already fully up to date with your domain files."; \
	fi

	@echo "🧹 Step 5: Cleaning up local temporary tracking contexts..."
	@docker stop $(SHADOW_CONTAINER) > /dev/null
	@rm -f migration_changes.sql
