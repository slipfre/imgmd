package collectable

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMarkdownFile(t *testing.T) {
	var collectableFile FileOperator

	collectableFile = NewMarkdownFile("", TestMDPath)
	require.Nil(t, collectableFile.FileError())
	require.Equal(t, "", collectableFile.GetParent())
	require.Equal(t, Markdown, collectableFile.GetFileType())

	expectMDURI, err := filepath.Abs(TestMDPath)
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

	err = collectableFile.To(TestMDTargetPath)
	require.Nil(t, err)

	_, err = os.Stat(TestMDTargetPath)
	require.Nil(t, err)

	for i := 1; i <= 13; i++ {
		err = os.Remove(fmt.Sprintf(newDependencyURIPrefix, i))
		require.Nil(t, err)
	}

	err = os.Remove(TestMDTargetPath)
	require.Nil(t, err)
}
