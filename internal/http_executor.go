package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/rs/zerolog"
)

type HTTPExecutorOpts struct {
	Logger zerolog.Logger
	Client *http.Client
}

func NewHTTPExecutor(opts HTTPExecutorOpts) *HTTPExecutor {
	client := http.DefaultClient
	if opts.Client != nil {
		client = opts.Client
	}

	return &HTTPExecutor{
		log: opts.Logger,
		client: client,
	}
}

var _ Executor = (*HTTPExecutor)(nil)

type HTTPExecutor struct {
	log    zerolog.Logger
	client *http.Client
}

func (h *HTTPExecutor) Execute(call Call) (*ExecuteResult, error) {
	var inBody io.Reader
	if call.Body != nil {
		bodyBytes, err := json.Marshal(call.Body)
		if err != nil {
			return nil, fmt.Errorf("error marshalling body as JSON: %w", err)
		}
		inBody = strings.NewReader(string(bodyBytes))
	}

	req, err := http.NewRequest(call.Method, call.Url, inBody)
	if err != nil {
		return nil, fmt.Errorf("error building request: %w", err)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	var outBody map[string]any
	err = json.NewDecoder(resp.Body).Decode(&outBody)
	if err != nil {
		return nil, fmt.Errorf("error decoding body as JSON: %w", err)
	}

	return &ExecuteResult{
		Body: outBody,
		StatusCode: resp.StatusCode,
	}, nil
}
