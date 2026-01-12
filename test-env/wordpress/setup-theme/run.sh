#!/bin/bash

# Check required environment variables
: "${HTTP_HOST:?HTTP_HOST is not set}"
: "${WORDPRESS_ADMIN_USER:?WORDPRESS_ADMIN_USER is not set}"
: "${WORDPRESS_ADMIN_PASSWORD:?WORDPRESS_ADMIN_PASSWORD is not set}"

python3 main.py --url "$HTTP_HOST" \
                --user "$WORDPRESS_ADMIN_USER" \
                --password "$WORDPRESS_ADMIN_PASSWORD" \
                --headless "$@"
