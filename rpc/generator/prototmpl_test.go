package generator

import (
	"github.com/zhangxiaohan228/zero-admin-goctl/util/pathx"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProtoTmpl(t *testing.T) {
	_ = Clean()
	// exists dir
	err := ProtoTmpl(pathx.MustTempDir())
	assert.Nil(t, err)

	// not exist dir
	dir := filepath.Join(pathx.MustTempDir(), "test")
	err = ProtoTmpl(dir)
	assert.Nil(t, err)
}
