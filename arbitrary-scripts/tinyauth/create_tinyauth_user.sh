#!/bin/bash

# Color definitions
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}TinyAuth User Creation Script${NC}"
echo -e "${YELLOW}This script helps you create properly formatted user credentials for TinyAuth${NC}"
echo ""

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
  echo -e "${RED}Error: Docker is not running or you don't have permission to use it.${NC}"
  exit 1
fi

# Ask for username
read -p "Enter username: " username
if [ -z "$username" ]; then
  echo -e "${RED}Error: Username cannot be empty${NC}"
  exit 1
fi

# Ask for password
read -s -p "Enter password: " password
echo ""
if [ -z "$password" ]; then
  echo -e "${RED}Error: Password cannot be empty${NC}"
  exit 1
fi

# Confirm password
read -s -p "Confirm password: " password_confirm
echo ""
if [ "$password" != "$password_confirm" ]; then
  echo -e "${RED}Error: Passwords do not match${NC}"
  exit 1
fi

echo -e "${YELLOW}Generating bcrypt hash using TinyAuth container...${NC}"
# Run the tinyauth user creation command
user_string=$(docker run --rm ghcr.io/steveiliop56/tinyauth:v3 user create --username "$username" --password "$password" --docker)

if [ $? -ne 0 ]; then
  echo -e "${RED}Error: Failed to create user${NC}"
  exit 1
fi

echo -e "${GREEN}Success! Your user string is:${NC}"
echo "$user_string"
echo ""
echo -e "${YELLOW}Add this to your .env file as TINYAUTH_USERS=${NC}"
echo -e "${YELLOW}For multiple users, separate with commas (no spaces)${NC}"
echo ""
echo -e "${BLUE}Example .env entry:${NC}"
echo "TINYAUTH_USERS=$user_string"

# Generate TOTP if required
read -p "Do you want to add TOTP (2FA) for this user? (y/n): " add_totp
if [[ "$add_totp" == "y" || "$add_totp" == "Y" ]]; then
  echo -e "${YELLOW}Generating TOTP secret...${NC}"
  totp_user_string=$(docker run --rm -i ghcr.io/steveiliop56/tinyauth:v3 totp generate --user "$user_string")
  
  if [ $? -ne 0 ]; then
    echo -e "${RED}Error: Failed to generate TOTP${NC}"
    exit 1
  fi
  
  echo -e "${GREEN}Success! Your user string with TOTP is:${NC}"
  echo "$totp_user_string"
  echo ""
  echo -e "${YELLOW}Update your .env file with this new value for TINYAUTH_USERS=${NC}"
  echo "TINYAUTH_USERS=$totp_user_string"
  echo ""
  echo -e "${YELLOW}Scan the QR code above with your authenticator app${NC}"
fi

echo ""
echo -e "${GREEN}Done!${NC}" 