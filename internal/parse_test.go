package internal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFSParser(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		parser := NewFSParser(FSParserOpts{
			Root: "./testdata/parser/happy",
		})

		got, err := parser.ParseSequences()
		require.NoError(t, err)

		require.Equal(
			t,
			map[string]Sequence{
				"foo_seq.yaml": {Calls: []Call{{Url: "https://foo.bar.com/foo_seq_top"}}},
				"subdir/foo_seq.yml": {Calls: []Call{{Url: "https://foo.bar.com/foo_seq_inner"}}},
			},
			got,
		)
	})
}
