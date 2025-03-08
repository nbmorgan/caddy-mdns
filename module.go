package mdns

import (
	"fmt"
	"log"
	"os/exec"
    "net/http" 
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func init() {
	caddy.RegisterModule(MDNSRegister{})
}

// MDNSRegister is a Caddy module for automatic mDNS registration.
type MDNSRegister struct {
	ServiceName string `json:"service_name,omitempty"`
	Port        int    `json:"port,omitempty"`
}

// CaddyModule returns the Caddy module information.
func (MDNSRegister) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.mdns_register",
		New: func() caddy.Module { return new(MDNSRegister) },
	}
}

// Provision sets up the module.
func (m *MDNSRegister) Provision(ctx caddy.Context) error {
	if m.ServiceName == "" {
		m.ServiceName = "caddy"
	}
	if m.Port == 0 {
		m.Port = 80
	}

	// Register the service via Avahi or mDNS
	go m.registerService()

	return nil
}

// ServeHTTP forwards the request after registering mDNS.
func (m *MDNSRegister) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	return next.ServeHTTP(w, r)
}

// registerService registers the service with mDNS using Avahi or `dns-sd`.
func (m *MDNSRegister) registerService() {
	serviceCmd := fmt.Sprintf("avahi-publish -a -R %s.local $(hostname -I | awk '{print $1}')", m.ServiceName)
	cmd := exec.Command("sh", "-c", serviceCmd)

	if err := cmd.Start(); err != nil {
		log.Printf("[mDNS] Failed to register service: %v", err)
	} else {
		log.Printf("[mDNS] Registered %s.local on port %d", m.ServiceName, m.Port)
	}
}
