#!/bin/sh -e

#############################################
# STREMIO WEB SERVICE RUN SCRIPT
#############################################

# CONFIGURATION
#############################################
# Set the configuration folder path
CONFIG_FOLDER="${APP_PATH:-${HOME}/.stremio-server/}"
echo "[CONFIG] Using configuration folder: ${CONFIG_FOLDER}"

#############################################
# HELPER FUNCTIONS
#############################################

# Function to get public IP address from various services
get_public_ip() {
    # Try multiple services to get public IP in case one fails
    PUBLIC_IP=""
    
    # Try icanhazip.com first
    if [ -z "${PUBLIC_IP}" ]; then
        echo "[IP LOOKUP] Attempting to get public IP from icanhazip.com..." >&2
        IP_RESULT=$(curl -s --connect-timeout 5 https://icanhazip.com)
        if [ -n "${IP_RESULT}" ] && [ "${IP_RESULT}" != "curl: "* ]; then
            echo "[IP LOOKUP] Successfully obtained IP from icanhazip.com: ${IP_RESULT}" >&2
            PUBLIC_IP="${IP_RESULT}"
        else
            echo "[IP LOOKUP] Failed to get IP from icanhazip.com" >&2
        fi
    fi
    
    # If icanhazip failed, try ifconfig.me
    if [ -z "${PUBLIC_IP}" ]; then
        echo "[IP LOOKUP] Attempting to get public IP from ifconfig.me..." >&2
        IP_RESULT=$(curl -s --connect-timeout 5 https://ifconfig.me)
        if [ -n "${IP_RESULT}" ] && [ "${IP_RESULT}" != "curl: "* ]; then
            echo "[IP LOOKUP] Successfully obtained IP from ifconfig.me: ${IP_RESULT}" >&2
            PUBLIC_IP="${IP_RESULT}"
        else
            echo "[IP LOOKUP] Failed to get IP from ifconfig.me" >&2
        fi
    fi
    
    # Return only the IP address, no debug messages
    echo "${PUBLIC_IP}"
}

# Function to start HTTP server with specified options
start_http_server() {
    echo "[SERVER] Starting HTTP server on port 8080 with options: $*"
    http-server build/ -p 8080 -d false "$@"
}

#############################################
# INITIAL SERVER CONFIGURATION
#############################################

# Update paths in server-settings.json if it exists
if [ -f "${CONFIG_FOLDER}server-settings.json" ]; then
    echo "[CONFIG] Found server-settings.json at ${CONFIG_FOLDER}server-settings.json"
    # Remove trailing slash from CONFIG_FOLDER for consistency in the JSON file
    CONFIG_PATH=$(echo "${CONFIG_FOLDER}" | sed 's:/$::')
    
    # Get current values for logging
    CURRENT_APP_PATH=$(grep -o '"appPath": "[^"]*"' "${CONFIG_FOLDER}server-settings.json" | cut -d'"' -f4)
    CURRENT_CACHE_ROOT=$(grep -o '"cacheRoot": "[^"]*"' "${CONFIG_FOLDER}server-settings.json" | cut -d'"' -f4)
    
    echo "[CONFIG] Updating paths in server-settings.json:"
    echo "  - appPath: '${CURRENT_APP_PATH}' -> '${CONFIG_PATH}'"
    echo "  - cacheRoot: '${CURRENT_CACHE_ROOT}' -> '${CONFIG_PATH}'"
    
    sed -i "s|\"appPath\": \"[^\"]*\"|\"appPath\": \"${CONFIG_PATH}\"|g" "${CONFIG_FOLDER}server-settings.json"
    sed -i "s|\"cacheRoot\": \"[^\"]*\"|\"cacheRoot\": \"${CONFIG_PATH}\"|g" "${CONFIG_FOLDER}server-settings.json"
fi

# Patch server.js with necessary fixes
echo "[PATCH] Checking server.js for required patches..."

# Disable proxy streams if not already disabled
if ! grep -q 'self.proxyStreamsEnabled = false,' server.js; then
    echo "[PATCH] Adding 'self.proxyStreamsEnabled = false' to server.js after 'self.allTranscodeProfiles = []' line"
    sed -i '/self.allTranscodeProfiles = \[\]/a \ \ \ \ \ \ \ \ self.proxyStreamsEnabled = false,' server.js
    echo "[PATCH] ✓ Proxy streams disabled"
else
    echo "[PATCH] ✓ Proxy streams already disabled"
fi

# Fix disk space check to work in all environments
if grep -q 'df -k' server.js; then
    echo "[PATCH] Replacing 'df -k' with 'df -Pk' in server.js to ensure consistent output format across systems"
    sed -i 's/df -k/df -Pk/g' server.js
    echo "[PATCH] ✓ Disk space check fixed"
else
    echo "[PATCH] ✓ Disk space check already using correct command"
fi

#############################################
# SERVER URL CONFIGURATION
#############################################

if [ -n "${SERVER_URL}" ]; then    
    echo "[URL CONFIG] Processing SERVER_URL: ${SERVER_URL}"
    ORIGINAL_SERVER_URL="${SERVER_URL}"
    
    # Handle special cases in SERVER_URL (0.0.0.0 or stremio.rocks domains)
    if echo "${SERVER_URL}" | grep -q "0\.0\.0\.0"; then
        echo "[URL CONFIG] Found '0.0.0.0' in SERVER_URL, will replace with actual IP"
        PUBLIC_IP=$(get_public_ip)
        
        if [ -n "${PUBLIC_IP}" ]; then
            # Replace 0.0.0.0 with the public IP while preserving protocol and port
            SERVER_URL=$(echo "${SERVER_URL}" | sed "s/0\.0\.0\.0/${PUBLIC_IP}/g")
            echo "[URL CONFIG] Replaced 0.0.0.0 with public IP: ${ORIGINAL_SERVER_URL} -> ${SERVER_URL}"
        else
            echo "[URL CONFIG] Failed to obtain public IP, keeping original SERVER_URL with 0.0.0.0"
        fi
    elif echo "${SERVER_URL}" | grep -q "0-0-0-0\.519b6502d940\.stremio\.rocks"; then
        echo "[URL CONFIG] Found 'stremio.rocks' test domain in SERVER_URL, will replace with actual IP"
        PUBLIC_IP=$(get_public_ip)
        
        if [ -n "${PUBLIC_IP}" ]; then
            # Replace 0-0-0-0.519b6502d940.stremio.rocks with the public IP
            SERVER_URL=$(echo "${SERVER_URL}" | sed "s/0-0-0-0\.519b6502d940\.stremio\.rocks/${PUBLIC_IP}/g")
            echo "[URL CONFIG] Replaced stremio.rocks domain with public IP: ${ORIGINAL_SERVER_URL} -> ${SERVER_URL}"
        else
            echo "[URL CONFIG] Failed to obtain public IP, keeping original SERVER_URL with stremio.rocks domain"
        fi
    fi
    
    # Ensure URL has trailing slash for consistency
    SERVER_URL_BEFORE="${SERVER_URL}"
    SERVER_URL=$(echo "${SERVER_URL}" | sed 's:/*$:/:' )
    
    if [ "${SERVER_URL}" != "${SERVER_URL_BEFORE}" ]; then
        echo "[URL CONFIG] Added trailing slash for consistency: ${SERVER_URL_BEFORE} -> ${SERVER_URL}"
    fi
    
    echo "[URL CONFIG] Final server URL: ${SERVER_URL}"
    
    # Update localStorage for the web client to use our server
    echo "[CONFIG] Updating localStorage.json to use configured server URL"
    echo "[CONFIG] Replacing 'http://127.0.0.1:11470/' with '${SERVER_URL}'"
    cp localStorage.json build/localStorage.json
    sed -i "s|http://127.0.0.1:11470/|${SERVER_URL}|g" build/localStorage.json
    sed -i "s|http://127.0.0.1:11470/|${SERVER_URL}|g" server.js
fi

#############################################
# WEB UI CONFIGURATION
#############################################

# Configure Web UI redirection if WEBUI_LOCATION is set
if [ -n "${WEBUI_LOCATION}" ]; then
    echo "[WEB UI] Configuring custom Web UI location: ${WEBUI_LOCATION}"
    
    # Ensure WEBUI_LOCATION ends with a trailing slash for consistency
    WEBUI_LOCATION_BEFORE="${WEBUI_LOCATION}"
    WEBUI_LOCATION=$(echo "${WEBUI_LOCATION}" | sed 's:/*$:/:')
    
    if [ "${WEBUI_LOCATION}" != "${WEBUI_LOCATION_BEFORE}" ]; then
        echo "[WEB UI] Added trailing slash to WEBUI_LOCATION: ${WEBUI_LOCATION_BEFORE} -> ${WEBUI_LOCATION}"
    fi
    
    # Escape forward slashes in the URL for sed
    ESCAPED_URL=$(echo "${WEBUI_LOCATION}" | sed 's/\//\\\//g')
    
    echo "[WEB UI] Updating server.js to redirect to custom Web UI location"
    echo "[WEB UI] Looking for patterns to replace with: ${WEBUI_LOCATION}"
    
    # Replace all variations of the default redirect URL patterns in server.js
    REPLACEMENTS_MADE=0
    
    # Search for occurrences before replacement for better logging
    PATTERN1_COUNT=$(grep -c "https://app\.strem\.io/shell-v4\.4/" server.js || echo 0)
    PATTERN2_COUNT=$(grep -c "https://app\.strem\.io/shell-v4\.4[^/]" server.js || echo 0)
    PATTERN3_COUNT=$(grep -c "https://app\.strem\.io/shell-v[0-9]\+\.[0-9]\+/" server.js || echo 0)
    PATTERN4_COUNT=$(grep -c "app\.strem\.io/shell-v[0-9.]\+" server.js || echo 0)
    
    # Log what we found before making changes
    echo "[WEB UI] Found URL patterns in server.js:"
    echo "  - Pattern 'https://app.strem.io/shell-v4.4/': ${PATTERN1_COUNT} occurrences"
    echo "  - Pattern 'https://app.strem.io/shell-v4.4' (no trailing slash): ${PATTERN2_COUNT} occurrences"
    echo "  - Pattern 'https://app.strem.io/shell-v*.*/': ${PATTERN3_COUNT} occurrences" 
    echo "  - Pattern 'app.strem.io/shell-v*.*': ${PATTERN4_COUNT} occurrences"
    
    # Pattern 1: https://app.strem.io/shell-v4.4/
    if [ "${PATTERN1_COUNT}" -gt 0 ]; then
        echo "[WEB UI] Replacing 'https://app.strem.io/shell-v4.4/' with '${WEBUI_LOCATION}'"
        sed -i "s/https:\/\/app\.strem\.io\/shell-v4\.4\//${ESCAPED_URL}/g" server.js
        REPLACEMENTS_MADE=$((REPLACEMENTS_MADE + PATTERN1_COUNT))
    fi
    
    # Pattern 2: https://app.strem.io/shell-v4.4 (no trailing slash)
    if [ "${PATTERN2_COUNT}" -gt 0 ]; then
        echo "[WEB UI] Replacing 'https://app.strem.io/shell-v4.4' (no trailing slash) with '${WEBUI_LOCATION}'"
        sed -i "s/https:\/\/app\.strem\.io\/shell-v4\.4([^\/])/${ESCAPED_URL}\1/g" server.js
        REPLACEMENTS_MADE=$((REPLACEMENTS_MADE + PATTERN2_COUNT))
    fi
    
    # Pattern 3: https://app.strem.io/shell-v (with any version number)
    if [ "${PATTERN3_COUNT}" -gt 0 ]; then
        echo "[WEB UI] Replacing 'https://app.strem.io/shell-v*.* with trailing slash' with '${WEBUI_LOCATION}'"
        sed -i "s/https:\/\/app\.strem\.io\/shell-v[0-9]\+\.[0-9]\+\//${ESCAPED_URL}/g" server.js
        REPLACEMENTS_MADE=$((REPLACEMENTS_MADE + PATTERN3_COUNT))
    fi
    
    # Pattern 4: Generic app.strem.io with shell-v pattern
    if [ "${PATTERN4_COUNT}" -gt 0 ]; then
        DOMAIN_PART=$(echo "${WEBUI_LOCATION}" | sed 's/^https\?:\/\///' | sed 's/\//\\\//g')
        echo "[WEB UI] Replacing 'app.strem.io/shell-v*.*' with '${DOMAIN_PART}'"
        sed -i "s/app\.strem\.io\/shell-v[0-9.]\+/${DOMAIN_PART}/g" server.js
        REPLACEMENTS_MADE=$((REPLACEMENTS_MADE + PATTERN4_COUNT))
    fi
    
    # Report results
    if [ $REPLACEMENTS_MADE -gt 0 ]; then
        echo "[WEB UI] ✓ Successfully updated server.js with custom redirect URLs (${REPLACEMENTS_MADE} replacements)"
    else
        echo "[WEB UI] ⚠ WARNING: No replacements made in server.js. The default URL patterns may be different than expected."
        echo "[WEB UI] Searching for strem.io patterns in server.js to help with debugging:"
        grep -n "app.strem.io" server.js || echo "[WEB UI] No instances of 'app.strem.io' found in server.js"
    fi
    
    # Check and update the streamingServer parameter if found
    if grep -q "streamingServer=" server.js; then
        echo "[WEB UI] Found streamingServer parameter in server.js"
        echo "[WEB UI] Note: Update for streamingServer parameter is currently disabled in this script"
        # Note: Original commented out line preserved
        # sed -i "s/\(streamingServer=\)[^&]*&/\1${ESCAPED_URL}&/g" server.js
    fi
fi

#############################################
# SERVER STARTUP
#############################################

# Echo startup message
echo "[STARTUP] Starting Stremio server at $(date)"
echo "[STARTUP] Configuration summary:"
echo "  - Config folder: ${CONFIG_FOLDER}"
echo "  - IP Address: ${IPADDRESS}"
echo "  - Server URL: ${SERVER_URL}"

# Handle different startup modes based on configuration
if [ -n "${IPADDRESS}" ]; then
    # Start the server process
    echo "[STARTUP] Starting Stremio server process"
    node server.js &
    SERVER_PID=$!
    echo "[STARTUP] Server started with PID: ${SERVER_PID}"

    # Check if IPADDRESS is an "any address" value
    if [ "${IPADDRESS}" = "0.0.0.0" ]; then
        echo "[STARTUP] IPADDRESS is set to '0.0.0.0' (any address)"
        PUBLIC_IP=$(get_public_ip)
        if [ -n "${PUBLIC_IP}" ]; then
            echo "[STARTUP] Using discovered public IP: ${PUBLIC_IP}"
            IPADDRESS="${PUBLIC_IP}"
        else
            echo "[STARTUP] ⚠ Failed to obtain public IP address, using original value: ${IPADDRESS}"
            echo "[STARTUP] ⚠ This may cause issues with certificate generation"
        fi
    elif [ "${IPADDRESS}" = "*" ]; then
        echo "[STARTUP] IPADDRESS is set to '*' (any address)"
        PUBLIC_IP=$(get_public_ip)
        if [ -n "${PUBLIC_IP}" ]; then
            echo "[STARTUP] Using discovered public IP: ${PUBLIC_IP}"
            IPADDRESS="${PUBLIC_IP}"
        else
            echo "[STARTUP] ⚠ Failed to obtain public IP address, using original value: ${IPADDRESS}"
            echo "[STARTUP] ⚠ This may cause issues with certificate generation"
        fi
    elif [ "${IPADDRESS}" = "any" ]; then
        echo "[STARTUP] IPADDRESS is set to 'any' (any address)"
        PUBLIC_IP=$(get_public_ip)
        if [ -n "${PUBLIC_IP}" ]; then
            echo "[STARTUP] Using discovered public IP: ${PUBLIC_IP}"
            IPADDRESS="${PUBLIC_IP}"
        else
            echo "[STARTUP] ⚠ Failed to obtain public IP address, using original value: ${IPADDRESS}"
            echo "[STARTUP] ⚠ This may cause issues with certificate generation"
        fi
    else
        # For specific IP addresses, just use them directly
        echo "[STARTUP] Using specific IP address: ${IPADDRESS}"
    fi
    
    # Certificate setup for HTTPS
    CERT_URL="http://localhost:11470/get-https?authKey=&ipAddress=${IPADDRESS}"
    echo "[HTTPS] Attempting to fetch HTTPS certificate for IP: ${IPADDRESS}"
    echo "[HTTPS] Using URL: ${CERT_URL}"
    
    # Use set -x to show the exact curl command being executed
    set -x
    curl --connect-timeout 5 \
         --retry-all-errors \
         --retry 10 \
         --retry-delay 1 \
         --verbose \
         "${CERT_URL}"
    CURL_STATUS="$?"
    set +x
    
    if [ "${CURL_STATUS}" -ne 0 ]; then
        echo "[HTTPS] ⚠ Failed to fetch HTTPS certificate. Curl exited with status: ${CURL_STATUS}"
    else
        echo "[HTTPS] ✓ Successfully requested HTTPS certificate"
    fi

    # Extract certificate information
    echo "[HTTPS] Extracting certificate information from ${CONFIG_FOLDER}httpsCert.json"
    IMPORTED_DOMAIN="$(node certificate.js --action extract --json-path "${CONFIG_FOLDER}httpsCert.json")"
    EXTRACT_STATUS="$?"
    IMPORTED_CERT_FILE="${CONFIG_FOLDER}${IMPORTED_DOMAIN}.pem"
    
    if [ "${EXTRACT_STATUS}" -eq 0 ]; then
        echo "[HTTPS] ✓ Successfully extracted domain from certificate: ${IMPORTED_DOMAIN}"
        echo "[HTTPS] Certificate file location: ${IMPORTED_CERT_FILE}"
    else
        echo "[HTTPS] ⚠ Failed to extract domain from certificate, exit code: ${EXTRACT_STATUS}"
    fi

    # Setup hosts file and start HTTP server with appropriate options
    if [ "${EXTRACT_STATUS}" -eq 0 ] && [ -n "${IMPORTED_DOMAIN}" ] && [ -f "${IMPORTED_CERT_FILE}" ]; then
        echo "[HOSTS] Adding entry to /etc/hosts: '${IPADDRESS} ${IMPORTED_DOMAIN}'"
        echo "${IPADDRESS} ${IMPORTED_DOMAIN}" >> /etc/hosts
        
        if [ "${FORCE_HTTP_FOR_WEBUI_SERVER}" = "true" ] || [ "${FORCE_HTTP_FOR_WEBUI_SERVER}" = "1" ] || [ "${FORCE_HTTP_FOR_WEBUI_SERVER}" = true ]; then
            echo "[SERVER] FORCE_HTTP_FOR_WEBUI_SERVER is set, starting Web UI server without HTTPS"
            echo "[SERVER] Using certificate for TLS options but not enabling HTTPS"
            start_http_server -C "${IMPORTED_CERT_FILE}" -K "${IMPORTED_CERT_FILE}"
        else
            echo "[SERVER] Starting Web UI server with HTTPS enabled"
            echo "[SERVER] Using certificate: ${IMPORTED_CERT_FILE}"
            start_http_server -S -C "${IMPORTED_CERT_FILE}" -K "${IMPORTED_CERT_FILE}"
        fi
    else
        echo "[HTTPS] ⚠ Failed to setup HTTPS due to missing or invalid certificate"
        echo "[SERVER] Falling back to HTTP for Web UI server"
        start_http_server
    fi
elif [ -n "${CERT_FILE}" ] && [ -n "${DOMAIN}" ]; then
    echo "[HTTPS] Using pre-configured certificate"
    echo "[HTTPS] Certificate file: ${CONFIG_FOLDER}${CERT_FILE}"
    echo "[HTTPS] Domain: ${DOMAIN}"
    
    # Start with pre-configured certificate
    node certificate.js --action load --pem-path "${CONFIG_FOLDER}${CERT_FILE}" --domain "${DOMAIN}" --json-path "${CONFIG_FOLDER}httpsCert.json"
    CERT_LOAD_STATUS="$?"
    
    if [ "${CERT_LOAD_STATUS}" -eq 0 ]; then
        echo "[HTTPS] ✓ Successfully loaded certificate"
        echo "[STARTUP] Starting Stremio server process"
        node server.js &
        SERVER_PID=$!
        echo "[STARTUP] Server started with PID: ${SERVER_PID}"
        
        if [ "${FORCE_HTTP_FOR_WEBUI_SERVER}" = "true" ] || [ "${FORCE_HTTP_FOR_WEBUI_SERVER}" = "1" ] || [ "${FORCE_HTTP_FOR_WEBUI_SERVER}" = true ]; then
            echo "[SERVER] FORCE_HTTP_FOR_WEBUI_SERVER is set, starting Web UI server without HTTPS"
            echo "[SERVER] Using certificate for TLS options but not enabling HTTPS"
            start_http_server -C "${CONFIG_FOLDER}${CERT_FILE}" -K "${CONFIG_FOLDER}${CERT_FILE}"
        else
            echo "[SERVER] Starting Web UI server with HTTPS enabled"
            echo "[SERVER] Using certificate: ${CONFIG_FOLDER}${CERT_FILE}"
            start_http_server -S -C "${CONFIG_FOLDER}${CERT_FILE}" -K "${CONFIG_FOLDER}${CERT_FILE}"
        fi
    else
        echo "[HTTPS] ⚠ Failed to load certificate, exit code: ${CERT_LOAD_STATUS}"
        echo "[SERVER] Falling back to HTTP for Web UI server"
        echo "[STARTUP] Starting Stremio server process"
        node server.js &
        SERVER_PID=$!
        echo "[STARTUP] Server started with PID: ${SERVER_PID}"
        start_http_server
    fi
else
    # Simple HTTP startup
    echo "[STARTUP] No HTTPS configuration found, using basic HTTP setup"
    echo "[STARTUP] Starting Stremio server process"
    node server.js &
    SERVER_PID=$!
    echo "[STARTUP] Server started with PID: ${SERVER_PID}"
    echo "[SERVER] Starting Web UI server with HTTP"
    start_http_server
fi

# Note: Original commented out line preserved
# wait ${SERVER_PID}
