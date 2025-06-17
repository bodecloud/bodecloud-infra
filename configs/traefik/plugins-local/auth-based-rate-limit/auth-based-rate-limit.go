// Package authbasedratelmit is a middleware plugin for Traefik
// that applies different rate limit middlewares based on auth status.
package authbasedratelmit

import (
	"context"
	"net/http"
)

// Config holds the plugin configuration.
type Config struct {
	// AuthHeader is the header to check for auth status
	AuthHeader string `json:"authHeader,omitempty"`
	// AuthenticatedValue is the value in the auth header that indicates a user is authenticated
	AuthenticatedValue string `json:"authenticatedValue,omitempty"`
	// AuthenticatedMiddleware is the middleware to use for authenticated users
	AuthenticatedMiddleware string `json:"authenticatedMiddleware,omitempty"`
	// UnauthenticatedMiddleware is the middleware to use for unauthenticated users
	UnauthenticatedMiddleware string `json:"unauthenticatedMiddleware,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		AuthHeader:                "X-Auth-Status",
		AuthenticatedValue:        "authenticated",
		AuthenticatedMiddleware:   "chain-authenticated-default@file",
		UnauthenticatedMiddleware: "chain-public-default@file",
	}
}

// AuthBasedRateLimit is a middleware that applies different rate limits
// based on authentication status
type AuthBasedRateLimit struct {
	next                      http.Handler
	authHeader                string
	authenticatedValue        string
	authenticatedMiddleware   string
	unauthenticatedMiddleware string
	name                      string
}

// New creates a new auth-based rate limit middleware.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &AuthBasedRateLimit{
		next:                      next,
		authHeader:                config.AuthHeader,
		authenticatedValue:        config.AuthenticatedValue,
		authenticatedMiddleware:   config.AuthenticatedMiddleware,
		unauthenticatedMiddleware: config.UnauthenticatedMiddleware,
		name:                      name,
	}, nil
}

// ServeHTTP implements http.Handler
func (a *AuthBasedRateLimit) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// Check if user is authenticated based on header
	authValue := req.Header.Get(a.authHeader)
	
	// Set header to indicate which middleware should be applied
	if authValue == a.authenticatedValue {
		// User is authenticated, use authenticated middleware
		req.Header.Set("X-Use-Middleware", a.authenticatedMiddleware)
	} else {
		// User is not authenticated, use public middleware
		req.Header.Set("X-Use-Middleware", a.unauthenticatedMiddleware)
	}
	
	// Pass to next middleware
	a.next.ServeHTTP(rw, req)
} 