# Nginx Multi-Middleware Authentication System

This nginx container provides **multiple middleware endpoints** in a single service, eliminating the need for separate nginx containers for different middleware types.

## 🎯 Problem Solved

Instead of needing separate nginx containers for different middleware types, you can now use **one nginx container** that provides multiple middleware endpoints while keeping the original authentication logic.

## 📋 Available Middlewares

| Middleware | Endpoint | Description | Use Case |
|------------|----------|-------------|----------|
| `nginx-auth@file` | `/auth` | **Standard Auth**: API key OR IP whitelist OR tinyauth | Main authentication (original logic) |
| `nginx-ratelimit@file` | `/ratelimit` | **Rate Limiting Only**: Just rate limiting, no auth | Public APIs, DDoS protection |
| `nginx-auth-with-ratelimit@file` | `/auth-with-ratelimit` | **Auth + Rate Limiting**: Standard auth with rate limiting | High-traffic authenticated services |

## 🔍 Authentication Logic

### Main Auth Endpoint (`/auth`) - Original Logic

This is your **main authentication endpoint** with the original OR logic:

1. **Check API Key**: If valid `X-API-Key` header → ✅ Allow immediately
2. **Check IP Whitelist**: If client IP is whitelisted → ✅ Allow immediately  
3. **Fallback to TinyAuth**: If neither passes → Proxy to tinyauth service

**Only ONE** of these methods needs to pass - they use OR logic, not AND logic.

## 🚀 Usage Examples

### Standard Authentication (Most Common)

```yaml
# docker-compose.yml
services:
  my-service:
    # ... service config ...
    labels:
      traefik.enable: "true"
      traefik.http.routers.my-service.rule: Host(`my-service.example.com`)
      traefik.http.routers.my-service.middlewares: nginx-auth@file
      traefik.http.services.my-service.loadbalancer.server.port: 8080
```

### Public API with Rate Limiting Only

```yaml
# docker-compose.yml
services:
  public-api:
    # ... service config ...
    labels:
      traefik.enable: "true"
      traefik.http.routers.public-api.rule: Host(`public-api.example.com`)
      traefik.http.routers.public-api.middlewares: nginx-ratelimit@file
      traefik.http.services.public-api.loadbalancer.server.port: 8080
```

### High-Traffic Service (Auth + Rate Limiting)

```yaml
# docker-compose.yml
services:
  high-traffic-service:
    # ... service config ...
    labels:
      traefik.enable: "true"
      traefik.http.routers.high-traffic-service.rule: Host(`app.example.com`)
      traefik.http.routers.high-traffic-service.middlewares: nginx-auth-with-ratelimit@file
      traefik.http.services.high-traffic-service.loadbalancer.server.port: 3000
```

### Multiple Middlewares (Chain)

```yaml
# docker-compose.yml
services:
  complex-service:
    # ... service config ...
    labels:
      traefik.enable: "true"
      traefik.http.routers.complex-service.rule: Host(`complex.example.com`)
      # Chain multiple middlewares together
      traefik.http.routers.complex-service.middlewares: nginx-auth@file,compress@file
      traefik.http.services.complex-service.loadbalancer.server.port: 8080
```

## 🔧 Configuration

### API Keys

Edit `configs/nginx-auth/nginx.conf` and update the `map $http_x_api_key $api_key_valid` section:

```nginx
map $http_x_api_key $api_key_valid {
    default 0;
    
    # Valid API keys
    "your-api-key-1" 1;
    "your-api-key-2" 1;
    # Add more API keys here as needed
}
```

### IP Whitelist

Edit `configs/nginx-auth/nginx.conf` and update the `geo $ip_whitelisted` section:

```nginx
geo $ip_whitelisted {
    default 0;
    
    # Trusted network ranges
    127.0.0.1/32    1;  # localhost (ipv4)
    ::1/128         1;  # localhost (ipv6)
    172.16.0.0/12   1;  # docker network
    10.76.0.0/16    1;  # docker 'publicnet' network
    100.64.0.0/10   1;  # tailscale network
    192.168.1.100/32 1; # specific whitelisted IP
}
```

### Rate Limiting

Rate limits are configured in the nginx.conf:

```nginx
# Rate limiting zones
limit_req_zone $binary_remote_addr zone=auth:10m rate=10r/s;
limit_req_zone $binary_remote_addr zone=ratelimit:10m rate=5r/s;
```

## 📊 Monitoring & Debugging

### Health Check

```bash
curl http://nginx-auth.example.com/health
# Response: nginx multi-middleware service healthy
```

### Available Middlewares

```bash
curl http://nginx-auth.example.com/middlewares
# Returns JSON with all available middlewares and usage instructions
```

### Response Headers

All middlewares add these headers for debugging:

- `X-Auth-Method`: How the request was authenticated (`api_key`, `ip_whitelist`, `tinyauth`, `ratelimit`)
- `X-Auth-Passed`: Whether authentication passed (`true`, `false`)
- `X-Middleware-Name`: Which middleware was used (`auth`, `ratelimit`, `auth-with-ratelimit`)

### Logs

Check nginx logs for authentication details:

```bash
# View auth logs
docker exec nginx-auth tail -f /var/log/nginx/auth.log

# View access logs  
docker exec nginx-auth tail -f /var/log/nginx/access.log
```

## 🔄 Testing Authentication Methods

### 1. Test API Key Authentication

```bash
# Should work with valid API key
curl -H "X-API-Key: your-api-key" https://protected.yourdomain.com/
# → Immediate access, header: X-Auth-Method: api_key

# Should fail with invalid API key, fall back to TinyAuth
curl -H "X-API-Key: invalid-key" https://protected.yourdomain.com/
# → Redirects to TinyAuth login
```

### 2. Test IP Whitelist

```bash
# From whitelisted IP (should work immediately)
curl https://protected.yourdomain.com/
# → header: X-Auth-Method: ip_whitelist

# From non-whitelisted IP (should fall back to TinyAuth)
# → Redirects to TinyAuth login
```

### 3. Test TinyAuth Fallback

```bash
# No API key, not whitelisted IP - should redirect to TinyAuth
curl https://protected.yourdomain.com/
# → Redirects to TinyAuth login
# → After login: X-Auth-Method: tinyauth
```

### 4. Test Rate Limiting

```bash
# Test rate limiting endpoint
for i in {1..20}; do curl https://ratelimited.yourdomain.com/; done
# → Should start returning 429 after hitting limits
```

## 🛠️ Adding New Middlewares

If you need additional middleware types in the future:

1. **Add new location block** in `configs/nginx-auth/nginx.conf`:

```nginx
location /custom-middleware {
    set $middleware_name "custom-middleware";
    # ... your custom logic ...
    return 200 "OK";
}
```

2. **Add new middleware** in `configs/traefik/config/dynamic.yaml`:

```yaml
nginx-custom:
  forwardAuth:
    address: "http://nginx-auth:80/custom-middleware"
    trustForwardHeader: true
    authResponseHeaders:
      - "X-Auth-Method"
      - "X-Auth-Passed"
      - "X-Middleware-Name"
```

3. **Use in services**:

```yaml
traefik.http.routers.my-service.middlewares: nginx-custom@file
```

## 🚨 Troubleshooting

### Common Issues

1. **403 Forbidden**: Check IP whitelist configuration
2. **401 Unauthorized**: Check API key configuration or tinyauth service
3. **429 Rate Limited**: Adjust rate limiting zones
4. **502 Bad Gateway**: Check if tinyauth service is running

## 📈 Benefits

- **Single Container**: Reduced resource usage compared to multiple nginx containers
- **Original Logic Preserved**: Main auth endpoint keeps the OR logic (API key OR IP whitelist OR tinyauth)
- **Flexible**: Add new middleware types without new containers
- **Service Agnostic**: No hardcoded service names anywhere
- **Easy Management**: Single configuration file to maintain

This approach gives you the flexibility to add new middleware types while keeping your original authentication logic intact in the main `/auth` endpoint.
