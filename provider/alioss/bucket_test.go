package alioss

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

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

type PutObjectWrapper func(bucket provider.Bucket, objectKey string, options ...provider.ObjectOption) (url string, err error)

func testPutObject(t *testing.T, put PutObjectWrapper) {
	testObjectKeyName := "test/test-bucket-put-and-remove-object-object"
	testBucketName := "test-bucket-put-and-remove-object-bucket"

	client, err := getClient()
	require.Nil(t, err)

	bucket, err := getBucket(client, testBucketName)
	require.Nil(t, err)

	defer cleanBucket(client, testBucketName)

	isExist, err := bucket.IsObjectExist(testObjectKeyName)
	require.Nil(t, err)
	require.Falsef(t, isExist, "Object %s should not exist", testObjectKeyName)

	url, err := put(
		bucket,
		testObjectKeyName,
		provider.WithACL(provider.PublicRead),
		provider.WithRedundancyType(provider.LRS),
		provider.WithStorage(provider.Standard),
	)
	require.Nil(t, err)
	require.NotNil(t, url)

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

func TestBucket_PutAndRemoveObjectFromFile(t *testing.T) {
	localFilePath := TestImgPath
	putObjectFromFileWrapper := func(bucket provider.Bucket, objectKey string, options ...provider.ObjectOption) (url string, err error) {
		return bucket.PutObjectFromFile(objectKey, localFilePath, options...)
	}
	testPutObject(t, putObjectFromFileWrapper)
	return
}

func TestBucket_PutAndRemoveObjectFromBytes(t *testing.T) {
	reader, err := os.Open(TestImgPath)
	require.Nil(t, err)
	reader.Close()

	bytes, err := ioutil.ReadAll(reader)
	require.Nil(t, err)

	putObjectFromBytesWrapper := func(bucket provider.Bucket, objectKey string, options ...provider.ObjectOption) (url string, err error) {
		return bucket.PutObjectFromBytes(objectKey, bytes, options...)
	}
	testPutObject(t, putObjectFromBytesWrapper)
	return
}

func TestBuket_GetObjectLastModifiedTime(t *testing.T) {
	testObjectKeyName := "test/test_object_last_modified_time_object"
	testBucketName := "test-object-last-modified-time-bucket"
	localFilePath := TestImgPath

	client, err := getClient()
	require.Nil(t, err)

	bucket, err := getBucket(client, testBucketName)
	require.Nil(t, err)

	defer cleanBucket(client, testBucketName)

	isExist, err := bucket.IsObjectExist(testObjectKeyName)
	require.Nil(t, err)

	if isExist {
		err = bucket.DeleteObject(testObjectKeyName)
		require.Nil(t, err)
	}

	beforePutObject := time.Now()
	_, err = bucket.PutObjectFromFile(
		testObjectKeyName,
		localFilePath,
		provider.WithACL(provider.PublicRead),
		provider.WithRedundancyType(provider.LRS),
		provider.WithStorage(provider.Standard),
	)
	require.Nil(t, err)

	lastModifiedTime, err := bucket.GetObjectLastModified(testObjectKeyName)
	require.Nil(t, err)
	require.NotNil(t, lastModifiedTime)
	beforePutObject.Before(*lastModifiedTime)

	err = bucket.DeleteObject(testObjectKeyName)
	require.Nil(t, err)
}
