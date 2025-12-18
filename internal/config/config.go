package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	OutputDir   string                `mapstructure:"output_dir"`
	Verbose     bool                  `mapstructure:"verbose"`
	Timeout     int                   `mapstructure:"timeout"`
	Concurrency ConcurrencyConfig     `mapstructure:"concurrency"`
	Tools       map[string]ToolConfig `mapstructure:"tools"`
	Modes       map[string]ModeConfig `mapstructure:"modes"`
	RateLimit   RateLimitConfig       `mapstructure:"rate_limit"`
}

type ConcurrencyConfig struct {
	MaxWorkers int `mapstructure:"max_workers"`
	QueueSize  int `mapstructure:"queue_size"`
}

type ToolConfig struct {
	Enabled  bool              `mapstructure:"enabled"`
	Priority int               `mapstructure:"priority"`
	Timeout  int               `mapstructure:"timeout"`
	Args     map[string]string `mapstructure:"args"`
}

type ModeConfig struct {
	Description string   `mapstructure:"description"`
	Tools       []string `mapstructure:"tools"`
}

type RateLimitConfig struct {
	RequestsPerSecond int `mapstructure:"requests_per_second"`
	Burst             int `mapstructure:"burst"`
}

func LoadConfig(cfgFile string) (*Config, error) {
	v := viper.New()

	v.SetDefault("output_dir", "./output")
	v.SetDefault("verbose", false)
	v.SetDefault("timeout", 3600)
	v.SetDefault("concurrency.max_workers", 10)
	v.SetDefault("concurrency.queue_size", 100)
	v.SetDefault("rate_limit.requests_per_second", 100)
	v.SetDefault("rate_limit.burst", 200)

	setDefaultModes(v)
	setDefaultTools(v)

	if cfgFile != "" {
		v.SetConfigFile(cfgFile)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("$HOME/.autoenum")
		v.AddConfigPath("/etc/autoenum")
	}

	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	return &config, nil
}

func setDefaultModes(v *viper.Viper) {
	v.SetDefault("modes.quick", map[string]interface{}{
		"description": "Fast subdomain and port discovery",
		"tools":       []string{"subfinder", "dnsx", "naabu", "httpx"},
	})

	v.SetDefault("modes.standard", map[string]interface{}{
		"description": "Full enumeration with web discovery",
		"tools": []string{
			"subfinder", "assetfinder", "dnsx", "naabu", "httpx",
			"katana", "waybackurls", "nuclei", "gowitness",
		},
	})

	v.SetDefault("modes.deep", map[string]interface{}{
		"description": "Comprehensive scan with vulnerability detection",
		"tools": []string{
			"subfinder", "assetfinder", "amass", "dnsx", "naabu",
			"httpx", "katana", "gospider", "waybackurls", "gau",
			"nuclei", "nikto", "gowitness", "tlsx", "webanalyze",
		},
	})
}

func setDefaultTools(v *viper.Viper) {
	defaultTools := map[string]map[string]interface{}{
		"subfinder":   {"enabled": true, "priority": 1, "timeout": 300},
		"assetfinder": {"enabled": true, "priority": 1, "timeout": 300},
		"amass":       {"enabled": true, "priority": 1, "timeout": 600},
		"dnsx":        {"enabled": true, "priority": 2, "timeout": 300},
		"naabu":       {"enabled": true, "priority": 3, "timeout": 600},
		"httpx":       {"enabled": true, "priority": 4, "timeout": 300},
		"katana":      {"enabled": true, "priority": 5, "timeout": 600},
		"nuclei":      {"enabled": true, "priority": 6, "timeout": 900},
		"gowitness":   {"enabled": true, "priority": 7, "timeout": 300},
	}

	for tool, settings := range defaultTools {
		v.SetDefault(fmt.Sprintf("tools.%s", tool), settings)
	}
}

func GenerateDefaultConfig(outputPath string) error {
	defaultConfig := `# AutoEnumeration Configuration File

output_dir: ./output
verbose: false
timeout: 3600

concurrency:
  max_workers: 10
  queue_size: 100

rate_limit:
  requests_per_second: 100
  burst: 200

modes:
  quick:
    description: "Fast subdomain and port discovery"
    tools:
      - subfinder
      - dnsx
      - naabu
      - httpx

  standard:
    description: "Full enumeration with web discovery"
    tools:
      - subfinder
      - assetfinder
      - dnsx
      - naabu
      - httpx
      - katana
      - waybackurls
      - nuclei
      - gowitness

  deep:
    description: "Comprehensive scan with vulnerability detection"
    tools:
      - subfinder
      - assetfinder
      - amass
      - dnsx
      - naabu
      - httpx
      - katana
      - gospider
      - waybackurls
      - gau
      - nuclei
      - nikto
      - gowitness
      - tlsx
      - webanalyze

tools:
  subfinder:
    enabled: true
    priority: 1
    timeout: 300
    args:
      threads: "10"

  assetfinder:
    enabled: true
    priority: 1
    timeout: 300

  amass:
    enabled: true
    priority: 1
    timeout: 600
    args:
      passive: "true"

  dnsx:
    enabled: true
    priority: 2
    timeout: 300

  naabu:
    enabled: true
    priority: 3
    timeout: 600
    args:
      top_ports: "1000"

  httpx:
    enabled: true
    priority: 4
    timeout: 300
    args:
      threads: "50"

  katana:
    enabled: true
    priority: 5
    timeout: 600

  nuclei:
    enabled: true
    priority: 6
    timeout: 900
    args:
      severity: "critical,high,medium"

  gowitness:
    enabled: true
    priority: 7
    timeout: 300
`

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}

	return os.WriteFile(outputPath, []byte(defaultConfig), 0644)
}
