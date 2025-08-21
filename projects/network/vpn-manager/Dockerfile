FROM docker:25-dind

# Install required packages
RUN apk add --no-cache \
    bash \
    curl \
    jq \
    iproute2 \
    iptables \
    openrc \
    supervisor \
    gettext \
    docker-compose \
    yq

# Create necessary directories
RUN mkdir -p /app /var/log/supervisor /etc/supervisor/conf.d

# Copy our management scripts
COPY scripts/ /app/
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf

# Make scripts executable
RUN chmod +x /app/*.sh

# Environment variables for configuration
ENV VPN_COMPOSE_FILE=""
ENV VPN_COMPOSE=""
ENV DOCKER_NETWORK_NAME="vpn-network"
ENV DOCKER_NETWORK_CIDR="10.45.0.0/16"
ENV HEALTH_CHECK_INTERVAL="60"
ENV MAX_FAILURES="3"

# Create entrypoint script that handles VPN_COMPOSE_FILE
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

WORKDIR /app

# Use our entrypoint to handle compose file loading
ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"] 