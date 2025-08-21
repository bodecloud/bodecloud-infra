# VPN Fallback Solution

A Docker-in-Docker (DinD) based VPN failover system that provides automatic failover between multiple VPN configurations. The system uses a single privileged container that manages multiple VPN services internally and handles failover, while the host manages routing to direct Docker network traffic through the active VPN.

## Architecture

- **Host Script** (`vpn-solution.sh`): Runs on the host, manages the DinD container and sets up `ip route`/`ip rule` routing
- **DinD Container**: Runs Docker internally, manages VPN containers defined in your docker-compose.yml, handles failover
- **VPN Network**: Docker network that other containers can connect to for VPN access

## Features

- **Multiple VPN Support**: Works with Gluetun, WARP, and any containerized VPN
- **Automatic Failover**: Round-robin failover when VPNs fail health checks
- **Flexible Configuration**: Use any docker-compose.yml with VPN services
- **Host-level Routing**: Uses `iproute2` for transparent routing (like your original `vpn-up.sh`)
- **Easy Integration**: Other containers just connect to the `vpn-network`

## Quick Start

### 1. Build the VPN Manager Image

```bash
cd vpn-manager
docker build -t vpn-fallback:latest .
```

### 2. Create Your VPN Configuration

Create a `docker-compose.yml` file with your VPN services. Services will be tried in the order they appear (top to bottom). See `example-vpn-compose.yml` for a complete example.

```yaml
services:
  gluetun-de:
    image: ghcr.io/qdm12/gluetun
    cap_add: [NET_ADMIN]
    devices: ["/dev/net/tun:/dev/net/tun"]
    environment:
      VPN_SERVICE_PROVIDER: custom
      # ... your VPN config
  
  warp:
    image: caomingjun/warp:latest
    cap_add: [NET_ADMIN, MKNOD, AUDIT_WRITE]
    # ... your WARP config
```

### 3. Start the VPN Solution

```bash
# Make the script executable
chmod +x vpn-solution.sh

# Start with your compose file (requires root for routing)
sudo ./vpn-solution.sh start --compose-file ./my-vpns.yml
```

### 4. Connect Other Containers

```bash
# Run any container through the VPN
docker run --rm --network vpn-network curlimages/curl:latest https://api.ipify.org

# Or in docker-compose:
services:
  my-app:
    image: my-app:latest
    networks:
      - vpn-network

networks:
  vpn-network:
    external: true
```

## Usage Examples

### Basic Usage
```bash
# Start with default settings
sudo ./vpn-solution.sh start --compose-file ./vpns.yml

# Start with custom network
sudo ./vpn-solution.sh start --compose-file ./vpns.yml --network-name my-vpn --network-cidr 10.50.0.0/16

# Check status
sudo ./vpn-solution.sh status

# View logs
sudo ./vpn-solution.sh logs

# Stop
sudo ./vpn-solution.sh stop
```

### Using Compose Content Directly
```bash
# Pass compose content as environment variable
sudo ./vpn-solution.sh start --compose-content "$(cat my-vpns.yml)"
```

### Advanced Configuration
```bash
sudo ./vpn-solution.sh start \
  --compose-file ./vpns.yml \
  --network-name production-vpn \
  --network-cidr 172.20.0.0/16 \
  --image my-custom-vpn-fallback:latest
```

## How It Works

1. **Host Script**: Creates a Docker network and starts the DinD container
2. **DinD Container**: 
   - Reads your docker-compose.yml file
   - Extracts service names as VPN options (in order)
   - Starts the first VPN service
   - Sets up internal routing to forward traffic through the active VPN
3. **Host Routing**: Uses `iproute2` to route traffic from the Docker network through the DinD container
4. **Health Monitoring**: Continuously checks VPN health and switches to the next VPN on failure
5. **Failover**: Round-robin through VPN services when failures occur

## Environment Variables

You can set these in your docker-compose.yml or pass them to the container:

- `PREMIUMIZE_CUSTOMER_ID`: Your Premiumize customer ID
- `PREMIUMIZE_API_KEY`: Your Premiumize API key
- `WARP_LICENSE_KEY`: Your Cloudflare WARP license key
- `PREMIUMIZE_DE_IPV4_ADDRESS`: Custom German server IP
- `PREMIUMIZE_NL_IPV4_ADDRESS`: Custom Netherlands server IP
- `PREMIUMIZE_US_IPV4_ADDRESS`: Custom US server IP

## Commands

### vpn-solution.sh Commands

- `start`: Start the VPN solution
- `stop`: Stop the VPN solution  
- `restart`: Restart the VPN solution
- `status`: Show current status
- `logs`: Show VPN manager logs

### Required Options for Start

- `--compose-file <path>`: Path to docker-compose.yml with VPN services
- OR `--compose-content <yaml>`: Docker compose content as string

### Optional Options

- `--network-name <name>`: Docker network name (default: vpn-network)
- `--network-cidr <cidr>`: Docker network CIDR (default: 10.45.0.0/16)
- `--network-gateway <ip>`: Docker network gateway (default: 10.45.0.1)
- `--image <image>`: VPN manager image (default: vpn-fallback:latest)

## Troubleshooting

### Check Status
```bash
sudo ./vpn-solution.sh status
```

### View Logs
```bash
sudo ./vpn-solution.sh logs
```

### Test VPN Connection
```bash
# Test that traffic goes through VPN
docker run --rm --network vpn-network curlimages/curl:latest https://api.ipify.org
```

### Common Issues

1. **Permission Denied**: Script must be run as root for `ip route`/`ip rule` commands
2. **VPN Not Working**: Check that your docker-compose.yml VPN configurations are correct
3. **Container Can't Start**: Ensure the vpn-fallback image is built
4. **Network Issues**: Check that the Docker network was created correctly

## Requirements

- Linux host with Docker
- Root privileges (for `ip route`/`ip rule` commands)
- Docker Compose support in the DinD container

## Security Notes

- The DinD container runs privileged to manage VPN containers
- Host routing requires root privileges
- VPN credentials should be properly secured in your compose file

## License

This solution is provided as-is for educational and personal use. 