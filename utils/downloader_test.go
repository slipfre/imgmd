package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDownloader_downloadLocalFile(t *testing.T) {
	err := DownloadFile(TestImgPath, TestImgTargetPath)
	require.Nil(t, err)

	_, err = os.Stat(TestImgTargetPath)
	require.Nil(t, err)

	err = os.Remove(TestImgTargetPath)
	require.Nil(t, err)
}

func TestDownloader_downloadHTTPHTTPSFile(t *testing.T) {
	err := DownloadFile(TestHTTPSrcURI, TestImgTargetPath)
	require.Nil(t, err)

	_, err = os.Stat(TestImgTargetPath)
	require.Nil(t, err)

	err = os.Remove(TestImgTargetPath)
	require.Nil(t, err)
}
