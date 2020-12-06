package collector

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/slipfre/imgmd/collectable"
	"github.com/slipfre/imgmd/provider"
	"github.com/slipfre/imgmd/provider/alioss"
	"github.com/stretchr/testify/require"
)

func TestAsyncCollector_testCollect(t *testing.T) {
	md := collectable.NewMarkdownFile("", TestMDPath)
	require.Nil(t, md.FileError())

	mdCollector, err := LocalCollectorGenerator(
		md,
		filepath.Dir(TestMDTargetPath),
		filepath.Base(TestMDTargetPath),
		LocalCollectorGenerator,
	)
	require.Nil(t, err)

	select {
	case err := <-mdCollector.Collect(context.Background()):
		require.Nil(t, err)
	}

	err = os.Remove(TestMDTargetPath)
	require.Nil(t, err)

	suffix := filepath.Ext(TestMDTargetPath)
	dirName := strings.TrimSuffix(TestMDTargetPath, suffix) + "_medias"
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

	md := collectable.NewMarkdownFile("", TestMDPath)
	require.Nil(t, md.FileError())

	require.Nil(t, err)
	mdCollector, err := LocalCollectorGenerator(
		md,
		filepath.Dir(TestMDTargetPath),
		filepath.Base(TestMDTargetPath),
		GetOBSCollectorGenerator(bucket),
	)
	require.Nil(t, err)

	select {
	case err := <-mdCollector.Collect(context.Background()):
		require.Nil(t, err)
	}

	err = os.Remove(TestMDTargetPath)
	require.Nil(t, err)

	objectKeyPrefix := "u_good_i_good_test_medias/img%d.png"
	for i := 1; i <= 13; i++ {
		err := bucket.DeleteObject(fmt.Sprintf(objectKeyPrefix, i))
		require.Nil(t, err)
	}
}
