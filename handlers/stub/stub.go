package stub

type Stubs map[string]Stub

type Stub struct {
	Request  Request  `json:"request" toml:"request" yaml:"request"`
	Response Response `json:"response" toml:"response" yaml:"response"`
}

type Request struct {
	Methods []string          `json:"methods" toml:"methods" yaml:"methods"`
	Path    string            `json:"path" toml:"path" yaml:"path"`
	Headers map[string]string `json:"headers" toml:"headers" yaml:"headers"`
	Query   map[string]string `json:"query" toml:"query" yaml:"query"`
	Body    string            `json:"body" toml:"body" yaml:"body"`
}

type Response struct {
	Code    int               `json:"code" toml:"code" yaml:"code"`
	Headers map[string]string `json:"headers" toml:"headers" yaml:"headers"`
	Body    string            `json:"body" toml:"body" yaml:"body"`
}
