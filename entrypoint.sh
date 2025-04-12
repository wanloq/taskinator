#!/bin/bash

# Load secrets into environment variables
export DATABASE_URL=$(cat /run/secrets/db_url)

# Optionally, load other secrets if needed in Go:
export SMTP_USERNAME=$(cat /run/secrets/smtp_username)
export SMTP_PASSWORD=$(cat /run/secrets/smtp_password)
export JWT_SECRET=$(cat /run/secrets/jwt_key)
export PORT=$(cat /run/secrets/port)

# Start the Go app
echo "🚀 Starting Taskinator..."
./main
