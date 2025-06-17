# Nomad Variables for Media Stack
# This file defines all variables used across the Nomad job files
# Based on the environment variables from the original Docker Compose configuration

# Common Environment Variables
variable "tz" {
  description = "Timezone"
  type        = string
  default     = "America/Chicago"
}

variable "puid" {
  description = "Process User ID"
  type        = string
  default     = "1002"
}

variable "pgid" {
  description = "Process Group ID"
  type        = string
  default     = "988"
}

variable "umask" {
  description = "File creation mask"
  type        = string
  default     = "002"
}

# Path Configuration
variable "config_path" {
  description = "Configuration path"
  type        = string
  default     = "./configs"
}

variable "certs_path" {
  description = "Certificates path"
  type        = string
  default     = "./certs"
}

variable "root_dir" {
  description = "Root directory"
  type        = string
  default     = "."
}

# Domain Configuration
variable "domain" {
  description = "Primary domain"
  type        = string
  default     = "example.com"
}

variable "duckdns_subdomain" {
  description = "DuckDNS subdomain"
  type        = string
  default     = "example"
}

variable "ts_hostname" {
  description = "Tailscale hostname"
  type        = string
  default     = "example"
}

# Network Configuration
variable "publicnet_subnet" {
  description = "Public network subnet"
  type        = string
  default     = "10.76.0.0/16"
}

variable "publicnet_gateway" {
  description = "Public network gateway"
  type        = string
  default     = "10.76.0.1"
}

variable "publicnet_ip_range" {
  description = "Public network IP range"
  type        = string
  default     = "10.76.0.0/16"
}

variable "tailscale_cidr" {
  description = "Tailscale CIDR"
  type        = string
  default     = "100.64.0.0/10"
}

# Service IP Addresses
variable "mongodb_ipv4_address" {
  description = "MongoDB IPv4 address"
  type        = string
  default     = "10.76.0.50"
}

variable "redis_ipv4_address" {
  description = "Redis IPv4 address"
  type        = string
  default     = "10.76.128.87"
}

variable "qdrant_ipv4_address" {
  description = "Qdrant IPv4 address"
  type        = string
  default     = "10.76.128.44"
}

variable "traefik_ipv4_address" {
  description = "Traefik IPv4 address"
  type        = string
  default     = "10.76.128.85"
}

variable "traefik_error_pages_ipv4_address" {
  description = "Traefik Error Pages IPv4 address"
  type        = string
  default     = "10.76.128.84"
}

variable "watchtower_ipv4_address" {
  description = "Watchtower IPv4 address"
  type        = string
  default     = "10.76.128.83"
}

variable "tinyauth_ipv4_address" {
  description = "TinyAuth IPv4 address"
  type        = string
  default     = "10.76.128.82"
}

variable "whoami_ipv4_address" {
  description = "Whoami IPv4 address"
  type        = string
  default     = "10.76.128.81"
}

variable "code_demo_ipv4_address" {
  description = "Code Demo IPv4 address"
  type        = string
  default     = "10.76.128.80"
}

variable "searxng_ipv4_address" {
  description = "SearxNG IPv4 address"
  type        = string
  default     = "10.76.128.90"
}

variable "dozzle_ipv4_address" {
  description = "Dozzle IPv4 address"
  type        = string
  default     = "10.76.128.89"
}

variable "homepage_ipv4_address" {
  description = "Homepage IPv4 address"
  type        = string
  default     = "10.76.128.88"
}

variable "speedtest_ipv4_address" {
  description = "Speedtest IPv4 address"
  type        = string
  default     = "10.76.128.86"
}

variable "code_dev_ipv4_address" {
  description = "Code Dev IPv4 address"
  type        = string
  default     = "10.76.128.92"
}

variable "flaresolverr_ipv4_address" {
  description = "FlareSolverr IPv4 address"
  type        = string
  default     = "10.76.128.93"
}

variable "nginx_auth_ipv4_address" {
  description = "Nginx Auth IPv4 address"
  type        = string
  default     = "10.76.128.94"
}

variable "warp_ipv4_address" {
  description = "WARP IPv4 address"
  type        = string
  default     = "10.76.128.97"
}

variable "gpt_researcher_ipv4_address" {
  description = "GPT Researcher IPv4 address"
  type        = string
  default     = "10.76.128.43"
}

variable "lobechat_ipv4_address" {
  description = "LobeChat IPv4 address"
  type        = string
  default     = "10.76.128.46"
}

variable "dash_ipv4_address" {
  description = "Dash IPv4 address"
  type        = string
  default     = "10.76.128.23"
}

# Service Hostnames
variable "mongodb_hostname" {
  description = "MongoDB hostname"
  type        = string
  default     = "mongodb"
}

variable "searxng_hostname" {
  description = "SearxNG hostname"
  type        = string
  default     = "searxng"
}

variable "gpt_researcher_hostname" {
  description = "GPT Researcher hostname"
  type        = string
  default     = "gptr"
}

variable "lobechat_hostname" {
  description = "LobeChat hostname"
  type        = string
  default     = "lobechat"
}

# SSL/TLS Configuration
variable "lets_encrypt_email" {
  description = "Let's Encrypt email"
  type        = string
  default     = "admin@example.com"
}

variable "cloudflare_email" {
  description = "Cloudflare email"
  type        = string
  default     = ""
}

variable "cloudflare_dns_api_token" {
  description = "Cloudflare DNS API token"
  type        = string
  default     = ""
}

variable "cloudflare_zone_api_token" {
  description = "Cloudflare Zone API token"
  type        = string
  default     = ""
}

variable "duckdns_token" {
  description = "DuckDNS token"
  type        = string
  default     = ""
}

# Watchtower Configuration
variable "watchtower_cleanup" {
  description = "Watchtower cleanup"
  type        = string
  default     = "true"
}

variable "watchtower_schedule" {
  description = "Watchtower schedule"
  type        = string
  default     = "0 0 6 * * *"
}

variable "watchtower_notification_url" {
  description = "Watchtower notification URL"
  type        = string
  default     = ""
}

variable "watchtower_notification_report" {
  description = "Watchtower notification report"
  type        = string
  default     = "true"
}

# DeUnhealth Configuration
variable "deunhealth_log_level" {
  description = "DeUnhealth log level"
  type        = string
  default     = "debug"
}

variable "deunhealth_health_server_address" {
  description = "DeUnhealth health server address"
  type        = string
  default     = "127.0.0.1:9999"
}

# SearxNG Configuration
variable "searxng_url" {
  description = "SearxNG base URL"
  type        = string
  default     = "http://searxng:8080"
}

# Homepage Configuration
variable "homepage_allowed_hosts" {
  description = "Homepage allowed hosts"
  type        = string
  default     = "*"
}

variable "homepage_var_title" {
  description = "Homepage title"
  type        = string
  default     = "Bolabaden"
}

variable "homepage_var_search_provider" {
  description = "Homepage search provider"
  type        = string
  default     = "google"
}

variable "homepage_var_header_style" {
  description = "Homepage header style"
  type        = string
  default     = ""
}

variable "homepage_var_weather_city" {
  description = "Homepage weather city"
  type        = string
  default     = "Chicago"
}

variable "homepage_var_weather_lat" {
  description = "Homepage weather latitude"
  type        = string
  default     = "41.8781"
}

variable "homepage_var_weather_long" {
  description = "Homepage weather longitude"
  type        = string
  default     = "-87.6298"
}

variable "homepage_var_weather_unit" {
  description = "Homepage weather unit"
  type        = string
  default     = "fahrenheit"
}

# Speedtest Tracker Configuration
variable "speedtest_tracker_admin_email" {
  description = "Speedtest Tracker admin email"
  type        = string
  default     = ""
}

variable "speedtest_tracker_admin_name" {
  description = "Speedtest Tracker admin name"
  type        = string
  default     = ""
}

variable "speedtest_tracker_admin_password" {
  description = "Speedtest Tracker admin password"
  type        = string
  default     = "b00tstr4p"
}

variable "speedtest_tracker_api_rate_limit" {
  description = "Speedtest Tracker API rate limit"
  type        = string
  default     = "60"
}

variable "speedtest_tracker_app_key" {
  description = "Speedtest Tracker app key"
  type        = string
  default     = ""
}

variable "speedtest_tracker_app_name" {
  description = "Speedtest Tracker app name"
  type        = string
  default     = "Speedtest Tracker"
}

variable "speedtest_tracker_app_timezone" {
  description = "Speedtest Tracker app timezone"
  type        = string
  default     = "America/Chicago"
}

variable "speedtest_tracker_app_url" {
  description = "Speedtest Tracker app URL"
  type        = string
  default     = ""
}

variable "speedtest_tracker_asset_url" {
  description = "Speedtest Tracker asset URL"
  type        = string
  default     = ""
}

variable "speedtest_tracker_chart_begin_at_zero" {
  description = "Speedtest Tracker chart begin at zero"
  type        = string
  default     = "true"
}

variable "speedtest_tracker_chart_datetime_format" {
  description = "Speedtest Tracker chart datetime format"
  type        = string
  default     = "j/m G:i"
}

variable "speedtest_tracker_content_width" {
  description = "Speedtest Tracker content width"
  type        = string
  default     = "7xl"
}

variable "speedtest_tracker_datetime_format" {
  description = "Speedtest Tracker datetime format"
  type        = string
  default     = "j M Y, G:i:s"
}

variable "speedtest_tracker_db_connection" {
  description = "Speedtest Tracker DB connection"
  type        = string
  default     = "sqlite"
}

variable "speedtest_tracker_display_timezone" {
  description = "Speedtest Tracker display timezone"
  type        = string
  default     = "America/Chicago"
}

variable "speedtest_tracker_prune_results_older_than" {
  description = "Speedtest Tracker prune results older than"
  type        = string
  default     = "0"
}

variable "speedtest_tracker_public_dashboard" {
  description = "Speedtest Tracker public dashboard"
  type        = string
  default     = "true"
}

variable "speedtest_tracker_blocked_servers" {
  description = "Speedtest Tracker blocked servers"
  type        = string
  default     = ""
}

variable "speedtest_tracker_interface" {
  description = "Speedtest Tracker interface"
  type        = string
  default     = ""
}

variable "speedtest_tracker_schedule" {
  description = "Speedtest Tracker schedule"
  type        = string
  default     = "0 * * * *"
}

variable "speedtest_tracker_servers" {
  description = "Speedtest Tracker servers"
  type        = string
  default     = ""
}

variable "speedtest_tracker_skip_ips" {
  description = "Speedtest Tracker skip IPs"
  type        = string
  default     = ""
}

variable "speedtest_tracker_threshold_download" {
  description = "Speedtest Tracker threshold download"
  type        = string
  default     = "900"
}

variable "speedtest_tracker_threshold_enabled" {
  description = "Speedtest Tracker threshold enabled"
  type        = string
  default     = "true"
}

variable "speedtest_tracker_threshold_ping" {
  description = "Speedtest Tracker threshold ping"
  type        = string
  default     = "25"
}

variable "speedtest_tracker_threshold_upload" {
  description = "Speedtest Tracker threshold upload"
  type        = string
  default     = "900"
}

# TinyAuth Configuration
variable "tinyauth_secret" {
  description = "TinyAuth secret"
  type        = string
  default     = ""
}

variable "tinyauth_app_url" {
  description = "TinyAuth app URL"
  type        = string
  default     = "https://auth.example.com"
}

variable "tinyauth_users" {
  description = "TinyAuth users"
  type        = string
  default     = ""
}

variable "tinyauth_google_client_id" {
  description = "TinyAuth Google client ID"
  type        = string
  default     = ""
}

variable "tinyauth_google_client_secret" {
  description = "TinyAuth Google client secret"
  type        = string
  default     = ""
}

variable "tinyauth_github_client_id" {
  description = "TinyAuth GitHub client ID"
  type        = string
  default     = ""
}

variable "tinyauth_github_client_secret" {
  description = "TinyAuth GitHub client secret"
  type        = string
  default     = ""
}

variable "tinyauth_session_expiry" {
  description = "TinyAuth session expiry"
  type        = string
  default     = "604800"
}

variable "tinyauth_cookie_secure" {
  description = "TinyAuth cookie secure"
  type        = string
  default     = "true"
}

variable "tinyauth_app_title" {
  description = "TinyAuth app title"
  type        = string
  default     = "Bolabaden"
}

variable "tinyauth_login_max_retries" {
  description = "TinyAuth login max retries"
  type        = string
  default     = "15"
}

variable "tinyauth_login_timeout" {
  description = "TinyAuth login timeout"
  type        = string
  default     = "300"
}

variable "tinyauth_oauth_auto_redirect" {
  description = "TinyAuth OAuth auto redirect"
  type        = string
  default     = "none"
}

variable "tinyauth_oauth_whitelist" {
  description = "TinyAuth OAuth whitelist"
  type        = string
  default     = "boden.crouch@gmail.com,halomastar@gmail.com,athenajaguiar@gmail.com,bolabaden.duckdns@gmail.com"
}

# Code Server Configuration
variable "codeserver_password" {
  description = "Code Server password"
  type        = string
  default     = ""
}

variable "codeserver_sudo_password" {
  description = "Code Server sudo password"
  type        = string
  default     = ""
}

variable "codeserver_default_workspace" {
  description = "Code Server default workspace"
  type        = string
  default     = "/workspace"
}

# FlareSolverr Configuration
variable "flaresolverr_port" {
  description = "FlareSolverr port"
  type        = string
  default     = "8191"
}

variable "flaresolverr_log_level" {
  description = "FlareSolverr log level"
  type        = string
  default     = "info"
}

variable "flaresolverr_log_html" {
  description = "FlareSolverr log HTML"
  type        = string
  default     = "false"
}

variable "flaresolverr_captcha_solver" {
  description = "FlareSolverr captcha solver"
  type        = string
  default     = "none"
}

variable "flaresolverr_host" {
  description = "FlareSolverr host"
  type        = string
  default     = "0.0.0.0"
}

variable "flaresolverr_headless" {
  description = "FlareSolverr headless"
  type        = string
  default     = "true"
}

variable "flaresolverr_browser_timeout" {
  description = "FlareSolverr browser timeout"
  type        = string
  default     = "120000"
}

variable "flaresolverr_test_url" {
  description = "FlareSolverr test URL"
  type        = string
  default     = "https://www.google.com"
}

variable "flaresolverr_prometheus_enabled" {
  description = "FlareSolverr Prometheus enabled"
  type        = string
  default     = "false"
}

variable "prometheus_port" {
  description = "Prometheus port"
  type        = string
  default     = "9090"
}

# VPN Configuration
variable "warp_tun_device" {
  description = "WARP TUN device"
  type        = string
  default     = "/dev/net/tun"
}

variable "ts_authkey" {
  description = "Tailscale auth key"
  type        = string
  default     = ""
}

variable "ts_state_dir" {
  description = "Tailscale state directory"
  type        = string
  default     = "/var/lib/tailscale"
}

variable "ts_routes" {
  description = "Tailscale routes"
  type        = string
  default     = "10.76.0.0/16,172.17.0.0/16,100.64.0.0/10"
}

variable "aiostreams_port" {
  description = "AIOStreams port"
  type        = string
  default     = "3005"
}

variable "comet_port" {
  description = "Comet port"
  type        = string
  default     = "2020"
}

# AI Service Configuration
variable "lobechat_access_code" {
  description = "LobeChat access code"
  type        = string
  default     = "brunner56"
}

# AI API Keys
variable "anthropic_api_key" {
  description = "Anthropic API key"
  type        = string
  default     = ""
}

variable "brave_api_key" {
  description = "Brave API key"
  type        = string
  default     = ""
}

variable "deepseek_api_key" {
  description = "DeepSeek API key"
  type        = string
  default     = ""
}

variable "exa_api_key" {
  description = "Exa API key"
  type        = string
  default     = ""
}

variable "firecrawl_api_key" {
  description = "Firecrawl API key"
  type        = string
  default     = ""
}

variable "fire_crawl_api_key" {
  description = "Fire Crawl API key"
  type        = string
  default     = ""
}

variable "gemini_api_key" {
  description = "Gemini API key"
  type        = string
  default     = ""
}

variable "glama_api_key" {
  description = "Glama API key"
  type        = string
  default     = ""
}

variable "groq_api_key" {
  description = "Groq API key"
  type        = string
  default     = ""
}

variable "hf_token" {
  description = "Hugging Face token"
  type        = string
  default     = ""
}

variable "huggingface_access_token" {
  description = "Hugging Face access token"
  type        = string
  default     = ""
}

variable "huggingface_api_token" {
  description = "Hugging Face API token"
  type        = string
  default     = ""
}

variable "langchain_api_key" {
  description = "LangChain API key"
  type        = string
  default     = ""
}

variable "mistral_api_key" {
  description = "Mistral API key"
  type        = string
  default     = ""
}

variable "mistralai_api_key" {
  description = "MistralAI API key"
  type        = string
  default     = ""
}

variable "openai_api_key" {
  description = "OpenAI API key"
  type        = string
  default     = ""
}

variable "openrouter_api_key" {
  description = "OpenRouter API key"
  type        = string
  default     = ""
}

variable "perplexity_api_key" {
  description = "Perplexity API key"
  type        = string
  default     = ""
}

variable "perplexityai_api_key" {
  description = "PerplexityAI API key"
  type        = string
  default     = ""
}

variable "replicate_api_key" {
  description = "Replicate API key"
  type        = string
  default     = ""
}

variable "revid_api_key" {
  description = "Revid API key"
  type        = string
  default     = ""
}

variable "sambanova_api_key" {
  description = "SambaNova API key"
  type        = string
  default     = ""
}

variable "search1api_key" {
  description = "Search1API key"
  type        = string
  default     = ""
}

variable "serpapi_api_key" {
  description = "SerpAPI key"
  type        = string
  default     = ""
}

variable "tavily_api_key" {
  description = "Tavily API key"
  type        = string
  default     = ""
}

variable "togetherai_api_key" {
  description = "TogetherAI API key"
  type        = string
  default     = ""
}

variable "unify_api_key" {
  description = "Unify API key"
  type        = string
  default     = ""
}

variable "upstage_api_key" {
  description = "Upstage API key"
  type        = string
  default     = ""
}

variable "upstageai_api_key" {
  description = "UpstageAI API key"
  type        = string
  default     = ""
}

variable "you_api_key" {
  description = "You API key"
  type        = string
  default     = ""
}

variable "next_public_ga_measurement_id" {
  description = "Next.js Google Analytics measurement ID"
  type        = string
  default     = ""
} 