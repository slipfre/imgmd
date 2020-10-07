package media

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testLocalSrcURI string
	testTargetPath  string
	testHTTPSrcURI  string
)

func TestMain(m *testing.M) {
	testLocalSrcURI = "C:\\Users\\Happy\\Desktop\\Resources\\griddle.png"
	testHTTPSrcURI = "http://testsdaf.oss-cn-hangzhou.aliyuncs.com/griddle.png"
	testTargetPath = "C:\\Users\\Happy\\Desktop\\code\\imgmd\\griddle.pdf"

	if testLocalSrcURI != "" && testTargetPath != "" {
		code := m.Run()
		os.Exit(code)
	}
}

func TestDownloader_downloadLocalMedia(t *testing.T) {
	err := DownloadMedia(testLocalSrcURI, testTargetPath)
	require.Nil(t, err)

	_, err = os.Stat(testTargetPath)
	require.Nil(t, err)

	err = os.Remove(testTargetPath)
	require.Nil(t, err)
}

func TestDownloader_downloadHTTPHTTPSMedia(t *testing.T) {
	err := DownloadMedia(testHTTPSrcURI, testTargetPath)
	require.Nil(t, err)

	_, err = os.Stat(testTargetPath)
	require.Nil(t, err)

	err = os.Remove(testTargetPath)
	require.Nil(t, err)
}
