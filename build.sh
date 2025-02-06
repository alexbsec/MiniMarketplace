#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Define the migration directory
MIGRATION_DIR="src/db/migrations"

# Define the environment
ENV="gorm"

echo "Checking and applying any pending migrations..."
# Apply any pending migrations to avoid a dirty state
atlas migrate apply --env $ENV

# Capture the number of files in the migration directory before running `diff`
initial_count=$(ls -A $MIGRATION_DIR | wc -l)

echo "Checking for schema changes using Atlas..."

# Generate a diff to check if any changes exist
atlas migrate diff --env $ENV --dir "file://$MIGRATION_DIR"

# Capture the number of files in the migration directory after running `diff`
final_count=$(ls -A $MIGRATION_DIR | wc -l)

# Compare the number of files before and after
if [ "$initial_count" -ne "$final_count" ]; then
    echo "Schema changes detected. Applying migrations..."

    # Apply the migration
    atlas migrate apply --env $ENV

    echo "Migrations applied successfully."
else
    echo "No schema changes detected. Skipping migration application."
fi

echo "Building the app..."
make build

echo "Starting the application..."

while [ 1 -le 2 ]; do
    :
done

./bin/main


