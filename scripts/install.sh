#!/bin/bash

set -e

BLUE='\033[0;34m'
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${BLUE}╔═══════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║                                                       ║${NC}"
echo -e "${BLUE}║          AutoEnumeration Installation Script         ║${NC}"
echo -e "${BLUE}║                                                       ║${NC}"
echo -e "${BLUE}╚═══════════════════════════════════════════════════════╝${NC}"
echo ""

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

check_command() {
    if command -v "$1" &> /dev/null; then
        return 0
    else
        return 1
    fi
}

install_go_tool() {
    local tool_name=$1
    local install_cmd=$2
    
    if check_command "$tool_name"; then
        log_success "$tool_name is already installed"
        return 0
    fi
    
    log_info "Installing $tool_name..."
    if eval "$install_cmd" &> /dev/null; then
        log_success "$tool_name installed successfully"
        return 0
    else
        log_error "Failed to install $tool_name"
        return 1
    fi
}

log_info "Checking prerequisites..."

if ! check_command "go"; then
    log_error "Go is not installed. Please install Go 1.21 or later."
    exit 1
fi
log_success "Go is installed"

log_info "Installing ProjectDiscovery tools..."

install_go_tool "subfinder" "go install -v github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest"
install_go_tool "httpx" "go install -v github.com/projectdiscovery/httpx/cmd/httpx@latest"
install_go_tool "naabu" "go install -v github.com/projectdiscovery/naabu/v2/cmd/naabu@latest"
install_go_tool "nuclei" "go install -v github.com/projectdiscovery/nuclei/v3/cmd/nuclei@latest"
install_go_tool "katana" "go install github.com/projectdiscovery/katana/cmd/katana@latest"
install_go_tool "dnsx" "go install -v github.com/projectdiscovery/dnsx/cmd/dnsx@latest"
install_go_tool "tlsx" "go install github.com/projectdiscovery/tlsx/cmd/tlsx@latest"
install_go_tool "cloudlist" "go install github.com/projectdiscovery/cloudlist/cmd/cloudlist@latest"

log_info "Installing community tools..."

install_go_tool "assetfinder" "go install github.com/tomnomnom/assetfinder@latest"
install_go_tool "httprobe" "go install github.com/tomnomnom/httprobe@latest"
install_go_tool "waybackurls" "go install github.com/tomnomnom/waybackurls@latest"
install_go_tool "amass" "go install -v github.com/owasp-amass/amass/v4/...@master"
install_go_tool "gospider" "go install github.com/jaeles-project/gospider@latest"
install_go_tool "hakrawler" "go install github.com/hakluke/hakrawler@latest"
install_go_tool "gau" "go install github.com/lc/gau/v2/cmd/gau@latest"
install_go_tool "gowitness" "go install github.com/sensepost/gowitness@latest"
install_go_tool "webanalyze" "go install github.com/rverton/webanalyze/cmd/webanalyze@latest"
install_go_tool "ffuf" "go install github.com/ffuf/ffuf/v2@latest"
install_go_tool "gobuster" "go install github.com/OJ/gobuster/v3@latest"

log_info "Building AutoEnumeration..."
cd "$(dirname "$0")/.."
go build -o autoenum main.go

if [ $? -eq 0 ]; then
    log_success "AutoEnumeration built successfully"
    
    log_info "Installing to /usr/local/bin..."
    if sudo mv autoenum /usr/local/bin/; then
        log_success "AutoEnumeration installed to /usr/local/bin/autoenum"
    else
        log_warning "Could not install to /usr/local/bin. Binary is in current directory."
    fi
else
    log_error "Build failed"
    exit 1
fi

log_info "Generating default configuration..."
mkdir -p ~/.autoenum
if [ ! -f ~/.autoenum/config.yaml ]; then
    cat > ~/.autoenum/config.yaml << 'EOF'
output_dir: ./output
verbose: false
timeout: 3600

concurrency:
  max_workers: 10
  queue_size: 100

rate_limit:
  requests_per_second: 100
  burst: 200
EOF
    log_success "Configuration file created at ~/.autoenum/config.yaml"
else
    log_info "Configuration file already exists"
fi

echo ""
log_success "Installation complete!"
echo ""
echo -e "${GREEN}Run 'autoenum --help' to get started${NC}"
echo ""
