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

		mockEx := NewMockExecutor(t)
		mockEx.EXPECT().Execute(call1).Return(nil, nil)
		mockEx.EXPECT().Execute(call2).Return(nil, nil)
		mockEx.EXPECT().Execute(call3).Return(nil, nil)
		mockEx.EXPECT().Execute(call4).Return(nil, nil)

		runner := NewRunner(RunnerOpts{
			HttpExecutor: mockEx,
			Parser: mockParser,
		})

		err := runner.Run("./some/path")
		require.NoError(t, err)
	})
}
