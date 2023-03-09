package internal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFSParser(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		parser := NewFSParser(FSParserOpts{})

		got, err := parser.Parse("./testdata/parser/happy")
		require.NoError(t, err)

		require.Equal(
			t,
			SequenceMap{
				"foo_seq.yaml": {
					Calls: []Call{{Url: "https://foo.bar.com/foo_seq_top"}},
					path: "testdata/parser/happy",
				},
				"subdir/foo_seq.yml": {
					Calls: []Call{{Url: "https://foo.bar.com/foo_seq_inner"}},
					path: "testdata/parser/happy/subdir",
				},
			},
			got,
		)
	})

	t.Run("unknown key errors", func(t *testing.T) {
		parser := NewFSParser(FSParserOpts{})

		_, err := parser.Parse("./testdata/parser/unknown_key.yaml")
		require.Error(t, err)
		require.ErrorContains(t, err, "field assertions not found")
	})
}
