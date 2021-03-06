package config

import (
	"testing"

	"os"

	"io/ioutil"

	"github.com/drud/ddev/pkg/cms/model"
	"github.com/drud/ddev/pkg/testcommon"
	"github.com/stretchr/testify/assert"
)

func TestWriteDrupalConfig(t *testing.T) {
	dir := testcommon.CreateTmpDir("example")

	file, err := ioutil.TempFile(dir, "file")
	assert.NoError(t, err)

	err = os.Chmod(file.Name(), 0444)
	assert.NoError(t, err)

	drupalConfig := model.NewDrupalConfig()
	err = WriteDrupalConfig(drupalConfig, file.Name())
	assert.NoError(t, err)

	err = os.Chmod(dir, 0755)
	assert.NoError(t, err)
	err = os.Chmod(file.Name(), 0666)
	assert.NoError(t, err)

	err = os.RemoveAll(dir)
	assert.NoError(t, err)
}

func TestWriteDrushConfig(t *testing.T) {
	dir := testcommon.CreateTmpDir("example")

	file, err := ioutil.TempFile(dir, "file")
	assert.NoError(t, err)

	err = os.Chmod(file.Name(), 0444)
	assert.NoError(t, err)

	drushConfig := model.NewDrushConfig()
	err = WriteDrushConfig(drushConfig, file.Name())
	assert.NoError(t, err)

	err = os.Chmod(dir, 0755)
	assert.NoError(t, err)
	err = os.Chmod(file.Name(), 0666)
	assert.NoError(t, err)

	err = os.RemoveAll(dir)
	assert.NoError(t, err)
}

func TestWriteWordpressConfig(t *testing.T) {
	dir := testcommon.CreateTmpDir("example")

	file, err := ioutil.TempFile(dir, "file")
	assert.NoError(t, err)

	err = os.Chmod(file.Name(), 0444)
	assert.NoError(t, err)

	wpConfig := model.NewWordpressConfig()
	err = WriteWordpressConfig(wpConfig, file.Name())
	assert.NoError(t, err)

	err = os.Chmod(dir, 0755)
	assert.NoError(t, err)
	err = os.Chmod(file.Name(), 0666)
	assert.NoError(t, err)

	err = os.RemoveAll(dir)
	assert.NoError(t, err)
}
