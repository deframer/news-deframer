#!/bin/bash

# To deduplicate logs (keep last occurrence of same feed/event/pubdate):
#   jq -s 'group_by(.feed, .event, .new_pub_date) | map(last)[]' log.json
# To filter by feed and only show unique dates (preserving order):
#   jq -r 'select(.feed == "spiegel") | .new_pub_date' log.json | awk '!seen[$0]++'

# Default Configuration
CONFIG_FILE="feeds.json"
LOG_DIR="."
CHECK_INTERVAL=300

usage() {
    echo "Usage: $0 [-c config_file] [-d log_dir] [-i interval_seconds]"
    echo "  -c  Path to JSON config file (default: $CONFIG_FILE)"
    echo "  -d  Directory to store RSS logs (default: $LOG_DIR)"
    echo "  -i  Check interval in seconds (default: $CHECK_INTERVAL)"
    echo ""
    echo "Example JSON config:"
    echo '['
    echo '  { "name": "xxx", "url": "https://www.xxx/index.rss" },'
    echo '  { "name": "yyy", "url": "https://www.yyy/rss2" }'
    echo ']'
    exit 1
}

while getopts "c:d:i:h" opt; do
    case $opt in
        c) CONFIG_FILE="$OPTARG" ;;
        d) LOG_DIR="$OPTARG" ;;
        i) CHECK_INTERVAL="$OPTARG" ;;
        h) usage ;;
        *) usage ;;
    esac
done

# Check dependencies
if ! command -v jq &> /dev/null; then
    echo "Error: 'jq' is required for JSON parsing."
    exit 1
fi

# Associative array to store last pub dates (requires Bash 4.0+)
declare -A last_pub_dates

echo "Starting RSS Watcher (Bash)..."
echo "Config:   $CONFIG_FILE"
echo "Log Dir:  $LOG_DIR"
echo "Interval: $CHECK_INTERVAL seconds"

# Ensure log directory exists
mkdir -p "$LOG_DIR"

while true; do
    if [ ! -f "$CONFIG_FILE" ]; then
        echo "Error: $CONFIG_FILE not found."
        sleep "$CHECK_INTERVAL"
        continue
    fi

    # Read feeds from JSON using process substitution to keep variables in scope
    while read -r item; do
        name=$(echo "$item" | jq -r '.name')
        url=$(echo "$item" | jq -r '.url')

        # Skip invalid entries
        if [[ -z "$name" || -z "$url" || "$name" == "null" ]]; then
            continue
        fi

        # Fetch the feed
        # -s: Silent mode
        # -L: Follow redirects
        content=$(curl -sL "$url")

        # Check if curl failed
        if [ $? -ne 0 ]; then
            echo "[$name] Error: Failed to fetch URL."
        else
            # Extract the first <pubDate> (Channel date)
            current_pub_date=$(echo "$content" | grep "<pubDate>" | head -n 1 | sed -e 's/.*<pubDate>\(.*\)<\/pubDate>.*/\1/')

            # Check if we found a date
            if [ -n "$current_pub_date" ]; then
                last_date="${last_pub_dates[$name]}"

                # Compare with last known date
                if [ "$current_pub_date" != "$last_date" ]; then
                    timestamp=$(date '+%H:%M:%S')
                    echo "[$timestamp] [$name] 🚨 PubDate changed! New: $current_pub_date"

                    # Append machine-readable JSON log to file
                    jq -n -c --arg ts "$(date -Iseconds)" --arg name "$name" --arg pubdate "$current_pub_date" \
                        '{timestamp: $ts, feed: $name, event: "pub_date_changed", new_pub_date: $pubdate}' >> "$LOG_DIR/log.json"

                    # Update state
                    last_pub_dates[$name]="$current_pub_date"
                fi
            fi
        fi
    done < <(jq -c '.[]' "$CONFIG_FILE")

    # Wait for next cycle
    sleep "$CHECK_INTERVAL"
done
