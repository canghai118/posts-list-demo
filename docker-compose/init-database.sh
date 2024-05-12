#!/bin/bash

DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASS=demopassword

# Command to check if the database is ready
CHECK_CMD="mysql -h$DB_HOST -u $DB_USER -p$DB_PASS -e 'SELECT 1'"

# Wait for database to be up
while ! timeout 1 bash -c "echo > /dev/tcp/$DB_HOST/$DB_PORT" 2>/dev/null
do
  echo "Waiting for database to be up..."
  sleep 1
done

echo "Database is up!"

set -e
mysql -h$DB_HOST -P$DB_PORT -u$DB_USER -p$DB_PASS < /sqls/init-schema.sql

echo "Database initialized!"