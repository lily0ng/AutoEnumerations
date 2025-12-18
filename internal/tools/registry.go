package tools

import (
	"context"
	"fmt"
	"os/exec"
	"sync"
)

type Tool struct {
	Name        string
	Category    string
	Description string
	InstallCmd  string
	ExecuteFunc func(ctx context.Context, target, outputDir string) (interface{}, error)
}

type Registry struct {
	mu    sync.RWMutex
	tools map[string]*Tool
}

var (
	registry *Registry
	once     sync.Once
)

func GetRegistry() *Registry {
	once.Do(func() {
		registry = &Registry{
			tools: make(map[string]*Tool),
		}
		registry.registerAllTools()
	})
	return registry
}

func (r *Registry) Register(tool *Tool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tools[tool.Name] = tool
}

func (r *Registry) GetTool(name string) *Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.tools[name]
}

func (r *Registry) GetCategories() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	categoryMap := make(map[string]bool)
	for _, tool := range r.tools {
		categoryMap[tool.Category] = true
	}

	categories := make([]string, 0, len(categoryMap))
	for cat := range categoryMap {
		categories = append(categories, cat)
	}
	return categories
}

func (r *Registry) GetToolsByCategory(category string) []*Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var tools []*Tool
	for _, tool := range r.tools {
		if tool.Category == category {
			tools = append(tools, tool)
		}
	}
	return tools
}

func (r *Registry) GetAllTools() []*Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tools := make([]*Tool, 0, len(r.tools))
	for _, tool := range r.tools {
		tools = append(tools, tool)
	}
	return tools
}

func (t *Tool) IsInstalled() bool {
	_, err := exec.LookPath(t.Name)
	return err == nil
}

func (t *Tool) Execute(ctx context.Context, target, outputDir string) (interface{}, error) {
	if !t.IsInstalled() {
		return nil, fmt.Errorf("tool %s is not installed", t.Name)
	}

	if t.ExecuteFunc == nil {
		return nil, fmt.Errorf("execute function not implemented for %s", t.Name)
	}

	return t.ExecuteFunc(ctx, target, outputDir)
}

func (r *Registry) registerAllTools() {
	r.registerSubdomainTools()
	r.registerPortScanningTools()
	r.registerHTTPProbingTools()
	r.registerDNSTools()
	r.registerVulnerabilityTools()
	r.registerScreenshotTools()
	r.registerTechnologyTools()
	r.registerCrawlingTools()
	r.registerCloudTools()
	r.registerSSLTools()
}
