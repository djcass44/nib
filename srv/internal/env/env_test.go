package env_test

import (
	"github.com/djcass44/nib/srv/internal/env"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGet(t *testing.T) {
	assert.NoError(t, os.Setenv("FOO", "a:b"))
	assert.NoError(t, os.Setenv("BAR", "a/b"))

	var cases = []struct {
		key    string
		values []string
	}{
		{
			"FOO",
			[]string{"a", "b"},
		},
		{
			"BAR",
			[]string{"a/b"},
		},
	}
	for _, tt := range cases {
		t.Run(tt.key, func(t *testing.T) {
			assert.ElementsMatch(t, tt.values, env.Get(tt.key, os.Getenv))
		})
	}
}

func TestGetFirstAndLast(t *testing.T) {
	assert.NoError(t, os.Unsetenv("BAR"))
	assert.NoError(t, os.Setenv("FOO", "a:b:c"))

	assert.EqualValues(t, "a", env.GetFirst("FOO", os.Getenv))
	assert.EqualValues(t, "c", env.GetLast("FOO", os.Getenv))

	assert.Empty(t, env.GetFirst("BAR", os.Getenv))
	assert.Empty(t, env.GetLast("BAR", os.Getenv))
}
