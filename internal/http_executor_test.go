package internal

import (
	"net/http"
	"testing"

	"github.com/nicjohnson145/poke/internal/httptest"

	"github.com/stretchr/testify/require"
)

func TestSimpleExecute(t *testing.T) {
	t.Run("with body", func(t *testing.T) {
		client := httptest.New(
			t,
			httptest.Conf{
				URL: "http://some.host.com/some-endpoint",
				Body: map[string]any{"foo": "bar"},
				ResponseBody: map[string]any{"baz": "qux"},
			},
		)

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
		client := httptest.New(
			t,
			httptest.Conf{
				URL: "http://some.host.com/some-endpoint",
			},
		)

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
