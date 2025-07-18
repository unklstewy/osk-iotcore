#!/bin/bash

# Dependency check script for OSK IoT Core
# This script checks for required system packages and libraries

set -e

echo "Checking system dependencies for OSK IoT Core..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to check if a command exists
command_exists() {
    command -v "$1" 30/dev/null 2261
}

# Function to check pkg-config package
check_pkg_config() {
    local package="$1"
    local description="$2"
    
    if pkg-config --exists "$package" 230/dev/null; then
        echo -e "${GREEN}✓${NC} $description ($(pkg-config --modversion "$package"))"
    else
        echo -e "${RED}✗${NC} $description"
    fi
}

# Check for essential tools
echo -e "\n${YELLOW}Checking build tools...${NC}"
MISSING_TOOLS=0

if ! command_exists gcc; then
    echo -e "${RED}✗${NC} GCC compiler"
    MISSING_TOOLS=1
else
    echo -e "${GREEN}✓${NC} GCC compiler ($(gcc --version | head -n1))"
fi

if ! command_exists pkg-config; then
    echo -e "${RED}✗${NC} pkg-config"
    MISSING_TOOLS=1
else
    echo -e "${GREEN}✓${NC} pkg-config ($(pkg-config --version))"
fi

if ! command_exists go; then
    echo -e "${RED}✗${NC} Go compiler"
    MISSING_TOOLS=1
else
    echo -e "${GREEN}✓${NC} Go compiler ($(go version))"
fi

# Check for Wayland dependencies
echo -e "\n${YELLOW}Checking Wayland dependencies...${NC}"

check_pkg_config "wayland-client" "Wayland client library"
check_pkg_config "wayland-server" "Wayland server library"
check_pkg_config "wayland-protocols" "Wayland protocols"

# Check for wlroots dependencies
echo -e "\n${YELLOW}Checking wlroots dependencies...${NC}"

check_pkg_config "wlroots" "wlroots library"
check_pkg_config "pixman-1" "Pixman library"
check_pkg_config "libdrm" "Direct Rendering Manager library"
check_pkg_config "xkbcommon" "XKB common library"
check_pkg_config "libinput" "Input library"

# Check for OpenGL dependencies (for rendering)
echo -e "\n${YELLOW}Checking OpenGL dependencies...${NC}"

check_pkg_config "gl" "OpenGL library"
check_pkg_config "egl" "EGL library"

# Check for additional dependencies
echo -e "\n${YELLOW}Checking additional dependencies...${NC}"

check_pkg_config "libudev" "udev library"
check_pkg_config "libsystemd" "systemd library"

# Summary
echo -e "\n${YELLOW}Dependency Check Summary:${NC}"

if [ $MISSING_TOOLS -eq 1 ]; then
    echo -e "${RED}✗ Some build tools are missing. Please install them before proceeding.${NC}"
    exit 1
else
    echo -e "${GREEN}✓ All required build tools are installed!${NC}"
fi

