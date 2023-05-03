package dotenv

import (
	"context"
	_ "embed"
	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

//go:embed testdata/kv_expected.js
var kvExpected string

//go:embed testdata/empty_expected.js
var emptyExpected string

func TestNewReader(t *testing.T) {
	ctx := logr.NewContext(context.TODO(), testr.NewWithOptions(t, testr.Options{Verbosity: 10}))
	err := NewReader(ctx, "./testdata/kv.env", filepath.Join(t.TempDir(), "env-config.js"))
	assert.NoError(t, err)
}

func TestReader_GetLines(t *testing.T) {
	ctx := logr.NewContext(context.TODO(), testr.NewWithOptions(t, testr.Options{Verbosity: 10}))
	var cases = []struct {
		path  string
		lines []string
		err   bool
	}{
		{
			"./testdata/kv.env",
			[]string{"NODE_ENV=production", "FOO=bar"},
			false,
		},
		{
			"./testdata/this-file-does-not-exist",
			nil,
			true,
		},
	}
	r := new(Reader)
	for _, tt := range cases {
		t.Run(tt.path, func(t *testing.T) {
			res, err := r.GetLines(ctx, tt.path)
			assert.EqualValues(t, tt.err, err != nil)
			assert.ElementsMatch(t, tt.lines, res)
		})
	}
}

func TestReader_Parse(t *testing.T) {
	ctx := logr.NewContext(context.TODO(), testr.NewWithOptions(t, testr.Options{Verbosity: 10}))
	var cases = []struct {
		path     string
		expected string
		err      bool
	}{
		{
			"./testdata/kv.env",
			kvExpected,
			false,
		},
		{
			"./testdata/empty.env",
			emptyExpected,
			false,
		},
		{
			"./testdata/this-file-does-not-exist",
			"",
			true,
		},
	}
	r := new(Reader)
	for _, tt := range cases {
		t.Run(tt.path, func(t *testing.T) {
			res, err := r.GetLines(ctx, tt.path)
			assert.EqualValues(t, tt.err, err != nil)
			if err != nil {
				return
			}

			data := r.Parse(ctx, res)
			assert.EqualValues(t, tt.expected, data)
		})
	}
}

func TestReader_Write(t *testing.T) {
	ctx := logr.NewContext(context.TODO(), testr.NewWithOptions(t, testr.Options{Verbosity: 10}))
	r := new(Reader)
	err := r.Write(ctx, kvExpected, filepath.Join(t.TempDir(), "test.js"))
	assert.NoError(t, err)
}
