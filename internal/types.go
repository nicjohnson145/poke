package internal

//go:generate go-enum --file $GOFILE --marshal --names

/*
ENUM(
http
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
}

type Sequence struct {
	Calls []Call `yaml:"calls"`
}
