#!/bin/bash

# Script to update the main docker-compose.yml file with Keycloak SSO configuration

# Check if we're in the right directory
if [ ! -f ./realm.json ]; then
  echo "Error: Please run this script from the keycloak-new directory where realm.json is located"
  exit 1
fi

# Create backup of existing docker-compose.yml
echo "Creating backup of existing docker-compose.yml..."
cp ../../docker-compose.yml ../../docker-compose.yml.bak.$(date +%Y%m%d%H%M%S)

# Function to check if keycloak is already in docker-compose.yml
check_keycloak_exists() {
  grep -q "  keycloak:" ../../docker-compose.yml
  return $?
}

# Function to add environment variables to .env file
update_env_file() {
  echo "Updating .env file with Keycloak variables..."
  
  # Check if .env exists
  if [ ! -f ../../.env ]; then
    echo "Error: .env file not found in the parent directory"
    exit 1
  fi
  
  # Create backup
  cp ../../.env ../../.env.bak.$(date +%Y%m%d%H%M%S)
  
  # Add the variables if they don't exist
  grep -q "GOOGLE_CLIENT_ID" ../../.env || echo "GOOGLE_CLIENT_ID=" >> ../../.env
  grep -q "GOOGLE_CLIENT_SECRET" ../../.env || echo "GOOGLE_CLIENT_SECRET=" >> ../../.env
  grep -q "FACEBOOK_CLIENT_ID" ../../.env || echo "FACEBOOK_CLIENT_ID=" >> ../../.env
  grep -q "FACEBOOK_CLIENT_SECRET" ../../.env || echo "FACEBOOK_CLIENT_SECRET=" >> ../../.env
  grep -q "TRAEFIK_AUTH_SECRET" ../../.env || echo "TRAEFIK_AUTH_SECRET=change-me-please" >> ../../.env
  grep -q "FORWARD_AUTH_SECRET" ../../.env || echo "FORWARD_AUTH_SECRET=change-me-please" >> ../../.env
  grep -q "KEYCLOAK_ADMIN_PASSWORD" ../../.env || echo "KEYCLOAK_ADMIN_PASSWORD=admin" >> ../../.env
  grep -q "KEYCLOAK_DB_PASSWORD" ../../.env || echo "KEYCLOAK_DB_PASSWORD=keycloak" >> ../../.env
  
  echo ".env file updated successfully"
}

# Copy realm.json to the main keycloak directory
copy_realm_json() {
  echo "Copying realm.json to the main keycloak directory..."
  mkdir -p ../../configs/keycloak
  cp ./realm.json ../../configs/keycloak/realm.json
  echo "realm.json copied successfully"
}

# Main execution
echo "Starting update process..."

# Update environment variables
update_env_file

# Copy realm.json
copy_realm_json

# Check if Keycloak is already in docker-compose.yml
if check_keycloak_exists; then
  echo "Keycloak service already exists in docker-compose.yml"
  echo "To update it, manually edit the file or replace the Keycloak section with the configuration from docker-compose.yml in this directory"
else
  echo "Adding Keycloak services to docker-compose.yml..."
  
  # Create a temporary file with the services to add
  cat > /tmp/keycloak-services.yml << 'EOF'
  keycloak:
    image: quay.io/keycloak/keycloak:latest
    container_name: keycloak
    hostname: keycloak
    networks:
      - infranet
    ports:
      - "${KEYCLOAK_PORT:-8074}:8443"
      - "${KEYCLOAK_PORT2:-8783}:9000"
    volumes:
      - ${CONFIG_DIR:-./configs}/keycloak/realm.json:/opt/keycloak/data/import/realm.json:ro
    environment:
      <<: *common-env
      KEYCLOAK_ADMIN: admin
      KEYCLOAK_ADMIN_PASSWORD: ${KEYCLOAK_ADMIN_PASSWORD:-admin}
      KC_DB: postgres
      KC_DB_URL: jdbc:postgresql://keycloak-db:5432/keycloak
      KC_DB_USERNAME: ${KEYCLOAK_DB_USER:-keycloak}
      KC_DB_PASSWORD: ${KEYCLOAK_DB_PASSWORD:-keycloak}
      KC_HOSTNAME: auth.${DOMAIN}
      KC_HOSTNAME_STRICT: "false"
      KC_HTTP_ENABLED: "true"
      KC_PROXY: edge
      GOOGLE_CLIENT_ID: ${GOOGLE_CLIENT_ID:-}
      GOOGLE_CLIENT_SECRET: ${GOOGLE_CLIENT_SECRET:-}
      FACEBOOK_CLIENT_ID: ${FACEBOOK_CLIENT_ID:-}
      FACEBOOK_CLIENT_SECRET: ${FACEBOOK_CLIENT_SECRET:-}
      CLIENT_SECRET: ${CLIENT_SECRET:-change-me}
      TRAEFIK_AUTH_SECRET: ${TRAEFIK_AUTH_SECRET:-change-me}
    depends_on:
      - keycloak-db
    labels:
      traefik.enable: "true"
      traefik.http.routers.keycloak.entrypoints: websecure
      traefik.http.routers.keycloak.tls: true
      traefik.http.routers.keycloak.rule: Host(`auth.${DOMAIN}`) || Host(`auth.${DUCKDNS_SUBDOMAIN}.duckdns.org`)
      traefik.http.routers.keycloak.tls.certresolver: default
      traefik.http.services.keycloak.loadbalancer.server.port: 8080
      traefik.http.middlewares.keycloak-rate-limit.ratelimit.average: 100
      traefik.http.middlewares.keycloak-rate-limit.ratelimit.burst: 50
      traefik.http.middlewares.keycloak-rate-limit.ratelimit.period: 1m
      traefik.http.routers.keycloak.middlewares: keycloak-rate-limit
    command:
      - start-dev
      - --import-realm
    restart: unless-stopped

  keycloak-db:
    image: postgres:14-alpine
    container_name: keycloak-db
    hostname: keycloak-db
    networks:
      - infranet
    environment:
      POSTGRES_DB: keycloak
      POSTGRES_USER: ${KEYCLOAK_DB_USER:-keycloak}
      POSTGRES_PASSWORD: ${KEYCLOAK_DB_PASSWORD:-keycloak}
    volumes:
      - keycloak_db_data:/var/lib/postgresql/data
    restart: unless-stopped
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U keycloak"]
      interval: 10s
      timeout: 5s
      retries: 5

  traefik-forward-auth:
    image: thomseddon/traefik-forward-auth:latest
    container_name: traefik-forward-auth
    hostname: traefik-forward-auth
    networks:
      - infranet
    environment:
      DEFAULT_PROVIDER: oidc
      PROVIDERS_OIDC_ISSUER_URL: https://auth.${DOMAIN}/realms/MediaStack
      PROVIDERS_OIDC_CLIENT_ID: traefik-forward-auth
      PROVIDERS_OIDC_CLIENT_SECRET: ${TRAEFIK_AUTH_SECRET:-change-me}
      SECRET: ${FORWARD_AUTH_SECRET:-change-me}
      INSECURE_COOKIE: "false"
      COOKIE_DOMAIN: ${DOMAIN}
      AUTH_HOST: auth.${DOMAIN}
      URL_PATH: /_oauth
      COOKIE_NAME: mediastack_auth
      LIFETIME: 3600
      LOG_LEVEL: debug
    restart: unless-stopped
    labels:
      traefik.enable: "true"
      traefik.http.routers.forward-auth.rule: Host(`auth.${DOMAIN}`) && PathPrefix(`/_oauth`)
      traefik.http.routers.forward-auth.entrypoints: websecure
      traefik.http.routers.forward-auth.tls: true
      traefik.http.services.forward-auth.loadbalancer.server.port: 4181
      traefik.docker.network: infranet
EOF

  # Append the services to docker-compose.yml
  cat /tmp/keycloak-services.yml >> ../../docker-compose.yml
  
  # Add the volume for keycloak-db
  if ! grep -q "keycloak_db_data:" ../../docker-compose.yml; then
    echo "Adding volumes for Keycloak..."
    echo "volumes:" >> ../../docker-compose.yml
    echo "  keycloak_db_data:" >> ../../docker-compose.yml
    echo "    name: keycloak_db_data" >> ../../docker-compose.yml
  fi
  
  echo "Keycloak services added successfully to docker-compose.yml"
fi

echo "Update process completed successfully!"
echo "Please review the changes and start the services with: docker-compose up -d keycloak keycloak-db traefik-forward-auth" 