#!/bin/sh

# CloudGram-GO Docker entrypoint script
# This script reads environment variables and converts them to command line arguments
# for the cloudgram-go binary.

set -e

# String to store command line arguments (POSIX shell compatible)
ARGS=""

# Add listen address if specified
if [ -n "$LISTEN" ]; then
    ARGS="$ARGS -L '$LISTEN'"
fi

# Add authentication user if specified
if [ -n "$AUTH_USER" ]; then
    ARGS="$ARGS -U '$AUTH_USER'"
fi

# Add authentication password if specified
if [ -n "$AUTH_PASSWORD" ]; then
    ARGS="$ARGS -P '$AUTH_PASSWORD'"
fi

# Add Telegram bot token if specified
if [ -n "$TELEGRAM_BOT_TOKEN" ]; then
    ARGS="$ARGS -T '$TELEGRAM_BOT_TOKEN'"
fi

# Add database type if specified
if [ -n "$DB_TYPE" ]; then
    ARGS="$ARGS --type '$DB_TYPE'"
fi

# Add database connection string if specified
if [ -n "$DB_DSN" ]; then
    ARGS="$ARGS --dsn '$DB_DSN'"
fi

# Add log file path if specified
if [ -n "$LOG_PATH" ]; then
    ARGS="$ARGS --log '$LOG_PATH'"
fi

# Add JWT secret key if specified
if [ -n "$JWT_SECRET_KEY" ]; then
    ARGS="$ARGS --jwt-secret '$JWT_SECRET_KEY'"
fi

# Enable debug mode if specified
if [ -n "$DEBUG" ] && ([ "$DEBUG" = "true" ] || [ "$DEBUG" = "1" ]); then
    ARGS="$ARGS -d"
fi

# Execute the cloudgram-go binary with the constructed arguments
# Use eval to properly handle the argument string
eval "exec /app/cloudgram-go $ARGS \"\$@\""