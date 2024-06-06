package generate

import (
	"github.com/zhangxiaohan228/zero-admin-goctl/util/pathx"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemplate(t *testing.T) {
	tempDir := t.TempDir()
	pathx.RegisterGoctlHome(tempDir)
	t.Cleanup(func() {
		pathx.RegisterGoctlHome("")
	})

	t.Run("Category", func(t *testing.T) {
		assert.Equal(t, category, Category())
	})

	t.Run("Clean", func(t *testing.T) {
		err := Clean()
		assert.NoError(t, err)
	})

	t.Run("Templates", func(t *testing.T) {
		err := Templates()
		assert.NoError(t, err)
		assert.True(t, pathx.FileExists(filepath.Join(tempDir, category, modelTemplateFile)))
	})

	t.Run("RevertTemplate", func(t *testing.T) {
		assert.NoError(t, RevertTemplate(modelTemplateFile))
		assert.Error(t, RevertTemplate("foo"))
	})

	t.Run("Update", func(t *testing.T) {
		assert.NoError(t, Update())
	})
}
