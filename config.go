package webmail

import (
	"net/url"
	"strings"
)

const (
	port = ":443"
	path = "/webmail/api/jsonrpc"
)

type parameters struct {
	JsonRpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	ID      int         `json:"id"`
	Token   *string     `json:"token,omitempty"`
	Params  interface{} `json:"params,omitempty"`
}

type Config struct {
	url string
	id  int
}

// NewConfig returns a pointer to structure with the configuration for connecting to the API server
// Parameters:
//      server - address without schema and port
func NewConfig(server string) *Config {
	if !strings.Contains(server, ":") {
		server += port
	}
	u := url.URL{
		Scheme: "https",
		Host:   server,
		Path:   path,
	}
	return &Config{
		url: u.String(),
	}
}

// NewApplication returns a pointer to structure with application data
func NewApplication(name, vendor, version string) *ApiApplication {
	if name == "" {
		name = "TempApp"
	}
	if vendor == "" {
		vendor = "TempVendor"
	}
	if version == "" {
		version = "v1.0.1"
	}
	return &ApiApplication{
		Name:    name,
		Vendor:  vendor,
		Version: version,
	}
}

func (c *Config) getID() int {
	c.id++
	return c.id
}
