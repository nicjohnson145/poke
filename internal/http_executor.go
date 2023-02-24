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
		h.log.Debug().Bytes("bodyBytes", bodyBytes).Msg("adding message body")
		if err != nil {
			return nil, fmt.Errorf("error marshalling body as JSON: %w", err)
		}
		inBody = strings.NewReader(string(bodyBytes))
	}

	method := call.Method
	if method == "" {
		method = http.MethodGet
	}
	h.log.Debug().Str("method", method).Str("url", call.Url).Msg("executing call")
	req, err := http.NewRequest(method, call.Url, inBody)
	if err != nil {
		return nil, fmt.Errorf("error building request: %w", err)
	}

	if call.Headers != nil {
		for k, v := range call.Headers {
			req.Header.Set(k, v)
		}
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	var outBody map[string]any
	err = json.NewDecoder(resp.Body).Decode(&outBody)
	if err != nil {
		if err != io.EOF {
			return nil, fmt.Errorf("error decoding body as JSON: %w", err)
		}
		// No body, move along
	}

	return &ExecuteResult{
		Body: outBody,
		StatusCode: resp.StatusCode,
	}, nil
}
