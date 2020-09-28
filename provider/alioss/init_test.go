package alioss

import (
	"os"
	"testing"
)

var (
	testAKID       string
	testAKS        string
	testEndpoint   string
	testBucketName string
)

func TestMain(m *testing.M) {
	testAKID = os.Getenv("AKID")
	testAKS = os.Getenv("AKS")
	testEndpoint = os.Getenv("Endpoint")
	testBucketName = os.Getenv("BucketName")

	// 本地获取环境变量为空不执行穿透测试
	if testAKID != "" {
		code := m.Run()
		os.Exit(code)
	}
}
