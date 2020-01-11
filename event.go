package melian

// https://docs.sentry.io/development/sdk-dev/event-payloads/
type event struct {
	Dist        string                 `json:"dist,omitempty"`
	Environment string                 `json:"environment,omitempty"`
	Release     string                 `json:"release,omitempty"`
	ServerName  string                 `json:"server_name,omitempty"`
	EventID     string                 `json:"event_id,omitempty"`
	Platform    string                 `json:"platform,omitempty"`
	Logger      string                 `json:"logger,omitempty"`
	Level       string                 `json:"level,omitempty"`
	Message     string                 `json:"message,omitempty"`
	Timestamp   int64                  `json:"timestamp,omitempty"`
	Exception   []exception            `json:"exception,omitempty"`
	Extra       map[string]interface{} `json:"extra,omitempty"`
	Tags        map[string]string      `json:"tags,omitempty"`
	Transaction string                 `json:"transaction,omitempty"`
	Request     request                `json:"request,omitempty"`
}

// https://docs.sentry.io/development/sdk-dev/event-payloads/exception/
type exception struct {
	Type          string      `json:"type,omitempty"`
	Value         string      `json:"value,omitempty"`
	Module        string      `json:"module,omitempty"`
	Stacktrace    *stacktrace `json:"stacktrace,omitempty"`
	RawStacktrace *stacktrace `json:"raw_stacktrace,omitempty"`
}
