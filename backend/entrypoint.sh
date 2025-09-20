#!/bin/sh

# Wait for MySQL to be ready
echo "Waiting for MySQL to be ready..."
while ! nc -z mysql 3306; do
  sleep 1
done
echo "MySQL is ready!"

# Run migrations only if migration was successful or if tables don't exist
echo "Running database migrations..."
mysql -h mysql -u root -ppassword --skip-ssl budget_book < /app/migrations/schema.sql || echo "Migration script executed (some tables may already exist)"

echo "Starting application..."
echo "Current directory: $(pwd)"
echo "Files in current directory:"
ls -la
echo "Looking for api file specifically:"
ls -la ./api || echo "api file not found!"
echo "Executing: $@"
exec "$@"