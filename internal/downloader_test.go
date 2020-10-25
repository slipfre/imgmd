package internal

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDownloader_downloadLocalFile(t *testing.T) {
	err := DownloadFile(testImgPath, testImgTargetPath)
	require.Nil(t, err)

	_, err = os.Stat(testImgTargetPath)
	require.Nil(t, err)

	err = os.Remove(testImgTargetPath)
	require.Nil(t, err)
}

func TestDownloader_downloadHTTPHTTPSFile(t *testing.T) {
	err := DownloadFile(testHTTPSrcURI, testImgTargetPath)
	require.Nil(t, err)

	_, err = os.Stat(testImgTargetPath)
	require.Nil(t, err)

	err = os.Remove(testImgTargetPath)
	require.Nil(t, err)
}
