# Traefik v2 OR Authentication Pattern
This solution implements an OR authentication pattern in Traefik v2, allowing access if ANY of these conditions are met:
- Request comes from a whitelisted IP
- Request contains a valid API key
- Request passes TinyAuth authentication

## Implementation (Docker Compose Labels)

```yaml
# Router 1: Default with TinyAuth (lowest priority)
- traefik.http.routers.myservice.rule=Host(`myservice.example.com`)
- traefik.http.routers.myservice.priority=1
- traefik.http.routers.myservice.middlewares=tinyauth@docker
- traefik.http.routers.myservice.service=myservice

# Router 2: IP Allowlist (higher priority)
