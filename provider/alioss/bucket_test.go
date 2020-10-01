package alioss

import (
	"testing"

	"github.com/slipfre/imgmd/provider"
	"github.com/stretchr/testify/require"
)

func getBucket(client provider.Client, bucketName string) (bucket provider.Bucket, err error) {
	bucket, err = client.GetOrCreateBucket(
		bucketName,
		provider.WithACL(provider.PublicRead),
		provider.WithRedundancyType(provider.LRS),
		provider.WithStorage(provider.Standard),
	)
	return
}

func cleanBucket(client provider.Client, bucketName string) (err error) {
	err = client.DeleteBucket(bucketName)
	return
}

func TestBucket_PutAndRemoveObject(t *testing.T) {
	testObjectKeyName := "test/test-bucket-put-and-remove-object-object"
	testBucketName := "test-bucket-put-and-remove-object-bucket"
	localFilePath := "C:\\Users\\Happy\\Desktop\\Resources\\griddle.png"

	client, err := getClient()
	require.Nil(t, err)

	bucket, err := getBucket(client, testBucketName)
	require.Nil(t, err)

	defer cleanBucket(client, testBucketName)

	isExist, err := bucket.IsObjectExist(testObjectKeyName)
	require.Nil(t, err)
	require.Falsef(t, isExist, "Object %s should not exist", testObjectKeyName)

	err = bucket.PutObjectFromFile(
		testObjectKeyName,
		localFilePath,
		provider.WithACL(provider.PublicRead),
		provider.WithRedundancyType(provider.LRS),
		provider.WithStorage(provider.Standard),
	)
	require.Nil(t, err)

	isExist, err = bucket.IsObjectExist(testObjectKeyName)
	require.Nil(t, err)
	require.Truef(t, isExist, "Objcet %s should exist", testObjectKeyName)

	err = bucket.DeleteObject(testObjectKeyName)
	require.Nil(t, err)

	isExist, err = bucket.IsObjectExist(testObjectKeyName)
	require.Nil(t, err)
	require.Falsef(t, isExist, "Object %s should not exist", testObjectKeyName)

	return
}
