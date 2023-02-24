package httptest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

type Conf struct {
	URL          string
	Method       string
	Body         map[string]any
	ResponseCode int
	ResponseBody map[string]any
}

type testClientKey struct {
	URL    string
	Method string
}

func New(t *testing.T, handlers ...Conf) *http.Client {
	tc := &TestClient{
		testingT:        t,
		handlers: make(map[testClientKey]Conf),
	}

	for _, handler := range handlers {
		h := handler
		require.NotEmpty(t, h.URL, "handler given with empty URL: %v", h)

		method := h.Method
		if method == "" {
			method = http.MethodGet
		}

		key := testClientKey{URL: h.URL, Method: method}
		if _, ok := tc.handlers[key]; ok {
			t.Fatalf("handler already registered for [%v]%v", method, h.URL)
		}

		tc.handlers[key] = h
	}

	return &http.Client{
		Transport: tc,
	}
}

type TestClient struct {
	testingT *testing.T
	handlers map[testClientKey]Conf
}

func (t *TestClient) RoundTrip(req *http.Request) (*http.Response, error) {
	conf, ok := t.handlers[testClientKey{URL: req.URL.String(), Method: req.Method}]
	if !ok {
		return nil, fmt.Errorf("no handler registered for [%v]%v", req.Method, req.URL.String())
	}


	var body map[string]any
	if req.Body != nil {
		defer req.Body.Close()
		err := json.NewDecoder(req.Body).Decode(&body)
		if err != nil {
			if err != io.EOF {
				t.testingT.Fatalf("got non-nil error decoding body: %v", err)
			}
			// Otherwise, there's just no body, move along
		}
	}

	// Make sure they gave us the right body
	require.Equal(t.testingT, conf.Body, body, "unexpected request body")
	// And if they did, then response with our status & body

	var outBody io.ReadCloser
	if conf.ResponseBody != nil {
		bodyBytes, err := json.Marshal(conf.ResponseBody)
		require.NoError(t.testingT, err, "marshalling response body")
		outBody = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	code := conf.ResponseCode
	if code == 0 {
		code = http.StatusOK
	}

	return &http.Response{
		StatusCode: code,
		Body: outBody,
		Header: make(http.Header),
	}, nil
}
