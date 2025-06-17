# Keycloak SSO Integration with Traefik

This directory contains the configuration files needed to set up Keycloak as an SSO provider for your Media Stack services, integrated with Traefik for authentication and rate limiting.

## Features

- Single Sign-On (SSO) with Google and Facebook
- Role-based access control
- Different rate limits for authenticated and unauthenticated users
- Custom rate limit error pages
- Integration with Traefik for centralized authentication

## Prerequisites

Before you begin, you'll need:

1. Google Developer Account and OAuth Credentials (for Google SSO)
2. Facebook Developer Account and OAuth App (for Facebook SSO)
3. A valid domain name with SSL/TLS certificates
4. Docker and Docker Compose installed

## Setup Instructions

### 1. Configure Environment Variables

1. Copy the `.env.example` file to `.env`:
   ```bash
   cp .env.example .env
   ```

2. Edit the `.env` file and update the following variables:
   - Domain settings
   - Keycloak admin credentials (use strong passwords)
   - Database credentials
   - Google and Facebook SSO credentials
   - SMTP settings for email notifications

### 2. Set Up Google SSO

1. Go to the [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Navigate to "APIs & Services" > "Credentials"
4. Click "Create Credentials" > "OAuth client ID"
5. Set the application type to "Web application"
6. Add authorized redirect URIs:
   ```
   https://auth.yourdomain.com/realms/MediaStack/broker/google/endpoint
   ```
7. Copy the Client ID and Client Secret to your `.env` file

### 3. Set Up Facebook SSO

1. Go to [Facebook Developers](https://developers.facebook.com/)
2. Create a new app or select an existing one
3. Add the "Facebook Login" product
4. In settings, add the OAuth redirect URI:
   ```
   https://auth.yourdomain.com/realms/MediaStack/broker/facebook/endpoint
   ```
5. Copy the App ID and App Secret to your `.env` file

### 4. Deploy the Stack

1. Start the services:
   ```bash
   docker-compose up -d
   ```

2. Verify Keycloak is running:
   ```bash
   docker-compose logs keycloak
   ```

3. Access Keycloak admin console at:
   ```
   https://auth.yourdomain.com/admin/
   ```

### 5. Configure Traefik Integration

The `dynamic.yaml` and `dynamic-routers.yaml` files in the Traefik config directory already contain the necessary middleware configurations. 

Key middleware chains:
- `chain-authenticated-default`: For authenticated users (higher rate limits)
- `chain-public-default`: For unauthenticated users (lower rate limits)

### 6. Testing the Configuration

1. Test authentication:
   ```
   https://auth.yourdomain.com/realms/MediaStack/account
   ```

2. Test rate limiting:
   - Make repeated requests to any protected endpoint
   - Authenticated users should get higher limits
   - Unauthenticated users will see the rate limit error page sooner

## Rate Limit Configuration

Rate limits are configured as follows:

- Unauthenticated users: 30 requests per minute with burst of 20
- Authenticated users: 100 requests per minute with burst of 50
- Premium users: Can have custom rate limits (configurable)

## Troubleshooting

1. Check container logs:
   ```bash
   docker-compose logs keycloak
   docker-compose logs traefik-forward-auth
   ```

2. Verify Keycloak realm import:
   ```bash
   docker-compose logs keycloak | grep "import"
   ```

3. Test authentication endpoints:
   ```bash
   curl -I https://auth.yourdomain.com/realms/MediaStack/.well-known/openid-configuration
   ```

## Security Recommendations

1. Always use strong, unique passwords for admin accounts
2. Keep your Google and Facebook OAuth credentials secure
3. Regularly audit user roles and permissions
4. Consider enabling MFA for administrative accounts
5. Regularly update container images for security patches

## Additional Resources

- [Keycloak Documentation](https://www.keycloak.org/documentation)
- [Traefik Forward Auth](https://github.com/thomseddon/traefik-forward-auth)
- [Google OAuth Setup Guide](https://developers.google.com/identity/protocols/oauth2/web-server)
- [Facebook Login Integration](https://developers.facebook.com/docs/facebook-login/web) 