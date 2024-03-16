package pathfinder_test

import (
	"github.com/djcass44/nib/cli/internal/pathfinder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestFindBuildDir(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(dir, "node_modules"), 0755))
	require.NoError(t, os.Mkdir(filepath.Join(dir, "public"), 0755))
	require.NoError(t, os.Mkdir(filepath.Join(dir, "dist"), 0755))

	t.Run("expected dir can be found", func(t *testing.T) {
		out, err := pathfinder.FindBuildDir(dir, pathfinder.DefaultBuildPaths)
		assert.NoError(t, err)
		t.Logf("dir: %s", out)
		assert.NotEmpty(t, out)
	})
	t.Run("dir not found", func(t *testing.T) {
		out, err := pathfinder.FindBuildDir(dir, []string{"export"})
		assert.ErrorIs(t, err, pathfinder.ErrDirNotFound)
		assert.Empty(t, out)
	})
}
