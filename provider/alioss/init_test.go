package alioss

import (
	"os"
	"testing"
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
	TestImgPath = "../../resources/icon_markdown.png"

	// 本地获取环境变量为空不执行穿透测试
	if TestAKID != "" {
		code := m.Run()
		os.Exit(code)
	}
}
