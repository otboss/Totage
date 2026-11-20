# --- CONFIGURATION ---
TMP_DB_NAME = migration_target_tmp
ENV_FILE = .env

# Local SQL schema tracking files (space-separated list)
# REGISTER SCHEMA FILES BELOW IN THE CORRECT ORDER DEPENDING ON THEIR RELATIONSHIPS WITH OTHER TABLES
SCHEMA_FILES = \
    ./shared/domain/entities/shared.schema.sql \
    ./modules/user/domain/entities/user.schema.sql \
    ./modules/payment/domain/entities/payment.schema.sql \
    ./modules/organization/domain/entities/organization.schema.sql \
    ./modules/logistic/domain/entities/logistic.schema.sql \
    ./modules/order/domain/entities/order.schema.sql

# --- DYNAMIC ENVIRONMENT LOADER (CROSS-PLATFORM) ---
# Uses Python to parse the .env file safely across Unix and Windows systems
DB_USER = $(shell python3 -c "import os; from pathlib import Path; [print(line.split('=')[1].strip()) for line in Path('$(ENV_FILE)').read_text().splitlines() if line.startswith('DATABASE_USER=')]")
DB_PASS = $(shell python3 -c "import os; from pathlib import Path; [print(line.split('=')[1].strip()) for line in Path('$(ENV_FILE)').read_text().splitlines() if line.startswith('DATABASE_PASSWORD=')]")
DB_HOST = $(shell python3 -c "import os; from pathlib import Path; vals = [line.split('=')[1].strip() for line in Path('$(ENV_FILE)').read_text().splitlines() if line.startswith('DATABASE_HOST=')]; print(vals[0] if vals else 'localhost')")
DB_PORT = $(shell python3 -c "import os; from pathlib import Path; vals = [line.split('=')[1].strip() for line in Path('$(ENV_FILE)').read_text().splitlines() if line.startswith('DATABASE_PORT=')]; print(vals[0] if vals else '5432')")
DB_NAME = $(shell python3 -c "import os; from pathlib import Path; [print(line.split('=')[1].strip()) for line in Path('$(ENV_FILE)').read_text().splitlines() if line.startswith('DATABASE_NAME=')]")

# Determine SSL parameter cleanly using Python
SSL_MODE = $(shell python3 -c "import os; from pathlib import Path; vals = [line.split('=')[1].strip().lower() for line in Path('$(ENV_FILE)').read_text().splitlines() if line.startswith('DATABASE_SSL_ENABLED=')]; print('require' if vals and vals[0] in ['true', '1'] else 'disable')")

# Assemble PostgreSQL connection URIs
LIVE_DB_URL = "postgresql://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(SSL_MODE)"
TMP_DB_URL  = "postgresql://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(TMP_DB_NAME)?sslmode=$(SSL_MODE)"

# Export password safely to the environment so Postgres binaries pick it up automatically
export PGPASSWORD=$(DB_PASS)

.PHONY: all check-deps init-db load-schemas generate-diff cleanup migrate

# Default target runs the full analysis pipeline
all: check-deps init-db load-schemas generate-diff cleanup

check-deps:
	@echo "🔑 Verifying environment configuration and system dependencies..."
	@if [ ! -f "$(ENV_FILE)" ]; then echo "❌ Error: $(ENV_FILE) file missing."; exit 1; fi
	@command -v migra >/dev/null 2>&1 || { echo "❌ Error: 'migra' utility is not installed. Run: pip install migra sqlbag"; exit 1; }
	@command -v psql >/dev/null 2>&1 || { echo "❌ Error: 'psql' client utilities are not installed."; exit 1; }

init-db:
	@echo "🔄 Resetting temporary sync database using template0..."
	@dropdb -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) --if-exists $(TMP_DB_NAME)
	@createdb -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) --template=template0 $(TMP_DB_NAME)
	@echo "🌐 Injecting PostGIS extension and namespaces..."
	@psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(TMP_DB_NAME) -q -c "CREATE EXTENSION IF NOT EXISTS postgis;"
	@psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(TMP_DB_NAME) -q -c "CREATE SCHEMA IF NOT EXISTS shared;"
	@psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(TMP_DB_NAME) -q -c "CREATE SCHEMA IF NOT EXISTS \"user\";"
	@psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(TMP_DB_NAME) -q -c "CREATE SCHEMA IF NOT EXISTS payment;"
	@psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(TMP_DB_NAME) -q -c "CREATE SCHEMA IF NOT EXISTS organization;"
	@psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(TMP_DB_NAME) -q -c "CREATE SCHEMA IF NOT EXISTS logistic;"
	@psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(TMP_DB_NAME) -q -c "CREATE SCHEMA IF NOT EXISTS _order;"

load-schemas:
	@echo "📂 Loading local schema sequence into temporary database..."
	@for file in $(SCHEMA_FILES); do \
		if [ -f "$$file" ]; then \
			echo "  ↳ Applying $$file..."; \
			psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(TMP_DB_NAME) -q -f "$$file"; \
		else \
			echo "⚠️ Warning: Schema file not found: $$file"; \
		fi \
	done

generate-diff:
	@echo "🔀 Calculating PostgreSQL schema differences..."
	@echo "--------------------------------------------------------"
	@# Capture migra's exit status code. Status 0 (no diff) and 2 (diff found) are both successes.
	@STATUS=0; migra --unsafe --schema=shared --schema=user --schema=payment --schema=organization --schema=logistic --schema=_order $(LIVE_DB_URL) $(TMP_DB_URL) > migration_output.sql || STATUS=$$?; \
	if [ $$STATUS -ne 0 ] && [ $$STATUS -ne 2 ]; then \
		echo "❌ Error: migra failed with unexpected crash exit code $$STATUS"; \
		exit 1; \
	fi
	@echo "✅ Migration script generated successfully!"
	@echo "📄 Output saved to: migration_output.sql"
	@echo "--------------------------------------------------------"
	@if [ ! -s migration_output.sql ]; then \
		echo "🎉 Databases are already in sync. No changes detected."; \
	else \
		echo "📋 Preview of ALTER commands to be executed:"; \
		cat migration_output.sql; \
	fi

cleanup:
	@echo "--------------------------------------------------------"
	@echo "🧹 Cleaning up temporary tracking entities..."
	@dropdb -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) $(TMP_DB_NAME)
	@echo "✨ Done."
