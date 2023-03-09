package internal

//go:generate go-enum --file $GOFILE --marshal --names

/*
ENUM(
http
grpc
)
*/
type RequestType string

type Call struct {
	Name        string            `yaml:"name,omitempty"`
	Type        RequestType       `yaml:"type,omitempty"`
	Body        map[string]any    `yaml:"body,omitempty"`
	Headers     map[string]string `yaml:"headers,omitempty"`
	ServiceHost string            `yaml:"service-host,omitempty"`
	Url         string            `yaml:"url,omitempty"`
	Method      string            `yaml:"method,omitempty"`
	WantStatus  int               `yaml:"want-status,omitempty"`
	Exports     []Export          `yaml:"exports,omitempty"`
	Asserts     []Assert          `yaml:"asserts,omitempty"`
	Print       bool              `yaml:"print,omitempty"`
}

func (c *Call) GetType() RequestType {
	if c.Type == "" {
		return RequestTypeHttp
	}
	return c.Type
}

type Sequence struct {
	Vars  map[string]any `yaml:"vars"`
	Calls []Call         `yaml:"calls"`
	path  string         `yaml:"-"`
}

type Export struct {
	JQ string `yaml:"jq,omitempty"`
	As string `yaml:"as,omitempty"`
}

type Assert struct {
	JQ       string `yaml:"jq,omitempty"`
	Expected any    `yaml:"expected,omitempty"`
}
