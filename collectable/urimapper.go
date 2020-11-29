package collectable

import (
	"errors"
	"path/filepath"

	"github.com/slipfre/imgmd/provider"
	"github.com/slipfre/imgmd/utils"
)

// LocalURIMapper Map the uri to 'targetDirPath/filename'
func LocalURIMapper(fileType FileType, uri []byte, base, objectKey string) []byte {
	destDirPath := utils.GetTargetResourcesDirPath(filepath.Join(base, objectKey))
	dirName := filepath.Base(destDirPath)
	fileName := filepath.Base(string(uri))
	newReferencePath := filepath.Join(dirName, fileName)
	return []byte(newReferencePath)
}

// GetOBSURIMapper Returns a OBSURIMapper which maps the uri to corresponding
// object under the bucket
func GetOBSURIMapper(bucket provider.Bucket) (URIMapper, error) {
	if bucket == nil {
		return nil, errors.New("bucket should not be nil")
	}
	return func(fileType FileType, originURI []byte, base, objectKey string) []byte {
		depObjDir := utils.GetTargetResourcesDirPath(objectKey)
		depObjKey := filepath.Join(depObjDir, filepath.Base(string(originURI)))
		return []byte(bucket.GetObjectURL(filepath.ToSlash(depObjKey)))
	}, nil
}
