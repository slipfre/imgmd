package collector

import (
	"errors"
	"path/filepath"

	"github.com/slipfre/imgmd/collectable"
	"github.com/slipfre/imgmd/provider"
)

// LocalMover Make files to specified local place
func LocalMover(cf collectable.FileOperator, base, objectKey string) error {
	return cf.To(filepath.Join(base, objectKey))
}

// GetOBSMover Return a OBSMover which make files to OBS with specified object key
func GetOBSMover(bucket provider.Bucket) (Mover, error) {
	if bucket == nil {
		return nil, errors.New("bucket should not be nil")
	}
	return func(cf collectable.FileOperator, base, objectKey string) error {
		return cf.ToOBS(bucket, objectKey)
	}, nil
}
