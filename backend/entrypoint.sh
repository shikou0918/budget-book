#!/bin/sh

# Wait for MySQL to be ready
echo "Waiting for MySQL to be ready..."
while ! nc -z mysql 3306; do
  sleep 1
done
echo "MySQL is ready!"

# Run migrations
echo "Running database migrations..."
mysql -h mysql -u root -ppassword --skip-ssl budget_book < /app/migrations/schema.sql

echo "Starting application..."
exec "$@"