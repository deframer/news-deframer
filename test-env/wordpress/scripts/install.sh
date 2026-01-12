n#!/bin/bash

WORDPRESS_TITLE=${WORDPRESS_TITLE:?WORDPRESS_TITLE is not set}
WORDPRESS_ADMIN_USER=${WORDPRESS_ADMIN_USER:?WORDPRESS_ADMIN_USER is not set}
WORDPRESS_ADMIN_PASSWORD=${WORDPRESS_ADMIN_PASSWORD:?WORDPRESS_ADMIN_PASSWORD is not set}
HTTP_HOST=${HTTP_HOST:?HTTP_HOST is not set}

# Check for --force argument
FORCE=false
if [[ "$1" == "--force" ]]; then
    FORCE=true
fi

# --- Install WordPress (if needed) ---
if ! wp core is-installed; then
    echo "Installing WordPress..."
    wp core install --url="$HTTP_HOST" --title="$WORDPRESS_TITLE" --admin_user="$WORDPRESS_ADMIN_USER" --admin_password="$WORDPRESS_ADMIN_PASSWORD" --admin_email="admin@example.com"

    # https://gist.github.com/JoeyBurzynski/34229834af1ac7a7e1e5007bc5b17021
    # Inject Dynamic Host/Protocol Detection into wp-config.php
    # This enables localhost:8000 AND wordpress:8000 AND https://reverse-proxy.com

    # This is insecure
    echo "Updating wp-config.php for dynamic hosts..."

    sed -i "1a \\
        // -- BEGIN DYNAMIC HOST LOGIC -- \\
        // Detect if SSL is used directly OR via a Reverse Proxy (Termination) \\
        \$is_secure = (isset(\$_SERVER['HTTPS']) && \$_SERVER['HTTPS'] === 'on') || \\
                    (isset(\$_SERVER['HTTP_X_FORWARDED_PROTO']) && \$_SERVER['HTTP_X_FORWARDED_PROTO'] === 'https'); \\
        \\
        \$protocol = \$is_secure ? 'https://' : 'http://'; \\
        \$current_host = \$_SERVER['HTTP_HOST']; \\
        \\
        // Override DB settings with the current request URL \\
        define('WP_HOME', \$protocol . \$current_host); \\
        define('WP_SITEURL', \$protocol . \$current_host); \\
        // -- END DYNAMIC HOST LOGIC --" wp-config.php

    # Create a mu-plugin to forcefully disable comments (our theme has comments still enabled)
    mkdir -p wp-content/mu-plugins
    cat <<EOF > wp-content/mu-plugins/force-disable-comments.php
<?php
add_filter('comments_open', '__return_false', 20, 2);
add_filter('pings_open', '__return_false', 20, 2);
add_filter('comments_array', '__return_empty_array', 10, 2);
EOF
else
    if [ "$FORCE" = true ]; then
        echo "WordPress is already installed. Force enabled, continuing..."
    else
        echo "WordPress is already installed. Assuming setup is complete. Use --force to override."
        exit 0
    fi
fi

echo "Starting setup..."

# --- Clean up Default Content ---
echo "Clean up Default Content..."

# Delete 'Hello World' post (ID 1)
wp post delete 1 --force 2>/dev/null || true

# Delete 'Sample Page'
SAMPLE_PAGE_ID=$(wp post list --post_type=page --post_title="Sample Page" --format=ids)
if [ -n "$SAMPLE_PAGE_ID" ]; then wp post delete $SAMPLE_PAGE_ID --force; fi

# --- Set Permalinks ---
echo "Set permalink to '/%postname%/' ..."
wp rewrite structure '/%postname%/'

# --- Disable Comments ---
echo "Disabling comments globally..."

# Also install this plugin
wp plugin install disable-comments --activate

# Configure settings for FUTURE posts
wp option set default_comment_status closed
wp option set default_ping_status closed
wp option set default_pingback_flag 0

# Close comments/pings on EXISTING posts and pages
# This removes the "Leave a Reply" form from existing content
# wp db query "UPDATE $(wp db prefix)posts SET comment_status = 'closed', ping_status = 'closed' WHERE post_type IN ('post', 'page');"
wp db query "UPDATE $(wp db prefix)posts SET comment_status = 'closed', ping_status = 'closed'"

# Wipe all comment data
wp db query "TRUNCATE TABLE $(wp db prefix)comments; TRUNCATE TABLE $(wp db prefix)commentmeta;"
wp db query "UPDATE $(wp db prefix)posts SET comment_count = 0;"

# Clear cache to make sure the "Leave a Reply" form disappears immediately
wp cache flush

# --- Setup RSS ---
echo "Setup RSS..."

##### print the feed usually $HOME/feed/
# wp option get home | awk '{print $1 "/feed/"}'
# http://localhost:8000/feed

# By default, WordPress shows the last 10 posts. To change this to 20:
wp option update posts_per_rss 20
# Set to 1 for Summary (Excerpt) / 0 for Full Text
wp option update rss_use_excerpt 1

# After the theme importer, the data of some streams is still
# https://ap0.qsandbox.cloud/site/colormag-demo - this need to be updated in DB

# --- Install theme ---

# Install the theme ColorMag and activate it - overwrite the existing theme if there is oe
echo "Installing ColorMag Theme..."
wp theme install colormag --activate --force

# Install we need the plugin to install some demo content
echo "Installing themegrill-demo-importer Plugin"
wp plugin install themegrill-demo-importer --activate

echo "crating done file"
touch "/var/www/html/done"

echo "--------------------------------------------------"
echo "SITE READY: $HTTP_HOST"
echo "Login: $WORDPRESS_ADMIN_USER / $WORDPRESS_ADMIN_PASSWORD"
echo "--------------------------------------------------"

exit 0
