package utils

import (
	"os"
	"testing"
)

const (
	ProjectRoot       string = "C:/Users/Happy/Desktop/code/imgmd/"
	TestImgPath       string = ProjectRoot + "resources/icon_markdown.png"
	TestImgTargetPath string = ProjectRoot + "resources/icon_markdown_test.png"
	TestMDPath        string = ProjectRoot + "resources/u_good_i_good.md"
	TestMDTargetPath  string = ProjectRoot + "resources/u_good_i_good_test.md"
	TestHTTPSrcURI    string = "http://testsdaf.oss-cn-hangzhou.aliyuncs.com/griddle.png"
)

var (
	TestAKID       string
	TestAKS        string
	TestEndpoint   string
	TestBucketName string
)

func TestMain(m *testing.M) {
	TestAKID = os.Getenv("AKID")
	TestAKS = os.Getenv("AKS")
	TestEndpoint = os.Getenv("Endpoint")
	TestBucketName = os.Getenv("BucketName")

	// 本地获取环境变量为空不执行穿透测试
	if TestImgPath != "" && TestImgTargetPath != "" && TestHTTPSrcURI != "" {
		code := m.Run()
		os.Exit(code)
	}
}
