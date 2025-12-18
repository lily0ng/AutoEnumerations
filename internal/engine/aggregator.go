package engine

import (
	"sync"
)

type ResultAggregator struct {
	mu              sync.RWMutex
	subdomains      map[string]bool
	liveHosts       map[string]bool
	openPorts       map[string]bool
	vulnerabilities []interface{}
	rawResults      map[string][]interface{}
}

func NewResultAggregator() *ResultAggregator {
	return &ResultAggregator{
		subdomains: make(map[string]bool),
		liveHosts:  make(map[string]bool),
		openPorts:  make(map[string]bool),
		rawResults: make(map[string][]interface{}),
	}
}

func (ra *ResultAggregator) Add(category string, data interface{}) {
	ra.mu.Lock()
	defer ra.mu.Unlock()

	ra.rawResults[category] = append(ra.rawResults[category], data)

	switch category {
	case "subdomain":
		if subdomains, ok := data.([]string); ok {
			for _, sub := range subdomains {
				ra.subdomains[sub] = true
			}
		}
	case "port_scanning":
		if ports, ok := data.([]string); ok {
			for _, port := range ports {
				ra.openPorts[port] = true
			}
		}
	case "http_probing":
		if hosts, ok := data.([]string); ok {
			for _, host := range hosts {
				ra.liveHosts[host] = true
			}
		}
	case "vulnerability":
		ra.vulnerabilities = append(ra.vulnerabilities, data)
	}
}

func (ra *ResultAggregator) GetSubdomains() []string {
	ra.mu.RLock()
	defer ra.mu.RUnlock()

	result := make([]string, 0, len(ra.subdomains))
	for sub := range ra.subdomains {
		result = append(result, sub)
	}
	return result
}

func (ra *ResultAggregator) GetLiveHosts() []string {
	ra.mu.RLock()
	defer ra.mu.RUnlock()

	result := make([]string, 0, len(ra.liveHosts))
	for host := range ra.liveHosts {
		result = append(result, host)
	}
	return result
}

func (ra *ResultAggregator) GetOpenPorts() []string {
	ra.mu.RLock()
	defer ra.mu.RUnlock()

	result := make([]string, 0, len(ra.openPorts))
	for port := range ra.openPorts {
		result = append(result, port)
	}
	return result
}

func (ra *ResultAggregator) GetVulnerabilityCount() int {
	ra.mu.RLock()
	defer ra.mu.RUnlock()

	return len(ra.vulnerabilities)
}

func (ra *ResultAggregator) GetRawResults(category string) []interface{} {
	ra.mu.RLock()
	defer ra.mu.RUnlock()

	return ra.rawResults[category]
}
