# Docker WARP Network Routing Script

A comprehensive bash script that sets up policy-based routing to route all traffic from a dedicated Docker network (`warpnet`) through a Cloudflare WARP container for enhanced privacy and security.

## Features

- **🔄 Automatic Container Deployment**: Automatically creates and configures the warp-with-nat container if it doesn't exist
- **🌐 Network Management**: Creates and manages Docker networks (`publicnet` and `warpnet`)
- **🛠️ Comprehensive Cleanup**: Thoroughly removes previous routing configurations
- **📊 Status Monitoring**: Real-time status checking and verification
- **🧪 Built-in Testing**: Automatic connectivity testing
- **⚙️ Highly Configurable**: Environment variable based configuration
- **🔒 Safety Features**: Confirmation prompts and rollback capabilities
- **📝 Detailed Logging**: Color-coded output with different verbosity levels

## Quick Start

```bash
# Make the script executable
chmod +x setup_warp_routing.sh

# Run with default settings (auto-creates container if needed)
sudo ./setup_warp_routing.sh

# Test the routing
docker run --rm --network=warpnet curlimages/curl:latest curl -s ifconfig.me
```

## Commands

| Command | Description |
|---------|-------------|
| `setup` | Set up routing (default action) |
| `cleanup` | Remove routing setup and optionally container |
| `status` | Show current configuration status |
| `test` | Test the routing setup |
| `help` | Show help message |

## Configuration

All aspects of the script can be configured via environment variables:

### Network Configuration

```bash
WARP_SUBNET="172.20.0.0/16"              # warpnet subnet
PUBLICNET_SUBNET="10.76.0.0/16"          # publicnet subnet  
PUBLICNET_GATEWAY="10.76.0.1"            # publicnet gateway
WARP_CONTAINER_IP="10.76.128.200"        # warp container IP
```

### Routing Configuration

```bash
ROUTING_TABLE_ID="200"                   # custom routing table ID
ROUTING_TABLE_NAME="warpnetvpn"          # custom routing table name
```

### Container Configuration

```bash
WARP_CONTAINER_NAME="warp-with-nat"      # container name
WARP_IMAGE="caomingjun/warp:latest"      # warp container image
CONFIG_PATH="./configs"                  # persistent config directory
WARP_SOCKS_PORT="5080"                   # SOCKS proxy port
```

### WARP Service Configuration

```bash
WARP_LICENSE_KEY="your_license_key"      # Cloudflare WARP license
TUNNEL_TOKEN="your_tunnel_token"         # Cloudflare tunnel token
```

### Script Behavior

```bash
DEBUG="false"                            # enable debug output
FORCE_RECREATE="false"                   # skip confirmation prompts
AUTO_START_CONTAINER="true"              # auto-create missing container
RUN_TESTS="true"                         # run comprehensive tests after setup
CLEANUP_ON_TEST_FAILURE="true"          # cleanup if tests fail
```

## Usage Examples

### Basic Setup

```bash
sudo ./setup_warp_routing.sh
```

### Setup with Custom Configuration

```bash
sudo WARP_SUBNET="192.168.100.0/24" WARP_CONTAINER_IP="10.76.128.201" ./setup_warp_routing.sh
```

### Force Setup (No Prompts)

```bash
sudo FORCE_RECREATE=true ./setup_warp_routing.sh
```

### Debug Mode

```bash
sudo DEBUG=true ./setup_warp_routing.sh
```

### Check Status

```bash
sudo ./setup_warp_routing.sh status
```

### Test Routing

```bash
sudo ./setup_warp_routing.sh test
```

### Setup Without Tests

```bash
sudo RUN_TESTS=false ./setup_warp_routing.sh
```

### Manual Testing

```bash
./test_warp_routing.sh  # Standalone test script (no sudo required)
```

### Cleanup Everything

```bash
sudo ./setup_warp_routing.sh cleanup
```

## How It Works

1. **Prerequisites Check**: Verifies all required tools are available
2. **Cleanup**: Removes any previous routing configurations
3. **Network Setup**: Creates/verifies Docker networks
4. **Container Management**: Deploys or verifies warp-with-nat container
5. **Routing Configuration**: Sets up policy-based routing using:
   - Custom routing table (`/etc/iproute2/rt_tables`)
   - Policy routing rules (`ip rule`)
   - Custom routes (`ip route`)
   - NAT masquerading (`iptables`)
6. **Verification**: Tests all components
7. **Testing**: Validates end-to-end connectivity

## Network Architecture

```mermaid
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   warpnet       │    │   publicnet      │    │   Internet      │
│ (172.20.0.0/16) │────│ (10.76.0.0/16)   │────│                 │
│                 │    │                  │    │                 │
│ Your containers │    │ warp-with-nat    │    │ WARP Egress     │
│                 │    │ (10.76.128.200)  │    │                 │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

### Traffic Flow

1. Container on `warpnet` sends traffic
2. Policy routing rule intercepts traffic from `172.20.0.0/16`
3. Traffic is routed via custom table to `warp-with-nat` container
4. NAT masquerading handles source translation
5. `warp-with-nat` forwards traffic through Cloudflare WARP
6. Return traffic follows the reverse path

## Container Deployment

The script automatically deploys the `warp-with-nat` container with the following configuration:

- **Image**: `caomingjun/warp:latest`
- **Network**: Connected to `publicnet` with static IP
- **Capabilities**: `MKNOD`, `AUDIT_WRITE`, `NET_ADMIN`
- **Sysctls**: Optimized for NAT forwarding
- **Ports**: SOCKS proxy on configurable port (default 5080)
- **Volumes**: Persistent WARP configuration
- **Environment**: Full WARP configuration with license key

## Testing Framework

The script includes a comprehensive testing suite that validates the routing setup without compromising host security. Tests run automatically after setup unless disabled.

### Automated Tests

The following tests run automatically during setup:

1. **Basic Connectivity Test**: Verifies containers can reach the internet through warpnet
2. **IP Routing Verification**: Confirms external IPs differ between host and containers
3. **Dynamic IP Assignment**: Validates containers receive unique IPs from the subnet
4. **Container Communication**: Tests inter-container communication within warpnet
5. **Routing Configuration**: Verifies all routing rules and tables are correctly configured
6. **WARP Container Health**: Checks the health and accessibility of the warp-with-nat container

### Test Configuration

```bash
# Run setup with tests (default)
sudo ./setup_warp_routing.sh

# Run setup without tests
sudo RUN_TESTS=false ./setup_warp_routing.sh

# Keep setup if tests fail (don't cleanup)
sudo CLEANUP_ON_TEST_FAILURE=false ./setup_warp_routing.sh

# Run only tests (requires existing setup)
sudo ./setup_warp_routing.sh test
```

### Manual Testing

Use the standalone test script for manual verification:

```bash
# Run comprehensive manual tests
./test_warp_routing.sh

# Quick manual test
docker run --rm --network=warpnet alpine sh -c 'apk add curl && curl ifconfig.me'
```

### Test Safety Features

- **Containerized Testing**: All tests run in isolated containers
- **No Host Modification**: Tests don't modify host networking
- **Automatic Cleanup**: Test containers are automatically removed
- **Failure Handling**: Optional cleanup on test failure
- **Non-Disruptive**: Tests don't affect existing containers

### Understanding Test Results

**✅ Success Indicators:**

- External IPs differ between host and containers
- Containers receive unique internal IPs from warpnet subnet
- All routing rules and tables are correctly configured
- DNS resolution works through warpnet

**❌ Failure Indicators:**

- Host and container external IPs are identical
- Containers cannot reach the internet
- Routing rules or NAT configuration missing
- IP assignment conflicts or subnet issues

## Troubleshooting

### Container Issues

```bash
# Check container status
docker ps | grep warp-with-nat

# View container logs
docker logs warp-with-nat

# Restart container
docker restart warp-with-nat
```

### Routing Issues

```bash
# Check routing rules
sudo ip rule list

# Check custom routing table
sudo ip route show table 200

# Check NAT rules
sudo iptables -t nat -L POSTROUTING -n
```

### Network Issues

```bash
# Test container connectivity
docker run --rm --network=warpnet curlimages/curl:latest curl -s ifconfig.me

# Compare with direct connection
curl -s ifconfig.me

# Check bridge interfaces
ip addr show | grep br-
```

### Debug Mode

Enable debug mode for detailed output:

```bash
sudo DEBUG=true ./setup_warp_routing.sh
```

## Safety Features

- **Confirmation Prompts**: Asks before destructive operations
- **Force Mode**: `FORCE_RECREATE=true` to skip prompts
- **Rollback**: `cleanup` command removes all changes
- **Verification**: Comprehensive validation of all components
- **Error Handling**: Graceful failure with informative messages

## Requirements

- **OS**: Linux with iproute2 and iptables
- **Docker**: Running Docker daemon
- **Privileges**: Root access (sudo)
- **Tools**: `ip`, `iptables`, `docker`, `ping`, `grep`, `awk`, `sed`

## License

This script is provided as-is for educational and operational purposes. Ensure you have proper licensing for Cloudflare WARP services.

## Contributing

Feel free to submit issues and enhancement requests. The script is designed to be modular and extensible.
