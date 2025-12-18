package engine

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/autoenumeration/autoenum/internal/config"
	"github.com/autoenumeration/autoenum/internal/logger"
	"github.com/autoenumeration/autoenum/internal/tools"
	"golang.org/x/sync/semaphore"
	"golang.org/x/time/rate"
)

type Engine struct {
	config      *config.Config
	registry    *tools.Registry
	rateLimiter *rate.Limiter
	semaphore   *semaphore.Weighted
	results     *ResultAggregator
}

type ScanConfig struct {
	Target    string
	Mode      string
	SkipTools []string
	OnlyTools []string
	Resume    bool
}

type ScanResults struct {
	Target      string                `json:"target"`
	StartTime   time.Time             `json:"start_time"`
	EndTime     time.Time             `json:"end_time"`
	Duration    string                `json:"duration"`
	Mode        string                `json:"mode"`
	ToolResults map[string]ToolResult `json:"tool_results"`
	Summary     Summary               `json:"summary"`
}

type ToolResult struct {
	ToolName  string      `json:"tool_name"`
	Status    string      `json:"status"`
	StartTime time.Time   `json:"start_time"`
	EndTime   time.Time   `json:"end_time"`
	Duration  string      `json:"duration"`
	Output    interface{} `json:"output"`
	Error     string      `json:"error,omitempty"`
}

type Summary struct {
	TotalTools      int      `json:"total_tools"`
	SuccessfulRuns  int      `json:"successful_runs"`
	FailedRuns      int      `json:"failed_runs"`
	SkippedRuns     int      `json:"skipped_runs"`
	Subdomains      []string `json:"subdomains"`
	LiveHosts       []string `json:"live_hosts"`
	OpenPorts       []string `json:"open_ports"`
	Vulnerabilities int      `json:"vulnerabilities"`
}

func NewEngine(cfg *config.Config) (*Engine, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	limiter := rate.NewLimiter(
		rate.Limit(cfg.RateLimit.RequestsPerSecond),
		cfg.RateLimit.Burst,
	)

	sem := semaphore.NewWeighted(int64(cfg.Concurrency.MaxWorkers))

	return &Engine{
		config:      cfg,
		registry:    tools.GetRegistry(),
		rateLimiter: limiter,
		semaphore:   sem,
		results:     NewResultAggregator(),
	}, nil
}

func (e *Engine) ExecuteScan(ctx context.Context, scanCfg *ScanConfig) (*ScanResults, error) {
	logger.Info("Initializing scan for target: %s", scanCfg.Target)

	results := &ScanResults{
		Target:      scanCfg.Target,
		StartTime:   time.Now(),
		Mode:        scanCfg.Mode,
		ToolResults: make(map[string]ToolResult),
	}

	toolsToRun, err := e.determineTools(scanCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to determine tools: %w", err)
	}

	logger.Info("Running %d tools in %s mode", len(toolsToRun), scanCfg.Mode)

	if err := e.executeTools(ctx, scanCfg.Target, toolsToRun, results); err != nil {
		return nil, fmt.Errorf("scan execution failed: %w", err)
	}

	results.EndTime = time.Now()
	results.Duration = results.EndTime.Sub(results.StartTime).String()
	results.Summary = e.generateSummary(results)

	return results, nil
}

func (e *Engine) determineTools(scanCfg *ScanConfig) ([]*tools.Tool, error) {
	var toolNames []string

	if len(scanCfg.OnlyTools) > 0 {
		toolNames = scanCfg.OnlyTools
	} else if modeConfig, exists := e.config.Modes[scanCfg.Mode]; exists {
		toolNames = modeConfig.Tools
	} else {
		return nil, fmt.Errorf("unknown scan mode: %s", scanCfg.Mode)
	}

	skipMap := make(map[string]bool)
	for _, skip := range scanCfg.SkipTools {
		skipMap[skip] = true
	}

	var selectedTools []*tools.Tool
	for _, name := range toolNames {
		if skipMap[name] {
			logger.Debug("Skipping tool: %s", name)
			continue
		}

		tool := e.registry.GetTool(name)
		if tool == nil {
			logger.Warning("Tool not found in registry: %s", name)
			continue
		}

		if toolCfg, exists := e.config.Tools[name]; exists && !toolCfg.Enabled {
			logger.Debug("Tool disabled in config: %s", name)
			continue
		}

		selectedTools = append(selectedTools, tool)
	}

	return selectedTools, nil
}

func (e *Engine) executeTools(ctx context.Context, target string, toolList []*tools.Tool, results *ScanResults) error {
	var wg sync.WaitGroup
	resultChan := make(chan ToolResult, len(toolList))
	errorChan := make(chan error, len(toolList))

	for _, tool := range toolList {
		wg.Add(1)

		go func(t *tools.Tool) {
			defer wg.Done()

			if err := e.semaphore.Acquire(ctx, 1); err != nil {
				errorChan <- fmt.Errorf("failed to acquire semaphore: %w", err)
				return
			}
			defer e.semaphore.Release(1)

			if err := e.rateLimiter.Wait(ctx); err != nil {
				errorChan <- fmt.Errorf("rate limiter error: %w", err)
				return
			}

			result := e.executeTool(ctx, t, target)
			resultChan <- result
		}(tool)
	}

	go func() {
		wg.Wait()
		close(resultChan)
		close(errorChan)
	}()

	for result := range resultChan {
		results.ToolResults[result.ToolName] = result

		if result.Status == "success" {
			logger.Success("Completed: %s (took %s)", result.ToolName, result.Duration)
		} else {
			logger.Error("Failed: %s - %s", result.ToolName, result.Error)
		}
	}

	for err := range errorChan {
		logger.Error("Execution error: %v", err)
	}

	return nil
}

func (e *Engine) executeTool(ctx context.Context, tool *tools.Tool, target string) ToolResult {
	result := ToolResult{
		ToolName:  tool.Name,
		StartTime: time.Now(),
		Status:    "running",
	}

	logger.Info("Executing: %s", tool.Name)

	timeout := time.Duration(e.config.Timeout) * time.Second
	if toolCfg, exists := e.config.Tools[tool.Name]; exists && toolCfg.Timeout > 0 {
		timeout = time.Duration(toolCfg.Timeout) * time.Second
	}

	toolCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	output, err := tool.Execute(toolCtx, target, e.config.OutputDir)

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime).String()

	if err != nil {
		result.Status = "failed"
		result.Error = err.Error()
		return result
	}

	result.Status = "success"
	result.Output = output

	e.results.Add(tool.Category, output)

	return result
}

func (e *Engine) generateSummary(results *ScanResults) Summary {
	summary := Summary{
		TotalTools: len(results.ToolResults),
	}

	for _, result := range results.ToolResults {
		switch result.Status {
		case "success":
			summary.SuccessfulRuns++
		case "failed":
			summary.FailedRuns++
		case "skipped":
			summary.SkippedRuns++
		}
	}

	summary.Subdomains = e.results.GetSubdomains()
	summary.LiveHosts = e.results.GetLiveHosts()
	summary.OpenPorts = e.results.GetOpenPorts()
	summary.Vulnerabilities = e.results.GetVulnerabilityCount()

	return summary
}
