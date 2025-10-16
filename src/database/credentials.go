package database

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// SaveCredentialsToFile saves admin credentials to a file
func SaveCredentialsToFile(creds *AdminCredentials, configDir, port string) error {
	serverURL := getAccessibleURL(port)

	content := fmt.Sprintf(`Quotes API - ADMIN CREDENTIALS
========================================
WEB UI LOGIN:
  URL:      %s/admin
  Username: %s

API ACCESS:
  URL:      %s/api/v1/admin
  Header:   Authorization: Bearer %s

CREDENTIALS:
  Username: %s
  Token:    %s

Created: %s
========================================
`, serverURL, creds.Username, serverURL, creds.Token,
		creds.Username, creds.Token, time.Now().Format("2006-01-02 15:04:05"))

	credFile := filepath.Join(configDir, "admin-credentials.txt")
	if err := os.WriteFile(credFile, []byte(content), 0600); err != nil {
		return fmt.Errorf("failed to write credentials file: %w", err)
	}

	return nil
}

// getAccessibleURL returns the most relevant URL for accessing the server
// Priority: FQDN > hostname > public IP > fallback
// NEVER shows localhost, 127.0.0.1, 0.0.0.0, or ::1
// Supports both IPv4 and IPv6
func getAccessibleURL(port string) string {
	// Try to get hostname
	hostname, err := os.Hostname()
	if err == nil && hostname != "" && hostname != "localhost" {
		// Try to resolve hostname to see if it's a valid FQDN
		if addrs, err := net.LookupHost(hostname); err == nil && len(addrs) > 0 {
			return fmt.Sprintf("http://%s:%s", hostname, port)
		}
	}

	// Try to get outbound IP (most likely accessible IP)
	if ip := getOutboundIP(); ip != "" {
		return formatURLWithIP(ip, port)
	}

	// Fallback to hostname if we have one
	if hostname != "" && hostname != "localhost" {
		return fmt.Sprintf("http://%s:%s", hostname, port)
	}

	// Last resort: use a generic message
	return fmt.Sprintf("http://<your-host>:%s", port)
}

// getOutboundIP gets the preferred outbound IP of this machine
// Tries IPv4 first, then IPv6
func getOutboundIP() string {
	// Try IPv4 first
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err == nil {
		defer conn.Close()
		localAddr := conn.LocalAddr().(*net.UDPAddr)
		return localAddr.IP.String()
	}

	// Try IPv6
	conn, err = net.Dial("udp", "[2001:4860:4860::8888]:80")
	if err == nil {
		defer conn.Close()
		localAddr := conn.LocalAddr().(*net.UDPAddr)
		return localAddr.IP.String()
	}

	return ""
}

// formatURLWithIP formats a URL with proper IPv6 bracket handling
func formatURLWithIP(ip, port string) string {
	// IPv6 addresses contain colons and need brackets
	if strings.Contains(ip, ":") {
		return fmt.Sprintf("http://[%s]:%s", ip, port)
	}
	return fmt.Sprintf("http://%s:%s", ip, port)
}
