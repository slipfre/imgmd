package internal

import (
	"os"
	"testing"
)

const (
	testImgPath       string = "../resources/icon_markdown.png"
	testImgTargetPath string = "../resources/icon_markdown_test.png"
	testMDPath        string = "../resources/u_good_i_good.md"
	testMDTargetPath  string = "../resources/u_good_i_good_test.md"
	testHTTPSrcURI    string = "http://testsdaf.oss-cn-hangzhou.aliyuncs.com/griddle.png"
)

func TestMain(m *testing.M) {

	// 本地获取环境变量为空不执行穿透测试
	if testImgPath != "" && testImgTargetPath != "" && testHTTPSrcURI != "" {
		code := m.Run()
		os.Exit(code)
	}
}
