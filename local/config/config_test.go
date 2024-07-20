package config

import (
	"strings"
	"testing"
)

func TestParseConfig(t *testing.T) {
	yamlData := `
targets:
  - name: home-assistant
    host_addr: http://192.168.1.20/api/webhook/someurl
  - name: hubitat
    host_addr: http://192.168.1.21/data
server:
  port: 8123
  path: /api/webhook/someurl
  verbose: true
  tls_insecure_skip_verify: true
`

	reader := strings.NewReader(yamlData)
	cfg, err := parseConfig(reader)
	if err != nil {
		t.Fatalf("Failed to parse config: %v", err)
	}

	if len(cfg.Targets) != 2 {
		t.Fatalf("Expected 2 targets, got %d", len(cfg.Targets))
	}

	if cfg.Server.Port != 8123 {
		t.Errorf("Expected server port 8123, got %d", cfg.Server.Port)
	}

	if cfg.Server.Verbose != true {
		t.Errorf("Expected verbose true, got %t", cfg.Server.Verbose)
	}
	if cfg.Server.TlsInsecureSkipVerify != true {
		t.Errorf("Expected tls_insecure_skip_verify true, got %t", cfg.Server.TlsInsecureSkipVerify)
	}

	if cfg.Targets[0].Name != "home-assistant" {
		t.Errorf("Expected first target name 'home-assistant', got '%s'", cfg.Targets[0].Name)
	}

	if cfg.Targets[1].HostAddr != "http://192.168.1.21/data" {
		t.Errorf("Expected second target host_addr 'http://192.168.1.21', got '%s'", cfg.Targets[1].HostAddr)
	}
}
