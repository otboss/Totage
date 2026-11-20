#!/usr/bin/env bash

# Exit immediately if a command exits with a non-zero status
set -eo pipefail

# --- CONFIGURATION ---
# Target schema files to read from (Order matters if there are dependencies)
SCHEMA_FILES=(
    "./schemas/users.sql"
    "./schemas/orders.sql"
    "./schemas/products.sql"
)

# Temporary Database name used to parse and compile the target schema files
TMP_DB_NAME="migration_target_tmp"
ENV_FILE=".env"

# --- ENVIRONMENT LOADER ---
if [ -f "$ENV_FILE" ]; then
    echo "🔑 Loading environment variables from $ENV_FILE..."
    # Export variables from .env while ignoring comments and empty lines
    set -a
    source <(grep -v '^#' "$ENV_FILE" | grep -v '^$')
    set +a
else
    echo "❌ Error: Configuration file '$ENV_FILE' not found in this directory." >&2
    exit 1
fi

# Map loaded variables to internal names (handling fallbacks if necessary)
DB_USER="${DATABASE_USER}"
DB_PASS="${DATABASE_PASSWORD}"
DB_HOST="${DATABASE_HOST:-localhost}"
DB_PORT="${DATABASE_PORT:-5432}"
DB_NAME="${DATABASE_NAME}"

# Determine SSL Mode based on DATABASE_SSL_ENABLED
SSL_MODE="disable"
if [ "${DATABASE_SSL_ENABLED,,}" = "true" ] || [ "${DATABASE_SSL_ENABLED}" = "1" ]; then
    SSL_MODE="require"
fi

# Build Connection Strings with SSL Mode
LIVE_DB_URL="postgresql://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${SSL_MODE}"
TMP_DB_URL="postgresql://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${TMP_DB_NAME}?sslmode=${SSL_MODE}"

# --- SANITY CHECKS ---
if ! command -v migra &> /dev/null; then
    echo "❌ Error: 'migra' command utility is not installed." >&2
    echo "Please run: pip install migra sqlbag" >&2
    exit 1
fi

if ! command -v psql &> /dev/null; then
    echo "❌ Error: 'psql' client is not installed or not in PATH." >&2
    exit 1
fi

# --- STEP 1: Prepare Temporary Database ---
echo "🔄 Resetting temporary sync database..."
# Drop the temp DB if it exists and create a fresh one
PGPASSWORD="$DB_PASS" dropdb -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" --if-exists "$TMP_DB_NAME"
PGPASSWORD="$DB_PASS" createdb -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" --template=template0 "$TMP_DB_NAME"

# --- STEP 2: Load Target Schema Array into Temp DB ---
echo "📂 Loading local schema files into temporary database..."
for file in "${SCHEMA_FILES[@]}"; do
    if [ -f "$file" ]; then
        echo "  ↳ Applying $file..."
        PGPASSWORD="$DB_PASS" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$TMP_DB_NAME" -q -f "$file"
    else
        echo "⚠️ Warning: Schema file not found: $file" >&2
    fi
done

# --- STEP 3: Generate Diff (Alter Commands) ---
echo "🔀 Calculating schema differences..."
echo "--------------------------------------------------------"

# Migra compares (Origin/Live DB) -> (Target/Temp DB) and outputs the necessary ALTER/CREATE statements
# --unsafe allows migra to generate DROP commands if columns/tables were removed in your files
if migra --unsafe "$LIVE_DB_URL" "$TMP_DB_URL" > migration_output.sql; then
    echo "✅ Migration script generated successfully!"
    echo "📄 Output saved to: migration_output.sql"
    echo "--------------------------------------------------------"
    
    # Check if the output file is empty (meaning schemas are fully synchronized)
    if [ ! -s migration_output.sql ]; then
        echo "🎉 Databases are already in sync. No changes detected."
    else
        echo "📋 Preview of ALTER commands to be executed:"
        cat migration_output.sql
    fi
else
    echo "❌ Error generating schema diff." >&2
    exit 1
fi

# --- STEP 4: Cleanup ---
echo "--------------------------------------------------------"
echo "🧹 Cleaning up temporary database..."
PGPASSWORD="$DB_PASS" dropdb -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" "$TMP_DB_NAME"
echo "✨ Done."
