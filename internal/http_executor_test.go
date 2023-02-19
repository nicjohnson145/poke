package internal

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

type roundTripFunc func(req *http.Request) *http.Response

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func newTestClient(fn roundTripFunc) *http.Client {
	return &http.Client{
		Transport: roundTripFunc(fn),
	}
}

func TestSimpleExecute(t *testing.T) {
	t.Run("with body", func(t *testing.T) {
		client := newTestClient(func(req *http.Request) *http.Response {
			// Check the host & path
			require.Equal(t, "http://some.host.com/some-endpoint", req.URL.String())
			// Assert we're posted a body
			var body map[string]string
			err := json.NewDecoder(req.Body).Decode(&body)
			require.NoError(t, err)
			require.Equal(t, map[string]string{"foo": "bar"}, body)

			// Write back
			return &http.Response{
				StatusCode: http.StatusOK,
				Body: ioutil.NopCloser(bytes.NewBufferString(`{"baz": "qux"}`)),
				Header: make(http.Header),
			}
		})

		ex := NewHTTPExecutor(HTTPExecutorOpts{
			Client: client,
		})

		got, err := ex.Execute(Call{
			Url: "http://some.host.com/some-endpoint",
			Body: map[string]any{"foo": "bar"},
		})
		require.NoError(t, err)
		require.Equal(
			t,
			&ExecuteResult{
				Body: map[string]any{"baz": "qux"},
				StatusCode: http.StatusOK,
			},
			got,
		)
	})
}
