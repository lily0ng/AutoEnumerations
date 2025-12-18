# AutoEnumeration - Usage Examples

## Table of Contents
- [Basic Usage](#basic-usage)
- [Scan Modes](#scan-modes)
- [Advanced Workflows](#advanced-workflows)
- [Custom Configurations](#custom-configurations)
- [Integration Examples](#integration-examples)

## Basic Usage

### Simple Domain Scan
```bash
autoenum scan -t example.com
```

### Scan with Verbose Output
```bash
autoenum scan -t example.com -v
```

### Custom Output Directory
```bash
autoenum scan -t example.com -o /tmp/recon/example
```

## Scan Modes

### Quick Mode - Fast Discovery
Perfect for initial reconnaissance:
```bash
autoenum scan -t example.com -m quick
```

**Tools used:**
- subfinder (subdomain discovery)
- dnsx (DNS resolution)
- naabu (port scanning)
- httpx (HTTP probing)

**Use case:** Quick asset discovery, initial scope mapping

### Standard Mode - Comprehensive Enumeration
Balanced approach with web discovery:
```bash
autoenum scan -t example.com -m standard
```

**Tools used:**
- subfinder, assetfinder (subdomain discovery)
- dnsx (DNS resolution)
- naabu (port scanning)
- httpx (HTTP probing)
- katana, waybackurls (web crawling)
- nuclei (vulnerability scanning)
- gowitness (screenshots)

**Use case:** Regular bug bounty reconnaissance, security assessments

### Deep Mode - Full Security Assessment
Comprehensive scan with all tools:
```bash
autoenum scan -t example.com -m deep
```

**Tools used:** All available tools including:
- Multiple subdomain finders
- Advanced crawling
- Technology detection
- SSL/TLS analysis
- Extensive vulnerability scanning

**Use case:** In-depth security audits, penetration testing

## Advanced Workflows

### Bug Bounty Workflow

#### Phase 1: Initial Discovery
```bash
# Quick scan to map the attack surface
autoenum scan -t target.com -m quick -o ./target/phase1

# Review subdomains
cat ./target/phase1/subfinder.txt
```

#### Phase 2: Deep Enumeration
```bash
# Deep scan on discovered assets
autoenum scan -t target.com -m deep -o ./target/phase2 --threads 20
```

#### Phase 3: Targeted Scanning
```bash
# Focus on specific tools
autoenum scan -t target.com \
  --only nuclei,katana,httpx \
  -o ./target/phase3
```

### Penetration Testing Workflow

```bash
# 1. Reconnaissance
autoenum scan -t client.com -m standard -o ./pentest/recon

# 2. Skip noisy tools for stealth
autoenum scan -t client.com \
  --skip masscan,rustscan,nikto \
  -o ./pentest/stealth

# 3. Vulnerability assessment
autoenum scan -t client.com \
  --only nuclei,tlsx,webanalyze \
  -o ./pentest/vulns \
  --threads 5
```

### Continuous Monitoring

```bash
#!/bin/bash
# monitor.sh - Daily reconnaissance script

TARGET="example.com"
DATE=$(date +%Y%m%d)
OUTPUT_DIR="./monitoring/$TARGET/$DATE"

# Run standard scan
autoenum scan -t $TARGET -m standard -o $OUTPUT_DIR

# Compare with previous day
if [ -d "./monitoring/$TARGET/$(date -d yesterday +%Y%m%d)" ]; then
    diff \
        "./monitoring/$TARGET/$(date -d yesterday +%Y%m%d)/subfinder.txt" \
        "$OUTPUT_DIR/subfinder.txt" \
        > "$OUTPUT_DIR/new_subdomains.txt"
    
    if [ -s "$OUTPUT_DIR/new_subdomains.txt" ]; then
        echo "New subdomains discovered!"
        cat "$OUTPUT_DIR/new_subdomains.txt"
    fi
fi
```

## Custom Configurations

### High-Speed Scanning
```yaml
# config-fast.yaml
concurrency:
  max_workers: 50
  queue_size: 500

rate_limit:
  requests_per_second: 500
  burst: 1000

tools:
  httpx:
    args:
      threads: "100"
  naabu:
    args:
      rate: "10000"
```

```bash
autoenum scan -t example.com --config config-fast.yaml
```

### Stealth Configuration
```yaml
# config-stealth.yaml
concurrency:
  max_workers: 2
  queue_size: 10

rate_limit:
  requests_per_second: 5
  burst: 10

tools:
  naabu:
    enabled: false
  masscan:
    enabled: false
  nuclei:
    args:
      rate_limit: "10"
```

```bash
autoenum scan -t example.com --config config-stealth.yaml
```

### Focused Web Application Testing
```yaml
# config-webapp.yaml
modes:
  webapp:
    description: "Web application focused testing"
    tools:
      - httpx
      - katana
      - gospider
      - waybackurls
      - nuclei
      - ffuf
      - gowitness
      - webanalyze
```

```bash
autoenum scan -t app.example.com -m webapp --config config-webapp.yaml
```

## Integration Examples

### Integration with Burp Suite
```bash
# 1. Discover endpoints
autoenum scan -t example.com -m standard -o ./burp_targets

# 2. Extract URLs for Burp
cat ./burp_targets/katana.txt \
    ./burp_targets/waybackurls.txt \
    ./burp_targets/gau.txt \
    | sort -u > burp_scope.txt

# 3. Import burp_scope.txt into Burp Suite
```

### Integration with Metasploit
```bash
# 1. Scan and identify services
autoenum scan -t example.com -m deep -o ./msf_recon

# 2. Parse nmap results
cat ./msf_recon/nmap.xml | grep -oP 'port="\K[0-9]+' > open_ports.txt

# 3. Use in Metasploit
# msfconsole
# > use auxiliary/scanner/portscan/tcp
# > set RHOSTS example.com
# > set PORTS file:open_ports.txt
```

### Integration with Custom Scripts
```bash
#!/bin/bash
# custom_pipeline.sh

TARGET=$1
OUTPUT="./pipeline/$TARGET"

# 1. Run AutoEnumeration
autoenum scan -t $TARGET -m standard -o $OUTPUT

# 2. Process subdomains
cat $OUTPUT/subfinder.txt | while read subdomain; do
    # Custom processing
    echo "Processing: $subdomain"
    
    # Example: Check for subdomain takeover
    subjack -w <(echo $subdomain) -t 100 -timeout 30 -o $OUTPUT/takeovers.txt
done

# 3. Analyze nuclei findings
jq '.info.severity' $OUTPUT/nuclei.json | sort | uniq -c

# 4. Generate summary
echo "=== Scan Summary ===" > $OUTPUT/summary.txt
echo "Subdomains: $(wc -l < $OUTPUT/subfinder.txt)" >> $OUTPUT/summary.txt
echo "Live Hosts: $(wc -l < $OUTPUT/httpx.txt)" >> $OUTPUT/summary.txt
echo "Vulnerabilities: $(jq -s length $OUTPUT/nuclei.json)" >> $OUTPUT/summary.txt
```

### CI/CD Integration
```yaml
# .github/workflows/security-scan.yml
name: Security Scan

on:
  schedule:
    - cron: '0 0 * * *'  # Daily
  workflow_dispatch:

jobs:
  scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      
      - name: Setup Go
        uses: actions/setup-go@v2
        with:
          go-version: '1.21'
      
      - name: Install AutoEnumeration
        run: |
          git clone https://github.com/yourusername/AutoEnumerations.git
          cd AutoEnumerations
          ./scripts/install.sh
      
      - name: Run Scan
        run: |
          autoenum scan -t ${{ secrets.TARGET_DOMAIN }} \
            -m standard \
            -o ./scan_results
      
      - name: Upload Results
        uses: actions/upload-artifact@v2
        with:
          name: scan-results
          path: ./scan_results/
      
      - name: Notify on New Findings
        run: |
          # Custom notification logic
          if [ -s ./scan_results/nuclei.json ]; then
            echo "Vulnerabilities found!"
          fi
```

## Tips and Best Practices

### 1. Resource Management
```bash
# Limit resource usage
autoenum scan -t example.com --threads 5 --timeout 1800
```

### 2. Incremental Scanning
```bash
# Day 1: Discovery
autoenum scan -t example.com -m quick -o ./day1

# Day 2: Deep dive on findings
autoenum scan -t example.com --only nuclei,katana -o ./day2
```

### 3. Handling Large Scopes
```bash
# Process subdomains in batches
cat subdomains.txt | split -l 100 - batch_
for batch in batch_*; do
    autoenum scan -t $(cat $batch) -m quick -o ./output_$(basename $batch)
done
```

### 4. Result Filtering
```bash
# Extract high-severity findings
jq '.[] | select(.info.severity == "high" or .info.severity == "critical")' \
    output/nuclei.json > critical_findings.json
```

### 5. Automation
```bash
# Create a wrapper script
#!/bin/bash
autoenum scan -t $1 -m ${2:-standard} -o ./scans/$(date +%Y%m%d)_$1
```

## Troubleshooting

### Scan Timeout
```bash
# Increase timeout for large targets
autoenum scan -t example.com --timeout 7200
```

### Rate Limiting Issues
```bash
# Reduce concurrency
autoenum scan -t example.com --threads 3
```

### Memory Issues
```bash
# Process in smaller batches
autoenum scan -t example.com --only subfinder,httpx
```
