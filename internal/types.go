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
	SkipVerify  bool              `yaml:"skip-verify,omitempty"`
	FromImport  *ImportedCall     `yaml:"from-import,omitempty"`
}

func (c *Call) GetType() RequestType {
	if c.Type == "" {
		return RequestTypeHttp
	}
	return c.Type
}

type Sequence struct {
	Vars          map[string]any             `yaml:"vars"`
	Imports       map[string]string          `yaml:"imports"`
	Calls         []Call                     `yaml:"calls"`
	path          string                     `yaml:"-"`
	importedCalls map[string]map[string]Call `yaml:"-"`
}

type Export struct {
	JQ string `yaml:"jq,omitempty"`
	As string `yaml:"as,omitempty"`
}

type Assert struct {
	JQ       string `yaml:"jq,omitempty"`
	Expected any    `yaml:"expected,omitempty"`
}

type ImportedCall struct {
	Name string `yaml:"name"`
	Call string `yaml:"call"`
}
