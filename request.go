package melian

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

// https://docs.sentry.io/development/sdk-dev/event-payloads/request/
type request struct {
	URL         string            `json:"url,omitempty"`
	Method      string            `json:"method,omitempty"`
	Data        string            `json:"data,omitempty"`
	QueryString string            `json:"query_string,omitempty"`
	Cookies     string            `json:"cookies,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	Env         map[string]string `json:"env,omitempty"`
}

func fromHTTPRequest(req *http.Request) request {
	r := request{}
	// Method
	r.Method = req.Method

	// URL
	protocol := "http"
	if req.TLS != nil || req.Header.Get("X-Forwarded-Proto") == "https" {
		protocol = "https"
	}
	r.URL = fmt.Sprintf("%s://%s%s", protocol, req.Host, req.URL.Path)

	// Headers
	headers := make(map[string]string, len(req.Header))
	for k, v := range req.Header {
		headers[k] = strings.Join(v, ",")
	}
	headers["Host"] = req.Host
	r.Headers = headers

	// Cookies
	r.Cookies = req.Header.Get("Cookie")

	// Env
	if addr, port, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		r.Env = map[string]string{"REMOTE_ADDR": addr, "REMOTE_PORT": port}
	}

	// QueryString
	r.QueryString = req.URL.RawQuery

	// Body
	if req.Body != nil {
		bodyBytes, err := ioutil.ReadAll(req.Body)
		_ = req.Body.Close()
		if err == nil {
			// We have to restore original state of *req.Body
			req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
			r.Data = string(bodyBytes)
		}
	}
	return r
}
