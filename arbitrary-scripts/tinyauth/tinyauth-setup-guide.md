# TinyAuth Setup Guide

This guide will walk you through setting up TinyAuth with Google and GitHub OAuth authentication for your media stack.

## 1. Creating Local Users

You can create local users using the provided `create_tinyauth_user.sh` script:

```bash
./create_tinyauth_user.sh
```

This will:

- Generate a properly hashed password
- Optionally set up TOTP (Two-Factor Authentication)
- Provide you with the string to add to your .env file

## 2. Setting Up Google OAuth

### Step 1: Create Google Cloud Project

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Navigate to "APIs & Services" > "Credentials"
4. Click "Create Credentials" and select "OAuth client ID"

### Step 2: Configure OAuth Consent Screen

1. Configure the OAuth consent screen with required information
2. Add your email address as a test user
3. Set the scopes to include "email" and "profile"

### Step 3: Create OAuth Client ID

1. Select "Web application" as the application type
2. Add a name for your client (e.g., "TinyAuth")
3. Add authorized redirect URIs:
   - `https://auth.yourdomain.com/api/auth/google/callback`
4. Click "Create"
5. Note the Client ID and Client Secret

### Step 4: Update Environment Variables

Add these to your `.env` file:

```shell
GOOGLE_CLIENT_ID=your_client_id
GOOGLE_CLIENT_SECRET=your_client_secret
```

## 3. Setting Up GitHub OAuth

### Step 1: Create GitHub OAuth App

1. Go to GitHub [Developer settings](https://github.com/settings/developers)
2. Click "New OAuth App"
3. Fill in the details:
   - Application name: e.g., "My Media Stack Auth"
   - Homepage URL: `https://auth.yourdomain.com`
   - Authorization callback URL: `https://auth.yourdomain.com/api/auth/github/callback`
4. Click "Register application"
5. Generate a new client secret

### Step 2: Update Environment Variables

Add these to your `.env` file:

```shell
GITHUB_CLIENT_ID=your_client_id
GITHUB_CLIENT_SECRET=your_client_secret
```

## 4. Additional Configuration Options

### Restricting OAuth Users

To restrict which users can log in with OAuth, use the `OAUTH_WHITELIST` variable:

```shell
# Only allow specific users (comma-separated, no spaces)
OAUTH_WHITELIST=user1@example.com,user2@example.com

# You can also use regex patterns with slashes
OAUTH_WHITELIST=/^admin.*/,/.*@yourdomain\.com$/
```

### Auto-Redirect to OAuth Provider

If you prefer to automatically redirect users to a specific OAuth provider, set:

```shell
OAUTH_AUTO_REDIRECT=google
```

Options include: `none`, `github`, `google`, or `generic`

### TOTP (Two-Factor Authentication)

TOTP can be enabled for local users by running:

```bash
./create_tinyauth_user.sh
```

And selecting "y" when asked about adding TOTP.

### Session Expiry

Adjust the session expiry time (in seconds):

```shell
SESSION_EXPIRY=86400  # 24 hours
```

## 5. Integrating with Services

To protect a service with TinyAuth authentication, add these labels in your docker-compose.yml:

```yaml
traefik.http.routers.your-service.middlewares: tinyauth
```

## 6. Troubleshooting

### Cookie Issues

If you're having cookie issues:

- Ensure APP_URL matches your actual domain
- Check that your services are on the same parent domain as TinyAuth
- Verify COOKIE_SECURE is set correctly for your environment

### OAuth Callback Errors

If OAuth callbacks are failing:

- Double-check the callback URLs match exactly
- Ensure your domain is accessible
- Check for any typos in client IDs or secrets

### Login Problems

If users can't log in:

- Verify the USERS variable is correctly formatted
- Check that OAuth providers are properly configured
- Look at the TinyAuth logs for specific error messages

## 7. Security Recommendations

- Use a strong SECRET value (generate with `openssl rand -base64 32 | tr -dc 'a-zA-Z0-9' | head -c 32`)
- Enable TOTP for added security
- Set COOKIE_SECURE=true if using HTTPS
- Configure LOGIN_MAX_RETRIES and LOGIN_TIMEOUT to prevent brute force attacks
- Use OAUTH_WHITELIST to restrict access to trusted users only
