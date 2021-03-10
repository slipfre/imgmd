package collectable

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLeafFile(t *testing.T) {
	var collectableFile FileOperator

	collectableFile = NewLeafFile("", TestImgPath)
	require.Nil(t, collectableFile.FileError())
	require.Equal(t, "", collectableFile.GetParent())
	require.Equal(t, Leaf, collectableFile.GetFileType())
	absTestImgPath, _ := filepath.Abs(TestImgPath)
	require.Equal(t, absTestImgPath, collectableFile.GetURI())

	dependencies, err := collectableFile.FindDependencies()
	require.Nil(t, err)
	require.NotNil(t, dependencies)
	require.Equal(t, 0, len(dependencies))

	err = collectableFile.To(TestImgTargetPath)
	require.Nil(t, err)

	_, err = os.Stat(TestImgTargetPath)
	require.Nil(t, err)

	err = os.Remove(TestImgTargetPath)
	require.Nil(t, err)
}
