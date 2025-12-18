package reporter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/autoenumeration/autoenum/internal/engine"
)

type Reporter struct {
	outputDir string
}

func NewReporter(outputDir string) *Reporter {
	return &Reporter{
		outputDir: outputDir,
	}
}

func (r *Reporter) GenerateJSON(results *engine.ScanResults) error {
	outputFile := filepath.Join(r.outputDir, fmt.Sprintf("scan_results_%s.json", time.Now().Format("20060102_150405")))

	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal results: %w", err)
	}

	if err := os.WriteFile(outputFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write JSON report: %w", err)
	}

	return nil
}

func (r *Reporter) GenerateHTML(results *engine.ScanResults) error {
	outputFile := filepath.Join(r.outputDir, fmt.Sprintf("scan_report_%s.html", time.Now().Format("20060102_150405")))

	html := r.generateHTMLContent(results)

	if err := os.WriteFile(outputFile, []byte(html), 0644); err != nil {
		return fmt.Errorf("failed to write HTML report: %w", err)
	}

	return nil
}

func (r *Reporter) generateHTMLContent(results *engine.ScanResults) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>AutoEnumeration Scan Report - %s</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            padding: 20px;
            color: #333;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: white;
            border-radius: 12px;
            box-shadow: 0 20px 60px rgba(0,0,0,0.3);
            overflow: hidden;
        }
        .header {
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            color: white;
            padding: 40px;
            text-align: center;
        }
        .header h1 {
            font-size: 2.5em;
            margin-bottom: 10px;
            text-shadow: 2px 2px 4px rgba(0,0,0,0.2);
        }
        .header .target {
            font-size: 1.5em;
            opacity: 0.9;
        }
        .summary {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 20px;
            padding: 40px;
            background: #f8f9fa;
        }
        .stat-card {
            background: white;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 2px 8px rgba(0,0,0,0.1);
            text-align: center;
            transition: transform 0.2s;
        }
        .stat-card:hover {
            transform: translateY(-5px);
            box-shadow: 0 4px 12px rgba(0,0,0,0.15);
        }
        .stat-number {
            font-size: 2.5em;
            font-weight: bold;
            color: #667eea;
            margin-bottom: 5px;
        }
        .stat-label {
            color: #666;
            font-size: 0.9em;
            text-transform: uppercase;
            letter-spacing: 1px;
        }
        .section {
            padding: 40px;
        }
        .section h2 {
            color: #667eea;
            margin-bottom: 20px;
            padding-bottom: 10px;
            border-bottom: 3px solid #667eea;
        }
        .tool-result {
            background: #f8f9fa;
            padding: 20px;
            margin-bottom: 15px;
            border-radius: 8px;
            border-left: 4px solid #667eea;
        }
        .tool-result.success {
            border-left-color: #28a745;
        }
        .tool-result.failed {
            border-left-color: #dc3545;
        }
        .tool-name {
            font-size: 1.2em;
            font-weight: bold;
            margin-bottom: 10px;
        }
        .tool-meta {
            display: flex;
            gap: 20px;
            margin-top: 10px;
            font-size: 0.9em;
            color: #666;
        }
        .badge {
            display: inline-block;
            padding: 4px 12px;
            border-radius: 12px;
            font-size: 0.85em;
            font-weight: 600;
        }
        .badge.success {
            background: #d4edda;
            color: #155724;
        }
        .badge.failed {
            background: #f8d7da;
            color: #721c24;
        }
        .list-items {
            background: white;
            padding: 15px;
            border-radius: 6px;
            margin-top: 10px;
            max-height: 300px;
            overflow-y: auto;
        }
        .list-items ul {
            list-style: none;
        }
        .list-items li {
            padding: 8px;
            border-bottom: 1px solid #eee;
            font-family: 'Courier New', monospace;
            font-size: 0.9em;
        }
        .list-items li:last-child {
            border-bottom: none;
        }
        .footer {
            background: #2c3e50;
            color: white;
            padding: 20px;
            text-align: center;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔍 AutoEnumeration Report</h1>
            <div class="target">Target: %s</div>
            <div style="margin-top: 10px; opacity: 0.8;">%s</div>
        </div>

        <div class="summary">
            <div class="stat-card">
                <div class="stat-number">%d</div>
                <div class="stat-label">Total Tools</div>
            </div>
            <div class="stat-card">
                <div class="stat-number">%d</div>
                <div class="stat-label">Successful</div>
            </div>
            <div class="stat-card">
                <div class="stat-number">%d</div>
                <div class="stat-label">Failed</div>
            </div>
            <div class="stat-card">
                <div class="stat-number">%d</div>
                <div class="stat-label">Subdomains</div>
            </div>
            <div class="stat-card">
                <div class="stat-number">%d</div>
                <div class="stat-label">Live Hosts</div>
            </div>
            <div class="stat-card">
                <div class="stat-number">%d</div>
                <div class="stat-label">Vulnerabilities</div>
            </div>
        </div>

        <div class="section">
            <h2>📊 Tool Execution Results</h2>
            %s
        </div>

        <div class="section">
            <h2>🌐 Discovered Subdomains</h2>
            <div class="list-items">
                <ul>%s</ul>
            </div>
        </div>

        <div class="section">
            <h2>✅ Live Hosts</h2>
            <div class="list-items">
                <ul>%s</ul>
            </div>
        </div>

        <div class="footer">
            <p>Generated by AutoEnumeration v1.0.0</p>
            <p>Duration: %s | Mode: %s</p>
        </div>
    </div>
</body>
</html>`,
		results.Target,
		results.Target,
		results.StartTime.Format("2006-01-02 15:04:05"),
		results.Summary.TotalTools,
		results.Summary.SuccessfulRuns,
		results.Summary.FailedRuns,
		len(results.Summary.Subdomains),
		len(results.Summary.LiveHosts),
		results.Summary.Vulnerabilities,
		r.generateToolResultsHTML(results),
		r.generateListHTML(results.Summary.Subdomains),
		r.generateListHTML(results.Summary.LiveHosts),
		results.Duration,
		results.Mode,
	)
}

func (r *Reporter) generateToolResultsHTML(results *engine.ScanResults) string {
	html := ""
	for _, result := range results.ToolResults {
		statusClass := result.Status
		badgeClass := result.Status

		html += fmt.Sprintf(`
            <div class="tool-result %s">
                <div class="tool-name">%s <span class="badge %s">%s</span></div>
                <div class="tool-meta">
                    <span>⏱️ Duration: %s</span>
                    <span>🕐 Started: %s</span>
                </div>
            </div>`,
			statusClass,
			result.ToolName,
			badgeClass,
			result.Status,
			result.Duration,
			result.StartTime.Format("15:04:05"),
		)
	}
	return html
}

func (r *Reporter) generateListHTML(items []string) string {
	if len(items) == 0 {
		return "<li>No items found</li>"
	}

	html := ""
	for _, item := range items {
		html += fmt.Sprintf("<li>%s</li>", item)
	}
	return html
}
