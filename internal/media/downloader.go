package media

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// DownloadMedia 根据 srcURI 下载 media 文件到 target 指示的位置
func DownloadMedia(srcURI, target string) (err error) {
	in, err := NewMediaReader(srcURI)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(target)
	if err != nil {
		return
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return
}

// NewMediaReader 根据 srcURI 的类型，创建 reader
func NewMediaReader(srcURI string) (reader io.ReadCloser, err error) {
	if strings.HasPrefix(srcURI, "http://") || strings.HasPrefix(srcURI, "https://") {
		return NewHTTPHTTPSMediaReader(srcURI)
	}
	return NewLocalMediaReader(srcURI)
}

// NewHTTPHTTPSMediaReader 创建 HTTP/HTTPS 类型的 MediaReader，从网络中读取 media 文件
func NewHTTPHTTPSMediaReader(srcURI string) (reader io.ReadCloser, err error) {
	resp, err := http.Get(srcURI)

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", resp.Status)
		return
	}

	reader = resp.Body
	return
}

// NewLocalMediaReader 创建本地文件类型的 MediaReader，从本地读取 media 文件
func NewLocalMediaReader(srcURI string) (reader io.ReadCloser, err error) {
	reader, err = os.Open(srcURI)
	return
}
