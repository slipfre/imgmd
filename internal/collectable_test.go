package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLeafFile(t *testing.T) {
	var collectableFile CollectableFileOperator

	collectableFile = NewLeafFile("", testImgPath)
	require.Nil(t, collectableFile.FileError())
	require.Equal(t, "", collectableFile.GetParent())
	require.Equal(t, Standalone, collectableFile.GetFileType())
	absTestImgPath, _ := filepath.Abs(testImgPath)
	require.Equal(t, absTestImgPath, collectableFile.GetURI())

	dependencies, err := collectableFile.FindDependencies()
	require.Nil(t, err)
	require.NotNil(t, dependencies)
	require.Equal(t, 0, len(dependencies))

	err = collectableFile.To(testImgTargetPath)
	require.Nil(t, err)

	_, err = os.Stat(testImgTargetPath)
	require.Nil(t, err)

	err = os.Remove(testImgTargetPath)
	require.Nil(t, err)
}

func TestMarkdownFile(t *testing.T) {
	var collectableFile CollectableFileOperator

	collectableFile = NewMarkdownFile("", testMDPath)
	require.Nil(t, collectableFile.FileError())
	require.Equal(t, "", collectableFile.GetParent())
	require.Equal(t, Markdown, collectableFile.GetFileType())

	expectMDURI, err := filepath.Abs(testMDPath)
	require.Nil(t, err)
	require.Equal(t, expectMDURI, collectableFile.GetURI())

	dependencies, err := collectableFile.FindDependencies()
	require.Nil(t, err)
	require.Equal(t, 13, len(dependencies))

	dependencyURIPrefix := "../resources/u_good_i_good_imgs/img%d.png"
	tempMDImgURIPrefix := "../resources/u_good_i_good_imgs/temp_%s"
	for i, dependency := range dependencies {
		require.Equal(t, expectMDURI, dependency.GetParent())
		require.Equal(t, Standalone, dependency.GetFileType())

		expectDependencyURI, err := filepath.Abs(fmt.Sprintf(dependencyURIPrefix, i+1))
		require.Nil(t, err)
		require.Equal(t, expectDependencyURI, dependency.GetURI())

		toPath := fmt.Sprintf(
			tempMDImgURIPrefix,
			filepath.Base(string(dependency.GetURI())),
		)
		err = dependency.To(toPath)
		require.Nil(t, err)
	}

	collectableFile.ReplaceDependencyURIs("", "", func(fileType FileType, uri []byte, base, objectKey string) []byte {
		filename := filepath.Base(string(uri))
		newReferencePath := fmt.Sprintf("u_good_i_good_imgs/temp_%s", filename)
		return []byte(newReferencePath)
	})

	newDependencyURIPrefix := "../resources/u_good_i_good_imgs/temp_img%d.png"
	dependencies, err = collectableFile.FindDependencies()
	require.Nil(t, err)
	for i, dependency := range dependencies {
		require.Equal(t, expectMDURI, dependency.GetParent())
		require.Equal(t, Standalone, dependency.GetFileType())

		expectDependencyURI, err := filepath.Abs(fmt.Sprintf(newDependencyURIPrefix, i+1))
		require.Nil(t, err)
		require.Equal(t, expectDependencyURI, dependency.GetURI())
	}

	err = collectableFile.To(testMDTargetPath)
	require.Nil(t, err)

	_, err = os.Stat(testMDTargetPath)
	require.Nil(t, err)

	for i := 1; i <= 13; i++ {
		err = os.Remove(fmt.Sprintf(newDependencyURIPrefix, i))
		require.Nil(t, err)
	}

	err = os.Remove(testMDTargetPath)
	require.Nil(t, err)
}
