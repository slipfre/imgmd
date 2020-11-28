package internal

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/slipfre/imgmd/provider"
	"github.com/slipfre/imgmd/provider/alioss"
	"github.com/stretchr/testify/require"
)

func TestAsyncCollector_testCollect(t *testing.T) {
	md := NewMarkdownFile("", testMDPath)
	require.Nil(t, md.FileError())

	mdCollector, err := NewLocalAsyncCollector(
		md,
		filepath.Dir(testMDTargetPath),
		filepath.Base(testMDTargetPath),
		LocalURIMapper,
	)
	require.Nil(t, err)

	select {
	case err := <-mdCollector.Collect(context.Background()):
		require.Nil(t, err)
	}

	err = os.Remove(testMDTargetPath)
	require.Nil(t, err)

	suffix := filepath.Ext(testMDTargetPath)
	dirName := strings.TrimSuffix(testMDTargetPath, suffix) + "_medias"
	err = os.RemoveAll(dirName)
	require.Nil(t, err)
}

func getBucket(bucketName string) (provider.Bucket, error) {
	client, err := alioss.NewClient(TestEndpoint, TestAKID, TestAKS)
	if err != nil {
		return nil, err
	}
	bucket, err := client.GetOrCreateBucket(
		bucketName,
		provider.WithACL(provider.PublicRead),
		provider.WithRedundancyType(provider.LRS),
		provider.WithStorage(provider.Standard),
	)
	if err != nil {
		return nil, err
	}
	return bucket, nil
}

func TestOBSAsyncCollector_testCollectMDWithMediasToOBS(t *testing.T) {
	testBucketName := "test-bucket-collect-md-with-medias-to-obs"
	bucket, err := getBucket(testBucketName)
	require.Nil(t, err)

	md := NewMarkdownFile("", testMDPath)
	require.Nil(t, md.FileError())

	obsCollectorGenerator := GetOBSCollectorGenerator(bucket)
	obsURIMapper, err := GetOBSURIMapper(bucket)
	require.Nil(t, err)
	mdCollector, err := NewLocalAsyncCollector(
		md,
		filepath.Dir(testMDTargetPath),
		filepath.Base(testMDTargetPath),
		obsURIMapper,
		WithDependencyCollectorGenerator(obsCollectorGenerator),
	)
	require.Nil(t, err)

	select {
	case err := <-mdCollector.Collect(context.Background()):
		require.Nil(t, err)
	}

	err = os.Remove(testMDTargetPath)
	require.Nil(t, err)

	objectKeyPrefix := "u_good_i_good_test_medias/img%d.png"
	for i := 1; i <= 13; i++ {
		err := bucket.DeleteObject(fmt.Sprintf(objectKeyPrefix, i))
		require.Nil(t, err)
	}
}
