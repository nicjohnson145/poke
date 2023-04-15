package internal

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSimpleExecute(t *testing.T) {
	jsonResponse := func(t *testing.T, data map[string]any, status int) *http.Response {
		t.Helper()
		outBody := ioutil.NopCloser(bytes.NewBuffer([]byte{}))
		if data != nil {
			bodyBytes, err := json.Marshal(data)
			require.NoError(t, err)
			outBody = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		return &http.Response{
			StatusCode: status,
			Body: outBody,
			Header: make(http.Header),
		}
	}

	expectBody := func(t *testing.T, data map[string]any) func(*http.Request) {
		return func(req *http.Request) {
			t.Helper()
			if data == nil {
				return
			}

			var body map[string]any
			if req.Body != nil {
				err := json.NewDecoder(req.Body).Decode(&body)
				if err != nil {
					if err != io.EOF {
						t.Fatalf("got non-nil error decoding body: %v", err)
					}
					// Otherwise, there's just no body, move along
				}
			}
			require.Equal(t, data, body, "unexpected request body")
		}
	}

	t.Run("with body", func(t *testing.T) {
		client := NewMockIHttpClient(t)
		client.EXPECT().
			Do(mock.Anything).
			RunAndReturn(func(r *http.Request) (*http.Response, error) {
				expectBody(t, map[string]any{"foo": "bar"})(r)

				return jsonResponse(
					t,
					map[string]any{
						"baz": "qux",
					},
					http.StatusOK,
				), nil
			})

		ex := NewHTTPExecutor(HTTPExecutorOpts{
			Client: client,
		})

		got, err := ex.Execute(Call{
			Url:  "http://some.host.com/some-endpoint",
			Body: map[string]any{"foo": "bar"},
		})
		require.NoError(t, err)
		require.Equal(
			t,
			&ExecuteResult{
				Body:       map[string]any{"baz": "qux"},
				StatusCode: http.StatusOK,
			},
			got,
		)
	})

	t.Run("no body", func(t *testing.T) {
		client := NewMockIHttpClient(t)
		client.EXPECT().
			Do(mock.Anything).
			RunAndReturn(func(r *http.Request) (*http.Response, error) {
				expectBody(t, nil)(r)
				return jsonResponse(t, nil, http.StatusOK), nil
			})
		ex := NewHTTPExecutor(HTTPExecutorOpts{
			Client: client,
		})

		got, err := ex.Execute(Call{
			Url:  "http://some.host.com/some-endpoint",
		})
		require.NoError(t, err)
		require.Equal(
			t,
			&ExecuteResult{
				StatusCode: http.StatusOK,
			},
			got,
		)
	})
}
