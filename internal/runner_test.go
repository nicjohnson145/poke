package internal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	call1 := Call{
		Name: "auth",
		Type: RequestTypeHttp,
		Url: "http://some.api.com/login",
	}
	call2 := Call{
		Name: "fetch",
		Type: RequestTypeHttp,
		Url: "http://some.api.com/objects/list",
	}
	seqA := Sequence{
		Calls: []Call{
			call1,
			call2,
		},
	}

	call3 := Call{
		Name: "auth",
		Type: RequestTypeHttp,
		Url: "http://other.api.com/login",
	}
	call4 := Call{
		Name: "fetch",
		Type: RequestTypeHttp,
		Url: "http://other.api.com/objects/list",
	}
	seqB := Sequence{
		Calls: []Call{
			call3,
			call4,
		},
	}


	t.Run("happy", func(t *testing.T) {
		mockParser := NewMockParser(t)
		mockParser.EXPECT().Parse("./some/path").Return(
			SequenceMap{
				"seqA.yaml": seqA,
				"seqB.yaml": seqB,
			},
			nil,
		)

		ok := &ExecuteResult{StatusCode: 200}
		mockEx := NewMockExecutor(t)
		mockEx.EXPECT().Execute(call1).Return(ok, nil)
		mockEx.EXPECT().Execute(call2).Return(ok, nil)
		mockEx.EXPECT().Execute(call3).Return(ok, nil)
		mockEx.EXPECT().Execute(call4).Return(ok, nil)

		runner := NewRunner(RunnerOpts{
			HttpExecutor: mockEx,
			Parser: mockParser,
		})

		err := runner.Run("./some/path")
		require.NoError(t, err)
	})
}

func TestVariablePassing(t *testing.T) {
	t.Run("happy jq body", func(t *testing.T) {
		call1 := Call{
			Name: "fetch",
			Url: "http://some.api.com/get",
			Export: []Export{
				{
					JQ: ".data.name",
					As: "fooVar",
				},
			},
		}

		call2 := Call{
			Name: "use",
			Url: "http://some.api.com/use",
			Body: map[string]any{
				"name": "{{ .fooVar }}",
			},
		}
		transformCall2 := Call{
			Name: "use",
			Url: "http://some.api.com/use",
			Body: map[string]any{
				"name": "fooNameActual",
			},
		}
		seqA := Sequence{Calls: []Call{call1, call2}}

		mockParser := NewMockParser(t)
		mockParser.EXPECT().Parse("./some/path").Return(SequenceMap{"seqA.yaml": seqA}, nil)

		mockEx := NewMockExecutor(t)
		mockEx.EXPECT().Execute(call1).Return(
			&ExecuteResult{
				StatusCode: 200,
				Body: map[string]any{
					"data": map[string]any{
						"name": "fooNameActual",
					},
				},
			},
			nil,
		)
		mockEx.EXPECT().Execute(transformCall2).Return(&ExecuteResult{StatusCode: 200}, nil)

		runner := NewRunner(RunnerOpts{
			HttpExecutor: mockEx,
			Parser: mockParser,
		})

		err := runner.Run("./some/path")
		require.NoError(t, err)
	})
}
