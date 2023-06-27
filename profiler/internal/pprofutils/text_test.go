

package pprofutils

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTextConvert(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		textIn := strings.TrimSpace(`
main;foo 5
main;foobar 4
main;foo;bar 3
`)
		proto, err := Text{}.Convert(strings.NewReader(textIn))
		require.NoError(t, err)
		textOut := bytes.Buffer{}
		require.NoError(t, Protobuf{}.Convert(proto, &textOut))
		require.Equal(t, textIn+"\n", textOut.String())
	})

	t.Run("headerWithOneSampleType", func(t *testing.T) {
		textIn := strings.TrimSpace(`
samples/count
main;foo 5
main;foobar 4
main;foo;bar 3
	`)
		proto, err := Text{}.Convert(strings.NewReader(textIn))
		require.NoError(t, err)
		textOut := bytes.Buffer{}
		require.NoError(t, Protobuf{SampleTypes: true}.Convert(proto, &textOut))
		require.Equal(t, textIn+"\n", textOut.String())
	})

	t.Run("headerWithMultipleSampleTypes", func(t *testing.T) {
		textIn := strings.TrimSpace(`
samples/count duration/nanoseconds
main;foo 5 50000000
main;foobar 4 40000000
main;foo;bar 3 30000000
	`)
		proto, err := Text{}.Convert(strings.NewReader(textIn))
		require.NoError(t, err)
		textOut := bytes.Buffer{}
		require.NoError(t, Protobuf{SampleTypes: true}.Convert(proto, &textOut))
		require.Equal(t, textIn+"\n", textOut.String())
	})
}
