package alioss

import (
	"testing"

	"github.com/slipfre/imgmd/provider"
)

func getBucket() (bucket provider.Bucket, err error) {
	client, err := getClient()
	if err != nil {
		return
	}

	testBucketBucketName := "test_bucket_bucket_name"
	bucket, err = client.GetOrCreateBucket(
		testBucketBucketName,
		provider.WithACL(provider.PublicRead),
		provider.WithRedundancyType(provider.LRS),
		provider.WithStorage(provider.Standard),
	)
	return
}

func TestBucket_PutAndRemoveObject(t *testing.T) {
	return
}
