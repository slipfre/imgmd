package internal

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// DownloadFile 根据 srcURI 下载 media 文件到 target 指示的位置
func DownloadFile(srcURI, target string) (err error) {
	in, err := NewFileReader(srcURI)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return
}

// NewFileReader 根据 srcURI 的类型，创建 reader
func NewFileReader(srcURI string) (reader io.ReadCloser, err error) {
	if strings.HasPrefix(srcURI, "http://") || strings.HasPrefix(srcURI, "https://") {
		return NewHTTPHTTPSFileReader(srcURI)
	}
	return NewLocalFileReader(srcURI)
}

// NewHTTPHTTPSFileReader 创建 HTTP/HTTPS 类型的 MediaReader，从网络中读取 media 文件
func NewHTTPHTTPSFileReader(srcURI string) (reader io.ReadCloser, err error) {
	resp, err := http.Get(srcURI)

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", resp.Status)
		return
	}

	reader = resp.Body
	return
}

// NewLocalFileReader 创建本地文件类型的 MediaReader，从本地读取 media 文件
func NewLocalFileReader(srcURI string) (reader io.ReadCloser, err error) {
	reader, err = os.Open(srcURI)
	return
}
