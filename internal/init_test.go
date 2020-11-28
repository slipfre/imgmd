package internal

import (
	"os"
	"testing"
)

const (
	testImgPath       string = "C:/Users/Happy/Desktop/code/imgmd/resources/icon_markdown.png"
	testImgTargetPath string = "C:/Users/Happy/Desktop/code/imgmd/resources/icon_markdown_test.png"
	testMDPath        string = "C:/Users/Happy/Desktop/code/imgmd/resources/u_good_i_good.md"
	testMDTargetPath  string = "C:/Users/Happy/Desktop/code/imgmd/resources/u_good_i_good_test.md"
	testHTTPSrcURI    string = "http://testsdaf.oss-cn-hangzhou.aliyuncs.com/griddle.png"
)

var (
	TestAKID       string
	TestAKS        string
	TestEndpoint   string
	TestBucketName string
	TestImgPath    string
)

func TestMain(m *testing.M) {
	TestAKID = os.Getenv("AKID")
	TestAKS = os.Getenv("AKS")
	TestEndpoint = os.Getenv("Endpoint")
	TestBucketName = os.Getenv("BucketName")

	// 本地获取环境变量为空不执行穿透测试
	if testImgPath != "" && testImgTargetPath != "" && testHTTPSrcURI != "" {
		code := m.Run()
		os.Exit(code)
	}
}
