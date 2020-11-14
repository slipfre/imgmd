package internal

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAsyncCollector_testCollect(t *testing.T) {
	md := NewMarkdownFile("", testMDPath)
	require.Nil(t, md.FileError())

	mdCollector := NewAsyncCollector(md, testMDTargetPath)

	select {
	case err := <-mdCollector.Collect(context.Background()):
		require.Nil(t, err)
	}

	err := os.Remove(testMDTargetPath)
	require.Nil(t, err)

	suffix := filepath.Ext(testMDTargetPath)
	dirName := strings.TrimSuffix(testMDTargetPath, suffix) + "_medias"
	err = os.RemoveAll(dirName)
	require.Nil(t, err)
}
