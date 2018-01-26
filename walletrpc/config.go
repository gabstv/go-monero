package walletrpc

import (
	"net/http"
)

type Config struct {
	Address       string
	CustomHeaders map[string]string
	Transport     http.RoundTripper
}
