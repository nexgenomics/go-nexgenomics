package fabric

// In a nexgenomics standard agent, we expect to get identifying marks
// (including the tenant and agent id) from the Linux kernel boot parameters.
// We support a fallback mechanism through env strings to support testing.

import (
	"os"
	"strings"
)

// get_tenant
func get_tenant(cfg *ServeCfg) string {
	a := cfg.Tenant
	if a == "" {
		a = getCmdlineValue("tenant")
	}
	if a == "" {
		a = os.Getenv("TENANT_ID")
	}
	return a
}

// get_agentid
func get_agentid(cfg *ServeCfg) string {
	a := cfg.AgentId
	if a == "" {
		a = getCmdlineValue("agent")
	}
	if a == "" {
		a = os.Getenv("AGENT_ID")
	}
	return a
}

// get_natsurl
func get_natsurl(cfg *ServeCfg) string {
	a := cfg.NatsUrl
	if a == "" {
		a = getCmdlineValue("nats")
	}
	if a == "" {
		a = os.Getenv("NATSURL")
	}
	return a
}

// getCmdlineValue
func getCmdlineValue(key string) string {
	data, err := os.ReadFile("/proc/cmdline")
	if err != nil {
		return ""
	}
	parts := strings.Fields(string(data))
	for _, p := range parts {
		if strings.HasPrefix(p, key+"=") {
			return strings.TrimPrefix(p, key+"=")
		}
	}
	return ""
}
