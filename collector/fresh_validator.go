package collector

import (
	"errors"
	"path/filepath"

	"github.com/slipfre/imgmd/collectable"
	"github.com/slipfre/imgmd/provider"
	"github.com/slipfre/imgmd/utils"
)

// LocalFileFreshValidator Validate whether the local file is up to date
func LocalFileFreshValidator(cf collectable.FileOperator, base, objectKey string) (bool, error) {
	targetPath := filepath.Join(base, objectKey)
	if utils.IsFileExist(targetPath) {
		updatedTime, _ := utils.GetUpdatedTime(targetPath)
		isUpdatedSince, err := cf.IsUpdatedSince(updatedTime)
		if err != nil {
			return false, err
		}
		if !isUpdatedSince {
			// No need for collecting
			return false, nil
		}
	}
	return true, nil
}

// GetOBSFileFreshValidator Get a FreshValidator which validates whether the obs file is up to date
func GetOBSFileFreshValidator(bucket provider.Bucket) (FreshValidator, error) {
	if bucket == nil {
		return nil, errors.New("bucket should not be nil")
	}
	return func(cf collectable.FileOperator, base, objectKey string) (bool, error) {
		exist, err := bucket.IsObjectExist(objectKey)
		if err != nil {
			return false, err
		}
		if exist {
			updatedTime, err := bucket.GetObjectLastModified(objectKey)
			if err != nil {
				return false, err
			}
			isUpdatedSince, err := cf.IsUpdatedSince(updatedTime)
			if err != nil {
				return false, err
			}
			if !isUpdatedSince {
				return false, nil
			}
		}
		return true, nil
	}, nil
}
