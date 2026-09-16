package api

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/coditary/wuji-core/pkg/config"
)

func applyAuth(req *http.Request, spec *config.APICapabilitySpec, apiKey string, vars map[string]string) {
	if strings.TrimSpace(apiKey) == "" {
		return
	}
	auth := spec.Auth
	typ := "bearer"
	if auth != nil && strings.TrimSpace(auth.Type) != "" {
		typ = strings.ToLower(strings.TrimSpace(auth.Type))
	}
	switch typ {
	case "none":
		return
	case "query":
		param := "api_key"
		if auth != nil && strings.TrimSpace(auth.QueryParam) != "" {
			param = auth.QueryParam
		}
		q := req.URL.Query()
		q.Set(param, apiKey)
		req.URL.RawQuery = q.Encode()
	case "header":
		header := "Authorization"
		prefix := ""
		if auth != nil {
			if strings.TrimSpace(auth.Header) != "" {
				header = auth.Header
			}
			prefix = auth.Prefix
		}
		req.Header.Set(header, strings.TrimSpace(prefix+" "+apiKey))
	default: // bearer
		prefix := "Bearer"
		if auth != nil && auth.Prefix != "" {
			prefix = auth.Prefix
		}
		header := "Authorization"
		if auth != nil && strings.TrimSpace(auth.Header) != "" {
			header = auth.Header
		}
		req.Header.Set(header, prefix+" "+apiKey)
	}
	_ = vars
}

func appendQuery(base string, query map[string]string, vars map[string]string) string {
	if len(query) == 0 {
		return base
	}
	u, err := url.Parse(base)
	if err != nil {
		return base
	}
	q := u.Query()
	for k, v := range expandMap(query, vars) {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String()
}
