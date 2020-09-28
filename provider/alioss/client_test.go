package alioss

import (
	"testing"

	"github.com/slipfre/imgmd/provider"
	"github.com/stretchr/testify/require"
)

func getClient() (client provider.Client, err error) {
	client, err = NewClient(testEndpoint, testAKID, testAKS)
	return
}

func TestClient_CreateBucketAndIsExistAndDelete(t *testing.T) {
	client, err := getClient()
	require.Nil(t, err)

	testClientBucketName := "test_oss_client"

	isExist, err := client.IsBucketExist(testClientBucketName)
	require.Nil(t, err)
	require.Falsef(t, isExist, "Bucket %s should not exist", testClientBucketName)

	err = client.CreateBucket(
		testClientBucketName,
		provider.WithACL(provider.PublicRead),
		provider.WithRedundancyType(provider.LRS),
		provider.WithStorage(provider.Standard),
	)
	require.Nil(t, err)
	isExist, err = client.IsBucketExist(testClientBucketName)
	require.Nil(t, err)
	require.Truef(t, isExist, "Bucket %s should exist", testClientBucketName)

	err = client.CreateBucket(
		testClientBucketName,
		provider.WithACL(provider.PublicRead),
		provider.WithRedundancyType(provider.LRS),
		provider.WithStorage(provider.Standard),
	)
	require.NotNil(t, err)

	err = client.DeleteBucket(testClientBucketName)
	require.Nil(t, err)

	isExist, err = client.IsBucketExist(testClientBucketName)
	require.Nil(t, err)
	require.Falsef(t, isExist, "Bucket %s should not exist", testClientBucketName)
}

func TestClient_GetOrCreateBucket(t *testing.T) {
	client, err := getClient()
	require.Nil(t, err)

	testClientGetOrCreateBucket := "test_client_get_or_create_bucket"

	isExist, err := client.IsBucketExist(testClientGetOrCreateBucket)
	require.Nil(t, isExist)
	require.Falsef(t, isExist, "Bucket %s should not exist", testClientGetOrCreateBucket)

	_, err = client.GetOrCreateBucket(
		testClientGetOrCreateBucket,
		provider.WithACL(provider.PublicRead),
		provider.WithStorage(provider.Standard),
		provider.WithRedundancyType(provider.LRS),
	)
	require.Nil(t, err)
	_, err = client.GetOrCreateBucket(
		testClientGetOrCreateBucket,
		provider.WithACL(provider.PublicRead),
		provider.WithStorage(provider.Standard),
		provider.WithRedundancyType(provider.LRS),
	)
	require.NotNil(t, err)

	isExist, err = client.IsBucketExist(testClientGetOrCreateBucket)
	require.Nil(t, err)
	require.Truef(t, isExist, "Bucket %s should exist", testClientGetOrCreateBucket)

	err = client.DeleteBucket(testClientGetOrCreateBucket)
	require.Nil(t, err)

	isExist, err = client.IsBucketExist(testClientGetOrCreateBucket)
	require.Nil(t, err)
	require.Falsef(t, isExist, "Bucket %s should not exist", testClientGetOrCreateBucket)
}
